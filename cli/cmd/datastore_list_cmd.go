package cmd

import (
	"context"
	"os"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var dataStoreListCmd = &cobra.Command{
	Use:    "list",
	Short:  "List data store configurations to your tracetest server",
	Long:   "List data store configurations to your tracetest server",
	PreRun: setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Track("Datastore List", "cmd", map[string]string{})

		ctx := context.Background()
		client := getAPIClient()

		applyDataStoreAction := actions.NewListDataStoreAction(cliConfig, cliLogger, client)
		actionArgs := actions.ListDataStoreConfig{}

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
	dataStoreCmd.AddCommand(dataStoreListCmd)
}
