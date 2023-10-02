package config

import (
	"fmt"
	"time"
)

type (
	serverConfig struct {
		Telemetry serverTelemetryConfig `yaml:",omitempty" mapstructure:"telemetry"`
	}

	serverTelemetryConfig struct {
		Exporter            string `yaml:",omitempty" mapstructure:"exporter"`
		ApplicationExporter string `yaml:",omitempty" mapstructure:"applicationExporter"`
		DataStore           string `yaml:",omitempty" mapstructure:"dataStore"`
	}

	telemetry struct {
		Exporters map[string]TelemetryExporterOption `yaml:",omitempty" mapstructure:"exporters"`
	}

	TelemetryExporterOption struct {
		ServiceName           string         `yaml:",omitempty" mapstructure:"serviceName"`
		Sampling              float64        `yaml:",omitempty" mapstructure:"sampling"`
		MetricsReaderInterval time.Duration  `yaml:",omitempty" mapstructure:"metricsReaderInterval"`
		Exporter              ExporterConfig `yaml:",omitempty" mapstructure:"exporter"`
	}

	ExporterConfig struct {
		Type                   string              `yaml:",omitempty" mapstructure:"type"`
		CollectorConfiguration OTELCollectorConfig `yaml:"collector,omitempty" mapstructure:"collector"`
	}

	OTELCollectorConfig struct {
		Endpoint string `yaml:",omitempty" mapstructure:"endpoint"`
	}
)

func (c *AppConfig) Exporter() (*TelemetryExporterOption, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.getExporter(c.config.Server.Telemetry.Exporter)
}

func (c *AppConfig) ApplicationExporter() (*TelemetryExporterOption, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.getExporter(c.config.Server.Telemetry.ApplicationExporter)
}

func (c *AppConfig) getExporter(name string) (*TelemetryExporterOption, error) {
	// Exporters are optional: if no name was provided we consider that the user don't want to have them enabled
	if name == "" {
		return nil, nil
	}

	exporterConfig, found := c.config.Telemetry.Exporters[name]
	if !found {
		availableOptions := mapKeys(c.config.Telemetry.Exporters)
		return nil, fmt.Errorf(`invalid exporter option: "%s". Available options: %v`, name, availableOptions)
	}

	return &exporterConfig, nil
}

func mapKeys[T any](m map[string]T) []string {
	keys := make([]string, 0)
	for key := range m {
		keys = append(keys, key)
	}

	return keys
}
