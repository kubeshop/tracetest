package config

import (
	"fmt"
	"io/ioutil"

	"github.com/mitchellh/mapstructure"
	"go.opentelemetry.io/collector/config/configgrpc"
	"gopkg.in/yaml.v2"
)

type Config struct {
	PostgresConnString     string                         `mapstructure:"postgresConnString"`
	JaegerConnectionConfig *configgrpc.GRPCClientSettings `mapstructure:"jaegerConnectionConfig"`
	TempoConnectionConfig  *configgrpc.GRPCClientSettings `mapstructure:"tempoConnectionConfig"`
	MaxWaitTimeForTrace    string                         `mapstructure:"maxWaitTimeForTrace"`
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
