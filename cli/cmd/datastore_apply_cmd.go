package cmd

import (
	"context"
	"os"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var dataStoreApplyFile string

var dataStoreApplyCmd = &cobra.Command{
	Use:    "apply",
	Short:  "Apply (create/update) data store configuration to your Tracetest server",
	Long:   "Apply (create/update) data store configuration to your Tracetest server",
	PreRun: setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Track("Datastore Apply", "cmd", map[string]string{})

		ctx := context.Background()
		client := utils.GetAPIClient(cliConfig)

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
	dataStoreApplyCmd.PersistentFlags().StringVarP(&dataStoreApplyFile, "file", "f", "", "file containing the data store configuration")
	dataStoreCmd.AddCommand(dataStoreApplyCmd)
}
