package config

import (
	"fmt"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
	"gopkg.in/yaml.v2"
)

var ErrInvalidTraceDBProvider = fmt.Errorf("invalid traceDB provider")

type (
	Config struct {
		Server               ServerConfig    `yaml:",omitempty" mapstructure:"server"`
		PostgresConnString   string          `yaml:",omitempty" mapstructure:"postgresConnString"`
		PoolingConfig        PoolingConfig   `yaml:",omitempty" mapstructure:"poolingConfig"`
		GA                   GoogleAnalytics `yaml:"googleAnalytics,omitempty" mapstructure:"googleAnalytics"`
		Telemetry            Telemetry       `yaml:",omitempty" mapstructure:"telemetry"`
		Demo                 Demo            `yaml:",omitempty" mapstructure:"demo"`
		ExperimentalFeatures []string        `yaml:",omitempty" mapstructure:"experimentalFeatures"`
	}

	Demo struct {
		Enabled   []string      `yaml:",omitempty" mapstructure:"enabled"`
		Endpoints DemoEndpoints `yaml:",omitempty" mapstructure:"endpoints"`
	}

	DemoEndpoints struct {
		PokeshopHttp       string `yaml:",omitempty" mapstructure:"pokeshopHttp"`
		PokeshopGrpc       string `yaml:",omitempty" mapstructure:"pokeshopGrpc"`
		OtelFrontend       string `yaml:",omitempty" mapstructure:"otelFrontend"`
		OtelProductCatalog string `yaml:",omitempty" mapstructure:"otelProductCatalog"`
		OtelCart           string `yaml:",omitempty" mapstructure:"otelCart"`
		OtelCheckout       string `yaml:",omitempty" mapstructure:"otelCheckout"`
	}

	GoogleAnalytics struct {
		Enabled bool `yaml:",omitempty" mapstructure:"enabled"`
	}

	PoolingConfig struct {
		MaxWaitTimeForTrace string `yaml:",omitempty" mapstructure:"maxWaitTimeForTrace"`
		RetryDelay          string `yaml:",omitempty" mapstructure:"retryDelay"`
	}

	ServerConfig struct {
		PathPrefix string                `yaml:",omitempty" mapstructure:"pathPrefix"`
		HttpPort   int                   `yaml:",omitempty" mapstructure:"httpPort"`
		Telemetry  ServerTelemetryConfig `yaml:",omitempty" mapstructure:"telemetry"`
	}

	ServerTelemetryConfig struct {
		Exporter            string `yaml:",omitempty" mapstructure:"exporter"`
		ApplicationExporter string `yaml:",omitempty" mapstructure:"applicationExporter"`
		DataStore           string `yaml:",omitempty" mapstructure:"dataStore"`
	}

	Telemetry struct {
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

	SignalFXDataStoreConfig struct {
		Realm string
		Token string
	}
)

func (c Config) PoolingRetryDelay() time.Duration {
	delay, err := time.ParseDuration(c.PoolingConfig.RetryDelay)
	if err != nil {
		return 1 * time.Second
	}

	return delay
}

func (c Config) MaxWaitTimeForTraceDuration() time.Duration {
	maxWaitTimeForTrace, err := time.ParseDuration(c.PoolingConfig.MaxWaitTimeForTrace)
	if err != nil {
		// use a default value
		maxWaitTimeForTrace = 30 * time.Second
	}
	return maxWaitTimeForTrace
}

func (c Config) DataStore() (*TracingBackendDataStoreConfig, error) {
	selectedStore := c.Server.Telemetry.DataStore
	dataStoreConfig, found := c.Telemetry.DataStores[selectedStore]

	if selectedStore != "" && !found {
		return nil, ErrInvalidTraceDBProvider
	}

	if !found {
		return nil, nil
	}

	return &dataStoreConfig, nil
}

func (c Config) IsDataStoreConfigured() bool {
	dataStore, _ := c.DataStore()

	return dataStore != nil
}

func (c Config) Exporter() (*TelemetryExporterOption, error) {
	return c.getExporter(c.Server.Telemetry.Exporter)
}

func (c Config) ApplicationExporter() (*TelemetryExporterOption, error) {
	return c.getExporter(c.Server.Telemetry.ApplicationExporter)
}

func (c Config) getExporter(name string) (*TelemetryExporterOption, error) {
	// Exporters are optional: if no name was provided we consider that the user don't want to have them enabled
	if name == "" {
		return nil, nil
	}

	exporterConfig, found := c.Telemetry.Exporters[name]
	if !found {
		availableOptions := mapKeys(c.Telemetry.DataStores)
		return nil, fmt.Errorf(`invalid exporter option: "%s". Available options: %v`, name, availableOptions)
	}

	return &exporterConfig, nil
}

func FromFile(file string) (Config, error) {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return Config{}, fmt.Errorf("read file: %w", err)
	}

	var m map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &m)
	if err != nil {
		return Config{}, fmt.Errorf("yaml unmarshal : %w", err)
	}

	var c Config
	err = mapstructure.Decode(m, &c)
	if err != nil {
		return Config{}, fmt.Errorf("yaml unmarshal : %w", err)
	}

	return c, nil
}

func mapKeys[T any](m map[string]T) []string {
	keys := make([]string, 0)
	for key := range m {
		keys = append(keys, key)
	}

	return keys
}
