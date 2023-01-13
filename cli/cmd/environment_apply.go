package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/spf13/cobra"
)

var environmentApplyFile string

var environmentApplyCmd = &cobra.Command{
	Use:    "apply",
	Short:  "Apply environment to tracetest",
	Long:   "Apply environment to tracetest",
	PreRun: setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Track("Environment apply", "cmd", map[string]string{})
		action := actions.NewApplyEnvironmentAction(cliLogger, getAPIClient())

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
	environmentApplyCmd.PersistentFlags().StringVarP(&environmentApplyFile, "file", "f", "", "--file environment.yaml")
	environmentCmd.AddCommand(environmentApplyCmd)
}
