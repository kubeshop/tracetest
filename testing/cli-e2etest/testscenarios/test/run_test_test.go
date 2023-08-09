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

func TestRunTestSuiteInsteadOfTest(t *testing.T) {
	t.Run("should fail if test suite resource is selected", func(t *testing.T) {
		// setup isolated e2e environment
		env := environment.CreateAndStart(t)
		defer env.Close(t)

		cliConfig := env.GetCLIConfigPath(t)

		// instantiate require with testing helper
		require := require.New(t)

		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And the datasource is already set

		// When I try to run a test suite
		// Then it should pass
		testFil := env.GetTestResourcePath(t, "import")

		command := fmt.Sprintf("run testsuite -f %s", testFil)
		result := tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 1)
		require.Contains(result.StdErr, "cannot apply Test to TestSuite resource")
	})
}

func TestRunTestWithHttpTriggerAndVariableSetFile(t *testing.T) {
	// setup isolated e2e environment
	env := environment.CreateAndStart(t, environment.WithDataStoreEnabled(), environment.WithPokeshop())
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	// instantiate require with testing helper
	require := require.New(t)

	t.Run("should pass when using variable set definition file", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And the datasource is already set

		// When I try to get a variable set
		// Then it should return a message saying that the variable set was not found
		result := tracetestcli.Exec(t, "get variableset --id pokeapi-env", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "Resource variableset with ID pokeapi-env not found")

		// When I try to run a test with a http trigger and a variable set file
		// Then it should pass
		environmentFile := env.GetTestResourcePath(t, "variableSet-file")
		testFile := env.GetTestResourcePath(t, "http-trigger-with-environment-file")

		command := fmt.Sprintf("run test -f %s --vars %s", testFile, environmentFile)
		result = tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "✔ It should add a Pokemon correctly")
		require.Contains(result.StdOut, "✔ It should save the correct data")

		// When I try to get the variable set created on the previous step
		// Then it should retrieve it correctly
		result = tracetestcli.Exec(t, "get variableset --id pokeapi-env --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		environmentVars := helpers.UnmarshalYAML[types.VariableSetResource](t, result.StdOut)
		require.Equal("VariableSet", environmentVars.Type)
		require.Equal("pokeapi-env", environmentVars.Spec.ID)
		require.Equal("pokeapi-env", environmentVars.Spec.Name)
		require.Len(environmentVars.Spec.Values, 2)
		require.Equal("POKEMON_NAME", environmentVars.Spec.Values[0].Key)
		require.Equal("snorlax", environmentVars.Spec.Values[0].Value)
		require.Equal("POKEMON_URL", environmentVars.Spec.Values[1].Key)
		require.Equal("https://assets.pokemon.com/assets/cms2/img/pokedex/full/143.png", environmentVars.Spec.Values[1].Value)
	})

	t.Run("should pass when using the deprecated environment definition file", func(t *testing.T) {
		result := tracetestcli.Exec(t, "get environment --id deprecated-pokeapi-env", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "The resource `environment` is deprecated and will be removed in a future version. Please use `variableset` instead.")
		require.Contains(result.StdOut, "Resource variableset with ID deprecated-pokeapi-env not found")

		// When I try to run a test with a http trigger and a variable set file
		// Then it should pass
		environmentFile := env.GetTestResourcePath(t, "deprecated-environment")
		testFile := env.GetTestResourcePath(t, "http-trigger-with-environment-file")

		command := fmt.Sprintf("run test -f %s --environment %s", testFile, environmentFile)
		result = tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "✔ It should add a Pokemon correctly")
		require.Contains(result.StdOut, "✔ It should save the correct data")

		// When I try to get the variable set created on the previous step
		// Then it should retrieve it correctly
		result = tracetestcli.Exec(t, "get environment --id deprecated-pokeapi-env", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		require.Contains(result.StdOut, "The resource `environment` is deprecated and will be removed in a future version. Please use `variableset` instead.")
		require.Contains(result.StdOut, "VariableSet")
		require.Contains(result.StdOut, "https://assets.pokemon.com/assets/cms2/img/pokedex/full/143.png")
	})

	t.Run("should pass when using variable set id", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And the datasource is already set

		// When I create a variable set
		// Then it should be created correctly
		environmentFile := env.GetTestResourcePath(t, "variableSet-file")

		result := tracetestcli.Exec(t, fmt.Sprintf("apply variableset --file %s", environmentFile), tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		environmentVars := helpers.UnmarshalYAML[types.VariableSetResource](t, result.StdOut)
		require.Equal("VariableSet", environmentVars.Type)
		require.Equal("pokeapi-env", environmentVars.Spec.ID)
		require.Equal("pokeapi-env", environmentVars.Spec.Name)
		require.Len(environmentVars.Spec.Values, 2)
		require.Equal("POKEMON_NAME", environmentVars.Spec.Values[0].Key)
		require.Equal("snorlax", environmentVars.Spec.Values[0].Value)
		require.Equal("POKEMON_URL", environmentVars.Spec.Values[1].Key)
		require.Equal("https://assets.pokemon.com/assets/cms2/img/pokedex/full/143.png", environmentVars.Spec.Values[1].Value)

		// When I try to run a test with a http trigger and a variable set id
		// Then it should pass

		testFile := env.GetTestResourcePath(t, "http-trigger-with-environment-file")

		command := fmt.Sprintf("run test -f %s --vars pokeapi-env", testFile)
		result = tracetestcli.Exec(t, command, tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "✔ It should add a Pokemon correctly")
		require.Contains(result.StdOut, "✔ It should save the correct data")
	})

	t.Run("should pass when using variable set id but using the old environment flag", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And the datasource is already set

		// When I create a variable set
		// Then it should be created correctly
		environmentFile := env.GetTestResourcePath(t, "variableSet-file")

		result := tracetestcli.Exec(t, fmt.Sprintf("apply variableset --file %s", environmentFile), tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		environmentVars := helpers.UnmarshalYAML[types.VariableSetResource](t, result.StdOut)
		require.Equal("VariableSet", environmentVars.Type)
		require.Equal("pokeapi-env", environmentVars.Spec.ID)
		require.Equal("pokeapi-env", environmentVars.Spec.Name)
		require.Len(environmentVars.Spec.Values, 2)
		require.Equal("POKEMON_NAME", environmentVars.Spec.Values[0].Key)
		require.Equal("snorlax", environmentVars.Spec.Values[0].Value)
		require.Equal("POKEMON_URL", environmentVars.Spec.Values[1].Key)
		require.Equal("https://assets.pokemon.com/assets/cms2/img/pokedex/full/143.png", environmentVars.Spec.Values[1].Value)

		// When I try to run a test with a http trigger and a variable set id
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
