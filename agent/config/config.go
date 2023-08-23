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
	Name      string `mapstructure:"agent_name"`
	ServerURL string `mapstructure:"server_url"`
}

func LoadConfig() (Config, error) {
	vp := viper.NewWithOptions(
		viper.EnvKeyReplacer(strings.NewReplacer(".", "_")),
	)

	tracetestFolder := getTracetestFolder()

	vp.SetEnvPrefix("tracetest")
	vp.AddConfigPath(tracetestFolder)
	vp.AddConfigPath("tracetest-agent.yaml")
	vp.SetConfigName("agent")
	vp.SetConfigType("env")
	vp.AutomaticEnv()

	vp.SetDefault("DEV_MODE", false)
	vp.SetDefault("AGENT_NAME", "")
	vp.SetDefault("API_KEY", "")
	vp.SetDefault("SERVER_URL", "https://cloud.tracetest.io")

	config := Config{}

	vp.ReadInConfig()

	err := vp.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("could not load config: %w", err)
	}

	return config, nil
}

func getTracetestFolder() string {
	homeFolder, err := os.UserHomeDir()
	if err != nil {
		// as a fallback, just return the current folder
		return "."
	}

	return path.Join(homeFolder, ".tracetest")
}
