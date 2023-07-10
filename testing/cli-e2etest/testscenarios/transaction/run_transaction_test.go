package transaction

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestRunTransaction(t *testing.T) {
	// setup isolated e2e environment
	env := environment.CreateAndStart(t, environment.WithDataStoreEnabled(), environment.WithPokeshop())
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("should pass", func(t *testing.T) {
		// instantiate require with testing helper
		require := require.New(t)

		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And the datasource is already set

		// When I try to run a transaction
		// Then it should pass
		transactionFile := env.GetTestResourcePath(t, "transaction-to-run")

		command := fmt.Sprintf("run transaction -f %s", transactionFile)
		result := tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "Transaction To Run") // transaction name
		require.Contains(result.StdOut, "Pokeshop - Add")     // first test
		require.Contains(result.StdOut, "✔ It should add a Pokemon correctly")
		require.Contains(result.StdOut, "✔ It should save the correct data")
		require.Contains(result.StdOut, "Pokeshop - Get") // second test
		require.Contains(result.StdOut, "✔ It should Get Pokemons correctly")
	})
}
