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

func TestListDatastore(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to set up a new datastore
	// Then it should be applied with success
	dataStorePath := env.GetManisfestResourcePath(t, "data-store")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply datastore --file %s", dataStorePath), tracetestcli.WithCLIConfig(cliConfig))
	require.Equal(0, result.ExitCode)

	// When I try to list datastore again on pretty mode
	// Then it should print a table with 4 lines printed: header, separator, data store item and empty line
	result = tracetestcli.Exec(t, "list datastore --output pretty", tracetestcli.WithCLIConfig(cliConfig))
	require.Equal(0, result.ExitCode)
	require.Contains(result.StdOut, "current")
	require.Contains(result.StdOut, env.Name())

	lines := strings.Split(result.StdOut, "\n")
	require.Len(lines, 4)

	// When I try to list datastore again on json mode
	// Then it should print a JSON list with one item
	result = tracetestcli.Exec(t, "list datastore --output json", tracetestcli.WithCLIConfig(cliConfig))
	require.Equal(0, result.ExitCode)

	dataStoresJSON := helpers.UnmarshalJSON[[]types.DataStoreResource](t, result.StdOut)

	require.Len(dataStoresJSON, 1)
	require.Equal("DataStore", dataStoresJSON[0].Type)
	require.Equal(env.Name(), dataStoresJSON[0].Spec.Name)
	require.True(dataStoresJSON[0].Spec.Default)

	// When I try to list datastore again on yaml mode
	// Then it should print a YAML list with one item
	result = tracetestcli.Exec(t, "list datastore --output yaml", tracetestcli.WithCLIConfig(cliConfig))
	require.Equal(0, result.ExitCode)

	dataStoresYAML := helpers.UnmarshalYAMLSequence[types.DataStoreResource](t, result.StdOut)

	require.Len(dataStoresYAML, 1)
	require.Equal("DataStore", dataStoresYAML[0].Type)
	require.Equal(env.Name(), dataStoresYAML[0].Spec.Name)
	require.True(dataStoresYAML[0].Spec.Default)
}
