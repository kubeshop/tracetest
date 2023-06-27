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
	overrideEndpoint   string
	cliExitInterceptor func(code int)
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
		ExitCLI(1)
	}
}

func ExitCLI(errorCode int) {
	if cliExitInterceptor != nil {
		cliExitInterceptor(errorCode)
		return
	}

	os.Exit(errorCode)
}

func RegisterCLIExitInterceptor(interceptor func(int)) {
	cliExitInterceptor = interceptor
}

var (
	cmdGroupConfig = &cobra.Group{
		ID:    "configuration",
		Title: "Configuration",
	}

	cmdGroupResources = &cobra.Group{
		ID:    "resources",
		Title: "Resources",
	}

	cmdGroupTests = &cobra.Group{
		ID:    "tests",
		Title: "Tests",
	}

	cmdGroupMisc = &cobra.Group{
		ID:    "misc",
		Title: "Misc",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", fmt.Sprintf("output format [%s]", outputFormatsString))
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "config.yml", "config file will be used by the CLI")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "display debug information")

	rootCmd.PersistentFlags().StringVarP(&overrideEndpoint, "server-url", "s", "", "server url")

	rootCmd.AddGroup(
		cmdGroupConfig,
		cmdGroupResources,
		cmdGroupTests,
		cmdGroupMisc,
	)

	rootCmd.SetCompletionCommandGroupID(cmdGroupConfig.ID)
	rootCmd.SetHelpCommandGroupID(cmdGroupMisc.ID)
}
