package actions

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type ConfigureConfig struct {
	Global bool
}

type configureAction struct {
	config config.Config
	logger *zap.Logger
	client *openapi.APIClient
}

var _ Action[ExportTestConfig] = &exportTestAction{}

func NewConfigureAction(config config.Config, logger *zap.Logger, client *openapi.APIClient) configureAction {
	return configureAction{
		config: config,
		logger: logger,
		client: client,
	}
}

func (a configureAction) Run(ctx context.Context, args ConfigureConfig) error {
	existingConfig := a.loadExistingConfig(args)
	var serverURL string
	if existingConfig.Scheme != "" && existingConfig.Endpoint != "" {
		fmt.Printf("Enter your Tracetest server URL (current value: %s://%s): ", existingConfig.Scheme, existingConfig.Endpoint)
	} else {
		fmt.Print("Enter your Tracetest server URL: ")
	}
	_, err := fmt.Fscanf(os.Stdin, "%s", &serverURL)
	if err != nil {
		return fmt.Errorf("could not read text from input: %w", err)
	}

	a.logger.Debug("user entered new server URL", zap.String("serverURL", serverURL))

	if !strings.HasPrefix(serverURL, "http://") && !strings.HasPrefix(serverURL, "https://") {
		return fmt.Errorf(`the server URL must start with the scheme, either "http://" or "https://"`)
	}

	urlParts := strings.Split(serverURL, "://")
	if len(urlParts) != 2 {
		return fmt.Errorf("invalid server url")
	}

	scheme := urlParts[0]
	endpoint := urlParts[1]

	config := config.Config{
		Scheme:   scheme,
		Endpoint: endpoint,
	}

	err = a.saveConfiguration(ctx, config, args)
	if err != nil {
		return fmt.Errorf("could not save configuration: %w", err)
	}

	return nil
}

func (a configureAction) loadExistingConfig(args ConfigureConfig) config.Config {
	configPath, err := a.getConfigurationPath(args)
	if err != nil {
		return config.Config{}
	}

	c, err := config.LoadConfig(configPath)
	if err != nil {
		return config.Config{}
	}

	return c
}

func (a configureAction) saveConfiguration(ctx context.Context, config config.Config, args ConfigureConfig) error {
	configPath, err := a.getConfigurationPath(args)
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

func (a configureAction) getConfigurationPath(args ConfigureConfig) (string, error) {
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
