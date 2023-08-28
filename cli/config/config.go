package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	Version          = "dev"
	Env              = "dev"
	FrontendEndpoint = "http://localhost:3001"
)

type ConfigureConfig struct {
	Global    bool
	SetValues ConfigureConfigSetValues
}

type ConfigureConfigSetValues struct {
	Endpoint string
}

type Config struct {
	Scheme           string  `yaml:"scheme"`
	Endpoint         string  `yaml:"endpoint"`
	ServerPath       *string `yaml:"serverPath,omitempty"`
	FrontendEndpoint string  `yaml:"-"`
	OrganizationID   string  `yaml:"organizationID,omitempty"`
	EnvironmentID    string  `yaml:"environmentID,omitempty"`
	Token            string  `yaml:"token,omitempty"`
	Jwt              string  `yaml:"jwt,omitempty"`
}

func (c Config) URL() string {
	if c.Scheme == "" || c.Endpoint == "" {
		return ""
	}

	return fmt.Sprintf("%s://%s", c.Scheme, strings.TrimSuffix(c.Endpoint, "/"))
}

func (c Config) Path() string {
	pathPrefix := "/api"
	if c.ServerPath != nil {
		pathPrefix = *c.ServerPath
	}

	return pathPrefix
}

func (c Config) IsEmpty() bool {
	thisConfigJson, _ := json.Marshal(c)
	emptyConfigJson, _ := json.Marshal(Config{})

	return string(thisConfigJson) == string(emptyConfigJson)
}

func LoadConfig(configFile string) (Config, error) {
	config, err := loadConfig(configFile)
	if err != nil {
		return Config{}, err
	}

	if !config.IsEmpty() {
		return config, nil
	}

	config.FrontendEndpoint = FrontendEndpoint

	homePath, err := os.UserHomeDir()
	if err != nil {
		return Config{}, fmt.Errorf("could not get user home path")
	}

	globalConfigPath := filepath.Join(homePath, ".tracetest/config.yml")
	return loadConfig(globalConfigPath)
}

func loadConfig(configFile string) (Config, error) {
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("tracetest")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// No configuration file found, return empty config
		return Config{}, nil
	}

	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("could not unmarshal config: %w", err)
	}

	config.FrontendEndpoint = FrontendEndpoint
	return config, nil
}

func ValidateServerURL(serverURL string) error {
	if !strings.HasPrefix(serverURL, "http://") && !strings.HasPrefix(serverURL, "https://") {
		return fmt.Errorf(`the server URL must start with the scheme, either "http://" or "https://"`)
	}

	return nil
}

func ParseServerURL(serverURL string) (scheme, endpoint string, err error) {
	urlParts := strings.Split(serverURL, "://")
	if len(urlParts) != 2 {
		return "", "", fmt.Errorf("invalid server url")
	}

	scheme = urlParts[0]
	endpoint = strings.TrimSuffix(urlParts[1], "/")

	return scheme, endpoint, nil
}

func Save(ctx context.Context, config Config, args ConfigureConfig) error {
	configPath, err := GetConfigurationPath(args)
	if err != nil {
		return fmt.Errorf("could not get configuration path: %w", err)
	}

	configYml, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("could not marshal configuration into yml: %w", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(configPath), 0700) // Ensure folder exists
	}
	err = os.WriteFile(configPath, configYml, 0755)
	if err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	return nil
}

func GetConfigurationPath(args ConfigureConfig) (string, error) {
	configPath := "./config.yml"
	if args.Global {
		homePath, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not get user home dir: %w", err)
		}
		configPath = path.Join(homePath, ".tracetest/config.yml")
	}

	return configPath, nil
}
