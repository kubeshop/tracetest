package starter

import (
	"context"
	"encoding/json"
	"fmt"

	agentConfig "github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/initialization"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/kubeshop/tracetest/cli/ui"
)

type Starter struct {
	configurator config.Configurator
	resources    *resourcemanager.Registry
	ui           ui.UI
}

func NewStarter(configurator config.Configurator, resources *resourcemanager.Registry) *Starter {
	ui := ui.DefaultUI
	return &Starter{configurator, resources, ui}
}

func (s *Starter) Run(ctx context.Context, cfg config.Config) error {
	flags := config.ConfigFlags{
		Endpoint: "http://localhost:8090/",
	}
	s.ui.Banner(config.Version)

	return s.configurator.WithOnFinish(s.onStartAgent).Start(ctx, cfg, flags)
}

func (s *Starter) onStartAgent(ctx context.Context, cfg config.Config) {
	env, err := s.getEnvironment(ctx, cfg)
	if err != nil {
		s.ui.Error(err.Error())
	}

	s.ui.Println(fmt.Sprintf("Connecting Agent to environment %s...", env.Name))
	err = startAgent(ctx, "localhost:8091", env.AgentApiKey)
	if err != nil {
		s.ui.Error(err.Error())
	}
}

type environment struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	AgentApiKey    string `json:"agentApiKey"`
	OrganizationID string `json:"organizationID"`
}

func (s *Starter) getEnvironment(ctx context.Context, cfg config.Config) (environment, error) {
	resource, err := s.resources.Get("env")
	if err != nil {
		return environment{}, err
	}

	resource = resource.
		WithHttpClient(config.SetupHttpClient(cfg)).
		WithOptions(resourcemanager.WithPrefixGetter(func() string {
			return fmt.Sprintf("/organizations/%s/", cfg.OrganizationID)
		}))

	resultFormat, err := resourcemanager.Formats.GetWithFallback("json", "json")
	if err != nil {
		return environment{}, err
	}

	raw, err := resource.Get(ctx, cfg.EnvironmentID, resultFormat)
	if err != nil {
		return environment{}, err
	}

	var env environment
	err = json.Unmarshal([]byte(raw), &env)
	if err != nil {
		return environment{}, err
	}

	return env, nil
}

func startAgent(ctx context.Context, endpoint, agentApiKey string) error {
	cfg := agentConfig.Config{
		ServerURL: endpoint,
		APIKey:    agentApiKey,
		Name:      "local",
	}

	return initialization.Start(ctx, cfg)
}
