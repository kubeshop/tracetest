package cmd

import (
	"context"

	"github.com/kubeshop/tracetest/cli/actions"
	"github.com/kubeshop/tracetest/cli/analytics"
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
		analytics.Track("Test Run", "cmd", map[string]string{})

		ctx := context.Background()
		client := utils.GetAPIClient(cliConfig)

		baseOptions := []actions.ResourceArgsOption{actions.WithLogger(cliLogger), actions.WithConfig(cliConfig)}
		environmentOptions := append(baseOptions, actions.WithClient(utils.GetResourceAPIClient("environments", cliConfig)))
		environmentActions := actions.NewEnvironmentsActions(environmentOptions...)

		runTestAction := actions.NewRunTestAction(cliConfig, cliLogger, client, environmentActions)
		actionArgs := actions.RunResourceArgs{
			DefinitionFile: runTestFileDefinition,
			EnvID:          runTestEnvID,
			WaitForResult:  runTestWaitForResult,
			JUnit:          runTestJUnit,
		}

		err := runTestAction.Run(ctx, actionArgs)
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
