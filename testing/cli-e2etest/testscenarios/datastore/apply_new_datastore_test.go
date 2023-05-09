package datastore

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestApplyNewDatastore(t *testing.T) {
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to get a datastore without any server setup
	// Then it should return an empty datastore
	result, err := tracetestcli.Exec("get datastore --id current", tracetestcli.WithCLIConfig(cliConfig))
	require.NoError(t, err)
	require.Equal(t, 0, result.ExitCode)
	require.Contains(t, result.StdOut, "type: DataStore")
	require.Contains(t, result.StdOut, "default: false")
	require.Contains(t, result.StdOut, "id: \"\"")
	require.Contains(t, result.StdOut, "name: \"\"")
	require.Contains(t, result.StdOut, "type: \"\"")

	// When I try to set up a new datastore
	// Then it should be applied with success
	dataStorePath := env.GetManisfestResourcePath(t, "data-store")

	result, err = tracetestcli.Exec(fmt.Sprintf("apply datastore --file %s", dataStorePath), tracetestcli.WithCLIConfig(cliConfig))
	require.NoError(t, err)
	require.Equal(t, 0, result.ExitCode)

	// When I try to get a datastore again
	// Then it should return the datastore applied on the last step
	result, err = tracetestcli.Exec("get datastore --id current", tracetestcli.WithCLIConfig(cliConfig))
	require.NoError(t, err)
	require.Equal(t, 0, result.ExitCode)
	require.Contains(t, result.StdOut, "type: DataStore")
	require.Contains(t, result.StdOut, "default: true")
	require.Contains(t, result.StdOut, "id: current")
	// TODO: we are not testing datastore specific fields because they are dependant of the environment
	// but we could read the dataStore file in the future and try to do a comparison
}
