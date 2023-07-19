package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ServerURL string `json:"serverURL"`
	DevMode   bool   `json:"devMode"`
	APIKey    string `json:"apiKey"`
}

func LoadConfig() (Config, error) {
	viper.SetEnvPrefix("tracetest")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	config := Config{}

	err := viper.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("could not load config: %w", err)
	}

	return config, nil
}
