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
	fmt.Println("Poll handled by agent")
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
		log.Printf("Invalid datastore: %s", err.Error())
		return err
	}

	log.Printf("Getting spans for trace %s", request.TraceID)
	trace, err := ds.GetTraceByID(ctx, request.TraceID)
	if err == connection.ErrTraceNotFound {
		// Nothing to send back
		log.Printf("Trace %s was not found", request.TraceID)
		return nil
	}

	spans := convertTraceInToProtoSpans(trace)

	err = w.client.SendTrace(ctx, request, spans...)
	if err != nil {
		log.Printf("cannot send trace to server: %s", err.Error())
		return err
	}

	return nil
}

func convertProtoToDataStore(r *proto.DataStore) (*datastore.DataStore, error) {
	var ds datastore.DataStore

	if r.Jaeger != nil && r.Jaeger.Grpc != nil {
		ds.Values.Jaeger = &datastore.GRPCClientSettings{}
		deepcopy.DeepCopy(r.Jaeger.Grpc, &ds.Values.Jaeger)
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
