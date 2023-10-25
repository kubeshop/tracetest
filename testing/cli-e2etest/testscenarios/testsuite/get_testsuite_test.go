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

func addGetTestSuitePreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to set up a new testsuite
	// Then it should be applied with success
	newTestSuitePath := env.GetTestResourcePath(t, "new-testsuite")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply testsuite --file %s", newTestSuitePath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
}

func TestGetTestSuite(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("get with no testsuite initialized", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And no testsuite registered

		// When I try to get a TestSuite on yaml mode
		// Then it should return a error message
		result := tracetestcli.Exec(t, "get testsuite --id no-id --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "Resource testsuite with ID no-id not found")
	})

	addGetTestSuitePreReqs(t, env)

	t.Run("get with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a TestSuite already set

		// When I try to get a TestSuite on yaml mode
		// Then it should print a YAML
		result := tracetestcli.Exec(t, "get testsuite --id Qti5R3_VR --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		testsuite := helpers.UnmarshalYAML[types.TestSuiteResource](t, result.StdOut)

		require.Equal("TestSuite", testsuite.Type)
		require.Equal("Qti5R3_VR", testsuite.Spec.ID)
		require.Equal("New TestSuite", testsuite.Spec.Name)
		require.Equal("a TestSuite", testsuite.Spec.Description)
		require.Len(testsuite.Spec.Steps, 2)
		require.Equal("9wtAH2_Vg", testsuite.Spec.Steps[0])
		require.Equal("ajksdkasjbd", testsuite.Spec.Steps[1])
	})

	t.Run("get with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a TestSuite already set

		// When I try to get a TestSuite on json mode
		// Then it should print a json
		result := tracetestcli.Exec(t, "get testsuite --id Qti5R3_VR --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		// it should return an augmented resource on get
		testsuite := helpers.UnmarshalJSON[types.AugmentedTestSuiteResource](t, result.StdOut)

		require.Equal("TestSuite", testsuite.Type)
		require.Equal("Qti5R3_VR", testsuite.Spec.ID)
		require.Equal("New TestSuite", testsuite.Spec.Name)
		require.Equal("a TestSuite", testsuite.Spec.Description)
		require.Len(testsuite.Spec.Steps, 2)
		require.Equal("9wtAH2_Vg", testsuite.Spec.Steps[0])
		require.Equal("ajksdkasjbd", testsuite.Spec.Steps[1])
		require.Equal(0, testsuite.Spec.Summary.Runs)
		require.Equal(0, testsuite.Spec.Summary.LastRun.Fails)
		require.Equal(0, testsuite.Spec.Summary.LastRun.Passes)
	})

	t.Run("get with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a TestSuite already set

		// When I try to get a TestSuite on pretty mode
		// Then it should print a table with 4 lines printed: header, separator, testsuite item and empty line
		result := tracetestcli.Exec(t, "get testsuite --id Qti5R3_VR --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		parsedTable := helpers.UnmarshalTable(t, result.StdOut)
		require.Len(parsedTable, 1)

		singleLine := parsedTable[0]

		require.Equal("Qti5R3_VR", singleLine["ID"])
		require.Equal("New TestSuite", singleLine["NAME"])
		require.Equal("1", singleLine["VERSION"])
		require.Equal("2", singleLine["STEPS"])
		require.Equal("0", singleLine["RUNS"])
		require.Equal("", singleLine["LAST RUN TIME"])
		require.Equal("0", singleLine["LAST RUN SUCCESSES"])
		require.Equal("0", singleLine["LAST RUN FAILURES"])
	})
}
