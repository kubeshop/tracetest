package workers

import (
	"context"
	"errors"
	"fmt"
	"log"

	gocache "github.com/Code-Hex/go-generics-cache"
	"github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
)

type PollerWorker struct {
	client            *client.Client
	tracer            trace.Tracer
	sentSpanIDs       *gocache.Cache[string, bool]
	inmemoryDatastore tracedb.TraceDB
}

type PollerOption func(*PollerWorker)

func WithInMemoryDatastore(datastore tracedb.TraceDB) PollerOption {
	return func(pw *PollerWorker) {
		pw.inmemoryDatastore = datastore
	}
}

func NewPollerWorker(client *client.Client, opts ...PollerOption) *PollerWorker {
	// TODO: use a real tracer
	tracer := trace.NewNoopTracerProvider().Tracer("noop")

	pollerWorker := &PollerWorker{
		client:      client,
		tracer:      tracer,
		sentSpanIDs: gocache.New[string, bool](),
	}

	for _, opt := range opts {
		opt(pollerWorker)
	}

	return pollerWorker
}

func (w *PollerWorker) Poll(ctx context.Context, request *proto.PollingRequest) error {
	err := w.poll(ctx, request)
	if err != nil {
		sendErr := w.client.SendTrace(ctx, &proto.PollingResponse{
			RequestID:           request.RequestID,
			AgentIdentification: w.client.SessionConfiguration().AgentIdentification,
			TestID:              request.GetTestID(),
			RunID:               request.GetRunID(),
			TraceID:             request.TraceID,
			TraceFound:          false,
			Error: &proto.Error{
				Message: err.Error(),
			},
		})

		if sendErr != nil {
			return fmt.Errorf("could not report polling error back to the server: %w. Original error: %s", sendErr, err.Error())
		}
	}

	return err
}

func (w *PollerWorker) poll(ctx context.Context, request *proto.PollingRequest) error {
	datastoreConfig, err := convertProtoToDataStore(request.Datastore)
	if err != nil {
		return err
	}

	if datastoreConfig == nil {
		return fmt.Errorf("invalid datastore: nil")
	}

	dsFactory := tracedb.Factory(nil)
	ds, err := dsFactory(*datastoreConfig)
	if err != nil {
		log.Printf("Invalid datastore: %s", err.Error())
		return err
	}

	if datastoreConfig.IsOTLPBasedProvider() && w.inmemoryDatastore != nil {
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
		if !errors.Is(err, connection.ErrTraceNotFound) {
			// let controlplane know we didn't find the trace
			log.Printf("cannot get trace from datastore: %s", err.Error())
			return err
		}

		pollingResponse.TraceFound = false
	} else {
		pollingResponse.TraceFound = true
		pollingResponse.Spans = convertTraceInToProtoSpans(trace)

		// remove sent spans
		newSpans := make([]*proto.Span, 0, len(pollingResponse.Spans))
		for _, span := range pollingResponse.Spans {
			runKey := fmt.Sprintf("%d-%s-%s", request.RunID, request.TestID, span.Id)
			_, alreadySent := w.sentSpanIDs.Get(runKey)
			if !alreadySent {
				newSpans = append(newSpans, span)
			}
		}
		pollingResponse.Spans = newSpans
	}

	err = w.client.SendTrace(ctx, pollingResponse)
	if err != nil {
		log.Printf("cannot send trace to server: %s", err.Error())
		return err
	}

	// mark spans as sent
	for _, span := range pollingResponse.Spans {
		runKey := fmt.Sprintf("%d-%s-%s", request.RunID, request.TestID, span.Id)
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
		attributes := make([]*proto.KeyValuePair, 0, len(span.Attributes))
		for name, value := range span.Attributes {
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
