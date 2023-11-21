package runner

import (
	"context"

	agentConfig "github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/ui"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
)

type Runner struct {
	configurator config.Configurator
	resources    *resourcemanager.Registry
	ui           ui.ConsoleUI
	mode         agentConfig.Mode
}

func NewRunner(configurator config.Configurator, resources *resourcemanager.Registry, ui ui.ConsoleUI) *Runner {
	return &Runner{
		configurator: configurator,
		resources:    resources,
		ui:           ui,
		mode:         agentConfig.Mode_Desktop,
	}
}

func (s *Runner) Run(ctx context.Context, cfg config.Config, flags agentConfig.Flags) error {
	s.ui.Banner(config.Version)
	s.ui.Println(`Tracetest start launches a lightweight agent. It enables you to run tests and collect traces with Tracetest.
Once started, Tracetest Agent exposes OTLP ports 4317 and 4318 to ingest traces via gRCP and HTTP.`)

	if flags.Token == "" || flags.AgentApiKey != "" {
		s.configurator = s.configurator.WithOnFinish(s.onStartAgent)
	}

	s.ui.Infof("Running in %s mode...", s.mode)

	s.mode = flags.Mode

	return s.configurator.Start(ctx, cfg, flags)
}

func (s *Runner) onStartAgent(ctx context.Context, cfg config.Config) {
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

func (s *Runner) StartAgent(ctx context.Context, endpoint, agentApiKey, uiEndpoint string) error {
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

	if s.mode == agentConfig.Mode_Desktop {
		return RunDesktopStrategy(ctx, cfg, s.ui, uiEndpoint)
	}

	// TODO: Add verbose strategy
	return nil
}
