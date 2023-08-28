package actions

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/ui"
)

type configureAction struct {
	config config.Config
}

func NewConfigureAction(config config.Config) configureAction {
	return configureAction{
		config: config,
	}
}

func (a configureAction) Run(ctx context.Context, args config.ConfigureConfig) error {
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

	cfg := config.Config{
		Scheme:   scheme,
		Endpoint: endpoint,
	}

	err = config.Save(ctx, cfg, args)
	if err != nil {
		return fmt.Errorf("could not save configuration: %w", err)
	}

	return nil
}

func (a configureAction) loadExistingConfig(args config.ConfigureConfig) config.Config {
	configPath, err := config.GetConfigurationPath(args)
	if err != nil {
		return config.Config{}
	}

	c, err := config.LoadConfig(configPath)
	if err != nil {
		return config.Config{}
	}

	return c
}
