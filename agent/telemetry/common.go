package telemetry

import (
	"context"
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/server/version"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var propagator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

func getAgentServiceName(serviceName string) string {
	return fmt.Sprintf("tracetest.agent-%s", serviceName)
}

func getResource(serviceName string) (*resource.Resource, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("could not get OS hostname: %w", err)
	}

	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.HostName(hostname),
			semconv.ServiceVersion(version.Version), // TODO: should we consider a version file for the agent?
		),
	)

	if err != nil {
		return nil, fmt.Errorf("could not merge resources: %w", err)
	}

	return resource, nil
}

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
