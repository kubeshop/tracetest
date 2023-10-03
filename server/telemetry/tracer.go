package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var tracerProviderInstance *sdktrace.TracerProvider

func ShutdownTracer(ctx context.Context) error {
	if tracerProviderInstance != nil {
		return tracerProviderInstance.Shutdown(ctx)
	}

	return nil
}

type exporterConfig interface {
	Exporter() (*config.TelemetryExporterOption, error)
}

func NewTracer(ctx context.Context, cfg exporterConfig) (trace.Tracer, error) {
	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(propagator)

	exporterConfig, err := cfg.Exporter()
	if err != nil {
		return nil, fmt.Errorf("could not get exporter config: %w", err)
	}

	tracerProvider, err := NewTracerProvider(ctx, exporterConfig)
	if err != nil {
		return nil, fmt.Errorf("could not get trace provider: %w", err)
	}

	tracerProviderInstance = tracerProvider

	otel.SetTracerProvider(tracerProvider)

	tracer := tracerProvider.Tracer("tracetest")
	return tracer, nil
}

func NewTracerProvider(ctx context.Context, exporterConfig *config.TelemetryExporterOption) (*sdktrace.TracerProvider, error) {
	if exporterConfig == nil {
		return sdktrace.NewTracerProvider(), nil
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(exporterConfig.ServiceName),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("could not get provider resource: %w", err)
	}

	exporter, err := getExporter(ctx, exporterConfig.Exporter)
	if err != nil {
		return nil, fmt.Errorf("could not create exporter: %w", err)
	}

	processor := sdktrace.NewBatchSpanProcessor(exporter, sdktrace.WithBatchTimeout(100*time.Millisecond))
	sampleRate := exporterConfig.Sampling / 100.0

	return sdktrace.NewTracerProvider(
		sdktrace.WithResource(r),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(sampleRate)),
		sdktrace.WithSpanProcessor(processor),
	), nil
}

func getExporter(ctx context.Context, exporterConfig config.ExporterConfig) (sdktrace.SpanExporter, error) {
	switch exporterConfig.Type {
	case "collector":
		return getOtelCollectorExporter(ctx, exporterConfig)
	}

	return nil, fmt.Errorf("invalid exporter type: %s", exporterConfig.Type)
}

func getOtelCollectorExporter(ctx context.Context, exporterConfig config.ExporterConfig) (sdktrace.SpanExporter, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, exporterConfig.CollectorConfiguration.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("could not create gRPC connection to collector: %w", err)
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("could not create trace exporter: %w", err)
	}

	return exporter, nil
}
