package deprecated

import (
	"github.com/kubeshop/tracetest/cli/global"
	global_middlewares "github.com/kubeshop/tracetest/cli/global/middlewares"
	"github.com/kubeshop/tracetest/cli/resources"
	"github.com/spf13/cobra"
)

var environmentApplyFile string

type EnvironmentLegacy []*cobra.Command

func NewEnvironmentLegacy(root global.Root) EnvironmentLegacy {
	postRun := global_middlewares.ComposeNoopRun(global_middlewares.WithTeardownMiddleware(root.Setup))
	run := global_middlewares.ComposeRun(global_middlewares.WithAnalytics("environments", "cmd"), global_middlewares.WithResultHandler(root.Setup))

	environmentCmd := &cobra.Command{
		GroupID:    global.GroupConfig.ID,
		Use:        "environment",
		Short:      "Manage your tracetest environments",
		Long:       "Manage your tracetest environments",
		Deprecated: "Please use `tracetest (apply|delete|list|get|export) environment` commands instead.",
		PreRun:     run(root.Setup.PreRun),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		PostRun: postRun(),
	}

	var environmentApplyCmd = &cobra.Command{
		Use:        "apply",
		Short:      "Create or update an environment to Tracetest",
		Long:       "Create or update an environment to Tracetest",
		Deprecated: "Please use `tracetest apply environment --file [path]` command instead.",
		PreRun:     run(root.Setup.PreRun),
		Run: func(cmd *cobra.Command, args []string) {
			// call new apply command
			apply := resources.NewApply(root)

			apply.Parameters.DefinitionFile = dataStoreApplyFile
			apply.Cmd.Run(apply.Cmd, []string{"environment"})
		},
		PostRun: postRun(),
	}

	root.Cmd.AddCommand(environmentCmd)

	environmentApplyCmd.PersistentFlags().StringVarP(&environmentApplyFile, "file", "f", "", "file containing the environment configuration")
	environmentCmd.AddCommand(environmentApplyCmd)

	return EnvironmentLegacy{environmentCmd, environmentApplyCmd}
}
