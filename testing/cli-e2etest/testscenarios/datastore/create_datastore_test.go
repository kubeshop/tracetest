package datastore

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestCreateDatastore(t *testing.T) {
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to get a datastore without any server setup
	// Then I should receive an error, telling there is no datastores registered
	result, err := tracetestcli.Exec("get datastore --id current", tracetestcli.WithCLIConfig(cliConfig))
	require.ErrorContains(t, err, "invalid datastores:")
	require.ErrorContains(t, err, "record not found")
	require.Nil(t, result)

	// When I try to set up a new datastore
	// Then I should receive success message
	dataStorePath := env.GetManisfestResourcePath(t, "data-store")

	result, err = tracetestcli.Exec(fmt.Sprintf("apply datastore --file %s", dataStorePath), tracetestcli.WithCLIConfig(cliConfig))
	require.NoError(t, err)
	require.Equal(t, 0, result.ExitCode)
}
