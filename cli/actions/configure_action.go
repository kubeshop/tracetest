package actions

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/ui"
	"gopkg.in/yaml.v3"
)

type ConfigureConfig struct {
	Global    bool
	SetValues ConfigureConfigSetValues
}

type ConfigureConfigSetValues struct {
	Endpoint string
}

type configureAction struct {
	config config.Config
}

func NewConfigureAction(config config.Config) configureAction {
	return configureAction{
		config: config,
	}
}

func (a configureAction) Run(ctx context.Context, args ConfigureConfig) error {
	ui := ui.DefaultUI
	existingConfig := a.loadExistingConfig(args)

	var serverURL string
	if args.SetValues.Endpoint != "" {
		serverURL = args.SetValues.Endpoint
	} else {
		serverURL = ui.TextInput("Enter your Tracetest server URL", existingConfig.URL())
	}

	if err := config.ValidateServerURL(serverURL); err != nil {
		return err
	}

	scheme, endpoint, err := config.ParseServerURL(serverURL)
	if err != nil {
		return err
	}

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
