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
		PostgresConnString      string                         `mapstructure:"postgresConnString"`
		JaegerConnectionConfig  *configgrpc.GRPCClientSettings `mapstructure:"jaegerConnectionConfig"`
		TempoConnectionConfig   *configgrpc.GRPCClientSettings `mapstructure:"tempoConnectionConfig"`
		PoolingConfig           PoolingConfig                  `mapstructure:"poolingConfig"`
		GA                      GoogleAnalytics                `mapstructure:"googleAnalytics"`
		PoolingRetryDelayString string                         `mapstructure:"poolingRetryDelay"`
		Telemetry               TelemetryConfig                `mapstructure:"telemetry"`
		Http                    HttpServerConfig               `mapstructure:"http"`
	}

	GoogleAnalytics struct {
		Enabled bool `mapstructure:"enabled"`
	}

	PoolingConfig struct {
		MaxWaitTimeForTrace string `mapstructure:"maxWaitTimeForTrace"`
		RetryDelay          string `mapstructure:"retryDelay"`
	}

	TelemetryConfig struct {
		ServiceName string                `mapstructure:"serviceName"`
		Sampling    float64               `mapstructure:"sampling"`
		Exporters   []string              `mapstructure:"exporters"`
		Jaeger      JaegerTelemetryConfig `mapstructure:"jaeger"`
	}

	JaegerTelemetryConfig struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}

	HttpServerConfig struct {
		Host   string `mapstructure:"host"`
		Port   int    `mapstructure:"port"`
		Prefix string `mapstructure:"prefix"`
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
