package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
)

var environmentApplyFile string

var environmentApplyCmd = &cobra.Command{
	Use:    "apply",
	Short:  "Create or update an environment to Tracetest",
	Long:   "Create or update an environment to Tracetest",
	PreRun: setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Track("Environment apply", "cmd", map[string]string{})
		client := utils.GetAPIClient(cliConfig)
		action := actions.NewApplyEnvironmentAction(cliLogger, client)

		err := action.Run(context.Background(), actions.ApplyEnvironmentConfig{
			File: environmentApplyFile,
		})
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
	PostRun: teardownCommand,
}

func init() {
	environmentApplyCmd.PersistentFlags().StringVarP(&environmentApplyFile, "file", "f", "", "file containing the environment configuration")
	environmentCmd.AddCommand(environmentApplyCmd)
}
