package config

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/mitchellh/mapstructure"
	"go.opentelemetry.io/collector/config/configgrpc"
	"gopkg.in/yaml.v2"
)

type (
	Config struct {
		Server             ServerConfig    `mapstructure:"server"`
		PostgresConnString string          `mapstructure:"postgresConnString"`
		PoolingConfig      PoolingConfig   `mapstructure:"poolingConfig"`
		GA                 GoogleAnalytics `mapstructure:"googleAnalytics"`
		Telemetry          Telemetry       `mapstructure:"telemetry"`
	}

	GoogleAnalytics struct {
		Enabled bool `mapstructure:"enabled"`
	}

	PoolingConfig struct {
		MaxWaitTimeForTrace string `mapstructure:"maxWaitTimeForTrace"`
		RetryDelay          string `mapstructure:"retryDelay"`
	}

	ServerConfig struct {
		PathPrefix string                `mapstructure:"pathPrefix"`
		HttpPort   int                   `mapstructure:"httpPort"`
		Telemetry  ServerTelemetryConfig `mapstructure:"telemetry"`
	}

	ServerTelemetryConfig struct {
		Exporter            string `mapstructure:"exporter"`
		ApplicationExporter string `mapstructure:"applicationExporter"`
		DataStore           string `mapstructure:"dataStore"`
	}

	Telemetry struct {
		DataStores map[string]TracingBackendDataStoreConfig `mapstructure:"dataStores"`
		Exporters  map[string]TelemetryExporterOption       `mapstructure:"exporters"`
	}

	TracingBackendDataStoreConfig struct {
		Type       string                        `mapstructure:"type"`
		Jaeger     configgrpc.GRPCClientSettings `mapstructure:"jaeger"`
		Tempo      configgrpc.GRPCClientSettings `mapstructure:"tempo"`
		OpenSearch OpensearchDataStoreConfig     `mapstructure:"opensearch"`
		SignalFX   SignalFXDataStoreConfig       `mapstructure:"signalfx"`
	}

	TelemetryExporterOption struct {
		ServiceName string         `mapstructure:"serviceName"`
		Sampling    float64        `mapstructure:"sampling"`
		Exporter    ExporterConfig `mapstructure:"exporter"`
	}

	ExporterConfig struct {
		Type                   string              `mapstructure:"type"`
		CollectorConfiguration OTELCollectorConfig `mapstructure:"collector"`
	}

	OTELCollectorConfig struct {
		Endpoint string `mapstructure:"endpoint"`
	}

	OpensearchDataStoreConfig struct {
		Addresses []string
		Username  string
		Password  string
		Index     string
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
	if !found {
		availableOptions := mapKeys(c.Telemetry.DataStores)
		return nil, fmt.Errorf(`invalid data store option: "%s". Available options: %v`, selectedStore, availableOptions)
	}

	return &dataStoreConfig, nil
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
	yamlFile, err := ioutil.ReadFile(file)
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
	for key, _ := range m {
		keys = append(keys, key)
	}

	return keys
}
