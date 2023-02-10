package config

import (
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
)

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

	TracingBackendDataStoreConfig struct {
		Type       string                        `yaml:",omitempty" mapstructure:"type"`
		Jaeger     configgrpc.GRPCClientSettings `yaml:",omitempty" mapstructure:"jaeger"`
		Tempo      BaseClientConfig              `yaml:",omitempty" mapstructure:"tempo"`
		OpenSearch ElasticSearchDataStoreConfig  `yaml:",omitempty" mapstructure:"opensearch"`
		SignalFX   SignalFXDataStoreConfig       `yaml:",omitempty" mapstructure:"signalfx"`
		ElasticApm ElasticSearchDataStoreConfig  `yaml:",omitempty" mapstructure:"elasticapm"`
	}

	TelemetryExporterOption struct {
		ServiceName string         `yaml:",omitempty" mapstructure:"serviceName"`
		Sampling    float64        `yaml:",omitempty" mapstructure:"sampling"`
		Exporter    ExporterConfig `yaml:",omitempty" mapstructure:"exporter"`
	}

	BaseClientConfig struct {
		Type string                        `yaml:",omitempty" mapstructure:"type"`
		Grpc configgrpc.GRPCClientSettings `yaml:",omitempty" mapstructure:"grpc"`
		Http HttpClientConfig              `yaml:",omitempty" mapstructure:"http"`
	}

	HttpClientConfig struct {
		Url        string                     `yaml:",omitempty" mapstructure:"url"`
		Headers    map[string]string          `yaml:",omitempty" mapstructure:"headers"`
		TLSSetting configtls.TLSClientSetting `yaml:",omitempty" mapstructure:"tls"`
	}

	ExporterConfig struct {
		Type                   string              `yaml:",omitempty" mapstructure:"type"`
		CollectorConfiguration OTELCollectorConfig `yaml:"collector,omitempty" mapstructure:"collector"`
	}

	OTELCollectorConfig struct {
		Endpoint string `yaml:",omitempty" mapstructure:"endpoint"`
	}

	ElasticSearchDataStoreConfig struct {
		Addresses          []string
		Username           string
		Password           string
		Index              string
		Certificate        string
		InsecureSkipVerify bool
	}

	SignalFXDataStoreConfig struct {
		Realm string
		Token string
	}
)
