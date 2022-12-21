package cmd

import (
	"context"
	"os"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var dataStoreApplyFile string

var dataStoreApplyCmd = &cobra.Command{
	Use:    "apply",
	Short:  "Apply data store configuration to your tracetest server",
	Long:   "Apply data store configuration to your tracetest server",
	PreRun: setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Track("Datastore Apply", "cmd", map[string]string{})

		ctx := context.Background()
		client := getAPIClient()

		applyDataStoreAction := actions.NewApplyDataStoreAction(cliLogger, client)
		actionArgs := actions.ApplyDataStoreConfig{
			File: dataStoreApplyFile,
		}

		err := applyDataStoreAction.Run(ctx, actionArgs)
		if err != nil {
			cliLogger.Error("failed to run test", zap.Error(err))
			os.Exit(1)
			return
		}

	},
	PostRun: teardownCommand,
}

func init() {
	dataStoreApplyCmd.PersistentFlags().StringVarP(&dataStoreApplyFile, "file", "f", "", "--file my-datastore.yaml")
	dataStoreCmd.AddCommand(dataStoreApplyCmd)
}
