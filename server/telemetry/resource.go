package telemetry

import (
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/version"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func getResource(cfg *config.TelemetryExporterOption) (*resource.Resource, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("could not get OS hostname: %w", err)
	}

	serviceName := "tracetest"
	if cfg != nil && cfg.ServiceName != "" {
		serviceName = cfg.ServiceName
	}

	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.HostName(hostname),
			semconv.ServiceVersion(version.Version),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("could not merge resources: %w", err)
	}

	return resource, nil
}
