package telemetry

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/config"
	"go.opentelemetry.io/otel/trace"
)

type appExporterConfig interface {
	ApplicationExporter() (*config.TelemetryExporterOption, error)
}

func GetApplicationTracer(ctx context.Context, cfg appExporterConfig) (trace.Tracer, error) {
	applicationSpanExporter, err := cfg.ApplicationExporter()
	if err != nil {
		return nil, fmt.Errorf("could not get application exporter config: %w", err)
	}

	triggerSpanTracerProvider, err := NewTracerProvider(ctx, applicationSpanExporter)
	if err != nil {
		return nil, fmt.Errorf("could not create tracer provider for application span: %w", err)
	}

	return triggerSpanTracerProvider.Tracer("trigger"), nil
}
