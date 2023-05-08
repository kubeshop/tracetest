package get_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	getCommand "github.com/kubeshop/tracetest/cli-e2etest/tracetestcli/get"
	"github.com/stretchr/testify/require"
)

func TestGetDatastoreCommand(t *testing.T) {
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfig(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to get a datastore without any server setup
	// Then I should receive an error, telling there is no datastores registered
	result, err := getCommand.Exec("datastore --id current", tracetestcli.WithCLIConfig(cliConfig))
	require.ErrorContains(t, err, "invalid datastores: \"record not found\"")
	require.Equal(t, 1, result.ExitCode)
}
