package testsuite

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/testing/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestRunTestSuite(t *testing.T) {
	t.Run("should fail if test resource is selected", func(t *testing.T) {
		// setup isolated e2e environment
		env := environment.CreateAndStart(t)
		defer env.Close(t)

		cliConfig := env.GetCLIConfigPath(t)

		// instantiate require with testing helper
		require := require.New(t)

		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And the datasource is already set

		// When I try to run a TestSuite
		// Then it should pass
		testsuiteFile := env.GetTestResourcePath(t, "testsuite-to-run")

		command := fmt.Sprintf("run test -f %s", testsuiteFile)
		result := tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 1)
		require.Contains(result.StdErr, "cannot apply TestSuite to Test resource")
	})

	t.Run("should pass", func(t *testing.T) {
		// setup isolated e2e environment
		env := environment.CreateAndStart(t, environment.WithDataStoreEnabled(), environment.WithPokeshop())
		defer env.Close(t)

		cliConfig := env.GetCLIConfigPath(t)

		// instantiate require with testing helper
		require := require.New(t)

		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And the datasource is already set

		// When I try to run a TestSuite
		// Then it should pass
		testsuiteFile := env.GetTestResourcePath(t, "testsuite-to-run")

		command := fmt.Sprintf("run testsuite -f %s", testsuiteFile)
		result := tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "TestSuite To Run") // testsuite name
		require.Contains(result.StdOut, "Pokeshop - Add")   // first test
		require.Contains(result.StdOut, "✔ It should add a Pokemon correctly")
		require.Contains(result.StdOut, "✔ It should save the correct data")
		require.Contains(result.StdOut, "Pokeshop - Get") // second test
		require.Contains(result.StdOut, "✔ It should Get Pokemons correctly")
	})

	t.Run("should run a legacy transaction", func(t *testing.T) {
		env := environment.CreateAndStart(t, environment.WithDataStoreEnabled(), environment.WithPokeshop())
		defer env.Close(t)

		cliConfig := env.GetCLIConfigPath(t)

		// instantiate require with testing helper
		require := require.New(t)

		testsuiteFile := env.GetTestResourcePath(t, "legacy-transaction")

		command := fmt.Sprintf("run transaction -f %s", testsuiteFile)
		result := tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))

		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "New Transaction") // testsuite name
		require.Contains(result.StdOut, "Pokeshop - Add")  // first test
		require.Contains(result.StdOut, "✔ It should add a Pokemon correctly")
		require.Contains(result.StdOut, "✔ It should save the correct data")
		require.Contains(result.StdOut, "Pokeshop - Get") // second test
		require.Contains(result.StdOut, "✔ It should Get Pokemons correctly")
	})
}
