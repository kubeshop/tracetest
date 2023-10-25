package testsuite

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/testing/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func addListTestSuitePreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to set up a new environment
	// Then it should be applied with success
	newTestSuitePath := env.GetTestResourcePath(t, "new-testsuite")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply testsuite --file %s", newTestSuitePath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	// When I try to set up a another environment
	// Then it should be applied with success
	anotherTestSuitePath := env.GetTestResourcePath(t, "another-testsuite")

	result = tracetestcli.Exec(t, fmt.Sprintf("apply testsuite --file %s", anotherTestSuitePath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	// When I try to set up a third environment
	// Then it should be applied with success
	oneMoreTestSuitePath := env.GetTestResourcePath(t, "one-more-testsuite")

	result = tracetestcli.Exec(t, fmt.Sprintf("apply testsuite --file %s", oneMoreTestSuitePath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
}

func TestListTestSuites(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup isolated e2e environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("list no testsuites", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And there is no envs
		result := tracetestcli.Exec(t, "list testsuite --sortBy name --sortDirection asc --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		testsuites := helpers.UnmarshalYAMLSequence[types.AugmentedTestSuiteResource](t, result.StdOut)
		require.Len(testsuites, 0)
	})

	addListTestSuitePreReqs(t, env)

	t.Run("list with invalid sortBy field", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these testsuites by an invalid field
		// Then I should receive an error
		result := tracetestcli.Exec(t, "list testsuite --sortBy id --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 1)
		require.Contains(result.StdErr, "invalid sort field: id") // TODO: think on how to improve this error handling
	})

	t.Run("list with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these testsuites by a valid field and in YAML format
		// Then I should receive three testsuites
		result := tracetestcli.Exec(t, "list testsuite --sortBy name --sortDirection asc --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		testsuites := helpers.UnmarshalYAMLSequence[types.AugmentedTestSuiteResource](t, result.StdOut)
		require.Len(testsuites, 3)

		anotherTestSuite := testsuites[0]
		require.Equal("TestSuite", anotherTestSuite.Type)
		require.Equal("asuhfdkj", anotherTestSuite.Spec.ID)
		require.Equal("Another TestSuite", anotherTestSuite.Spec.Name)
		require.Equal("another TestSuite", anotherTestSuite.Spec.Description)
		require.Equal(0, anotherTestSuite.Spec.Summary.Runs)
		require.Equal(0, anotherTestSuite.Spec.Summary.LastRun.Fails)
		require.Equal(0, anotherTestSuite.Spec.Summary.LastRun.Passes)
		require.Len(anotherTestSuite.Spec.Steps, 4)
		require.Equal("9wtAH2_Vg", anotherTestSuite.Spec.Steps[0])
		require.Equal("9wtAH2_Vg", anotherTestSuite.Spec.Steps[1])
		require.Equal("ajksdkasjbd", anotherTestSuite.Spec.Steps[2])
		require.Equal("ajksdkasjbd", anotherTestSuite.Spec.Steps[3])

		newTestSuite := testsuites[1]
		require.Equal("TestSuite", newTestSuite.Type)
		require.Equal("Qti5R3_VR", newTestSuite.Spec.ID)
		require.Equal("New TestSuite", newTestSuite.Spec.Name)
		require.Equal("a TestSuite", newTestSuite.Spec.Description)
		require.Equal(0, newTestSuite.Spec.Summary.Runs)
		require.Equal(0, newTestSuite.Spec.Summary.LastRun.Fails)
		require.Equal(0, newTestSuite.Spec.Summary.LastRun.Passes)
		require.Len(newTestSuite.Spec.Steps, 2)
		require.Equal("9wtAH2_Vg", newTestSuite.Spec.Steps[0])
		require.Equal("ajksdkasjbd", newTestSuite.Spec.Steps[1])

		oneMoreTestSuite := testsuites[2]
		require.Equal("TestSuite", oneMoreTestSuite.Type)
		require.Equal("i2ug34j", oneMoreTestSuite.Spec.ID)
		require.Equal("One More TestSuite", oneMoreTestSuite.Spec.Name)
		require.Equal("one more TestSuite", oneMoreTestSuite.Spec.Description)
		require.Equal(0, oneMoreTestSuite.Spec.Summary.Runs)
		require.Equal(0, oneMoreTestSuite.Spec.Summary.LastRun.Fails)
		require.Equal(0, oneMoreTestSuite.Spec.Summary.LastRun.Passes)
		require.Len(oneMoreTestSuite.Spec.Steps, 3)
		require.Equal("9wtAH2_Vg", oneMoreTestSuite.Spec.Steps[0])
		require.Equal("9wtAH2_Vg", oneMoreTestSuite.Spec.Steps[1])
		require.Equal("ajksdkasjbd", oneMoreTestSuite.Spec.Steps[2])
	})

	t.Run("list with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these testsuites by a valid field and in JSON format
		// Then I should receive three testsuites
		result := tracetestcli.Exec(t, "list testsuite --sortBy name --sortDirection asc --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		testsuites := helpers.UnmarshalJSON[types.ResourceList[types.AugmentedTestSuiteResource]](t, result.StdOut)
		require.Len(testsuites.Items, 3)
		require.Equal(len(testsuites.Items), testsuites.Count)

		anotherTestSuite := testsuites.Items[0]
		require.Equal("TestSuite", anotherTestSuite.Type)
		require.Equal("asuhfdkj", anotherTestSuite.Spec.ID)
		require.Equal("Another TestSuite", anotherTestSuite.Spec.Name)
		require.Equal("another TestSuite", anotherTestSuite.Spec.Description)
		require.Equal(0, anotherTestSuite.Spec.Summary.Runs)
		require.Equal(0, anotherTestSuite.Spec.Summary.LastRun.Fails)
		require.Equal(0, anotherTestSuite.Spec.Summary.LastRun.Passes)
		require.Len(anotherTestSuite.Spec.Steps, 4)
		require.Equal("9wtAH2_Vg", anotherTestSuite.Spec.Steps[0])
		require.Equal("9wtAH2_Vg", anotherTestSuite.Spec.Steps[1])
		require.Equal("ajksdkasjbd", anotherTestSuite.Spec.Steps[2])
		require.Equal("ajksdkasjbd", anotherTestSuite.Spec.Steps[3])

		newTestSuite := testsuites.Items[1]
		require.Equal("TestSuite", newTestSuite.Type)
		require.Equal("Qti5R3_VR", newTestSuite.Spec.ID)
		require.Equal("New TestSuite", newTestSuite.Spec.Name)
		require.Equal("a TestSuite", newTestSuite.Spec.Description)
		require.Equal(0, newTestSuite.Spec.Summary.Runs)
		require.Equal(0, newTestSuite.Spec.Summary.LastRun.Fails)
		require.Equal(0, newTestSuite.Spec.Summary.LastRun.Passes)
		require.Len(newTestSuite.Spec.Steps, 2)
		require.Equal("9wtAH2_Vg", newTestSuite.Spec.Steps[0])
		require.Equal("ajksdkasjbd", newTestSuite.Spec.Steps[1])

		oneMoreTestSuite := testsuites.Items[2]
		require.Equal("TestSuite", oneMoreTestSuite.Type)
		require.Equal("i2ug34j", oneMoreTestSuite.Spec.ID)
		require.Equal("One More TestSuite", oneMoreTestSuite.Spec.Name)
		require.Equal("one more TestSuite", oneMoreTestSuite.Spec.Description)
		require.Equal(0, oneMoreTestSuite.Spec.Summary.Runs)
		require.Equal(0, oneMoreTestSuite.Spec.Summary.LastRun.Fails)
		require.Equal(0, oneMoreTestSuite.Spec.Summary.LastRun.Passes)
		require.Len(oneMoreTestSuite.Spec.Steps, 3)
		require.Equal("9wtAH2_Vg", oneMoreTestSuite.Spec.Steps[0])
		require.Equal("9wtAH2_Vg", oneMoreTestSuite.Spec.Steps[1])
		require.Equal("ajksdkasjbd", oneMoreTestSuite.Spec.Steps[2])
	})

	t.Run("list with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these testsuites by a valid field and in pretty format
		// Then it should print a table with 6 lines printed: header, separator, three testsuites and empty line
		result := tracetestcli.Exec(t, "list testsuite --sortBy name --sortDirection asc --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		parsedTable := helpers.UnmarshalTable(t, result.StdOut)
		require.Len(parsedTable, 3)

		firstLine := parsedTable[0]
		require.Equal("asuhfdkj", firstLine["ID"])
		require.Equal("Another TestSuite", firstLine["NAME"])
		require.Equal("1", firstLine["VERSION"])
		require.Equal("4", firstLine["STEPS"])
		require.Equal("0", firstLine["RUNS"])
		require.Equal("", firstLine["LAST RUN TIME"])
		require.Equal("0", firstLine["LAST RUN SUCCESSES"])
		require.Equal("0", firstLine["LAST RUN FAILURES"])

		secondLine := parsedTable[1]
		require.Equal("Qti5R3_VR", secondLine["ID"])
		require.Equal("New TestSuite", secondLine["NAME"])
		require.Equal("1", secondLine["VERSION"])
		require.Equal("2", secondLine["STEPS"])
		require.Equal("0", secondLine["RUNS"])
		require.Equal("", secondLine["LAST RUN TIME"])
		require.Equal("0", secondLine["LAST RUN SUCCESSES"])
		require.Equal("0", secondLine["LAST RUN FAILURES"])

		thirdLine := parsedTable[2]
		require.Equal("i2ug34j", thirdLine["ID"])
		require.Equal("One More TestSuite", thirdLine["NAME"])
		require.Equal("1", thirdLine["VERSION"])
		require.Equal("3", thirdLine["STEPS"])
		require.Equal("0", thirdLine["RUNS"])
		require.Equal("", thirdLine["LAST RUN TIME"])
		require.Equal("0", thirdLine["LAST RUN SUCCESSES"])
		require.Equal("0", thirdLine["LAST RUN FAILURES"])
	})

	t.Run("list with YAML format skipping the first and taking two items", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list these testsuites by a valid field, paging options and in YAML format
		// Then I should receive two testsuites
		result := tracetestcli.Exec(t, "list testsuite --sortBy name --sortDirection asc --skip 1 --take 2 --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		testsuites := helpers.UnmarshalYAMLSequence[types.AugmentedTestSuiteResource](t, result.StdOut)
		require.Len(testsuites, 2)

		newTestSuite := testsuites[0]
		require.Equal("TestSuite", newTestSuite.Type)
		require.Equal("Qti5R3_VR", newTestSuite.Spec.ID)
		require.Equal("New TestSuite", newTestSuite.Spec.Name)
		require.Equal("a TestSuite", newTestSuite.Spec.Description)
		require.Equal(0, newTestSuite.Spec.Summary.Runs)
		require.Equal(0, newTestSuite.Spec.Summary.LastRun.Fails)
		require.Equal(0, newTestSuite.Spec.Summary.LastRun.Passes)
		require.Len(newTestSuite.Spec.Steps, 2)
		require.Equal("9wtAH2_Vg", newTestSuite.Spec.Steps[0])
		require.Equal("ajksdkasjbd", newTestSuite.Spec.Steps[1])

		oneMoreTestSuite := testsuites[1]
		require.Equal("TestSuite", oneMoreTestSuite.Type)
		require.Equal("i2ug34j", oneMoreTestSuite.Spec.ID)
		require.Equal("One More TestSuite", oneMoreTestSuite.Spec.Name)
		require.Equal("one more TestSuite", oneMoreTestSuite.Spec.Description)
		require.Equal(0, oneMoreTestSuite.Spec.Summary.Runs)
		require.Equal(0, oneMoreTestSuite.Spec.Summary.LastRun.Fails)
		require.Equal(0, oneMoreTestSuite.Spec.Summary.LastRun.Passes)
		require.Len(oneMoreTestSuite.Spec.Steps, 3)
		require.Equal("9wtAH2_Vg", oneMoreTestSuite.Spec.Steps[0])
		require.Equal("9wtAH2_Vg", oneMoreTestSuite.Spec.Steps[1])
		require.Equal("ajksdkasjbd", oneMoreTestSuite.Spec.Steps[2])
	})
}
