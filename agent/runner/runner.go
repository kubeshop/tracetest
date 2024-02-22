package runner

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
	agentConfig "github.com/kubeshop/tracetest/agent/config"
	"github.com/kubeshop/tracetest/agent/event"
	"github.com/kubeshop/tracetest/agent/ui"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/pkg/oauth"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Runner struct {
	configurator config.Configurator
	resources    *resourcemanager.Registry
	ui           ui.ConsoleUI
	mode         agentConfig.Mode
	logger       *zap.Logger
	loggerLevel  *zap.AtomicLevel
	claims       jwt.MapClaims
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

func (s *Runner) Run(ctx context.Context, cfg config.Config, flags agentConfig.Flags, verbose bool) error {
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
		atom := zap.NewAtomicLevel()
		if verbose {
			atom.SetLevel(zapcore.DebugLevel)
		}

		logger, err = zap.NewDevelopment()
		if err != nil {
			return fmt.Errorf("could not create logger: %w", err)
		}

		logger = logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
				zapcore.Lock(os.Stdout),
				atom,
			)
		}))

		s.loggerLevel = &atom
	}

	s.logger = logger
	s.configurator = s.configurator.WithLogger(logger)
	oauth.SetLogger(logger)

	s.logger.Debug("Starting agent with flags", zap.Any("flags", flags))

	return s.configurator.Start(ctx, &cfg, flags)
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

	if env.AgentApiKey == "" {
		s.ui.Error("You are attempting to start the agent in an environment you do not have admin rights to. Please ask the administrator of this environment to grant you admin rights.")
		return
	}

	err = s.StartAgent(ctx, cfg.AgentEndpoint, env.AgentApiKey, cfg.UIEndpoint)
	if err != nil {
		s.ui.Error(err.Error())
	}
}

func (s *Runner) StartAgent(ctx context.Context, endpoint, agentApiKey, uiEndpoint string) error {
	cfg, err := agentConfig.LoadConfig()
	s.logger.Debug("Loaded agent config", zap.Any("config", cfg))
	if err != nil {
		s.logger.Error("Could not load agent config", zap.Error(err))
		return err
	}

	if endpoint != "" {
		s.logger.Debug("Overriding agent endpoint", zap.String("endpoint", endpoint))
		cfg.ServerURL = endpoint
	}
	s.logger.Debug("Agent endpoint", zap.String("endpoint", cfg.ServerURL))

	if agentApiKey != "" {
		s.logger.Debug("Overriding agent api key", zap.String("apiKey", agentApiKey))
		cfg.APIKey = agentApiKey
	}
	s.logger.Debug("Agent api key", zap.String("apiKey", cfg.APIKey))

	if s.mode == agentConfig.Mode_Desktop {
		s.logger.Debug("Starting agent in desktop mode")
		return s.RunDesktopStrategy(ctx, cfg, uiEndpoint)
	}

	s.logger.Debug("Starting agent in verbose mode")
	return s.RunVerboseStrategy(ctx, cfg)
}

func enableLogging(logLevel string) bool {
	return os.Getenv("TRACETEST_DEV") == "true" && logLevel == "debug"
}

func (s *Runner) authenticate(ctx context.Context, cfg agentConfig.Config, observer event.Observer) (*Session, jwt.MapClaims, error) {
	isStarted := false
	session := &Session{}

	var err error

	for !isStarted {
		session, err = StartSession(ctx, cfg, observer, s.logger)
		if err != nil && errors.Is(err, ErrOtlpServerStart) {
			s.ui.Error("Tracetest Agent binds to the OpenTelemetry ports 4317 and 4318 which are used to receive trace information from your system. The agent tried to bind to these ports, but failed.")
			shouldRetry := s.ui.Enter("Please stop the process currently listening on these ports and press enter to try again.")

			if !shouldRetry {
				s.ui.Finish()
				return nil, nil, err
			}

			continue
		}

		if err != nil {
			return nil, nil, err
		}

		isStarted = true
	}

	claims, err := config.GetTokenClaims(session.Token)
	if err != nil {
		return nil, nil, err
	}
	s.claims = claims
	return session, claims, nil
}

func (s *Runner) getCurrentSessionClaims() jwt.MapClaims {
	return s.claims
}
