package cmd

import (
	"context"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	exportTestId         string
	exportTestOutputFile string
)

var testExportCmd = &cobra.Command{
	Use:    "export",
	Short:  "export a test",
	Long:   "export a test",
	PreRun: setupCommand,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		cliLogger.Debug("Exporting test", zap.String("testID", exportTestId))
		client := getAPIClient()
		exportTestAction := actions.NewExportTestAction(cliConfig, cliLogger, client)

		actionArgs := actions.ExportTestConfig{
			TestId:     exportTestId,
			OutputFile: exportTestOutputFile,
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
	testExportCmd.PersistentFlags().StringVarP(&exportTestId, "id", "", "", "tracetest test export --id <id>")
	testExportCmd.PersistentFlags().StringVarP(&exportTestOutputFile, "output", "o", "", "tracetest test export --output file")

	testCmd.AddCommand(testExportCmd)
}
