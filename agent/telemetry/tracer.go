package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const spanExporterTimeout = 1 * time.Minute

func GetTracer(ctx context.Context, otelExporterEndpoint, serviceName string) (trace.Tracer, error) {
	realServiceName := fmt.Sprintf("tracetestAgent_%s", serviceName)

	spanExporter, err := newSpanExporter(ctx, otelExporterEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to setup span exporter: %w", err)
	}

	traceProvider, err := newTraceProvider(ctx, spanExporter, realServiceName)
	if err != nil {
		return nil, fmt.Errorf("failed to setup trace provider: %w", err)
	}

	return traceProvider.Tracer(realServiceName), nil
}

func newSpanExporter(ctx context.Context, otelExporterEndpoint string) (sdkTrace.SpanExporter, error) {
	ctx, cancel := context.WithTimeout(ctx, spanExporterTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, otelExporterEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	return traceExporter, nil
}

func newTraceProvider(ctx context.Context, spanExporter sdkTrace.SpanExporter, serviceName string) (*sdkTrace.TracerProvider, error) {
	defaultResource := resource.Default()

	mergedResource, err := resource.Merge(
		defaultResource,
		resource.NewWithAttributes(
			defaultResource.SchemaURL(),
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create otel resource: %w", err)
	}

	tp := sdkTrace.NewTracerProvider(
		sdkTrace.WithBatcher(spanExporter),
		sdkTrace.WithResource(mergedResource),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp, nil
}
