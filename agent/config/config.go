package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DevMode   bool   `mapstructure:"dev_mode"`
	APIKey    string `mapstructure:"api_key"`
	Name      string `mapstructure:"name"`
	ServerURL string `mapstructure:"server_url"`
}

func LoadConfig() (Config, error) {
	vp := viper.NewWithOptions(
		viper.EnvKeyReplacer(strings.NewReplacer(".", "_")),
	)

	homeFolder, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("could not get user home folder")
	}
	tracetestFolder := path.Join(homeFolder, ".tracetest")

	vp.SetEnvPrefix("tracetest")
	vp.AddConfigPath(tracetestFolder)
	vp.AddConfigPath("tracetest-agent.yaml")
	vp.SetConfigName("agent")
	vp.SetConfigType("env")
	vp.AutomaticEnv()

	vp.SetDefault("DEV_MODE", false)
	vp.SetDefault("SERVER_URL", "https://cloud.tracetest.io")

	config := Config{}

	vp.ReadInConfig()

	err = vp.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("could not load config: %w", err)
	}

	return config, nil
}
