package misc

import (
	"context"

	"github.com/kubeshop/tracetest/cli/analytics"
	misc_actions "github.com/kubeshop/tracetest/cli/misc/actions"
	misc_parameters "github.com/kubeshop/tracetest/cli/misc/parameters"
	resources_actions "github.com/kubeshop/tracetest/cli/resources/actions"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
)

type TestRun struct {
	args[*misc_parameters.TestRun]
}

func (d TestRun) Run(cmd *cobra.Command, args []string) (string, error) {
	analytics.Track("Test Run", "cmd", map[string]string{})
	ctx := context.Background()

	baseOptions := []resources_actions.ResourceArgsOption{resources_actions.WithLogger(d.Setup.Logger), resources_actions.WithConfig(*d.Setup.Config)}
	environmentOptions := append(baseOptions, resources_actions.WithClient(utils.GetResourceAPIClient("environments", *d.Setup.Config)))
	environmentActions := resources_actions.NewEnvironmentsActions(environmentOptions...)

	runTestAction := misc_actions.NewRunTestAction(*d.Setup.Config, d.Setup.Logger, d.Setup.Client, environmentActions)
	actionArgs := misc_actions.RunTestConfig{
		DefinitionFile: d.Parameters.RunTestFileDefinition,
		EnvID:          d.Parameters.RunTestEnvID,
		WaitForResult:  d.Parameters.RunTestWaitForResult,
		JUnit:          d.Parameters.RunTestJUnit,
	}

	err := runTestAction.Run(ctx, actionArgs)
	return "", err
}

func NewTestRun(root Test) TestRun {
	parameters := misc_parameters.NewTestRun()
	defaults := NewDefaults("test export", root.Setup.Setup)

	testRun := TestRun{
		args: NewArgs(defaults, parameters),
	}

	testRun.Cmd = &cobra.Command{
		Use:     "run",
		Short:   "Run a test on your Tracetest server",
		Long:    "Run a test on your Tracetest server",
		PreRun:  defaults.PreRun,
		Run:     defaults.Run(testRun.Run),
		PostRun: defaults.PostRun,
	}

	testRun.Cmd.PersistentFlags().StringVarP(&parameters.RunTestEnvID, "environment", "e", "", "id of the environment to be used")
	testRun.Cmd.PersistentFlags().StringVarP(&parameters.RunTestFileDefinition, "definition", "d", "", "path to the definition file to be run")
	testRun.Cmd.PersistentFlags().BoolVarP(&parameters.RunTestWaitForResult, "wait-for-result", "w", false, "wait for the test result to print it's result")
	testRun.Cmd.PersistentFlags().StringVarP(&parameters.RunTestJUnit, "junit", "j", "", "path to the junit file that will be generated")

	root.Cmd.AddCommand(testRun.Cmd)

	return testRun
}
