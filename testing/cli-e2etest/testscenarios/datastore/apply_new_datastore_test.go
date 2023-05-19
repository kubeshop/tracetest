package datastore

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestApplyNewDatastore(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup isolated e2e environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to get a datastore without any server setup
	// Then it should return an empty datastore
	result := tracetestcli.Exec(t, "get datastore --id current", tracetestcli.WithCLIConfig(cliConfig))
	// TODO: we haven't defined a valid output to tell to the user that we are on `no-tracing mode`
	require.Equal(0, result.ExitCode)

	dataStore := helpers.UnmarshalYAML[types.DataStoreResource](t, result.StdOut)
	require.Equal("DataStore", dataStore.Type)
	require.False(dataStore.Spec.Default)

	// When I try to set up a new datastore
	// Then it should be applied with success
	dataStorePath := env.GetManisfestResourcePath(t, "data-store")

	result = tracetestcli.Exec(t, fmt.Sprintf("apply datastore --file %s", dataStorePath), tracetestcli.WithCLIConfig(cliConfig))
	require.Equal(0, result.ExitCode)

	// When I try to get a datastore again
	// Then it should return the datastore applied on the last step
	result = tracetestcli.Exec(t, "get datastore --id current", tracetestcli.WithCLIConfig(cliConfig))
	require.Equal(0, result.ExitCode)

	dataStore = helpers.UnmarshalYAML[types.DataStoreResource](t, result.StdOut)
	require.Equal("DataStore", dataStore.Type)
	require.Equal("current", dataStore.Spec.ID)
	require.Equal(env.Name(), dataStore.Spec.Name)
	require.True(dataStore.Spec.Default)
}
