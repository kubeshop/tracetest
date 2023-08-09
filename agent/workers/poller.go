package workers

import (
	"context"
	"fmt"
	"log"

	"github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/workers/datastores"
	"github.com/kubeshop/tracetest/agent/workers/datastores/connection"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
)

type PollerWorker struct {
	client *client.Client
	tracer trace.Tracer
}

func NewPollerWorker(client *client.Client) *PollerWorker {
	// TODO: use a real tracer
	tracer := trace.NewNoopTracerProvider().Tracer("noop")

	return &PollerWorker{client, tracer}
}

func (w *PollerWorker) Poll(ctx context.Context, request *proto.PollingRequest) error {
	datastoreConfig, err := convertProtoToDataStore(request.Datastore)
	if err != nil {
		return err
	}

	if datastoreConfig == nil {
		return fmt.Errorf("invalid datastore: nil")
	}

	dsFactory := datastores.Factory()
	ds, err := dsFactory(*datastoreConfig)
	if err != nil {
		return err
	}

	trace, err := ds.GetTraceByID(ctx, request.TraceID)
	if err == connection.ErrTraceNotFound {
		// Nothing to send back
		log.Printf("Trace %s was not found", request.TraceID)
		return nil
	}

	spans := convertTraceInToProtoSpans(trace)

	err = w.client.SendTrace(ctx, request, spans...)
	if err != nil {
		return err
	}

	return nil
}

func convertProtoToDataStore(request *proto.DataStore) (*datastore.DataStore, error) {
	var ds datastore.DataStore
	err := deepcopy.DeepCopy(request, &ds.Values)
	if err != nil {
		return nil, fmt.Errorf("could not deep copy datastore: %w", err)
	}

	ds.Type = datastore.DataStoreType(request.Type)
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
			StartTime:  span.StartTime.UnixMicro(),
			EndTime:    span.EndTime.UnixMicro(),
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
