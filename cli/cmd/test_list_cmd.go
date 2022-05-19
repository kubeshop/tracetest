package cmd

import (
	"context"
	"log"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var testListCmd = &cobra.Command{
	Use:    "list",
	Short:  "list all test",
	Long:   "list all test",
	PreRun: setupCommand,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		cliLogger.Debug("Retrieving list of tests", zap.String("endpoint", cliConfig.Endpoint))
		client := getAPIClient()
		listTestsAction := actions.NewListTestsAction(cliConfig, cliLogger, client)

		err := listTestsAction.Run(ctx, args)
		if err != nil {
			log.Fatal(err)
		}
	},
	PostRun: teardownCommand,
}

func init() {
	testCmd.AddCommand(testListCmd)
}
