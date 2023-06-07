package datastore

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

func addGetDatastorePreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// When I try to set up a new datastore
	// Then it should be applied with success
	dataStorePath := env.GetEnvironmentResourcePath(t, "data-store")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply datastore --file %s", dataStorePath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
}

func TestGetDatastore(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("get with no datastore initialized", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And no datastores registered

		// When I try to get a datastore on yaml mode
		// Then it should print a YAML list with one item
		result := tracetestcli.Exec(t, "list datastore --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		require.Equal(0, result.ExitCode)

		dataStoresYAML := helpers.UnmarshalYAMLSequence[types.DataStoreResource](t, result.StdOut)

		require.Len(dataStoresYAML, 1)
		require.Equal("DataStore", dataStoresYAML[0].Type)
	})

	addGetDatastorePreReqs(t, env)

	t.Run("get with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a Datastore already set

		// When I try to get a datastore on yaml mode
		// Then it should print a YAML
		result := tracetestcli.Exec(t, "get datastore --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		require.Equal(0, result.ExitCode)

		dataStore := helpers.UnmarshalYAML[types.DataStoreResource](t, result.StdOut)

		require.Equal("DataStore", dataStore.Type)
		require.Equal(env.Name(), dataStore.Spec.Name)
		require.True(dataStore.Spec.Default)
	})

	t.Run("get with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a Datastore already set

		// When I try to get a datastore on json mode
		// Then it should print a json
		result := tracetestcli.Exec(t, "get datastore --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		dataStores := helpers.UnmarshalJSON[types.DataStoreResource](t, result.StdOut)

		require.Equal("DataStore", dataStores.Type)
		require.Equal(env.Name(), dataStores.Spec.Name)
		require.True(dataStores.Spec.Default)
	})

	t.Run("get with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a Datastore already set

		// When I try to list datastore again on pretty mode
		// Then it should print a table with 4 lines printed: header, separator, data store item and empty line
		result := tracetestcli.Exec(t, "get datastore --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "current")
		require.Contains(result.StdOut, env.Name())

		lines := strings.Split(result.StdOut, "\n")
		require.Len(lines, 4)
	})
}
