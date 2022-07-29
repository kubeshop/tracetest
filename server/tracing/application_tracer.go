package tracing

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/config"
	"go.opentelemetry.io/otel/trace"
)

func GetApplicationTracer(ctx context.Context, config config.Config) (trace.Tracer, error) {
	applicationSpanExporter, err := config.ApplicationExporter()
	if err != nil {
		return nil, fmt.Errorf("could not get application exporter config: %w", err)
	}

	triggerSpanTracerProvider, err := NewTracerProvider(ctx, applicationSpanExporter)
	if err != nil {
		return nil, fmt.Errorf("could not create tracer provider for application span: %w", err)
	}

	return triggerSpanTracerProvider.Tracer("trigger"), nil
}
