package datastore

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestDeleteDatastore(t *testing.T) {
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to set up a new datastore
	// Then it should be applied with success
	dataStorePath := env.GetManisfestResourcePath(t, "data-store")

	result, err := tracetestcli.Exec(fmt.Sprintf("apply datastore --file %s", dataStorePath), tracetestcli.WithCLIConfig(cliConfig))
	require.NoError(t, err)
	require.Equal(t, 0, result.ExitCode)

	// When I try to get a datastore
	// Then it should return the datastore applied on the last step
	result, err = tracetestcli.Exec("get datastore --id current", tracetestcli.WithCLIConfig(cliConfig))
	require.NoError(t, err)
	require.Equal(t, 0, result.ExitCode)
	require.Contains(t, result.StdOut, "type: DataStore")
	require.Contains(t, result.StdOut, "default: true")
	require.Contains(t, result.StdOut, "id: current")

	// When I try to delete the datastore
	// Then it should delete with success
	result, err = tracetestcli.Exec("delete datastore --id current", tracetestcli.WithCLIConfig(cliConfig))
	require.NoError(t, err)
	require.Equal(t, 0, result.ExitCode)
	//TODO: on the future we should tell to the user that we will use an empty datastore again

	// When I try to get a datastore again
	// Then it should return an empty datastore
	result, err = tracetestcli.Exec("get datastore --id current", tracetestcli.WithCLIConfig(cliConfig))
	// TODO: we haven't defined a valid output to tell to the user that we are on `no-tracing mode`
	require.NoError(t, err)
	require.Equal(t, 0, result.ExitCode)
	require.Contains(t, result.StdOut, "type: DataStore")
	require.Contains(t, result.StdOut, "default: false")
}
