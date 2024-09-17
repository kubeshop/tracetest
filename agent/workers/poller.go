package workers

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/collector"
	"github.com/kubeshop/tracetest/agent/event"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/telemetry"
	"github.com/kubeshop/tracetest/agent/tracedb"
	"github.com/kubeshop/tracetest/agent/tracedb/connection"
	"github.com/kubeshop/tracetest/agent/workers/envelope"
	"github.com/kubeshop/tracetest/agent/workers/poller"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const maxEnvelopeSize = 3 * 1024 * 1024 // 3MB

type PollerWorker struct {
	client                 *client.Client
	inmemoryDatastore      tracedb.TraceDB
	sentSpansCache         *poller.SentSpansCache
	logger                 *zap.Logger
	observer               event.Observer
	stoppableProcessRunner StoppableProcessRunner
	tracer                 trace.Tracer
	meter                  metric.Meter
	traceCache             collector.TraceCache
}

type PollerOption func(*PollerWorker)

func WithInMemoryDatastore(datastore tracedb.TraceDB) PollerOption {
	return func(pw *PollerWorker) {
		pw.inmemoryDatastore = datastore
	}
}

func WithPollerObserver(observer event.Observer) PollerOption {
	return func(pw *PollerWorker) {
		pw.observer = observer
	}
}

func WithPollerTraceCache(cache collector.TraceCache) PollerOption {
	return func(pw *PollerWorker) {
		pw.traceCache = cache
	}
}

func WithPollerStoppableProcessRunner(stoppableProcessRunner StoppableProcessRunner) PollerOption {
	return func(pw *PollerWorker) {
		pw.stoppableProcessRunner = stoppableProcessRunner
	}
}

func WithPollerLogger(logger *zap.Logger) PollerOption {
	return func(pw *PollerWorker) {
		pw.logger = logger
	}
}

func WithPollerTracer(tracer trace.Tracer) PollerOption {
	return func(pw *PollerWorker) {
		pw.tracer = tracer
	}
}

func WithPollerMeter(meter metric.Meter) PollerOption {
	return func(pw *PollerWorker) {
		pw.meter = meter
	}
}

func NewPollerWorker(client *client.Client, opts ...PollerOption) *PollerWorker {
	pollerWorker := &PollerWorker{
		client:         client,
		logger:         zap.NewNop(),
		sentSpansCache: poller.NewSentSpansCache(),
		observer:       event.NewNopObserver(),
		tracer:         telemetry.GetNoopTracer(),
		meter:          telemetry.GetNoopMeter(),
	}

	for _, opt := range opts {
		opt(pollerWorker)
	}

	return pollerWorker
}

func (w *PollerWorker) Poll(ctx context.Context, request *proto.PollingRequest) error {
	ctx, span := w.tracer.Start(ctx, "PollingRequest Worker operation")
	defer span.End()

	runCounter, _ := w.meter.Int64Counter("tracetest.agent.pollerworker.runs")
	runCounter.Add(ctx, 1)

	errorCounter, _ := w.meter.Int64Counter("tracetest.agent.pollerworker.errors")

	w.logger.Debug("Received polling request", zap.Any("request", request))
	w.observer.StartTracePoll(request)

	var err error
	w.stoppableProcessRunner(ctx, request.TestID, request.RunID, func(ctx context.Context) {
		err = w.poll(ctx, request)
	}, func(cause string) {
		err = executor.ErrUserCancelled
		if cause == string(executor.UserRequestTypeSkipTraceCollection) {
			err = executor.ErrSkipTraceCollection
		}
	})

	if err != nil {
		w.logger.Error("Error polling", zap.Error(err))
		errorResponse := &proto.PollingResponse{
			RequestID:           request.RequestID,
			AgentIdentification: w.client.SessionConfiguration().AgentIdentification,
			TestID:              request.GetTestID(),
			RunID:               request.GetRunID(),
			TraceID:             request.TraceID,
			TraceFound:          false,
			Error: &proto.Error{
				Message: err.Error(),
			},
		}

		w.observer.EndTracePoll(request, err)
		w.logger.Debug("Sending polling error", zap.Any("response", errorResponse))
		sendErr := w.client.SendTrace(ctx, errorResponse)

		if sendErr != nil {
			w.logger.Error("Error sending polling error", zap.Error(sendErr))
			w.observer.Error(sendErr)

			formattedErr := fmt.Errorf("could not report polling error back to the server: %w. Original error: %s", sendErr, err.Error())
			span.RecordError(formattedErr)
			errorCounter.Add(ctx, 1)

			return formattedErr
		}

		span.RecordError(err)
		errorCounter.Add(ctx, 1)
	}

	w.observer.EndTracePoll(request, nil)
	return err
}

func (w *PollerWorker) poll(ctx context.Context, request *proto.PollingRequest) error {
	w.logger.Debug("Received polling request", zap.Any("request", request))
	datastoreConfig, err := convertProtoToDataStore(request.Datastore)
	if err != nil {
		w.logger.Error("Invalid datastore", zap.Error(err))
		return err
	}
	w.logger.Debug("Converted datastore", zap.Any("datastore", datastoreConfig), zap.Any("originalDatastore", request.Datastore))

	if datastoreConfig == nil {
		w.logger.Error("Invalid datastore: nil")
		return fmt.Errorf("invalid datastore: nil")
	}

	dsFactory := tracedb.Factory(nil)
	ds, err := dsFactory(*datastoreConfig)
	if err != nil {
		w.logger.Error("Invalid datastore", zap.Error(err))
		log.Printf("Invalid datastore: %s", err.Error())
		return err
	}
	defer ds.Close()
	w.logger.Debug("Created datastore", zap.Any("datastore", ds), zap.Bool("isOTLPBasedProvider", datastoreConfig.IsOTLPBasedProvider()))

	if datastoreConfig.IsOTLPBasedProvider() && w.inmemoryDatastore != nil {
		w.logger.Debug("Using in-memory datastore")
		ds = w.inmemoryDatastore
	}

	pollingResponse := &proto.PollingResponse{
		RequestID: request.RequestID,
		TestID:    request.TestID,
		RunID:     request.RunID,
		TraceID:   request.TraceID,
	}

	trace, err := ds.GetTraceByID(ctx, request.TraceID)
	if err != nil {
		w.logger.Error("Cannot get trace from datastore", zap.Error(err))
		if !errors.Is(err, connection.ErrTraceNotFound) {
			w.logger.Debug("error was %s was not %s", zap.Error(err), zap.Error(connection.ErrTraceNotFound))
			log.Printf("cannot get trace from datastore: %s", err.Error())
			return err
		}

		w.logger.Debug("Trace not found")
		pollingResponse.TraceFound = false
	} else {
		w.logger.Debug("Trace found")
		pollingResponse.TraceFound = true
		pollingResponse.Spans = convertTraceInToProtoSpans(trace)
		w.logger.Debug("Converted trace", zap.Any("trace", trace), zap.Any("pollingResponse", spew.Sdump(pollingResponse)))

		// remove sent spans
		newSpans := make([]*proto.Span, 0, len(pollingResponse.Spans))
		for _, span := range pollingResponse.Spans {
			runKey := fmt.Sprintf("%d-%s-%s", request.RunID, request.TestID, span.Id)
			w.logger.Debug("Checking if span was already sent", zap.String("runKey", runKey))
			alreadySent := w.sentSpansCache.Get(request.TraceID, runKey)
			if !alreadySent {
				w.logger.Debug("Span was not sent", zap.String("runKey", runKey))
				newSpans = append(newSpans, span)
			} else {
				w.logger.Debug("Span was already sent", zap.String("runKey", runKey))
			}
		}
		pollingResponse.Spans = envelope.EnvelopeSpans(newSpans, maxEnvelopeSize)

		w.logger.Debug("Filtered spans", zap.Any("pollingResponse", spew.Sdump(pollingResponse)))
	}

	err = w.client.SendTrace(ctx, pollingResponse)
	if err != nil {
		w.logger.Error("Cannot send trace to server", zap.Error(err))
		log.Printf("cannot send trace to server: %s", err.Error())
		return err
	}

	spanIDs := make([]string, 0, len(pollingResponse.Spans))
	for _, span := range pollingResponse.Spans {
		spanIDs = append(spanIDs, span.Id)

		// mark span as sent
		runKey := fmt.Sprintf("%d-%s-%s", request.RunID, request.TestID, span.Id)
		w.logger.Debug("Marking span as sent", zap.String("runKey", runKey))
		// TODO: we can set the expiration for this key to be
		// 1 second after the pollingProfile max waiting time
		// but we need to get that info here from controlplane
		w.sentSpansCache.Set(request.TraceID, runKey)
	}

	if w.traceCache != nil {
		w.traceCache.RemoveSpans(request.TraceID, spanIDs)
	}

	return nil
}

func convertProtoToDataStore(r *proto.DataStore) (*datastore.DataStore, error) {
	var ds datastore.DataStore

	if r.Jaeger != nil && r.Jaeger.Grpc != nil {
		ds.Values.Jaeger = &datastore.GRPCClientSettings{}
		deepcopy.DeepCopy(r.Jaeger.Grpc, &ds.Values.Jaeger)
	}

	if r.Tempo != nil {
		ds.Values.Tempo = &datastore.MultiChannelClientConfig{}
		if r.Tempo.Grpc != nil {
			ds.Values.Tempo.Type = datastore.MultiChannelClientTypeGRPC
			ds.Values.Tempo.Grpc = &datastore.GRPCClientSettings{}
			deepcopy.DeepCopy(r.Tempo.Grpc, &ds.Values.Tempo.Grpc)
		}
		if r.Tempo.Http != nil {
			ds.Values.Tempo.Type = datastore.MultiChannelClientTypeHTTP
			ds.Values.Tempo.Http = &datastore.HttpClientConfig{}
			deepcopy.DeepCopy(r.Tempo.Http, &ds.Values.Tempo.Http)
		}
	}

	if r.Opensearch != nil {
		ds.Values.OpenSearch = &datastore.ElasticSearchConfig{}
		deepcopy.DeepCopy(r.Opensearch, &ds.Values.OpenSearch)
	}

	if r.Elasticapm != nil {
		ds.Values.ElasticApm = &datastore.ElasticSearchConfig{}
		deepcopy.DeepCopy(r.Elasticapm, &ds.Values.ElasticApm)
	}

	if r.Signalfx != nil {
		ds.Values.SignalFx = &datastore.SignalFXConfig{}
		deepcopy.DeepCopy(r.Signalfx, &ds.Values.SignalFx)
	}

	if r.Awsxray != nil {
		ds.Values.AwsXRay = &datastore.AWSXRayConfig{}
		deepcopy.DeepCopy(r.Awsxray, &ds.Values.AwsXRay)
	}

	if r.Azureappinsights != nil {
		ds.Values.AzureAppInsights = &datastore.AzureAppInsightsConfig{}
		deepcopy.DeepCopy(r.Azureappinsights, &ds.Values.AzureAppInsights)
	}

	if r.Sumologic != nil {
		ds.Values.SumoLogic = &datastore.SumoLogicConfig{}
		deepcopy.DeepCopy(r.Sumologic, &ds.Values.SumoLogic)
	}

	ds.Type = datastore.DataStoreType(r.Type)
	return &ds, nil
}

func convertTraceInToProtoSpans(trace traces.Trace) []*proto.Span {
	spans := make([]*proto.Span, 0, len(trace.Flat))
	for _, span := range trace.Flat {
		attributes := make([]*proto.KeyValuePair, 0, span.Attributes.Len())
		for name, value := range span.Attributes.Values() {
			attributes = append(attributes, &proto.KeyValuePair{
				Key:   name,
				Value: value,
			})
		}

		protoSpan := proto.Span{
			Id:         span.ID.String(),
			ParentId:   getParentID(span),
			Name:       span.Name,
			Kind:       string(span.Kind),
			StartTime:  span.StartTime.UnixNano(),
			EndTime:    span.EndTime.UnixNano(),
			Attributes: attributes,
		}
		spans = append(spans, &protoSpan)
	}

	// hack to prevent the "Temporary root span" to be sent alone to the server.
	// This causes the server to be confused when evaluating the trace
	if len(spans) == 1 && spans[0].Name == traces.TemporaryRootSpanName {
		return []*proto.Span{}
	}

	return spans
}

func getParentID(span *traces.Span) string {
	if span.Parent != nil {
		return span.Parent.ID.String()
	}

	return ""
}
