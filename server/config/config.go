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
		PostgresConnString     string                         `mapstructure:"postgresConnString"`
		JaegerConnectionConfig *configgrpc.GRPCClientSettings `mapstructure:"jaegerConnectionConfig"`
		TempoConnectionConfig  *configgrpc.GRPCClientSettings `mapstructure:"tempoConnectionConfig"`
		MaxWaitTimeForTrace    string                         `mapstructure:"maxWaitTimeForTrace"`
		GA                     GoogleAnalytics                `mapstructure:"googleAnalytics"`
	}

	GoogleAnalytics struct {
		MeasurementID string `mapstructure:"measurementId"`
		SecretKey     string `mapstructure:"secretKey"`
		Enabled       bool   `mapstructure:"enabled"`
	}
)

func (c Config) MaxWaitTimeForTraceDuration() time.Duration {
	maxWaitTimeForTrace, err := time.ParseDuration(c.MaxWaitTimeForTrace)
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
