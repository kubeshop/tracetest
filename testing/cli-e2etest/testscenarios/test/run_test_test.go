package test

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestRunTransactionInsteadOfTest(t *testing.T) {
	t.Run("should fail if transaction resource is selected", func(t *testing.T) {
		// setup isolated e2e environment
		env := environment.CreateAndStart(t)
		defer env.Close(t)

		cliConfig := env.GetCLIConfigPath(t)

		// instantiate require with testing helper
		require := require.New(t)

		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And the datasource is already set

		// When I try to run a transaction
		// Then it should pass
		testFil := env.GetTestResourcePath(t, "import")

		command := fmt.Sprintf("run transaction -f %s", testFil)
		result := tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 1)
		require.Contains(result.StdErr, "cannot apply Test to Transaction resource")
	})
}

func TestRunTestWithHttpTriggerAndEnvironmentFile(t *testing.T) {
	// setup isolated e2e environment
	env := environment.CreateAndStart(t, environment.WithDataStoreEnabled(), environment.WithPokeshop())
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	// instantiate require with testing helper
	require := require.New(t)

	t.Run("should pass when using environment definition file", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And the datasource is already set

		// When I try to get an environment
		// Then it should return a message saying that the environment was not found
		result := tracetestcli.Exec(t, "get environment --id pokeapi-env", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "Resource environment with ID pokeapi-env not found")

		// When I try to run a test with a http trigger and a environment file
		// Then it should pass
		environmentFile := env.GetTestResourcePath(t, "environment-file")
		testFile := env.GetTestResourcePath(t, "http-trigger-with-environment-file")

		command := fmt.Sprintf("run test -f %s --environment %s", testFile, environmentFile)
		result = tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "✔ It should add a Pokemon correctly")
		require.Contains(result.StdOut, "✔ It should save the correct data")

		// When I try to get the environment created on the previous step
		// Then it should retrieve it correctly
		result = tracetestcli.Exec(t, "get environment --id pokeapi-env --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		environmentVars := helpers.UnmarshalYAML[types.EnvironmentResource](t, result.StdOut)
		require.Equal("Environment", environmentVars.Type)
		require.Equal("pokeapi-env", environmentVars.Spec.ID)
		require.Equal("pokeapi-env", environmentVars.Spec.Name)
		require.Len(environmentVars.Spec.Values, 2)
		require.Equal("POKEMON_NAME", environmentVars.Spec.Values[0].Key)
		require.Equal("snorlax", environmentVars.Spec.Values[0].Value)
		require.Equal("POKEMON_URL", environmentVars.Spec.Values[1].Key)
		require.Equal("https://assets.pokemon.com/assets/cms2/img/pokedex/full/143.png", environmentVars.Spec.Values[1].Value)
	})

	t.Run("should pass when using environment id", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And the datasource is already set

		// When I create an environment
		// Then it should be created correctly
		environmentFile := env.GetTestResourcePath(t, "environment-file")

		result := tracetestcli.Exec(t, fmt.Sprintf("apply environment --file %s", environmentFile), tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		environmentVars := helpers.UnmarshalYAML[types.EnvironmentResource](t, result.StdOut)
		require.Equal("Environment", environmentVars.Type)
		require.Equal("pokeapi-env", environmentVars.Spec.ID)
		require.Equal("pokeapi-env", environmentVars.Spec.Name)
		require.Len(environmentVars.Spec.Values, 2)
		require.Equal("POKEMON_NAME", environmentVars.Spec.Values[0].Key)
		require.Equal("snorlax", environmentVars.Spec.Values[0].Value)
		require.Equal("POKEMON_URL", environmentVars.Spec.Values[1].Key)
		require.Equal("https://assets.pokemon.com/assets/cms2/img/pokedex/full/143.png", environmentVars.Spec.Values[1].Value)

		// When I try to run a test with a http trigger and a environment id
		// Then it should pass

		testFile := env.GetTestResourcePath(t, "http-trigger-with-environment-file")

		command := fmt.Sprintf("run test -f %s --environment pokeapi-env", testFile)
		result = tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "✔ It should add a Pokemon correctly")
		require.Contains(result.StdOut, "✔ It should save the correct data")
	})
}

func TestRunTestWithGrpcTrigger(t *testing.T) {
	// setup isolated e2e environment
	env := environment.CreateAndStart(t, environment.WithDataStoreEnabled(), environment.WithPokeshop())
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("should pass when using an embedded protobuf string in the test", func(t *testing.T) {
		// instantiate require with testing helper
		require := require.New(t)

		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And the datasource is already set

		// When I try to run a test with a gRPC trigger with embedded protobuf
		// Then it should pass
		testFile := env.GetTestResourcePath(t, "grpc-trigger-embedded-protobuf")

		command := fmt.Sprintf("run test -f %s", testFile)
		result := tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "✔ It calls Pokeshop correctly") // checks if the assertion was succeeded
	})

	t.Run("should pass when referencing a protobuf file in the test", func(t *testing.T) {
		// instantiate require with testing helper
		require := require.New(t)

		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And the datasource is already set

		// When I try to run a test with a gRPC trigger with a reference to a protobuf file
		// Then it should pass
		testFile := env.GetTestResourcePath(t, "grpc-trigger-reference-protobuf")

		command := fmt.Sprintf("run test -f %s", testFile)
		result := tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "✔ It calls Pokeshop correctly") // checks if the assertion was succeeded
	})
}
