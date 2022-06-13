package tracing

import (
	"context"
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/server/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

var tracerProviderInstance *sdktrace.TracerProvider

func ShutdownTracer(ctx context.Context) error {
	if tracerProviderInstance != nil {
		return tracerProviderInstance.Shutdown(ctx)
	}

	return nil
}

func NewTracer(ctx context.Context, config config.Config) (trace.Tracer, error) {
	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(propagator)

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.Telemetry.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("could not get provider resource: %w", err)
	}

	tracerOptions := make([]sdktrace.TracerProviderOption, 0, len(config.Telemetry.Exporters))
	for _, exporterName := range config.Telemetry.Exporters {
		exporter, err := getExporter(ctx, config, exporterName)
		if err != nil {
			return nil, fmt.Errorf(`could not create "%s" exporter: %w`, exporterName, err)
		}

		processor := sdktrace.NewSimpleSpanProcessor(exporter)
		tracerOptions = append(tracerOptions, sdktrace.WithSpanProcessor(processor))
	}

	sampleRate := config.Telemetry.Sampling / 100.0
	tracerOptions = append(tracerOptions, sdktrace.WithResource(r))
	tracerOptions = append(tracerOptions, sdktrace.WithSampler(sdktrace.TraceIDRatioBased(sampleRate)))

	tracerProvider := sdktrace.NewTracerProvider(
		tracerOptions...,
	)

	tracerProviderInstance = tracerProvider

	otel.SetTracerProvider(tracerProvider)

	tracer := tracerProvider.Tracer("tracetest")
	return tracer, nil
}

func getExporter(ctx context.Context, config config.Config, name string) (sdktrace.SpanExporter, error) {
	switch name {
	case "jaeger":
		return newJaegerExporter(ctx, config)
	case "console":
		return newConsoleExporter(ctx, config)
	}

	return nil, nil
}

func newJaegerExporter(ctx context.Context, config config.Config) (sdktrace.SpanExporter, error) {
	exporter, err := jaeger.New(
		jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost(config.Telemetry.Jaeger.Host),
			jaeger.WithAgentPort(fmt.Sprintf("%d", config.Telemetry.Jaeger.Port)),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("could not get jaeger exporter: %w", err)
	}

	return exporter, nil
}

func newConsoleExporter(ctx context.Context, config config.Config) (sdktrace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(os.Stdout),
	)
}
