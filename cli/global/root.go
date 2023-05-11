package global

import (
	"fmt"

	global_formatters "github.com/kubeshop/tracetest/cli/global/formatters"
	global_middlewares "github.com/kubeshop/tracetest/cli/global/middlewares"
	global_parameters "github.com/kubeshop/tracetest/cli/global/parameters"
	global_setup "github.com/kubeshop/tracetest/cli/global/setup"
	"github.com/spf13/cobra"
)

var (
	GroupConfig = &cobra.Group{
		ID:    "configuration",
		Title: "Configuration",
	}

	GroupResources = &cobra.Group{
		ID:    "resources",
		Title: "Resources",
	}

	GroupTests = &cobra.Group{
		ID:    "tests",
		Title: "Tests",
	}

	GroupMisc = &cobra.Group{
		ID:    "misc",
		Title: "Misc",
	}
)

type Root struct {
	Parameters *global_parameters.Global
	Setup      *global_setup.Setup
	Cmd        *cobra.Command
}

func NewRoot() Root {
	parameters := global_parameters.NewGlobal()
	setup := global_setup.NewSetup(parameters, global_setup.WithConfig(), global_setup.WithLogger(), global_setup.WithVersionText(), global_setup.WithInitAnalytics())
	run := global_middlewares.ComposeRun(global_middlewares.WithAnalytics("root", "cmd"), global_middlewares.WithResultHandler(setup))
	postRun := global_middlewares.ComposeNoopRun(global_middlewares.WithTeardownMiddleware(setup))

	root := Root{
		Setup:      setup,
		Parameters: parameters,
	}

	root.Cmd = &cobra.Command{
		Use:     "tracetest",
		Short:   "CLI to configure, install and execute tests on a Tracetest server",
		Long:    `CLI to configure, install and execute tests on a Tracetest server`,
		PreRun:  run(setup.PreRun),
		PostRun: postRun(),
	}

	root.Cmd.PersistentFlags().StringVarP(&root.Parameters.Output, "output", "o", "", fmt.Sprintf("output format [%s]", global_formatters.OutputFormatsString))
	root.Cmd.PersistentFlags().StringVarP(&root.Parameters.ConfigFile, "config", "c", "config.yml", "config file will be used by the CLI")
	root.Cmd.PersistentFlags().BoolVarP(&root.Parameters.Verbose, "verbose", "v", false, "display debug information")
	root.Cmd.PersistentFlags().StringVarP(&root.Parameters.OverrideEndpoint, "server-url", "s", "", "server url")

	root.Cmd.AddGroup(
		GroupConfig,
		GroupResources,
		GroupTests,
		GroupMisc,
	)

	root.Cmd.SetCompletionCommandGroupID(GroupConfig.ID)
	root.Cmd.SetHelpCommandGroupID(GroupMisc.ID)

	return root
}
