package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
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
	resourceParams = &resourceParameters{}
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

var httpClient = &resourcemanager.HTTPClient{}
var resources = resourcemanager.NewRegistry().
	Register(
		resourcemanager.NewClient(
			httpClient,
			"config", "configs",
			resourcemanager.WithTableConfig(resourcemanager.TableConfig{
				Cells: []resourcemanager.TableCellConfig{
					{Header: "ID", Path: "spec.id"},
					{Header: "NAME", Path: "spec.name"},
					{Header: "ANALYTICS ENABLED", Path: "spec.analyticsEnabled"},
				},
			}),
		),
	).
	Register(
		resourcemanager.NewClient(
			httpClient,
			"analyzer", "analyzers",
			resourcemanager.WithTableConfig(resourcemanager.TableConfig{
				Cells: []resourcemanager.TableCellConfig{
					{Header: "ID", Path: "spec.id"},
					{Header: "NAME", Path: "spec.name"},
					{Header: "ENABLED", Path: "spec.enabled"},
					{Header: "MINIMUM SCORE", Path: "spec.minimumScore"},
					{Header: "PLUGINS", Path: "spec.total.plugins"},
				},
				ItemModifier: func(item *gabs.Container) error {
					item.SetP(len(item.Path("spec.plugins").Children()), "spec.total.plugins")

					return nil
				},
			}),
		),
	).
	Register(
		resourcemanager.NewClient(
			httpClient,
			"pollingprofile", "pollingprofiles",
			resourcemanager.WithTableConfig(resourcemanager.TableConfig{
				Cells: []resourcemanager.TableCellConfig{
					{Header: "ID", Path: "spec.id"},
					{Header: "NAME", Path: "spec.name"},
					{Header: "STRATEGY", Path: "spec.strategy"},
				},
			}),
		),
	).
	Register(
		resourcemanager.NewClient(
			httpClient,
			"demo", "demos",
			resourcemanager.WithTableConfig(resourcemanager.TableConfig{
				Cells: []resourcemanager.TableCellConfig{
					{Header: "ID", Path: "spec.id"},
					{Header: "NAME", Path: "spec.name"},
					{Header: "TYPE", Path: "spec.type"},
					{Header: "ENABLED", Path: "spec.enabled"},
				},
			}),
		),
	).
	Register(
		resourcemanager.NewClient(
			httpClient,
			"datastore", "datastores",
			resourcemanager.WithTableConfig(resourcemanager.TableConfig{
				Cells: []resourcemanager.TableCellConfig{
					{Header: "ID", Path: "spec.id"},
					{Header: "NAME", Path: "spec.name"},
					{Header: "DEFAULT", Path: "spec.default"},
				},
				ItemModifier: func(item *gabs.Container) error {
					isDefault := item.Path("spec.default").Data().(bool)
					if !isDefault {
						item.SetP("", "spec.default")
					} else {
						item.SetP("*", "spec.default")
					}
					return nil
				},
			}),
			resourcemanager.WithDeleteEnabled("DataStore removed. Defaulting back to no-tracing mode"),
		),
	).
	Register(
		resourcemanager.NewClient(
			httpClient,
			"environment", "environments",
			resourcemanager.WithTableConfig(resourcemanager.TableConfig{
				Cells: []resourcemanager.TableCellConfig{
					{Header: "ID", Path: "spec.id"},
					{Header: "NAME", Path: "spec.name"},
					{Header: "DESCRIPTION", Path: "spec.description"},
				},
			}),
		),
	).
	Register(
		resourcemanager.NewClient(
			httpClient,
			"transaction", "transactions",
			resourcemanager.WithTableConfig(resourcemanager.TableConfig{
				Cells: []resourcemanager.TableCellConfig{
					{Header: "ID", Path: "spec.id"},
					{Header: "NAME", Path: "spec.name"},
					{Header: "VERSION", Path: "spec.version"},
					{Header: "STEPS", Path: "spec.summary.steps"},
					{Header: "RUNS", Path: "spec.summary.runs"},
					{Header: "LAST RUN TIME", Path: "spec.summary.lastRun.time"},
					{Header: "LAST RUN SUCCESSES", Path: "spec.summary.lastRun.passes"},
					{Header: "LAST RUN FAILURES", Path: "spec.summary.lastRun.fails"},
				},
				ItemModifier: func(item *gabs.Container) error {
					// set spec.summary.steps to the number of steps in the transaction
					item.SetP(len(item.Path("spec.steps").Children()), "spec.summary.steps")

					// if lastRun.time is not empty, show it in a nicer format
					lastRunTime := item.Path("spec.summary.lastRun.time").Data().(string)
					if lastRunTime != "" {
						date, err := time.Parse(time.RFC3339, lastRunTime)
						if err != nil {
							return fmt.Errorf("failed to parse last run time: %s", err)
						}
						if date.IsZero() {
							item.SetP("", "spec.summary.lastRun.time")
						} else {
							item.SetP(date.Format(time.DateTime), "spec.summary.lastRun.time")
						}
					}
					return nil
				},
			}),
		),
	)

func resourceList() string {
	return strings.Join(resources.List(), "|")
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

		extraHeaders := http.Header{}
		extraHeaders.Set("x-client-id", analytics.ClientID())
		extraHeaders.Set("x-source", "cli")

		// To avoid a ciruclar reference initialization when setting up the registry and its resources,
		// we create the resources with a pointer to an unconfigured HTTPClient.
		// When each command is run, this function is run in the PreRun stage, before any of the actual `Run` code is executed
		// We take this chance to configure the HTTPClient with the correct URL and headers.
		// To make this configuration propagate to all the resources, we need to replace the pointer to the HTTPClient.
		// For more details, see https://github.com/kubeshop/tracetest/pull/2832#discussion_r1245616804
		hc := resourcemanager.NewHTTPClient(cliConfig.URL(), extraHeaders)
		*httpClient = *hc

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
