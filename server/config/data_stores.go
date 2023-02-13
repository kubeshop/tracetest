package config

import (
	"fmt"

	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
)

var ErrInvalidTraceDBProvider = fmt.Errorf("invalid traceDB provider")

type (
	TracingBackendDataStoreConfig struct {
		Type       string                        `yaml:",omitempty" mapstructure:"type"`
		Jaeger     configgrpc.GRPCClientSettings `yaml:",omitempty" mapstructure:"jaeger"`
		Tempo      BaseClientConfig              `yaml:",omitempty" mapstructure:"tempo"`
		OpenSearch ElasticSearchDataStoreConfig  `yaml:",omitempty" mapstructure:"opensearch"`
		SignalFX   SignalFXDataStoreConfig       `yaml:",omitempty" mapstructure:"signalfx"`
		ElasticApm ElasticSearchDataStoreConfig  `yaml:",omitempty" mapstructure:"elasticapm"`
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

func (c *Config) DataStore() (*TracingBackendDataStoreConfig, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	selectedStore := c.config.Server.Telemetry.DataStore
	dataStoreConfig, found := c.config.Telemetry.DataStores[selectedStore]

	if selectedStore != "" && !found {
		return nil, ErrInvalidTraceDBProvider
	}

	if !found {
		return nil, nil
	}

	return &dataStoreConfig, nil
}

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
