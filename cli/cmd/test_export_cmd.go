package cmd

import (
	"context"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	exportTestId         string
	exportTestOutputFile string
	version              int32
)

var testExportCmd = &cobra.Command{
	Use:    "export",
	Short:  "Exports a test into a file",
	Long:   "Exports a test into a file",
	PreRun: setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Track("Test Export", "cmd", map[string]string{})

		ctx := context.Background()
		cliLogger.Debug("Exporting test", zap.String("testID", exportTestId))
		client := getAPIClient()
		exportTestAction := actions.NewExportTestAction(cliConfig, cliLogger, client)

		actionArgs := actions.ExportTestConfig{
			TestId:     exportTestId,
			OutputFile: exportTestOutputFile,
			Version:    version,
		}

		err := exportTestAction.Run(ctx, actionArgs)
		if err != nil {
			cliLogger.Error("could not get tests", zap.Error(err))
			return
		}
	},
	PostRun: teardownCommand,
}

func init() {
	testExportCmd.PersistentFlags().StringVarP(&exportTestId, "id", "", "", "id of the test")
	testExportCmd.PersistentFlags().StringVarP(&exportTestOutputFile, "output", "o", "", "file to be created with definition")
	testExportCmd.PersistentFlags().Int32VarP(&version, "version", "", -1, "version of the test. Default is latest")

	testCmd.AddCommand(testExportCmd)
}
