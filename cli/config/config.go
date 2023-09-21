package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/goware/urlx"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	Version              = "dev"
	Env                  = "dev"
	DefaultCloudEndpoint = "http://localhost:3000/"
	DefaultCloudDomain   = "tracetest.io"
	DefaultCloudPath     = "/"
)

type ConfigFlags struct {
	Endpoint       string
	OrganizationID string
	EnvironmentID  string
	CI             bool
	AgentApiKey    string
}

type Config struct {
	Scheme         string  `yaml:"scheme"`
	Endpoint       string  `yaml:"endpoint"`
	ServerPath     *string `yaml:"serverPath,omitempty"`
	OrganizationID string  `yaml:"organizationID,omitempty"`
	EnvironmentID  string  `yaml:"environmentID,omitempty"`
	Token          string  `yaml:"token,omitempty"`
	Jwt            string  `yaml:"jwt,omitempty"`
	AgentApiKey    string  `yaml:"-"`

	// cloud config
	CloudAPIEndpoint string `yaml:"-"`
	AgentEndpoint    string `yaml:"agentEndpoint,omitempty"`
	UIEndpoint       string `yaml:"uIEndpoint,omitempty"`
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

	if pathPrefix == "/" {
		return ""
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

	if config.CloudAPIEndpoint == "" {
		config.CloudAPIEndpoint = DefaultCloudEndpoint
	}

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

	return config, nil
}

func ValidateServerURL(serverURL string) error {
	if !strings.HasPrefix(serverURL, "http://") && !strings.HasPrefix(serverURL, "https://") {
		return fmt.Errorf(`the server URL must start with the scheme, either "http://" or "https://"`)
	}

	return nil
}

func ParseServerURL(serverURL string) (scheme, endpoint string, serverPath *string, err error) {
	url, err := urlx.Parse(serverURL)
	if err != nil {
		return "", "", nil, fmt.Errorf("could not parse server URL: %w", err)
	}

	var path *string
	if url.Path != "" {
		path = &url.Path
	}

	return url.Scheme, url.Host, path, nil
}

func Save(config Config) error {
	configPath, err := GetConfigurationPath()
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

func GetConfigurationPath() (string, error) {
	configPath := "./config.yml"
	if _, err := os.Stat("config.yml"); os.IsNotExist(err) {
		homePath, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not get user home dir: %w", err)
		}

		configPath = path.Join(homePath, ".tracetest/config.yml")
	}

	return configPath, nil
}
