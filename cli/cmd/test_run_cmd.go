package cmd

import (
	"context"
	"os"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var runTestFileDefinition string
var runTestWaitForResult bool
var runTestJUnit string

var testRunCmd = &cobra.Command{
	Use:    "run",
	Short:  "run a test",
	Long:   "run a test using a definition or an id",
	PreRun: setupCommand,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		client := getAPIClient()

		runTestAction := actions.NewRunTestAction(cliConfig, cliLogger, client)
		actionArgs := actions.RunTestConfig{
			DefinitionFile: runTestFileDefinition,
			WaitForResult:  runTestWaitForResult,
			JUnit:          runTestJUnit,
		}

		err := runTestAction.Run(ctx, actionArgs)
		if err != nil {
			cliLogger.Error("failed to run test", zap.Error(err))
			os.Exit(1)
			return
		}
	},
	PostRun: teardownCommand,
}

func init() {
	testRunCmd.PersistentFlags().StringVarP(&runTestFileDefinition, "definition", "d", "", "--definition <definition-file.yml>")
	testRunCmd.PersistentFlags().BoolVarP(&runTestWaitForResult, "wait-for-result", "w", false, "")
	testRunCmd.PersistentFlags().StringVarP(&runTestJUnit, "junit", "j", "", "--junit <junit-output.xml>")
	testCmd.AddCommand(testRunCmd)
}
