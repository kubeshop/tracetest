package cmd

import (
	"context"
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

var (
	cliConfig      config.Config
	cliLogger      *zap.Logger
	versionText    string
	isVersionMatch bool
)

type setupConfig struct {
	shouldValidateConfig          bool
	shouldValidateVersionMismatch bool
}

type setupOption func(*setupConfig)

func SkipConfigValidation() setupOption {
	return func(sc *setupConfig) {
		sc.shouldValidateConfig = false
	}
}

func SkipVersionMismatchCheck() setupOption {
	return func(sc *setupConfig) {
		sc.shouldValidateVersionMismatch = false
	}
}

func setupCommand(options ...setupOption) func(cmd *cobra.Command, args []string) {
	config := setupConfig{
		shouldValidateConfig:          true,
		shouldValidateVersionMismatch: true,
	}
	for _, option := range options {
		option(&config)
	}

	return func(cmd *cobra.Command, args []string) {
		setupOutputFormat(cmd)
		setupLogger(cmd, args)
		loadConfig(cmd, args)
		overrideConfig()
		setupVersion()
		setupResources()

		if config.shouldValidateConfig {
			validateConfig(cmd, args)
		}

		if config.shouldValidateVersionMismatch {
			validateVersionMismatch()
		}

		analytics.Init()
	}
}

func overrideConfig() {
	if overrideEndpoint != "" {
		scheme, endpoint, err := config.ParseServerURL(overrideEndpoint)
		if err != nil {
			msg := fmt.Sprintf("cannot parse endpoint %s", overrideEndpoint)
			cliLogger.Error(msg, zap.Error(err))
			ExitCLI(1)
		}
		cliConfig.Scheme = scheme
		cliConfig.Endpoint = endpoint
	}
}

func setupOutputFormat(cmd *cobra.Command) {
	if cmd.GroupID != "resources" && output == string(formatters.Empty) {
		output = string(formatters.DefaultOutput)
	}

	o := formatters.Output(output)
	if !formatters.ValidOutput(o) {
		fmt.Fprintf(os.Stderr, "Invalid output format %s. Available formats are [%s]\n", output, outputFormatsString)
		ExitCLI(1)
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
}

func setupVersion() {
	versionText, isVersionMatch = actions.GetVersion(
		context.Background(),
		cliConfig,
		utils.GetAPIClient(cliConfig),
	)
}

func validateVersionMismatch() {
	if !isVersionMatch && os.Getenv("TRACETEST_DEV") == "" {
		fmt.Fprintf(os.Stderr, versionText+`
✖️ Error: Version Mismatch
The CLI version and the server version are not compatible. To fix this, you'll need to make sure that both your CLI and server are using compatible versions.
We recommend upgrading both of them to the latest available version. Check out our documentation https://docs.tracetest.io/configuration/upgrade for simple instructions on how to upgrade.
Thank you for using Tracetest! We apologize for any inconvenience caused.
`)
		ExitCLI(1)
	}
}
