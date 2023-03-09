package cmd

import (
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var cliConfig config.Config
var cliLogger *zap.Logger
var resourceRegistry = actions.NewResourceRegistry()

type setupConfig struct {
	shouldValidateConfig bool
}

type setupOption func(*setupConfig)

func SkipConfigValidation() setupOption {
	return func(sc *setupConfig) {
		sc.shouldValidateConfig = false
	}
}

func setupCommand(options ...setupOption) func(cmd *cobra.Command, args []string) {
	config := setupConfig{
		shouldValidateConfig: true,
	}
	for _, option := range options {
		option(&config)
	}

	return func(cmd *cobra.Command, args []string) {
		setupOutputFormat()
		setupLogger(cmd, args)
		loadConfig(cmd, args)
		overrideConfig()

		apiClient := utils.GetAPIClient(cliConfig)

		options := []actions.ResourceArgsOption{actions.WithClient(apiClient), actions.WithLogger(cliLogger), actions.WithConfig(cliConfig)}
		configActions := actions.NewConfigActions(options...)
		resourceRegistry.Register(actions.SupportedResourceConfig, configActions)

		if config.shouldValidateConfig {
			validateConfig(cmd, args)
		}

		analytics.Init(cliConfig)
	}
}

func overrideConfig() {
	if overrideEndpoint != "" {
		scheme, endpoint, err := config.ParseServerURL(overrideEndpoint)
		if err != nil {
			msg := fmt.Sprintf("cannot parse endpoint %s", overrideEndpoint)
			cliLogger.Error(msg, zap.Error(err))
			os.Exit(1)
		}
		cliConfig.Scheme = scheme
		cliConfig.Endpoint = endpoint
	}
}

func setupOutputFormat() {
	o := formatters.Output(output)
	if !formatters.ValidOutput(o) {
		fmt.Fprintf(os.Stderr, "Invalid output format %s. Available formats are [%s]\n", output, outputFormatsString)
		os.Exit(1)
	}
	formatters.SetOutput(o)
}

func loadConfig(cmd *cobra.Command, args []string) {
	config, err := config.LoadConfig(configFile)
	if err != nil {
		cliLogger.Fatal("could not load config", zap.Error(err))
	}

	cliConfig = config
}

func validateConfig(cmd *cobra.Command, args []string) {
	if cliConfig.IsEmpty() {
		cliLogger.Warn("You haven't configured your CLI, some commands might fail!")
		cliLogger.Warn("Run 'tracetest configure' to configure your CLI")
	}
}

func setupLogger(cmd *cobra.Command, args []string) {
	atom := zap.NewAtomicLevel()
	if verbose {
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

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))

	cliLogger = logger
}

func teardownCommand(cmd *cobra.Command, args []string) {
	cliLogger.Sync()
	analytics.Close()
}
