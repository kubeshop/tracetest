package tracing

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
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
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

	exporter, err := getExporter(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("could not create exporter: %w", err)
	}

	processor := sdktrace.NewBatchSpanProcessor(exporter)
	sampleRate := config.Telemetry.Sampling / 100.0

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(r),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(sampleRate)),
		sdktrace.WithSpanProcessor(processor),
	)

	tracerProviderInstance = tracerProvider

	otel.SetTracerProvider(tracerProvider)

	tracer := tracerProvider.Tracer("tracetest")
	return tracer, nil
}

func getExporter(ctx context.Context, config config.Config) (sdktrace.SpanExporter, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, config.Telemetry.OTelCollectorEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("could not create gRPC connection to collector: %w", err)
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("could not create trace exporter: %w", err)
	}

	return exporter, nil
}
