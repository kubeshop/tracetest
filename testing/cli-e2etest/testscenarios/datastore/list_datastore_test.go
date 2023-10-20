package datastore

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/testing/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func addListDatastorePreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// When I try to set up a new datastore
	// Then it should be applied with success
	dataStorePath := env.GetEnvironmentResourcePath(t, "data-store")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply datastore --file %s", dataStorePath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
}

func TestListDatastore(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("list with no datastore initialized", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list datastore on pretty mode and there is no datastore
		// Then it should print an empty table
		result := tracetestcli.Exec(t, "list datastore --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.NotContains(result.StdOut, "current")
	})

	addListDatastorePreReqs(t, env)

	t.Run("list with invalid sortBy field", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a datastore created

		// When I try to list a datastore by an invalid field
		// Then I should receive an error
		result := tracetestcli.Exec(t, "list datastore --sortBy id --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 1)
		require.Contains(result.StdErr, "invalid sort field: id") // TODO: think on how to improve this error handling
	})

	t.Run("list with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a datastore created

		// When I try to list datastore again on yaml mode
		// Then it should print a YAML list with one item
		result := tracetestcli.Exec(t, "list datastore --sortBy name --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		dataStoresYAML := helpers.UnmarshalYAMLSequence[types.DataStoreResource](t, result.StdOut)

		require.Len(dataStoresYAML, 1)
		require.Equal("DataStore", dataStoresYAML[0].Type)
		require.Equal(env.Name(), dataStoresYAML[0].Spec.Name)
		require.True(dataStoresYAML[0].Spec.Default)
	})

	t.Run("list with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a datastore created

		// When I try to list datastore again on json mode
		// Then it should print a JSON list with one item
		result := tracetestcli.Exec(t, "list datastore --sortBy name --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		dataStoresList := helpers.UnmarshalJSON[types.ResourceList[types.DataStoreResource]](t, result.StdOut)
		require.Len(dataStoresList.Items, 1)
		require.Equal(len(dataStoresList.Items), dataStoresList.Count)

		require.Equal("DataStore", dataStoresList.Items[0].Type)
		require.Equal(env.Name(), dataStoresList.Items[0].Spec.Name)
		require.True(dataStoresList.Items[0].Spec.Default)
	})

	t.Run("list with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a datastore created

		// When I try to list datastore again on pretty mode
		// Then it should print a table with 4 lines printed: header, separator, data store item and empty line
		result := tracetestcli.Exec(t, "list datastore --sortBy name --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		parsedTable := helpers.UnmarshalTable(t, result.StdOut)
		require.Len(parsedTable, 1)

		singleLine := parsedTable[0]

		require.Equal("current", singleLine["ID"])
		require.Equal(env.Name(), singleLine["NAME"])
		require.Equal("*", singleLine["DEFAULT"])
	})
}
