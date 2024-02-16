package config

import (
	"context"
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
	DefaultCloudEndpoint = "http://app.tracetest.io"
	DefaultCloudDomain   = "tracetest.io"
	DefaultCloudPath     = "/"
)

type Config struct {
	Scheme            string `yaml:"scheme"`
	Endpoint          string `yaml:"endpoint"`
	ServerPath        string `yaml:"serverPath,omitempty"`
	OrganizationID    string `yaml:"organizationID,omitempty"`
	EnvironmentID     string `yaml:"environmentID,omitempty"`
	Token             string `yaml:"token,omitempty"`
	Jwt               string `yaml:"jwt,omitempty"`
	AgentApiKey       string `yaml:"-"`
	EndpointOverriden bool   `yaml:"-"`

	// cloud config
	CloudAPIEndpoint string `yaml:"-"`
	AgentEndpoint    string `yaml:"agentEndpoint,omitempty"`
	UIEndpoint       string `yaml:"uIEndpoint,omitempty"`
}

func (c Config) OAuthEndpoint() string {
	return fmt.Sprintf("%s%s", c.URL(), c.Path())

}

func (c Config) URL() string {
	if c.Scheme == "" || c.Endpoint == "" {
		return ""
	}

	return fmt.Sprintf("%s://%s", c.Scheme, strings.TrimSuffix(c.Endpoint, "/"))
}

func (c Config) FullURL() string {
	return fmt.Sprintf("%s%s", c.URL(), c.ServerPath)
}

func (c Config) UI() string {
	if c.UIEndpoint != "" && !c.EndpointOverriden {
		return fmt.Sprintf("%s/organizations/%s/environments/%s", strings.TrimSuffix(c.UIEndpoint, "/"), c.OrganizationID, c.EnvironmentID)
	}

	return c.URL()
}

func (c Config) Path() string {
	pathPrefix := "/api"
	if c.ServerPath != "" {
		pathPrefix = c.ServerPath
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

func ParseServerURL(serverURL string) (scheme, endpoint, serverPath string, err error) {
	url, err := urlx.Parse(serverURL)
	if err != nil {
		return "", "", "", fmt.Errorf("could not parse server URL: %w", err)
	}

	return url.Scheme, url.Host, url.Path, nil
}

type orgIDKeyType struct{}
type envIDKeyType struct{}

var orgIDKey = orgIDKeyType{}
var envIDKey = envIDKeyType{}

func ContextWithOrganizationID(ctx context.Context, orgID string) context.Context {
	return context.WithValue(ctx, orgIDKey, orgID)
}

func ContextWithEnvironmentID(ctx context.Context, envID string) context.Context {
	return context.WithValue(ctx, envIDKey, envID)
}

func ContextGetOrganizationID(ctx context.Context) string {
	v := ctx.Value(orgIDKey)
	if v == nil {
		return ""
	}
	return v.(string)
}

func ContextGetEnvironmentID(ctx context.Context) string {
	v := ctx.Value(envIDKey)
	if v == nil {
		return ""
	}
	return v.(string)
}

func Save(ctx context.Context, config Config) (context.Context, error) {
	configPath, err := GetConfigurationPath()
	if err != nil {
		return ctx, fmt.Errorf("could not get configuration path: %w", err)
	}

	configYml, err := yaml.Marshal(config)
	if err != nil {
		return ctx, fmt.Errorf("could not marshal configuration into yml: %w", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(configPath), 0700) // Ensure folder exists
	}
	err = os.WriteFile(configPath, configYml, 0755)
	if err != nil {
		return ctx, fmt.Errorf("could not write file: %w", err)
	}

	ctx = ContextWithOrganizationID(ctx, config.OrganizationID)
	ctx = ContextWithEnvironmentID(ctx, config.EnvironmentID)

	return ctx, nil
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
