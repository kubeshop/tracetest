package environment

import (
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
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
	require.Equal(0, result.ExitCode)
	require.Contains(result.StdOut, "Resource environment with ID .noenv not found")

	// When I try to set up a new environment
	// Then it should be applied with success
	// newEnvironmentPath := env.GetTestResourcePath(t, "new-environment")

	// result = tracetestcli.Exec(t, fmt.Sprintf("apply environment --file %s", newEnvironmentPath), tracetestcli.WithCLIConfig(cliConfig))
	// require.Equal(0, result.ExitCode)

	// // When I try to get the environment applied on the last step
	// // Then it should return it
	// result = tracetestcli.Exec(t, "get environment --id .env", tracetestcli.WithCLIConfig(cliConfig))
	// require.Equal(0, result.ExitCode)

	// environmentVars := helpers.UnmarshalYAML[types.EnvironmentResource](t, result.StdOut)
	// require.Equal("Environment", environmentVars.Type)
	// require.Equal(".env", environmentVars.Spec.ID)
	// require.Equal(".env", environmentVars.Spec.Name)
	// require.Equal("some-value", environmentVars.Spec.Values["FIRST_VALUE"])
	// require.Equal("another_value", environmentVars.Spec.Values["SECOND_VALUE"])
}
