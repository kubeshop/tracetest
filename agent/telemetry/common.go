package telemetry

import (
	"fmt"
	"os"

<<<<<<< Updated upstream
	"github.com/kubeshop/tracetest/server/version"
=======
	"github.com/kubeshop/tracetest/agent/version"
	"go.opentelemetry.io/otel/propagation"
>>>>>>> Stashed changes
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func getAgentServiceName(serviceName string) string {
	return fmt.Sprintf("tracetest.agent-%s", serviceName)
}

func getResource(serviceName string) (*resource.Resource, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("could not get OS hostname: %w", err)
	}

	resource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.HostName(hostname),
			semconv.ServiceVersion(version.Version), // TODO: should we consider a version file for the agent?
		),
	)

	if err != nil {
		return nil, fmt.Errorf("could not merge resources: %w", err)
	}

	return resource, nil
}
