package environment

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func addGetEnvironmentPreReqs(t *testing.T, env environment.Manager) {
	// instantiate require with testing helper
	require := require.New(t)

	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to set up a new environment
	// Then it should be applied with success
	newEnvironmentPath := env.GetTestResourcePath(t, "new-environment")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply environment --file %s", newEnvironmentPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	environmentVars := helpers.UnmarshalYAML[types.EnvironmentResource](t, result.StdOut)
	require.Equal("Environment", environmentVars.Type)
	require.Equal(".env", environmentVars.Spec.ID)
	require.Equal(".env", environmentVars.Spec.Name)
}

func TestGetEnvironment(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("get with no environment initialized", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And no environment registered

		// When I try to get a environment on yaml mode
		// Then it should return a error message
		result := tracetestcli.Exec(t, "get environment --id .env --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "Resource environment with ID .env not found")
	})

	addGetEnvironmentPreReqs(t, env)

	t.Run("get with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a environment already set

		// When I try to get an environment on yaml mode
		// Then it should print a YAML
		result := tracetestcli.Exec(t, "get environment --id .env --output yaml", tracetestcli.WithCLIConfig(cliConfig))
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
	})

	t.Run("get with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a environment already set

		// When I try to get an environment on json mode
		// Then it should print a json
		result := tracetestcli.Exec(t, "get environment --id .env --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		environmentVars := helpers.UnmarshalJSON[types.EnvironmentResource](t, result.StdOut)

		require.Equal("Environment", environmentVars.Type)
		require.Equal(".env", environmentVars.Spec.ID)
		require.Equal(".env", environmentVars.Spec.Name)
		require.Len(environmentVars.Spec.Values, 2)
		require.Equal("FIRST_VAR", environmentVars.Spec.Values[0].Key)
		require.Equal("some-value", environmentVars.Spec.Values[0].Value)
		require.Equal("SECOND_VAR", environmentVars.Spec.Values[1].Key)
		require.Equal("another_value", environmentVars.Spec.Values[1].Value)
	})

	t.Run("get with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a environment already set

		// When I try to get an environment on pretty mode
		// Then it should print a table with 4 lines printed: header, separator, environment item and empty line
		result := tracetestcli.Exec(t, "get environment --id .env --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, ".env")

		lines := strings.Split(result.StdOut, "\n")
		require.Len(lines, 4)
	})
}
