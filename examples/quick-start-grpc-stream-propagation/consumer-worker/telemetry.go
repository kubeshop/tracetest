package main

import (
	"context"
	"fmt"
	"log"
	"os"
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

func getEnvVar(envVarName, defaultValue string) string {
	envVarValue := os.Getenv(envVarName)
	if envVarValue == "" {
		return defaultValue
	}

	return envVarValue
}

func setupOpenTelemetry(ctx context.Context, otelExporterEndpoint, serviceName string) (trace.Tracer, error) {
	log.Printf("Setting up OpenTelemetry with exporter endpoint %s and service name %s", otelExporterEndpoint, serviceName)

	spanExporter, err := getSpanExporter(ctx, otelExporterEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to setup span exporter: %w", err)
	}

	traceProvider, err := getTraceProvider(spanExporter, serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to setup trace provider: %w", err)
	}

	return traceProvider.Tracer(serviceName), nil
}

func getSpanExporter(ctx context.Context, otelExporterEndpoint string) (sdkTrace.SpanExporter, error) {
	ctx, cancel := context.WithTimeout(ctx, spanExporterTimeout)
	defer cancel()

	conn, err := grpc.NewClient(
		otelExporterEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	return traceExporter, nil
}

func getTraceProvider(spanExporter sdkTrace.SpanExporter, serviceName string) (*sdkTrace.TracerProvider, error) {
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

func injectMetadataIntoContext(ctx context.Context, metadata map[string]string) context.Context {
	propagator := otel.GetTextMapPropagator()

	return propagator.Extract(
		ctx,
		propagation.MapCarrier(metadata),
	)
}

func extractMetadataFromContext(ctx context.Context) map[string]string {
	propagator := otel.GetTextMapPropagator()

	metadata := map[string]string{}
	propagator.Inject(
		ctx,
		propagation.MapCarrier(metadata),
	)

	return metadata
}
