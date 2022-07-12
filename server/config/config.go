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
		Server                  ServerConfig    `mapstructure:"server"`
		TracingBackend          TracingBackend  `mapstructure:"tracingBackend"`
		PostgresConnString      string          `mapstructure:"postgresConnString"`
		PoolingConfig           PoolingConfig   `mapstructure:"poolingConfig"`
		GA                      GoogleAnalytics `mapstructure:"googleAnalytics"`
		PoolingRetryDelayString string          `mapstructure:"poolingRetryDelay"`
		Telemetry               TelemetryConfig `mapstructure:"telemetry"`
	}

	TracingBackend struct {
		DataStore TracingBackendDataStoreConfig `mapstructure:"dataStore"`
	}

	TracingBackendDataStoreConfig struct {
		Type   string                        `mapstructure:"type"`
		Jaeger configgrpc.GRPCClientSettings `mapstructure:"jaeger"`
		Tempo  configgrpc.GRPCClientSettings `mapstructure:"tempo"`
	}

	ServerConfig struct {
		PathPrefix string `mapstructure:"pathPrefix"`
		HttpPort   int    `mapstructure:"httpPort"`
	}

	GoogleAnalytics struct {
		Enabled bool `mapstructure:"enabled"`
	}

	PoolingConfig struct {
		MaxWaitTimeForTrace string `mapstructure:"maxWaitTimeForTrace"`
		RetryDelay          string `mapstructure:"retryDelay"`
	}

	TelemetryConfig struct {
		Enabled               bool    `mapstructure:"enabled"`
		ServiceName           string  `mapstructure:"serviceName"`
		Sampling              float64 `mapstructure:"sampling"`
		OTelCollectorEndpoint string  `mapstructure:"otelCollectorEndpoint"`
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
