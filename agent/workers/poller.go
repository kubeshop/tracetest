package workers

import (
	"context"
	"errors"
	"fmt"
	"log"

	gocache "github.com/Code-Hex/go-generics-cache"
	"github.com/davecgh/go-spew/spew"
	"github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/event"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type PollerWorker struct {
	client            *client.Client
	tracer            trace.Tracer
	sentSpanIDs       *gocache.Cache[string, bool]
	inmemoryDatastore tracedb.TraceDB
	logger            *zap.Logger
	observer          event.Observer
}

type PollerOption func(*PollerWorker)

func WithInMemoryDatastore(datastore tracedb.TraceDB) PollerOption {
	return func(pw *PollerWorker) {
		pw.inmemoryDatastore = datastore
	}
}

func WithObserver(observer event.Observer) PollerOption {
	return func(pw *PollerWorker) {
		pw.observer = observer
	}
}

func NewPollerWorker(client *client.Client, opts ...PollerOption) *PollerWorker {
	// TODO: use a real tracer
	tracer := trace.NewNoopTracerProvider().Tracer("noop")

	pollerWorker := &PollerWorker{
		client:      client,
		tracer:      tracer,
		sentSpanIDs: gocache.New[string, bool](),
		logger:      zap.NewNop(),
		observer:    event.NewNopObserver(),
	}

	for _, opt := range opts {
		opt(pollerWorker)
	}

	return pollerWorker
}

func (w *PollerWorker) SetLogger(logger *zap.Logger) {
	w.logger = logger
}

func (w *PollerWorker) Poll(ctx context.Context, request *proto.PollingRequest) error {
	w.logger.Debug("Received polling request", zap.Any("request", request))
	w.observer.StartTracePoll(request)

	err := w.poll(ctx, request)
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

			return fmt.Errorf("could not report polling error back to the server: %w. Original error: %s", sendErr, err.Error())
		}
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
	w.logger.Debug("Converted datastore", zap.Any("datastore", datastoreConfig))

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
			_, alreadySent := w.sentSpanIDs.Get(runKey)
			if !alreadySent {
				w.logger.Debug("Span was not sent", zap.String("runKey", runKey))
				newSpans = append(newSpans, span)
			} else {
				w.logger.Debug("Span was already sent", zap.String("runKey", runKey))
			}
		}
		pollingResponse.Spans = newSpans
		w.logger.Debug("Filtered spans", zap.Any("pollingResponse", spew.Sdump(pollingResponse)))
	}

	err = w.client.SendTrace(ctx, pollingResponse)
	if err != nil {
		w.logger.Error("Cannot send trace to server", zap.Error(err))
		log.Printf("cannot send trace to server: %s", err.Error())
		return err
	}

	// mark spans as sent
	for _, span := range pollingResponse.Spans {
		runKey := fmt.Sprintf("%d-%s-%s", request.RunID, request.TestID, span.Id)
		w.logger.Debug("Marking span as sent", zap.String("runKey", runKey))
		// TODO: we can set the expiration for this key to be
		// 1 second after the pollingProfile max waiting time
		// but we need to get that info here from controlplane
		w.sentSpanIDs.Set(runKey, true)
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

	return spans
}

func getParentID(span *traces.Span) string {
	if span.Parent != nil {
		return span.Parent.ID.String()
	}

	return ""
}
