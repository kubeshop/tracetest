package starter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	agentConfig "github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/ui"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
)

type Starter struct {
	configurator config.Configurator
	resources    *resourcemanager.Registry
	ui           ui.UI
}

func NewStarter(configurator config.Configurator, resources *resourcemanager.Registry, ui ui.UI) *Starter {
	return &Starter{configurator, resources, ui}
}

func (s *Starter) Run(ctx context.Context, cfg config.Config, flags agentConfig.Flags) error {
	s.ui.Banner(config.Version)
	s.ui.Println(`Tracetest start launches a lightweight agent. It enables you to run tests and collect traces with Tracetest.
Once started, Tracetest Agent exposes OTLP ports 4317 and 4318 to ingest traces via gRCP and HTTP.`)

	if flags.Token == "" || flags.AgentApiKey != "" {
		s.configurator = s.configurator.WithOnFinish(s.onStartAgent)
	}

	return s.configurator.Start(ctx, cfg, flags)
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

	isStarted := false
	session := &Session{}
	for !isStarted {
		session, err = Start(ctx, cfg)
		if err != nil && errors.Is(err, ErrOtlpServerStart) {
			s.ui.Error("Tracetest Agent binds to the OpenTelemetry ports 4317 and 4318 which are used to receive trace information from your system. The agent tried to bind to these ports, but failed.")
			shouldRetry := s.ui.Enter("Please stop the process currently listening on these ports and press enter to try again.")

			if !shouldRetry {
				s.ui.Finish()
				return err
			}

			continue
		}

		if err != nil {
			return err
		}

		isStarted = true
	}

	claims, err := config.GetTokenClaims(session.Token)
	if err != nil {
		return err
	}

	isOpen := true
	message := `Agent is started! Leave the terminal open so tests can be run and traces gathered from this environment.
You can`
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
			s.ui.Finish()
		},
	}}

	for isOpen {
		selected := s.ui.Select(message, options, 0)
		selected.Fn(s.ui)
	}
	return nil
}
