package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
)

var propagator = propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})

func InjectMetadataIntoContext(ctx context.Context, metadata map[string]string) context.Context {
	return propagator.Extract(
		ctx,
		propagation.MapCarrier(metadata),
	)
}

func ExtractMetadataFromContext(ctx context.Context) map[string]string {
	metadata := map[string]string{}
	propagator.Inject(
		ctx,
		propagation.MapCarrier(metadata),
	)

	return metadata
}
