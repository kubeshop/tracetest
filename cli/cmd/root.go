package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/spf13/cobra"
)

var (
	verbose    bool
	configFile string
	output     string

	outputFormats       = formatters.OuputsStr()
	outputFormatsString = strings.Join(outputFormats, "|")

	// overrides
	overrideEndpoint string
)

var rootCmd = &cobra.Command{
	Use:     "tracetest",
	Short:   "CLI to configure, install and execute tests on a Tracetest server",
	Long:    `CLI to configure, install and execute tests on a Tracetest server`,
	PreRun:  setupCommand(),
	PostRun: teardownCommand,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var (
	cmdGroupResources = &cobra.Group{
		ID:    "resources",
		Title: "Resources",
	}

	cmdGroupServer = &cobra.Group{
		ID:    "server",
		Title: "Server",
	}

	cmdGroupCLIConfig = &cobra.Group{
		ID:    "cli-config",
		Title: "CLI Config",
	}

	cmdGroupTests = &cobra.Group{
		ID:    "tests",
		Title: "Tests",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", string(formatters.DefaultOutput), fmt.Sprintf("output format [%s]", outputFormatsString))
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yml", "config file will be used by the CLI")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "display debug information")

	rootCmd.PersistentFlags().StringVarP(&overrideEndpoint, "server-url", "s", "", "server url")

	rootCmd.AddGroup(
		cmdGroupResources,
		cmdGroupServer,
		cmdGroupCLIConfig,
		cmdGroupTests,
	)
}
