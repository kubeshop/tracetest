package config

import "fmt"

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
		DataStores map[string]TracingBackendDataStoreConfig `yaml:",omitempty" mapstructure:"dataStores"`
		Exporters  map[string]TelemetryExporterOption       `yaml:",omitempty" mapstructure:"exporters"`
	}

	TelemetryExporterOption struct {
		ServiceName string         `yaml:",omitempty" mapstructure:"serviceName"`
		Sampling    float64        `yaml:",omitempty" mapstructure:"sampling"`
		Exporter    ExporterConfig `yaml:",omitempty" mapstructure:"exporter"`
	}

	ExporterConfig struct {
		Type                   string              `yaml:",omitempty" mapstructure:"type"`
		CollectorConfiguration OTELCollectorConfig `yaml:"collector,omitempty" mapstructure:"collector"`
	}
)

func (c *Config) Exporter() (*TelemetryExporterOption, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.getExporter(c.config.Server.Telemetry.Exporter)
}

func (c *Config) ApplicationExporter() (*TelemetryExporterOption, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.getExporter(c.config.Server.Telemetry.ApplicationExporter)
}

func (c *Config) getExporter(name string) (*TelemetryExporterOption, error) {
	// Exporters are optional: if no name was provided we consider that the user don't want to have them enabled
	if name == "" {
		return nil, nil
	}

	exporterConfig, found := c.config.Telemetry.Exporters[name]
	if !found {
		availableOptions := mapKeys(c.config.Telemetry.DataStores)
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
