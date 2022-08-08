package trigger

import (
	"context"
	"fmt"
	"io"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/tracing"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func getTracerProvider(config config.Config) (*sdktrace.TracerProvider, error) {
	appExporterConfig, err := config.ApplicationExporter()
	if err != nil {
		return nil, fmt.Errorf("could not create application exporter: %w", err)
	}

	return tracing.NewTracerProvider(context.Background(), appExporterConfig)
}

func noopTracerProvider() *sdktrace.TracerProvider {
	// Set standard attributes per semantic conventions
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("tracetest"),
	)

	// this is in fact a noop exporter, so we can ignore errors
	spanExporter, _ := stdouttrace.New(stdouttrace.WithWriter(io.Discard))

	return sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(spanExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
	)
}
