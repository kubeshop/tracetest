package cmd

import (
	"context"
	"os"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	runTestFileDefinition,
	runTestEnvID,
	runTestJUnit string
	runTestWaitForResult bool
)

var testRunCmd = &cobra.Command{
	Use:    "run",
	Short:  "Run a test on your Tracetest server",
	Long:   "Run a test on your Tracetest server",
	PreRun: setupCommand(),
	Run: func(cmd *cobra.Command, args []string) {
		analytics.Track("Test Run", "cmd", map[string]string{})

		ctx := context.Background()
		client := getAPIClient()

		runTestAction := actions.NewRunTestAction(cliConfig, cliLogger, client)
		actionArgs := actions.RunTestConfig{
			DefinitionFile: runTestFileDefinition,
			EnvID:          runTestEnvID,
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
	testRunCmd.PersistentFlags().StringVarP(&runTestEnvID, "environment", "e", "", "--environment <id>")
	testRunCmd.PersistentFlags().StringVarP(&runTestFileDefinition, "definition", "d", "", "--definition <definition-file.yml>")
	testRunCmd.PersistentFlags().BoolVarP(&runTestWaitForResult, "wait-for-result", "w", false, "")
	testRunCmd.PersistentFlags().StringVarP(&runTestJUnit, "junit", "j", "", "--junit <junit-output.xml>")
	testCmd.AddCommand(testRunCmd)
}
