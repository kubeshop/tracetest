package test

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func addListTestsPreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to set up a new test
	// Then it should be applied with success
	newTestPath := env.GetTestResourcePath(t, "list")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply test --file %s", newTestPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	// When I try to set up a another test
	// Then it should be applied with success
	anotherTestPath := env.GetTestResourcePath(t, "import")

	result = tracetestcli.Exec(t, fmt.Sprintf("apply test --file %s", anotherTestPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
}

func TestListTests(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)
	assert := assert.New(t)

	// setup isolated e2e test
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("list no tests", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And there is no envs
		result := tracetestcli.Exec(t, "list test --sortBy name --sortDirection asc --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		testVarsList := helpers.UnmarshalYAMLSequence[types.TestResource](t, result.StdOut)
		require.Len(testVarsList, 0)
	})

	addListTestsPreReqs(t, env)

	t.Run("list with invalid sortBy field", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these tests by an invalid field
		// Then I should receive an error
		result := tracetestcli.Exec(t, "list test --sortBy id --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 1)
		require.Contains(result.StdErr, "invalid sort field: id") // TODO: think on how to improve this error handling
	})

	t.Run("list with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have 2 existing tests

		// When I try to list these tests by a valid field and in YAML format
		// Then I should receive 2 tests
		result := tracetestcli.Exec(t, "list test --sortBy name --sortDirection desc --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		testVarsList := helpers.UnmarshalYAMLSequence[types.TestResource](t, result.StdOut)
		require.Len(testVarsList, 2)

		spew.Dump(testVarsList)

		listTest := testVarsList[0]
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

		importTest := testVarsList[1]
		assert.Equal("Test", importTest.Type)
		assert.Equal("RXrbV__4g", importTest.Spec.ID)
		assert.Equal("Pokeshop - Import", importTest.Spec.Name)
		assert.Equal("Import a Pokemon", importTest.Spec.Description)
		assert.Equal("http", importTest.Spec.Trigger.Type)
		assert.Equal("http://demo-api:8081/pokemon/import", importTest.Spec.Trigger.HTTPRequest.URL)
		assert.Equal("POST", importTest.Spec.Trigger.HTTPRequest.Method)
		assert.Equal(`{"id":52}`, importTest.Spec.Trigger.HTTPRequest.Body)
		require.Len(importTest.Spec.Trigger.HTTPRequest.Headers, 1)
		assert.Equal("Content-Type", importTest.Spec.Trigger.HTTPRequest.Headers[0].Key)
		assert.Equal("application/json", importTest.Spec.Trigger.HTTPRequest.Headers[0].Value)
	})

	// 	t.Run("list with JSON format", func(t *testing.T) {
	// 		// Given I am a Tracetest CLI user
	// 		// And I have my server recently created

	// 		// When I try to list these tests by a valid field and in JSON format
	// 		// Then I should receive three tests
	// 		result := tracetestcli.Exec(t, "list test --sortBy name --sortDirection asc --output json", tracetestcli.WithCLIConfig(cliConfig))
	// 		helpers.RequireExitCodeEqual(t, result, 0)

	// 		testVarsList := helpers.UnmarshalJSON[[]types.TestResource](t, result.StdOut)
	// 		require.Len(testVarsList, 3)

	// 		// due our database sorting algorithm, "another-env" comes in the front of ".env"
	// 		// ref https://wiki.postgresql.org/wiki/FAQ#Why_do_my_strings_sort_incorrectly.3F
	// 		anotherTestVars := testVarsList[0]
	// 		require.Equal("Test", anotherTestVars.Type)
	// 		require.Equal("another-env", anotherTestVars.Spec.ID)
	// 		require.Equal("another-env", anotherTestVars.Spec.Name)
	// 		require.Len(anotherTestVars.Spec.Values, 2)
	// 		require.Equal("Here", anotherTestVars.Spec.Values[0].Key)
	// 		require.Equal("We", anotherTestVars.Spec.Values[0].Value)
	// 		require.Equal("Come", anotherTestVars.Spec.Values[1].Key)
	// 		require.Equal("Again", anotherTestVars.Spec.Values[1].Value)

	// 		testVars := testVarsList[1]
	// 		require.Equal("Test", testVars.Type)
	// 		require.Equal(".env", testVars.Spec.ID)
	// 		require.Equal(".env", testVars.Spec.Name)
	// 		require.Len(testVars.Spec.Values, 2)
	// 		require.Equal("FIRST_VAR", testVars.Spec.Values[0].Key)
	// 		require.Equal("some-value", testVars.Spec.Values[0].Value)
	// 		require.Equal("SECOND_VAR", testVars.Spec.Values[1].Key)
	// 		require.Equal("another_value", testVars.Spec.Values[1].Value)

	// 		oneMoreTestVars := testVarsList[2]
	// 		require.Equal("Test", oneMoreTestVars.Type)
	// 		require.Equal("one-more-env", oneMoreTestVars.Spec.ID)
	// 		require.Equal("one-more-env", oneMoreTestVars.Spec.Name)
	// 		require.Len(oneMoreTestVars.Spec.Values, 2)
	// 		require.Equal("This", oneMoreTestVars.Spec.Values[0].Key)
	// 		require.Equal("Is", oneMoreTestVars.Spec.Values[0].Value)
	// 		require.Equal("The", oneMoreTestVars.Spec.Values[1].Key)
	// 		require.Equal("Third", oneMoreTestVars.Spec.Values[1].Value)
	// 	})

	// 	t.Run("list with pretty format", func(t *testing.T) {
	// 		// Given I am a Tracetest CLI user
	// 		// And I have my server recently created

	// 		// When I try to list these tests by a valid field and in pretty format
	// 		// Then it should print a table with 6 lines printed: header, separator, three envs and empty line
	// 		result := tracetestcli.Exec(t, "list test --sortBy name --sortDirection asc --output pretty", tracetestcli.WithCLIConfig(cliConfig))
	// 		helpers.RequireExitCodeEqual(t, result, 0)

	// 		lines := strings.Split(result.StdOut, "\n")
	// 		require.Len(lines, 6)

	// 		// due our database sorting algorithm, "another-env" comes in the front of ".env"
	// 		// ref https://wiki.postgresql.org/wiki/FAQ#Why_do_my_strings_sort_incorrectly.3F
	// 		require.Contains(lines[2], "another-env")
	// 		require.Contains(lines[3], ".env")
	// 		require.Contains(lines[4], "one-more-env")
	// 	})

	// 	t.Run("list with YAML format skipping the first and taking two items", func(t *testing.T) {
	// 		// Given I am a Tracetest CLI user
	// 		// And I have my server recently created

	// 		// When I try to list these tests by a valid field, paging options and in YAML format
	// 		// Then I should receive two tests
	// 		result := tracetestcli.Exec(t, "list test --sortBy name --sortDirection asc --skip 1 --take 2 --output yaml", tracetestcli.WithCLIConfig(cliConfig))
	// 		helpers.RequireExitCodeEqual(t, result, 0)

	// 		testVarsList := helpers.UnmarshalYAMLSequence[types.TestResource](t, result.StdOut)
	// 		require.Len(testVarsList, 2)

	// 		// due our database sorting algorithm, "another-env" comes in the front of ".env"
	// 		// ref https://wiki.postgresql.org/wiki/FAQ#Why_do_my_strings_sort_incorrectly.3F
	// 		testVars := testVarsList[0]
	// 		require.Equal("Test", testVars.Type)
	// 		require.Equal(".env", testVars.Spec.ID)
	// 		require.Equal(".env", testVars.Spec.Name)
	// 		require.Len(testVars.Spec.Values, 2)
	// 		require.Equal("FIRST_VAR", testVars.Spec.Values[0].Key)
	// 		require.Equal("some-value", testVars.Spec.Values[0].Value)
	// 		require.Equal("SECOND_VAR", testVars.Spec.Values[1].Key)
	// 		require.Equal("another_value", testVars.Spec.Values[1].Value)

	//		oneMoreTestVars := testVarsList[1]
	//		require.Equal("Test", oneMoreTestVars.Type)
	//		require.Equal("one-more-env", oneMoreTestVars.Spec.ID)
	//		require.Equal("one-more-env", oneMoreTestVars.Spec.Name)
	//		require.Len(oneMoreTestVars.Spec.Values, 2)
	//		require.Equal("This", oneMoreTestVars.Spec.Values[0].Key)
	//		require.Equal("Is", oneMoreTestVars.Spec.Values[0].Value)
	//		require.Equal("The", oneMoreTestVars.Spec.Values[1].Key)
	//		require.Equal("Third", oneMoreTestVars.Spec.Values[1].Value)
	//	})
}
