package testrunner

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kubeshop/tracetest/testing/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func addListConfigPreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// When I try to set up a new testrunner
	// Then it should be applied with success
	testRunnerPath := env.GetTestResourcePath(t, "new-testrunner")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply testrunner --file %s", testRunnerPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
}

func TestListConfig(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)
	assert := assert.New(t)

	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("list with no testrunner initialized", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list testrunner on pretty mode and there is no testrunner previously registered
		// Then it should print an empty table
		// Then it should print a table with 5 lines printed: header, separator, the default testrunner item, an entire line for the second gate and empty line
		result := tracetestcli.Exec(t, "list testrunner --sortBy name --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "current")

		lines := strings.Split(result.StdOut, "\n")
		require.Len(lines, 5)
	})

	addListConfigPreReqs(t, env)

	t.Run("list with invalid sortBy field", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a testrunner created

		// When I try to list a testrunner by an invalid field
		// Then I should receive an error
		result := tracetestcli.Exec(t, "list testrunner --sortBy id --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 1)
		require.Contains(result.StdErr, "invalid sort field: id") // TODO: think on how to improve this error handling
	})

	t.Run("list with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a testrunner created

		// When I try to list testrunner again on yaml mode
		// Then it should print a YAML list with one item
		result := tracetestcli.Exec(t, "list testrunner --sortBy name --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		testRunner := helpers.UnmarshalYAML[types.TestRunnerResource](t, result.StdOut)
		assert.Equal("TestRunner", testRunner.Type)
		assert.Equal("current", testRunner.Spec.ID)
		assert.Equal("default", testRunner.Spec.Name)
		require.Len(testRunner.Spec.RequiredGates, 2)
		assert.Equal("analyzer-score", testRunner.Spec.RequiredGates[0])
		assert.Equal("test-specs", testRunner.Spec.RequiredGates[1])
	})

	t.Run("list with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a testrunner created

		// When I try to list testrunner again on json mode
		// Then it should print a JSON list with one item
		result := tracetestcli.Exec(t, "list testrunner --sortBy name --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		testRunnerList := helpers.UnmarshalJSON[types.ResourceList[types.TestRunnerResource]](t, result.StdOut)
		require.Len(testRunnerList.Items, 1)
		require.Equal(len(testRunnerList.Items), testRunnerList.Count)

		testRunner := testRunnerList.Items[0]
		assert.Equal("TestRunner", testRunner.Type)
		assert.Equal("current", testRunner.Spec.ID)
		assert.Equal("default", testRunner.Spec.Name)
		require.Len(testRunner.Spec.RequiredGates, 2)
		assert.Equal("analyzer-score", testRunner.Spec.RequiredGates[0])
		assert.Equal("test-specs", testRunner.Spec.RequiredGates[1])
	})

	t.Run("list with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a testrunner created

		// When I try to list testrunner again on pretty mode
		// Then it should print a table with 4 lines printed: header, separator, testrunner item and empty line
		result := tracetestcli.Exec(t, "list testrunner --sortBy name --output pretty", tracetestcli.WithCLIConfig(cliConfig))
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
