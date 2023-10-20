package test

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/testing/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func addGetTestPreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to set up a new test
	// Then it should be applied with success
	newTestPath := env.GetTestResourcePath(t, "list")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply test --file %s", newTestPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
}

func TestGetTest(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("get with no test initialized", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And no test registered

		// When I try to get a test on yaml mode
		// Then it should return a error message
		result := tracetestcli.Exec(t, "get test --id fH_8AulVR --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "Resource test with ID fH_8AulVR not found")
	})

	addGetTestPreReqs(t, env)

	t.Run("get with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have an test already set

		// When I try to get an test on yaml mode
		// Then it should print a YAML
		result := tracetestcli.Exec(t, "get test --id fH_8AulVR --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		listTest := helpers.UnmarshalYAML[types.TestResource](t, result.StdOut)
		assert.Equal("Test", listTest.Type)
		assert.Equal("fH_8AulVR", listTest.Spec.ID)
		assert.Equal("Pokeshop - List", listTest.Spec.Name)
		assert.Equal("List Pokemon", listTest.Spec.Description)
		assert.Equal("http", listTest.Spec.Trigger.Type)
		assert.Equal("http://demo-api:8081/pokemon?take=20&skip=0", listTest.Spec.Trigger.HTTPRequest.URL)
		assert.Equal("GET", listTest.Spec.Trigger.HTTPRequest.Method)
		assert.Equal("", listTest.Spec.Trigger.HTTPRequest.Body)
		require.Len(listTest.Spec.Trigger.HTTPRequest.Headers, 1)
		assert.Equal("Content-Type", listTest.Spec.Trigger.HTTPRequest.Headers[0].Key)
		assert.Equal("application/json", listTest.Spec.Trigger.HTTPRequest.Headers[0].Value)
	})

	t.Run("get with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have an test already set

		// When I try to get an test on json mode
		// Then it should print a json
		result := tracetestcli.Exec(t, "get test --id fH_8AulVR --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		listTest := helpers.UnmarshalJSON[types.TestResource](t, result.StdOut)
		assert.Equal("Test", listTest.Type)
		assert.Equal("fH_8AulVR", listTest.Spec.ID)
		assert.Equal("Pokeshop - List", listTest.Spec.Name)
		assert.Equal("List Pokemon", listTest.Spec.Description)
		assert.Equal("http", listTest.Spec.Trigger.Type)
		assert.Equal("http://demo-api:8081/pokemon?take=20&skip=0", listTest.Spec.Trigger.HTTPRequest.URL)
		assert.Equal("GET", listTest.Spec.Trigger.HTTPRequest.Method)
		assert.Equal("", listTest.Spec.Trigger.HTTPRequest.Body)
		require.Len(listTest.Spec.Trigger.HTTPRequest.Headers, 1)
		assert.Equal("Content-Type", listTest.Spec.Trigger.HTTPRequest.Headers[0].Key)
		assert.Equal("application/json", listTest.Spec.Trigger.HTTPRequest.Headers[0].Value)
	})

	t.Run("get with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have an test already set

		// When I try to get an test on pretty mode
		// Then it should print a table with 4 lines printed: header, separator, test item and empty line
		result := tracetestcli.Exec(t, "get test --id fH_8AulVR --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		parsedTable := helpers.UnmarshalTable(t, result.StdOut)
		require.Len(parsedTable, 1)

		singleLine := parsedTable[0]

		require.Equal("fH_8AulVR", singleLine["ID"])
		require.Equal("Pokeshop - List", singleLine["NAME"])
		require.Equal("1", singleLine["VERSION"])
		require.Equal("http", singleLine["TRIGGER TYPE"])
		require.Equal("0", singleLine["RUNS"])
		require.Equal("", singleLine["LAST RUN TIME"])
		require.Equal("0", singleLine["LAST RUN SUCCESSES"])
		require.Equal("0", singleLine["LAST RUN FAILURES"])
		require.Equal("http://localhost:11633/test/fH_8AulVR", singleLine["URL"])
	})
}
