package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/cmdutil"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/oauth"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	cliLogger      = &zap.Logger{}
	cliConfig      config.Config
	openapiClient  = &openapi.APIClient{}
	versionText    string
	isVersionMatch bool
)

type setupConfig struct {
	shouldValidateConfig          bool
	shouldValidateVersionMismatch bool
	optionalResourceName          bool
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

func WithOptionalResourceName() setupOption {
	return func(sc *setupConfig) {
		sc.optionalResourceName = true
	}
}

func setupCommand(options ...setupOption) func(cmd *cobra.Command, args []string) {
	config := setupConfig{
		shouldValidateConfig:          true,
		shouldValidateVersionMismatch: true,
		optionalResourceName:          false,
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
		setupRunners()

		if config.shouldValidateConfig {
			validateConfig(cmd, args)
		}

		if config.shouldValidateVersionMismatch {
			validateVersionMismatch()
		}

		if config.optionalResourceName {
			resourceParams.optional = true
		}

		analytics.Init()
	}
}

func overrideConfig() {
	if overrideEndpoint != "" {
		scheme, endpoint, _, err := config.ParseServerURL(overrideEndpoint)
		if err != nil {
			msg := fmt.Sprintf("cannot parse endpoint %s", overrideEndpoint)
			cliLogger.Error(msg, zap.Error(err))
			ExitCLI(1)
		}
		cliConfig.Scheme = scheme
		cliConfig.Endpoint = endpoint
		cliConfig.EndpointOverriden = true
	}
}

func setupRunners() {
	c := config.GetAPIClient(cliConfig)
	*openapiClient = *c
}

func setupOutputFormat(cmd *cobra.Command) {
	o := formatters.Output(output)
	if output == "" {
		o = formatters.Pretty
	}
	if !formatters.ValidOutput(o) {
		fmt.Fprintf(os.Stderr, "Invalid output format %s. Available formats are [%s]\n", output, outputFormatsString)
		ExitCLI(1)
	}
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
	l := cmdutil.GetLogger(cmdutil.WithVerbose(verbose))
	*cliLogger = *l
	oauth.SetLogger(l)
}

func teardownCommand(cmd *cobra.Command, args []string) {
	cliLogger.Sync()
}

func setupVersion() {
	versionText, isVersionMatch = config.GetVersion(context.Background(), cliConfig)
}

func validateVersionMismatch() {
	if !isVersionMatch && os.Getenv("TRACETEST_DEV") == "" {
		fmt.Fprintf(os.Stderr, versionText+config.ErrVersionMismatch.Error())
		ExitCLI(1)
	}
}
