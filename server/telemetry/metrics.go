package telemetry

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kubeshop/tracetest/server/config"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
)

func NewMeter(ctx context.Context, cfg exporterConfig) (metric.Meter, error) {
	provider, err := newMeterProvider(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("could not create meter provider: %w", err)
	}

	return provider.Meter("tracetest"), nil
}

func newMeterProvider(ctx context.Context, cfg exporterConfig) (metric.MeterProvider, error) {
	exporterConfig, err := cfg.Exporter()
	if err != nil {
		return nil, fmt.Errorf("could not get exporter config: %w", err)
	}

	if exporterConfig == nil {
		log.Println("empty exporter config: falling back to noop meter provider")
		return noop.NewMeterProvider(), nil
	}

	resource, err := getResource(exporterConfig)
	if err != nil {
		return nil, fmt.Errorf("could not get resource: %w", err)
	}

	collectorExporter, err := getOtelMetricsCollectorExporter(ctx, exporterConfig)
	if err != nil {
		return nil, fmt.Errorf("could not get collector exporter: %w", err)
	}

	interval := 10 * time.Second
	if exporterConfig.MetricsReaderInterval != 0 {
		interval = exporterConfig.MetricsReaderInterval
	}

	periodicReader := metricsdk.NewPeriodicReader(collectorExporter,
		metricsdk.WithInterval(interval),
	)

	provider := metricsdk.NewMeterProvider(
		metricsdk.WithResource(resource),
		metricsdk.WithReader(periodicReader),
	)

	return provider, nil
}

func getOtelMetricsCollectorExporter(ctx context.Context, exporterConfig *config.TelemetryExporterOption) (metricsdk.Exporter, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(exporterConfig.Exporter.CollectorConfiguration.Endpoint),
		otlpmetricgrpc.WithCompressor("gzip"),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("could not create metric exporter: %w", err)
	}

	return exporter, nil
}
