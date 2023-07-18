package testrunner

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func addGetTestRunnerPreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// When I try to set up a testrunner
	// Then it should be applied with success
	testRunnerPath := env.GetTestResourcePath(t, "new-testrunner")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply testrunner --file %s", testRunnerPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
}

func TestGetTestRunner(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)
	assert := assert.New(t)

	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("get with no testrunner initialized", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And no testrunner previously registered

		// When I try to get a testrunner on yaml mode
		// Then it should print a YAML with the default testrunner
		result := tracetestcli.Exec(t, "get testrunner --id current --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		require.Equal(0, result.ExitCode)

		testRunner := helpers.UnmarshalYAML[types.TestRunnerResource](t, result.StdOut)
		assert.Equal("TestRunner", testRunner.Type)
		assert.Equal("current", testRunner.Spec.ID)
		assert.Equal("default", testRunner.Spec.Name)
		require.Len(testRunner.Spec.RequiredGates, 2)
		assert.Equal("analyzer-score", testRunner.Spec.RequiredGates[0])
		assert.Equal("test-specs", testRunner.Spec.RequiredGates[1])
	})

	addGetTestRunnerPreReqs(t, env)

	t.Run("get with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a testrunner already set

		// When I try to get a testrunner on yaml mode
		// Then it should print a YAML
		result := tracetestcli.Exec(t, "get testrunner --id current --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		require.Equal(0, result.ExitCode)

		testRunner := helpers.UnmarshalYAML[types.TestRunnerResource](t, result.StdOut)
		assert.Equal("TestRunner", testRunner.Type)
		assert.Equal("current", testRunner.Spec.ID)
		assert.Equal("default", testRunner.Spec.Name)
		require.Len(testRunner.Spec.RequiredGates, 2)
		assert.Equal("analyzer-score", testRunner.Spec.RequiredGates[0])
		assert.Equal("test-specs", testRunner.Spec.RequiredGates[1])
	})

	t.Run("get with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a testrunner already set

		// When I try to get a testrunner on json mode
		// Then it should print a json
		result := tracetestcli.Exec(t, "get testrunner --id current --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		testRunner := helpers.UnmarshalJSON[types.TestRunnerResource](t, result.StdOut)
		assert.Equal("TestRunner", testRunner.Type)
		assert.Equal("current", testRunner.Spec.ID)
		assert.Equal("default", testRunner.Spec.Name)
		require.Len(testRunner.Spec.RequiredGates, 2)
		assert.Equal("analyzer-score", testRunner.Spec.RequiredGates[0])
		assert.Equal("test-specs", testRunner.Spec.RequiredGates[1])
	})

	t.Run("get with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a testrunner already set

		// When I try to get a testrunner on pretty mode
		// Then it should print a table with 4 lines printed: header, separator, a testrunner item and empty line
		result := tracetestcli.Exec(t, "get testrunner --id current --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		parsedTable := helpers.UnmarshalTable(t, result.StdOut)
		// this output shows one gate per line, so the parser reads that as an entire new row
		require.Len(parsedTable, 2)

		singleLine := parsedTable[0]
		require.Equal("current", singleLine["ID"])
		require.Equal("default", singleLine["NAME"])
		require.Equal("- analyzer-score", parsedTable[0]["REQUIRED GATES"])
		require.Equal("- test-specs", parsedTable[1]["REQUIRED GATES"])
	})
}
