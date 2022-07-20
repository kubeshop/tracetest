package trigger

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/tracing"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func getTracerProvider(config config.Config) (*sdktrace.TracerProvider, error) {
	appExporterConfig, err := config.ApplicationExporter()
	if err != nil {
		return nil, fmt.Errorf("could not create application exporter: %w", err)
	}

	return tracing.NewTracerProvider(context.Background(), appExporterConfig)
}

func getTracer(spanContext *sdktrace.TracerProvider) trace.Tracer {
	return spanContext.Tracer("tracetest")
}
