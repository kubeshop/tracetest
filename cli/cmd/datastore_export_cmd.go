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

var (
	exportOutputFile string
	dataStoreID      string
)

var dataStoreExportCmd = &cobra.Command{
	Use:    "export",
	Short:  "Exports a data store configuration into a file",
	Long:   "Exports a data store configuration into a file",
	PreRun: setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Track("Datastore Export", "cmd", map[string]string{})

		ctx := context.Background()
		client := utils.GetAPIClient(cliConfig)

		exportDataStoreAction := actions.NewExportDataStoreAction(cliLogger, client)
		actionArgs := actions.ExportDataStoreConfig{
			ID:         dataStoreID,
			OutputFile: exportOutputFile,
		}

		err := exportDataStoreAction.Run(ctx, actionArgs)
		if err != nil {
			cliLogger.Error("failed to export data store", zap.Error(err))
			os.Exit(1)
			return
		}

	},
	PostRun: teardownCommand,
}

func init() {
	dataStoreExportCmd.PersistentFlags().StringVarP(&exportOutputFile, "output", "o", "", "file where data store configuration will be saved")
	dataStoreExportCmd.PersistentFlags().StringVarP(&dataStoreID, "id", "", "", "id of the data store that will be exported")

	dataStoreCmd.AddCommand(dataStoreExportCmd)
}
