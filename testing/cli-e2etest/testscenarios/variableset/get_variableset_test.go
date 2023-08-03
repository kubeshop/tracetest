package variableset

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func addGetVariableSetPreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to set up a new variable set
	// Then it should be applied with success
	newEnvironmentPath := env.GetTestResourcePath(t, "new-varSet")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply variableset --file %s", newEnvironmentPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
}

func TestGetVariableSet(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("get with no variable set initialized", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And no environment registered

		// When I try to get a environment on yaml mode
		// Then it should return a error message
		result := tracetestcli.Exec(t, "get variableset --id .env --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "Resource variableset with ID .env not found")
	})

	addGetVariableSetPreReqs(t, env)

	t.Run("get with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have an variable set already set

		// When I try to get an variable set on yaml mode
		// Then it should print a YAML
		result := tracetestcli.Exec(t, "get variableset --id .env --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		environmentVars := helpers.UnmarshalYAML[types.VariableSetResource](t, result.StdOut)

		require.Equal("VariableSet", environmentVars.Type)
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
		// And I have an variable set already set

		// When I try to get an variable set on json mode
		// Then it should print a json
		result := tracetestcli.Exec(t, "get variableset --id .env --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		environmentVars := helpers.UnmarshalJSON[types.VariableSetResource](t, result.StdOut)

		require.Equal("VariableSet", environmentVars.Type)
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
		// And I have an variable set already set

		// When I try to get an variable set on pretty mode
		// Then it should print a table with 4 lines printed: header, separator, variable set item and empty line
		result := tracetestcli.Exec(t, "get variableset --id .env --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		parsedTable := helpers.UnmarshalTable(t, result.StdOut)
		require.Len(parsedTable, 1)

		singleLine := parsedTable[0]

		require.Equal(".env", singleLine["ID"])
		require.Equal(".env", singleLine["NAME"])
		require.Equal("", singleLine["DESCRIPTION"])
	})

	t.Run("getting a variable set using the deprecated environment command", func(t *testing.T) {
		result := tracetestcli.Exec(t, "get environment --id .env", tracetestcli.WithCLIConfig(cliConfig))

		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "The resource `environment` is deprecated and will be removed in a future version. Please use `variableset` instead.")
		require.Contains(result.StdOut, "VariableSet")
		require.Contains(result.StdOut, ".env")
	})
}
