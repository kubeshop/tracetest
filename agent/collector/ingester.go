package collector

import (
	"context"
	"sync"
	"time"

	"github.com/kubeshop/tracetest/agent/event"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/events"
	"github.com/kubeshop/tracetest/agent/ui/dashboard/sensors"
	"github.com/kubeshop/tracetest/server/otlp"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opencensus.io/trace"
	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	v11 "go.opentelemetry.io/proto/otlp/common/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"go.uber.org/zap"
)

type stoppable interface {
	Stop()
}

type ingester interface {
	otlp.Ingester
	stoppable

	Statistics() Statistics
	ResetStatistics()

	SetSensor(sensors.Sensor)
}

type TraceModeForwarder interface {
	Export(ctx context.Context, request *pb.ExportTraceServiceRequest) error
}

func newForwardIngester(ctx context.Context, batchTimeout time.Duration, cfg remoteIngesterConfig, traceModeForwarder TraceModeForwarder, startRemoteServer bool) (ingester, error) {
	ingester := &forwardIngester{
		BatchTimeout:       batchTimeout,
		RemoteIngester:     cfg,
		traceIDs:           make(map[string]bool, 0),
		done:               make(chan bool),
		traceCache:         cfg.traceCache,
		logger:             cfg.logger,
		sensor:             cfg.sensor,
		traceModeForwarder: traceModeForwarder,
	}

	return ingester, nil
}

type Statistics struct {
	SpanCount         int64
	LastSpanTimestamp time.Time
}

// forwardIngester forwards all incoming spans to a remote ingester. It also batches those
// spans to reduce network traffic.
type forwardIngester struct {
	BatchTimeout       time.Duration
	RemoteIngester     remoteIngesterConfig
	mutex              sync.Mutex
	traceIDs           map[string]bool
	done               chan bool
	traceCache         TraceCache
	logger             *zap.Logger
	sensor             sensors.Sensor
	traceModeForwarder TraceModeForwarder

	statistics Statistics

	sync.Mutex
}

type remoteIngesterConfig struct {
	URL                string
	Token              string
	traceCache         TraceCache
	startRemoteServer  bool
	logger             *zap.Logger
	observer           event.Observer
	sensor             sensors.Sensor
	traceMode          bool
	traceModeForwarder TraceModeForwarder
}

func (i *forwardIngester) Statistics() Statistics {
	return i.statistics
}

func (i *forwardIngester) ResetStatistics() {
	i.statistics = Statistics{}
}

func (i *forwardIngester) SetSensor(sensor sensors.Sensor) {
	i.sensor = sensor
}

func (i *forwardIngester) Ingest(ctx context.Context, request *pb.ExportTraceServiceRequest, requestType otlp.RequestType) (*pb.ExportTraceServiceResponse, error) {
	go i.ingestSpans(request)
	if i.RemoteIngester.traceMode {
		err := i.traceModeForwarder.Export(ctx, request)
		if err != nil {
			i.logger.Error("failed to forward spans to trace mode", zap.Error(err))
		}
	}

	return &pb.ExportTraceServiceResponse{
		PartialSuccess: &pb.ExportTracePartialSuccess{
			RejectedSpans: 0,
		},
	}, nil
}

func (i *forwardIngester) ingestSpans(request *pb.ExportTraceServiceRequest) {
	spanCount := countSpans(request)
	i.mutex.Lock()

	i.statistics.SpanCount += int64(spanCount)
	i.statistics.LastSpanTimestamp = time.Now()
	realSpanCount := i.statistics.SpanCount

	i.mutex.Unlock()

	i.sensor.Emit(events.SpanCountUpdated, realSpanCount)
	i.logger.Debug("received spans", zap.Int("count", spanCount))

	if i.traceCache != nil {
		i.logger.Debug("caching test spans")
		// In case of OTLP datastore, those spans will be polled from this cache instead
		// of a real datastore
		i.cacheTestSpans(request.ResourceSpans)
		i.sensor.Emit(events.TraceCountUpdated, len(i.traceIDs))
	}
}

func countSpans(request *pb.ExportTraceServiceRequest) int {
	count := 0
	for _, resourceSpan := range request.ResourceSpans {
		for _, scopeSpan := range resourceSpan.ScopeSpans {
			count += len(scopeSpan.Spans)
		}
	}

	return count
}

func (i *forwardIngester) cacheTestSpans(resourceSpans []*v1.ResourceSpans) {
	i.logger.Debug("caching test spans")
	spans := make(map[string][]*v1.Span)
	for _, resourceSpan := range resourceSpans {
		for _, scopedSpan := range resourceSpan.ScopeSpans {
			for _, span := range scopedSpan.Spans {
				if scopedSpan.Scope != nil {
					span.Attributes = append(span.Attributes, &v11.KeyValue{
						Key:   traces.MetadataServiceName,
						Value: &v11.AnyValue{Value: &v11.AnyValue_StringValue{StringValue: scopedSpan.Scope.Name}},
					})

					// Add attributes from the resource
					span.Attributes = append(span.Attributes, scopedSpan.Scope.Attributes...)
				}

				// Add attributes from the resource
				if resourceSpan.Resource != nil {
					span.Attributes = append(span.Attributes, resourceSpan.Resource.Attributes...)
				}

				traceID := trace.TraceID(span.TraceId).String()
				spans[traceID] = append(spans[traceID], span)
			}
		}
	}

	i.logger.Debug("caching test spans", zap.Int("count", len(spans)))

	for traceID, spans := range spans {
		i.Lock()
		i.traceIDs[traceID] = true
		i.Unlock()
		if _, ok := i.traceCache.Get(traceID); !ok {
			i.logger.Debug("traceID is not part of a test", zap.String("traceID", traceID))
			// traceID is not part of a test
			continue
		}

		i.RemoteIngester.observer.StartSpanReceive(spans)

		i.traceCache.Append(traceID, spans)
		i.logger.Debug("caching test spans", zap.String("traceID", traceID), zap.Int("count", len(spans)))

		i.RemoteIngester.observer.EndSpanReceive(spans, nil)
	}
}

func (i *forwardIngester) Stop() {
	i.logger.Debug("stopping forward ingester")
	i.done <- true
}
