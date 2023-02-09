package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"go.opentelemetry.io/collector/config/configgrpc"
	"gopkg.in/yaml.v2"
)

type (
	config struct {
		Server               serverConfig    `yaml:",omitempty" mapstructure:"server"`
		PostgresConnString   string          `yaml:",omitempty" mapstructure:"postgresConnString"`
		PoolingConfig        PoolingConfig   `yaml:",omitempty" mapstructure:"poolingConfig"`
		GA                   googleAnalytics `yaml:"googleAnalytics,omitempty" mapstructure:"googleAnalytics"`
		Telemetry            telemetry       `yaml:",omitempty" mapstructure:"telemetry"`
		Demo                 demo            `yaml:",omitempty" mapstructure:"demo"`
		ExperimentalFeatures []string        `yaml:",omitempty" mapstructure:"experimentalFeatures"`
	}

	demo struct {
		Enabled   []string          `yaml:",omitempty" mapstructure:"enabled"`
		Endpoints map[string]string `yaml:",omitempty" mapstructure:"endpoints"`
	}

	googleAnalytics struct {
		Enabled bool `yaml:",omitempty" mapstructure:"enabled"`
	}

	PoolingConfig struct {
		MaxWaitTimeForTrace string `yaml:",omitempty" mapstructure:"maxWaitTimeForTrace"`
		RetryDelay          string `yaml:",omitempty" mapstructure:"retryDelay"`
	}

	serverConfig struct {
		PathPrefix string                `yaml:",omitempty" mapstructure:"pathPrefix"`
		HttpPort   int                   `yaml:",omitempty" mapstructure:"httpPort"`
		Telemetry  serverTelemetryConfig `yaml:",omitempty" mapstructure:"telemetry"`
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
		Tempo      configgrpc.GRPCClientSettings `yaml:",omitempty" mapstructure:"tempo"`
		OpenSearch ElasticSearchDataStoreConfig  `yaml:",omitempty" mapstructure:"opensearch"`
		SignalFX   SignalFXDataStoreConfig       `yaml:",omitempty" mapstructure:"signalfx"`
		ElasticApm ElasticSearchDataStoreConfig  `yaml:",omitempty" mapstructure:"elasticapm"`
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

func FromFile(file string) (*Config, error) {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return &Config{}, fmt.Errorf("read file: %w", err)
	}

	var m map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		return &Config{}, fmt.Errorf("yaml unmarshal : %w", err)
	}

	var c config
	err = mapstructure.Decode(m, &c)
	if err != nil {
		return &Config{}, fmt.Errorf("yaml unmarshal : %w", err)
	}

	cfg := New()
	cfg.config = &c

	return cfg, nil
}
