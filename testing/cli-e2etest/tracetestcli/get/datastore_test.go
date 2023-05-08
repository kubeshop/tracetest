package get_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	getCommand "github.com/kubeshop/tracetest/cli-e2etest/tracetestcli/get"
	"github.com/stretchr/testify/require"
)

func TestGetDatastoreCommand(t *testing.T) {
	// test case: get current datastore without proper setup

	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfig(t)

	_, exitCode, err := getCommand.Exec("datastore --id current", tracetestcli.WithCLIConfig(cliConfig))
	require.NoError(t, err)

	require.Equal(t, 0, exitCode)
}
