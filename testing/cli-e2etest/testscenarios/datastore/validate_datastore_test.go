package datastore

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestValidateDatastore(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup isolated e2e environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("validate file with correct datastore setup", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// And I have a valid datastore file
		dataStorePath := env.GetEnvironmentResourcePath(t, "data-store")

		// When I try to validate it
		// Then it should return successful checks
		command := fmt.Sprintf("validate datastore --file %s", dataStorePath)
		result := tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		require.Contains(result.StdOut, "✔ \x1b[32mPort checking:\x1b[0m")
		require.Contains(result.StdOut, "✔ \x1b[32mConnectivity:\x1b[0m")
		require.Contains(result.StdOut, "✔ \x1b[32mAuthentication:\x1b[0m")
		require.Contains(result.StdOut, "✔ \x1b[32mFetch traces:\x1b[0m")
	})

	t.Run("validate file with incorrect datastore setup", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// And I have an invalid datastore file
		dataStorePath := env.GetEnvironmentResourcePath(t, "invalid-data-store")

		// When I try to validate it
		// Then it should return failed checks
		command := fmt.Sprintf("validate datastore --file %s", dataStorePath)
		result := tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		require.Contains(result.StdOut, "✔ \x1b[32mPort checking:\x1b[0m")
		require.Contains(result.StdOut, "✘ \x1b[31mConnectivity:\x1b[0m")
		require.NotContains(result.StdOut, "Authentication")
		require.NotContains(result.StdOut, "Fetch traces")
	})
}
