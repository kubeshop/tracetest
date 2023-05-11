package global_setup

import (
	"context"
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/config"
	global_parameters "github.com/kubeshop/tracetest/cli/global/parameters"
	misc_actions "github.com/kubeshop/tracetest/cli/misc/actions"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Setup struct {
	Config      *config.Config
	Logger      *zap.Logger
	VersionText *string
	Parameters  *global_parameters.Global
	options     []SetupOption
}

type SetupOption = func(setup *Setup, args []string) error

func NewSetup(parameters *global_parameters.Global, options ...SetupOption) *Setup {
	return &Setup{
		Parameters: parameters,
		options:    options,
	}
}

func (s *Setup) PreRun(cmd *cobra.Command, args []string) (string, error) {
	for _, option := range s.options {
		err := option(s, args)
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

func WithInitAnalytics() SetupOption {
	return func(setup *Setup, _ []string) error {
		analytics.Init(*setup.Config)
		return nil
	}
}

func WithLogger() SetupOption {
	return func(setup *Setup, _ []string) error {
		atom := zap.NewAtomicLevel()
		if setup.Parameters.Verbose {
			atom.SetLevel(zap.DebugLevel)
		} else {
			atom.SetLevel(zap.WarnLevel)
		}

		encoderCfg := zapcore.EncoderConfig{
			TimeKey:        zapcore.OmitKey,
			LevelKey:       "level",
			NameKey:        zapcore.OmitKey,
			CallerKey:      zapcore.OmitKey,
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			StacktraceKey:  zapcore.OmitKey,
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		setup.Logger = zap.New(zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			atom,
		))

		return nil
	}
}

func WithConfig() SetupOption {
	return func(setup *Setup, _ []string) error {
		cfg, err := config.LoadConfig(setup.Parameters.ConfigFile)
		if err != nil {
			setup.Logger.Fatal("could not load config", zap.Error(err))
		}

		setup.Config = &cfg

		if setup.Parameters.OverrideEndpoint != "" {
			scheme, endpoint, err := config.ParseServerURL(setup.Parameters.OverrideEndpoint)
			if err != nil {
				msg := fmt.Sprintf("cannot parse endpoint %s", setup.Parameters.OverrideEndpoint)
				setup.Logger.Error(msg, zap.Error(err))
				os.Exit(1)
			}
			setup.Config.Scheme = scheme
			setup.Config.Endpoint = endpoint
		}

		return nil
	}
}

func WithVersionText() SetupOption {
	return func(setup *Setup, _ []string) error {
		ctx := context.Background()
		options := []misc_actions.ActionArgsOption{
			misc_actions.ActionWithClient(utils.GetAPIClient(*setup.Config)),
			misc_actions.ActionWithConfig(*setup.Config),
			misc_actions.ActionWithLogger(setup.Logger),
		}

		action := misc_actions.NewGetServerVersionAction(options...)
		version := action.GetVersionText(ctx)

		setup.VersionText = &version
		return nil
	}
}
