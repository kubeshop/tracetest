package mocks

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewTracer(ctx context.Context, endpoint string) (trace.Tracer, error) {
	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(propagator)

	tracerProvider, err := newTracerProvider(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not get trace provider: %w", err)
	}

	otel.SetTracerProvider(tracerProvider)

	tracer := tracerProvider.Tracer("tracetest")
	return tracer, nil
}

func newTracerProvider(ctx context.Context, endpoint string) (*sdktrace.TracerProvider, error) {
	exporter, err := getExporter(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not create exporter: %w", err)
	}

	processor := sdktrace.NewBatchSpanProcessor(exporter, sdktrace.WithBatchTimeout(100*time.Millisecond))
	sampleRate := 100.0

	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(sampleRate)),
		sdktrace.WithSpanProcessor(processor),
	), nil
}

func getExporter(ctx context.Context, endpoint string) (sdktrace.SpanExporter, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("could not create gRPC connection to collector: %w", err)
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("could not create trace exporter: %w", err)
	}

	return exporter, nil
}
