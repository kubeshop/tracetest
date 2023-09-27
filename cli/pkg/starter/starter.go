package starter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
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

func (s *Starter) Run(ctx context.Context, cfg config.Config, flags config.ConfigFlags) error {
	s.ui.Banner(config.Version)

	return s.configurator.WithOnFinish(s.onStartAgent).Start(ctx, cfg, flags)
}

func (s *Starter) onStartAgent(ctx context.Context, cfg config.Config) {
	if cfg.AgentApiKey != "" {
		err := s.StartAgent(ctx, cfg.AgentEndpoint, cfg.AgentApiKey, cfg.UIEndpoint)
		if err != nil {
			s.ui.Error(err.Error())
		}

		return
	}

	env, err := s.getEnvironment(ctx, cfg)
	if err != nil {
		s.ui.Error(err.Error())
	}

	s.ui.Info(fmt.Sprintf(`
Connecting Agent with name %s to Organization %s and Environment %s
`, "local", cfg.OrganizationID, env.Name))

	err = s.StartAgent(ctx, cfg.AgentEndpoint, env.AgentApiKey, cfg.UIEndpoint)
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

func (s *Starter) StartAgent(ctx context.Context, endpoint, agentApiKey, uiEndpoint string) error {
	cfg, err := agentConfig.LoadConfig()
	if err != nil {
		return err
	}

	if endpoint != "" {
		cfg.ServerURL = endpoint
	}

	if agentApiKey != "" {
		cfg.APIKey = agentApiKey
	}

	s.ui.Info(fmt.Sprintf("Starting Agent with name %s...", cfg.Name))
	session, err := initialization.Start(ctx, cfg)

	claims, err := s.getTokenClaims(session.Token)
	if err != nil {
		return err
	}

	isOpen := true
	message := fmt.Sprintf("Agent is started! Leave the terminal open so tests can be run and traces gathered from this environment (%s). You can:", claims["environment_id"])
	for isOpen {
		options := []ui.Option{{
			Text: "Open Tracetest in a browser to this environment",
			Fn: func(_ ui.UI) {
				s.ui.OpenBrowser(fmt.Sprintf("%sorganizations/%s/environments/%s/dashboard", uiEndpoint, claims["organization_id"], claims["environment_id"]))
			},
		}, {
			Text: "Stop this agent",
			Fn: func(_ ui.UI) {
				isOpen = false
				session.Close()
			},
		}}

		selected := s.ui.Select(message, options, 0)
		selected.Fn(s.ui)
	}
	return nil
}

func (s *Starter) getTokenClaims(tokenString string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return jwt.MapClaims{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return jwt.MapClaims{}, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
