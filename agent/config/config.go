package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	APIKey    string `mapstructure:"api_key"`
	Name      string `mapstructure:"agent_name"`
	ServerURL string `mapstructure:"server_url"`

	OTLPServer otlpServer `mapstructure:"otlp_server"`
}

type otlpServer struct {
	GRPCPort int `mapstructure:"grpc_port"`
	HTTPPort int `mapstructure:"http_port"`
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

	vp.SetDefault("AGENT_NAME", getHostname())
	vp.SetDefault("API_KEY", "")
	vp.SetDefault("SERVER_URL", "https://cloud.tracetest.io")
	vp.SetDefault("OTLP_SERVER.GRPC_PORT", 4317)
	vp.SetDefault("OTLP_SERVER.HTTP_PORT", 4318)

	config := Config{}

	vp.ReadInConfig()

	err := vp.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("could not load config: %w", err)
	}

	if config.Name == "" {
		return Config{}, fmt.Errorf("invalid host name, use the environment variable TRACETEST_AGENT_NAME to name your agent")
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

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		// Don't fail yet because user still can name the agent using the TRACETEST_AGENT_NAME
		// env variable.
		return ""
	}

	return hostname
}
