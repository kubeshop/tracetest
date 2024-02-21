package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

func GetMeter(ctx context.Context, otelExporterEndpoint, serviceName string) (metric.Meter, error) {
	if otelExporterEndpoint == "" {
		// fallback, return noop
		return noop.NewMeterProvider().Meter("noop"), nil
	}

	realServiceName := fmt.Sprintf("tracetestAgent_%s", serviceName)

	return nil, nil
}
