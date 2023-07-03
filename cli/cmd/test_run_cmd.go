package cmd

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
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
	Run: WithResultHandler(func(_ *cobra.Command, _ []string) (string, error) {
		ctx := context.Background()
		client := utils.GetAPIClient(cliConfig)

		envClient, err := resources.Get("environment")
		if err != nil {
			return "", fmt.Errorf("failed to get environment client: %w", err)
		}

		testClient, err := resources.Get("test")
		if err != nil {
			return "", fmt.Errorf("failed to get test client: %w", err)
		}

		runTestAction := actions.NewRunTestAction(cliConfig, cliLogger, client, testClient, envClient, ExitCLI)
		actionArgs := actions.RunResourceArgs{
			DefinitionFile: runTestFileDefinition,
			EnvID:          runTestEnvID,
			WaitForResult:  runTestWaitForResult,
			JUnit:          runTestJUnit,
		}

		err = runTestAction.Run(ctx, actionArgs)
		return "", err
	}),
	PostRun: teardownCommand,
}

func init() {
	testRunCmd.PersistentFlags().StringVarP(&runTestEnvID, "environment", "e", "", "id of the environment to be used")
	testRunCmd.PersistentFlags().StringVarP(&runTestFileDefinition, "definition", "d", "", "path to the definition file to be run")
	testRunCmd.PersistentFlags().BoolVarP(&runTestWaitForResult, "wait-for-result", "w", false, "wait for the test result to print it's result")
	testRunCmd.PersistentFlags().StringVarP(&runTestJUnit, "junit", "j", "", "path to the junit file that will be generated")
	testCmd.AddCommand(testRunCmd)
}
