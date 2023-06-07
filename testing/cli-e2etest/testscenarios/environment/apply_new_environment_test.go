package environment

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestApplyNewEnvironment(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup isolated e2e environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to get an environment that doesn't exists
	// Then it should return error message
	result := tracetestcli.Exec(t, "get environment --id .noenv", tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
	require.Contains(result.StdOut, "Resource environment with ID .noenv not found")

	// When I try to set up a new environment
	// Then it should be applied with success
	newEnvironmentPath := env.GetTestResourcePath(t, "new-environment")

	result = tracetestcli.Exec(t, fmt.Sprintf("apply environment --file %s", newEnvironmentPath), tracetestcli.WithCLIConfig(cliConfig))
	require.Equal(0, result.ExitCode)

	// When I try to get the environment applied on the last step
	// Then it should return it
	result = tracetestcli.Exec(t, "get environment --id .env", tracetestcli.WithCLIConfig(cliConfig))
	require.Equal(0, result.ExitCode)

	environmentVars := helpers.UnmarshalYAML[types.EnvironmentResource](t, result.StdOut)

	require.Equal("Environment", environmentVars.Type)
	require.Equal(".env", environmentVars.Spec.ID)
	require.Equal(".env", environmentVars.Spec.Name)
	require.Len(environmentVars.Spec.Values, 2)
	require.Equal("FIRST_VAR", environmentVars.Spec.Values[0].Key)
	require.Equal("some-value", environmentVars.Spec.Values[0].Value)
	require.Equal("SECOND_VAR", environmentVars.Spec.Values[1].Key)
	require.Equal("another_value", environmentVars.Spec.Values[1].Value)

	// When I try to update the last environment
	// Then it should be applied with success
	updatedNewEnvironmentPath := env.GetTestResourcePath(t, "updated-new-environment")

	result = tracetestcli.Exec(t, fmt.Sprintf("apply environment --file %s", updatedNewEnvironmentPath), tracetestcli.WithCLIConfig(cliConfig))
	require.Equal(0, result.ExitCode)

	// When I try to get the environment applied on the last step
	// Then it should return it
	result = tracetestcli.Exec(t, "get environment --id .env", tracetestcli.WithCLIConfig(cliConfig))
	require.Equal(0, result.ExitCode)

	updatedEnvironmentVars := helpers.UnmarshalYAML[types.EnvironmentResource](t, result.StdOut)

	require.Equal("Environment", updatedEnvironmentVars.Type)
	require.Equal(".env", updatedEnvironmentVars.Spec.ID)
	require.Equal(".env", updatedEnvironmentVars.Spec.Name)
	require.Len(updatedEnvironmentVars.Spec.Values, 3)
	require.Equal("FIRST_VAR", updatedEnvironmentVars.Spec.Values[0].Key)
	require.Equal("some-value", updatedEnvironmentVars.Spec.Values[0].Value)
	require.Equal("SECOND_VAR", updatedEnvironmentVars.Spec.Values[1].Key)
	require.Equal("updated_value", updatedEnvironmentVars.Spec.Values[1].Value) // this value has been updated
	require.Equal("THIRD_VAR", updatedEnvironmentVars.Spec.Values[2].Key)
	require.Equal("hello", updatedEnvironmentVars.Spec.Values[2].Value) // this value was added

}
