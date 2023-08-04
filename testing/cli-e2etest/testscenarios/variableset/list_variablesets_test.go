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

func addListVariableSetPreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to set up a new environment
	// Then it should be applied with success
	newEnvironmentPath := env.GetTestResourcePath(t, "new-varSet")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply variableset --file %s", newEnvironmentPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	// When I try to set up a another environment
	// Then it should be applied with success
	anotherEnvironmentPath := env.GetTestResourcePath(t, "another-varSet")

	result = tracetestcli.Exec(t, fmt.Sprintf("apply variableset --file %s", anotherEnvironmentPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	// When I try to set up a third environment
	// Then it should be applied with success
	oneMoreEnvironmentPath := env.GetTestResourcePath(t, "one-more-varSet")

	result = tracetestcli.Exec(t, fmt.Sprintf("apply variableset --file %s", oneMoreEnvironmentPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
}

func TestListVariableSets(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup isolated e2e environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("list no variable sets", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And there is no variable sets
		result := tracetestcli.Exec(t, "list variableset --sortBy name --sortDirection asc --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		environmentVarsList := helpers.UnmarshalYAMLSequence[types.VariableSetResource](t, result.StdOut)
		require.Len(environmentVarsList, 0)
	})

	addListVariableSetPreReqs(t, env)

	t.Run("list with invalid sortBy field", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these variable set by an invalid field
		// Then I should receive an error
		result := tracetestcli.Exec(t, "list variableset --sortBy id --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 1)
		require.Contains(result.StdErr, "invalid sort field: id") // TODO: think on how to improve this error handling
	})

	t.Run("list with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these variable sets by a valid field and in YAML format
		// Then I should receive three variable sets
		result := tracetestcli.Exec(t, "list variableset --sortBy name --sortDirection asc --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		environmentVarsList := helpers.UnmarshalYAMLSequence[types.VariableSetResource](t, result.StdOut)
		require.Len(environmentVarsList, 3)

		// due our database sorting algorithm, "another-env" comes in the front of ".env"
		// ref https://wiki.postgresql.org/wiki/FAQ#Why_do_my_strings_sort_incorrectly.3F
		anotherEnvironmentVars := environmentVarsList[0]
		require.Equal("VariableSet", anotherEnvironmentVars.Type)
		require.Equal("another-env", anotherEnvironmentVars.Spec.ID)
		require.Equal("another-env", anotherEnvironmentVars.Spec.Name)
		require.Len(anotherEnvironmentVars.Spec.Values, 2)
		require.Equal("Here", anotherEnvironmentVars.Spec.Values[0].Key)
		require.Equal("We", anotherEnvironmentVars.Spec.Values[0].Value)
		require.Equal("Come", anotherEnvironmentVars.Spec.Values[1].Key)
		require.Equal("Again", anotherEnvironmentVars.Spec.Values[1].Value)

		environmentVars := environmentVarsList[1]
		require.Equal("VariableSet", environmentVars.Type)
		require.Equal(".env", environmentVars.Spec.ID)
		require.Equal(".env", environmentVars.Spec.Name)
		require.Len(environmentVars.Spec.Values, 2)
		require.Equal("FIRST_VAR", environmentVars.Spec.Values[0].Key)
		require.Equal("some-value", environmentVars.Spec.Values[0].Value)
		require.Equal("SECOND_VAR", environmentVars.Spec.Values[1].Key)
		require.Equal("another_value", environmentVars.Spec.Values[1].Value)

		oneMoreEnvironmentVars := environmentVarsList[2]
		require.Equal("VariableSet", oneMoreEnvironmentVars.Type)
		require.Equal("one-more-env", oneMoreEnvironmentVars.Spec.ID)
		require.Equal("one-more-env", oneMoreEnvironmentVars.Spec.Name)
		require.Len(oneMoreEnvironmentVars.Spec.Values, 2)
		require.Equal("This", oneMoreEnvironmentVars.Spec.Values[0].Key)
		require.Equal("Is", oneMoreEnvironmentVars.Spec.Values[0].Value)
		require.Equal("The", oneMoreEnvironmentVars.Spec.Values[1].Key)
		require.Equal("Third", oneMoreEnvironmentVars.Spec.Values[1].Value)
	})

	t.Run("list with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these environments by a valid field and in JSON format
		// Then I should receive three environments
		result := tracetestcli.Exec(t, "list variableset --sortBy name --sortDirection asc --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		environmentVarsList := helpers.UnmarshalJSON[types.ResourceList[types.VariableSetResource]](t, result.StdOut)
		require.Equal(3, environmentVarsList.Count)
		require.Len(environmentVarsList.Items, 3)

		// due our database sorting algorithm, "another-env" comes in the front of ".env"
		// ref https://wiki.postgresql.org/wiki/FAQ#Why_do_my_strings_sort_incorrectly.3F
		anotherEnvironmentVars := environmentVarsList.Items[0]
		require.Equal("VariableSet", anotherEnvironmentVars.Type)
		require.Equal("another-env", anotherEnvironmentVars.Spec.ID)
		require.Equal("another-env", anotherEnvironmentVars.Spec.Name)
		require.Len(anotherEnvironmentVars.Spec.Values, 2)
		require.Equal("Here", anotherEnvironmentVars.Spec.Values[0].Key)
		require.Equal("We", anotherEnvironmentVars.Spec.Values[0].Value)
		require.Equal("Come", anotherEnvironmentVars.Spec.Values[1].Key)
		require.Equal("Again", anotherEnvironmentVars.Spec.Values[1].Value)

		environmentVars := environmentVarsList.Items[1]
		require.Equal("VariableSet", environmentVars.Type)
		require.Equal(".env", environmentVars.Spec.ID)
		require.Equal(".env", environmentVars.Spec.Name)
		require.Len(environmentVars.Spec.Values, 2)
		require.Equal("FIRST_VAR", environmentVars.Spec.Values[0].Key)
		require.Equal("some-value", environmentVars.Spec.Values[0].Value)
		require.Equal("SECOND_VAR", environmentVars.Spec.Values[1].Key)
		require.Equal("another_value", environmentVars.Spec.Values[1].Value)

		oneMoreEnvironmentVars := environmentVarsList.Items[2]
		require.Equal("VariableSet", oneMoreEnvironmentVars.Type)
		require.Equal("one-more-env", oneMoreEnvironmentVars.Spec.ID)
		require.Equal("one-more-env", oneMoreEnvironmentVars.Spec.Name)
		require.Len(oneMoreEnvironmentVars.Spec.Values, 2)
		require.Equal("This", oneMoreEnvironmentVars.Spec.Values[0].Key)
		require.Equal("Is", oneMoreEnvironmentVars.Spec.Values[0].Value)
		require.Equal("The", oneMoreEnvironmentVars.Spec.Values[1].Key)
		require.Equal("Third", oneMoreEnvironmentVars.Spec.Values[1].Value)
	})

	t.Run("list with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these environments by a valid field and in pretty format
		// Then it should print a table with 6 lines printed: header, separator, three envs and empty line
		result := tracetestcli.Exec(t, "list variableset --sortBy name --sortDirection asc --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		// due our database sorting algorithm, "another-env" comes in the front of ".env"
		// ref https://wiki.postgresql.org/wiki/FAQ#Why_do_my_strings_sort_incorrectly.3F
		parsedTable := helpers.UnmarshalTable(t, result.StdOut)
		require.Len(parsedTable, 3)

		firstLine := parsedTable[0]

		require.Equal("another-env", firstLine["ID"])
		require.Equal("another-env", firstLine["NAME"])
		require.Equal("", firstLine["DESCRIPTION"])

		secondLine := parsedTable[1]

		require.Equal(".env", secondLine["ID"])
		require.Equal(".env", secondLine["NAME"])
		require.Equal("", secondLine["DESCRIPTION"])

		thirdLine := parsedTable[2]

		require.Equal("one-more-env", thirdLine["ID"])
		require.Equal("one-more-env", thirdLine["NAME"])
		require.Equal("", thirdLine["DESCRIPTION"])
	})

	t.Run("list with YAML format skipping the first and taking two items", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these environments by a valid field, paging options and in YAML format
		// Then I should receive two environments
		result := tracetestcli.Exec(t, "list variableset --sortBy name --sortDirection asc --skip 1 --take 2 --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		environmentVarsList := helpers.UnmarshalYAMLSequence[types.VariableSetResource](t, result.StdOut)
		require.Len(environmentVarsList, 2)

		// due our database sorting algorithm, "another-env" comes in the front of ".env"
		// ref https://wiki.postgresql.org/wiki/FAQ#Why_do_my_strings_sort_incorrectly.3F
		environmentVars := environmentVarsList[0]
		require.Equal("VariableSet", environmentVars.Type)
		require.Equal(".env", environmentVars.Spec.ID)
		require.Equal(".env", environmentVars.Spec.Name)
		require.Len(environmentVars.Spec.Values, 2)
		require.Equal("FIRST_VAR", environmentVars.Spec.Values[0].Key)
		require.Equal("some-value", environmentVars.Spec.Values[0].Value)
		require.Equal("SECOND_VAR", environmentVars.Spec.Values[1].Key)
		require.Equal("another_value", environmentVars.Spec.Values[1].Value)

		oneMoreEnvironmentVars := environmentVarsList[1]
		require.Equal("VariableSet", oneMoreEnvironmentVars.Type)
		require.Equal("one-more-env", oneMoreEnvironmentVars.Spec.ID)
		require.Equal("one-more-env", oneMoreEnvironmentVars.Spec.Name)
		require.Len(oneMoreEnvironmentVars.Spec.Values, 2)
		require.Equal("This", oneMoreEnvironmentVars.Spec.Values[0].Key)
		require.Equal("Is", oneMoreEnvironmentVars.Spec.Values[0].Value)
		require.Equal("The", oneMoreEnvironmentVars.Spec.Values[1].Key)
		require.Equal("Third", oneMoreEnvironmentVars.Spec.Values[1].Value)
	})

	t.Run("list using deprecated environment command", func(t *testing.T) {
		result := tracetestcli.Exec(t, "list environment --sortBy name --sortDirection asc --skip 1 --take 2 --output yaml", tracetestcli.WithCLIConfig(cliConfig))

		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "The resource `environment` is deprecated and will be removed in a future version. Please use `variableset` instead.")
		require.Contains(result.StdOut, "VariableSet")
	})
}
