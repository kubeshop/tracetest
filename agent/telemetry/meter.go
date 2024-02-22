package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
)

const (
	metricsReaderInterval = 30 * time.Second
	metricExporterTimeout = 5 * time.Second
)

func GetNoopMeter() metric.Meter {
	return noop.NewMeterProvider().Meter("noop")
}

func GetMeter(ctx context.Context, otelExporterEndpoint, serviceName string) (metric.Meter, error) {
	if otelExporterEndpoint == "" {
		// fallback, return noop
		return GetNoopMeter(), nil
	}

	realServiceName := getAgentServiceName(serviceName)

	provider, err := newMeterProvider(ctx, otelExporterEndpoint, realServiceName)
	if err != nil {
		return nil, fmt.Errorf("could not create meter provider: %w", err)
	}

	return provider.Meter("tracetest.agent"), nil
}

func newMeterProvider(ctx context.Context, otelExporterEndpoint, serviceName string) (metric.MeterProvider, error) {
	resource, err := getResource(serviceName)
	if err != nil {
		return nil, fmt.Errorf("could not get resource: %w", err)
	}

	exporter, err := getMetricExporter(ctx, otelExporterEndpoint)
	if err != nil {
		return nil, fmt.Errorf("could not create metric exporter: %w", err)
	}

	periodicReader := metricsdk.NewPeriodicReader(
		exporter,
		metricsdk.WithInterval(metricsReaderInterval),
	)

	provider := metricsdk.NewMeterProvider(
		metricsdk.WithResource(resource),
		metricsdk.WithReader(periodicReader),
	)

	return provider, nil
}

func getMetricExporter(ctx context.Context, otelExporterEndpoint string) (*otlpmetricgrpc.Exporter, error) {
	ctx, cancel := context.WithTimeout(ctx, metricExporterTimeout)
	defer cancel()

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(otelExporterEndpoint),
		otlpmetricgrpc.WithCompressor("gzip"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create metric exporter: %w", err)
	}

	return exporter, nil
}
