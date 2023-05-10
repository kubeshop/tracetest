package datastore

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestListDatastore(t *testing.T) {
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

	// When I try to list datastore again
	// Then it should receive a error message saying that a datastore cannot be listed
	result, err = tracetestcli.Exec("list datastore", tracetestcli.WithCLIConfig(cliConfig))
	require.NoError(t, err)
	require.Contains(t, result.StdOut, "DataStore does not support listing. Try `tracetest get datastore`")
}
