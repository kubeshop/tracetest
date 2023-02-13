package config

type (
	config struct {
		Server    serverConfig    `yaml:",omitempty" mapstructure:"server"`
		GA        googleAnalytics `yaml:"googleAnalytics,omitempty" mapstructure:"googleAnalytics"`
		Telemetry telemetry       `yaml:",omitempty" mapstructure:"telemetry"`
		Demo      demo            `yaml:",omitempty" mapstructure:"demo"`
	}

	demo struct {
		Enabled   []string          `yaml:",omitempty" mapstructure:"enabled"`
		Endpoints map[string]string `yaml:",omitempty" mapstructure:"endpoints"`
	}

	googleAnalytics struct {
		Enabled bool `yaml:",omitempty" mapstructure:"enabled"`
	}

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
