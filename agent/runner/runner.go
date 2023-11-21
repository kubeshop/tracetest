package runner

import (
	"context"
	"fmt"
	"os"

	agentConfig "github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/ui"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"

	"go.uber.org/zap"
)

type Runner struct {
	configurator config.Configurator
	resources    *resourcemanager.Registry
	ui           ui.ConsoleUI
	mode         agentConfig.Mode
	logger       *zap.Logger
}

func NewRunner(configurator config.Configurator, resources *resourcemanager.Registry, ui ui.ConsoleUI) *Runner {
	return &Runner{
		configurator: configurator,
		resources:    resources,
		ui:           ui,
		mode:         agentConfig.Mode_Desktop,
		logger:       nil,
	}
}

func (s *Runner) Run(ctx context.Context, cfg config.Config, flags agentConfig.Flags) error {
	s.ui.Banner(config.Version)
	s.ui.Println(`Tracetest start launches a lightweight agent. It enables you to run tests and collect traces with Tracetest.
Once started, Tracetest Agent exposes OTLP ports 4317 and 4318 to ingest traces via gRCP and HTTP.`)
	s.ui.Println("") // print empty line

	if flags.Token == "" || flags.AgentApiKey != "" {
		s.configurator = s.configurator.WithOnFinish(s.onStartAgent)
	}

	s.mode = flags.Mode
	s.ui.Infof("Running in %s mode...", s.mode)

	logger := zap.NewNop()

	if enableLogging(flags.LogLevel) {
		var err error
		logger, err = zap.NewDevelopment()
		if err != nil {
			return fmt.Errorf("could not create logger: %w", err)
		}
	}

	s.logger = logger

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
		return s.RunDesktopStrategy(ctx, cfg, uiEndpoint)
	}

	return s.RunVerboseStrategy(ctx, cfg)
}

func enableLogging(logLevel string) bool {
	return os.Getenv("TRACETEST_DEV") == "true" && logLevel == "debug"
}
