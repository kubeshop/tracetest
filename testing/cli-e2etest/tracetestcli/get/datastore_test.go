package get_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	getCommand "github.com/kubeshop/tracetest/cli-e2etest/tracetestcli/get"
	"github.com/stretchr/testify/require"
)

func TestGetDatastoreCommand(t *testing.T) {
	// test case: get current datastore without proper setup

	env := environment.CreateAndStart(t)
	defer env.Close(t)

	_, exitCode, err := getCommand.Exec("datastore --id current")
	require.NoError(t, err)

	require.Equal(t, 0, exitCode)
}
