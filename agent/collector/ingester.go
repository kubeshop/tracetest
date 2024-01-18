package collector

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kubeshop/tracetest/agent/event"
	"github.com/kubeshop/tracetest/server/otlp"
	"go.opencensus.io/trace"
	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type stoppable interface {
	Stop()
}

type ingester interface {
	otlp.Ingester
	stoppable

	Statistics() Statistics
}

func newForwardIngester(ctx context.Context, batchTimeout time.Duration, cfg remoteIngesterConfig, startRemoteServer bool) (ingester, error) {
	ingester := &forwardIngester{
		BatchTimeout:   batchTimeout,
		RemoteIngester: cfg,
		buffer:         &buffer{},
		done:           make(chan bool),
		traceCache:     cfg.traceCache,
		logger:         cfg.logger,
	}

	if startRemoteServer {
		err := ingester.connectToRemoteServer(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not connect to remote server: %w", err)
		}

		go ingester.startBatchWorker()
	}

	return ingester, nil
}

type Statistics struct {
	spanCount         int
	lastSpanTimestamp time.Time
}

// forwardIngester forwards all incoming spans to a remote ingester. It also batches those
// spans to reduce network traffic.
type forwardIngester struct {
	BatchTimeout   time.Duration
	RemoteIngester remoteIngesterConfig
	client         pb.TraceServiceClient
	buffer         *buffer
	done           chan bool
	traceCache     TraceCache
	logger         *zap.Logger

	statistics Statistics
}

type remoteIngesterConfig struct {
	URL               string
	Token             string
	traceCache        TraceCache
	startRemoteServer bool
	logger            *zap.Logger
	observer          event.Observer
}

type buffer struct {
	mutex sync.Mutex
	spans []*v1.ResourceSpans
}

func (i *forwardIngester) Statistics() Statistics {
	return i.statistics
}

func (i *forwardIngester) ResetStatistics() {
	i.statistics = Statistics{}
}

func (i *forwardIngester) Ingest(ctx context.Context, request *pb.ExportTraceServiceRequest, requestType otlp.RequestType) (*pb.ExportTraceServiceResponse, error) {
	spanCount := countSpans(request)
	i.buffer.mutex.Lock()

	i.buffer.spans = append(i.buffer.spans, request.ResourceSpans...)
	i.statistics.spanCount += spanCount
	i.statistics.lastSpanTimestamp = time.Now()

	i.buffer.mutex.Unlock()
	i.logger.Debug("received spans", zap.Int("count", spanCount))

	if i.traceCache != nil {
		i.logger.Debug("caching test spans")
		// In case of OTLP datastore, those spans will be polled from this cache instead
		// of a real datastore
		i.cacheTestSpans(request.ResourceSpans)
	}

	return &pb.ExportTraceServiceResponse{
		PartialSuccess: &pb.ExportTracePartialSuccess{
			RejectedSpans: 0,
		},
	}, nil
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

func (i *forwardIngester) connectToRemoteServer(ctx context.Context) error {
	conn, err := grpc.DialContext(ctx, i.RemoteIngester.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		i.logger.Error("could not connect to remote server", zap.Error(err))
		return fmt.Errorf("could not connect to remote server: %w", err)
	}

	i.client = pb.NewTraceServiceClient(conn)
	return nil
}

func (i *forwardIngester) startBatchWorker() {
	i.logger.Debug("starting batch worker", zap.Duration("batch_timeout", i.BatchTimeout))
	ticker := time.NewTicker(i.BatchTimeout)
	done := make(chan bool)
	for {
		select {
		case <-done:
			i.logger.Debug("stopping batch worker")
			return
		case <-ticker.C:
			i.logger.Debug("executing batch")
			err := i.executeBatch(context.Background())
			if err != nil {
				i.logger.Error("could not execute batch", zap.Error(err))
				log.Println(err)
			}
		}
	}
}

func (i *forwardIngester) executeBatch(ctx context.Context) error {
	i.buffer.mutex.Lock()
	newSpans := i.buffer.spans
	i.buffer.spans = []*v1.ResourceSpans{}
	i.buffer.mutex.Unlock()

	if len(newSpans) == 0 {
		i.logger.Debug("no spans to forward")
		return nil
	}

	err := i.forwardSpans(ctx, newSpans)
	if err != nil {
		i.logger.Error("could not forward spans", zap.Error(err))
		return err
	}

	i.logger.Debug("successfully forwarded spans", zap.Int("count", len(newSpans)))
	return nil
}

func (i *forwardIngester) forwardSpans(ctx context.Context, spans []*v1.ResourceSpans) error {
	_, err := i.client.Export(ctx, &pb.ExportTraceServiceRequest{
		ResourceSpans: spans,
	})

	if err != nil {
		i.logger.Error("could not forward spans to remote server", zap.Error(err))
		return fmt.Errorf("could not forward spans to remote server: %w", err)
	}

	return nil
}

func (i *forwardIngester) cacheTestSpans(resourceSpans []*v1.ResourceSpans) {
	i.logger.Debug("caching test spans")
	spans := make(map[string][]*v1.Span)
	for _, resourceSpan := range resourceSpans {
		for _, scopedSpan := range resourceSpan.ScopeSpans {
			for _, span := range scopedSpan.Spans {
				traceID := trace.TraceID(span.TraceId).String()
				spans[traceID] = append(spans[traceID], span)
			}
		}
	}

	i.logger.Debug("caching test spans", zap.Int("count", len(spans)))

	for traceID, spans := range spans {
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
