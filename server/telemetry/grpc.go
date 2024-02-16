package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var propagator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

func InjectContextIntoStream(ctx context.Context, stream grpc.ServerStream) error {
	header := make(metadata.MD)
	propagator.Inject(ctx, &metadataSupplier{metadata: &header})

	err := stream.SetHeader(header)
	if err != nil {
		return fmt.Errorf("could not set header: %w", err)
	}

	return nil
}

func ExtractContextFromStream(stream grpc.ClientStream) (context.Context, error) {
	ctx := stream.Context()
	header, err := stream.Header()
	if err != nil {
		return ctx, fmt.Errorf("coult not get header from stream: %w", err)
	}

	ctx = propagator.Extract(ctx, &metadataSupplier{metadata: &header})
	return ctx, nil
}

type metadataSupplier struct {
	metadata *metadata.MD
}

// assert that metadataSupplier implements the TextMapCarrier interface.
var _ propagation.TextMapCarrier = &metadataSupplier{}

func (s *metadataSupplier) Get(key string) string {
	values := s.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (s *metadataSupplier) Set(key string, value string) {
	s.metadata.Set(key, value)
}

func (s *metadataSupplier) Keys() []string {
	out := make([]string, 0, len(*s.metadata))
	for key := range *s.metadata {
		out = append(out, key)
	}
	return out
}
