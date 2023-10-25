package testsuite

import (
	"fmt"
	"os"
	"testing"

	"github.com/kubeshop/tracetest/testing/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestApplyTestSuite(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup isolated e2e environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("should apply a test suite", func(t *testing.T) {

		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to set up a new testsuite
		// Then it should be applied with success
		newTestSuitePath := env.GetTestResourcePath(t, "new-testsuite")

		result := tracetestcli.Exec(t, fmt.Sprintf("apply testsuite --file %s", newTestSuitePath), tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		testsuite := helpers.UnmarshalYAML[types.TestSuiteResource](t, result.StdOut)

		require.Equal("TestSuite", testsuite.Type)
		require.Equal("Qti5R3_VR", testsuite.Spec.ID)
		require.Equal("New TestSuite", testsuite.Spec.Name)
		require.Equal("a TestSuite", testsuite.Spec.Description)
		require.Len(testsuite.Spec.Steps, 2)
		require.Equal("9wtAH2_Vg", testsuite.Spec.Steps[0])
		require.Equal("ajksdkasjbd", testsuite.Spec.Steps[1])

		// When I try to get the testsuite applied on the last step
		// Then it should return it
		result = tracetestcli.Exec(t, "get testsuite --id Qti5R3_VR --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		require.Equal("TestSuite", testsuite.Type)
		require.Equal("Qti5R3_VR", testsuite.Spec.ID)
		require.Equal("New TestSuite", testsuite.Spec.Name)
		require.Equal("a TestSuite", testsuite.Spec.Description)
		require.Len(testsuite.Spec.Steps, 2)
		require.Equal("9wtAH2_Vg", testsuite.Spec.Steps[0])
		require.Equal("ajksdkasjbd", testsuite.Spec.Steps[1])

		// When I try to update the last testsuite
		// Then it should be applied with success
		updatedNewTestSuitePath := env.GetTestResourcePath(t, "updated-new-testsuite")

		result = tracetestcli.Exec(t, fmt.Sprintf("apply testsuite --file %s", updatedNewTestSuitePath), tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		updatedTestSuite := helpers.UnmarshalYAML[types.TestSuiteResource](t, result.StdOut)
		require.Equal("TestSuite", updatedTestSuite.Type)
		require.Equal("Qti5R3_VR", updatedTestSuite.Spec.ID)
		require.Equal("Updated TestSuite", updatedTestSuite.Spec.Name)
		require.Equal("an updated TestSuite", updatedTestSuite.Spec.Description)
		require.Len(updatedTestSuite.Spec.Steps, 3)
		require.Equal("9wtAH2_Vg", updatedTestSuite.Spec.Steps[0])
		require.Equal("ajksdkasjbd", updatedTestSuite.Spec.Steps[1])
		require.Equal("ajksdkasjbd", updatedTestSuite.Spec.Steps[2])

		// When I try to get the testsuite applied on the last step
		// Then it should return it
		result = tracetestcli.Exec(t, "get testsuite --id Qti5R3_VR --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		updatedTestSuite = helpers.UnmarshalYAML[types.TestSuiteResource](t, result.StdOut)
		require.Equal("TestSuite", updatedTestSuite.Type)
		require.Equal("Qti5R3_VR", updatedTestSuite.Spec.ID)
		require.Equal("Updated TestSuite", updatedTestSuite.Spec.Name)
		require.Equal("an updated TestSuite", updatedTestSuite.Spec.Description)
		require.Len(updatedTestSuite.Spec.Steps, 3)
		require.Equal("9wtAH2_Vg", updatedTestSuite.Spec.Steps[0])
		require.Equal("ajksdkasjbd", updatedTestSuite.Spec.Steps[1])
		require.Equal("ajksdkasjbd", updatedTestSuite.Spec.Steps[2])

		// When I try to set up a new testsuite without any id
		// Then it should be applied with success and it should not update
		// the steps with its ids.
		testsuiteWithoutIDPath := env.GetTestResourcePath(t, "new-testsuite-without-id")
		helpers.Copy(testsuiteWithoutIDPath+".tpl", testsuiteWithoutIDPath)

		helpers.RemoveIDFromTestSuiteFile(t, testsuiteWithoutIDPath)

		testsuiteWithoutIDResult := tracetestcli.Exec(t, fmt.Sprintf("apply testsuite --file %s", testsuiteWithoutIDPath), tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, testsuiteWithoutIDResult, 0)

		content, err := os.ReadFile(testsuiteWithoutIDPath)
		require.NoError(err)

		testsuiteWithoutID := helpers.UnmarshalYAML[types.TestSuiteResource](t, string(content))

		require.Equal("TestSuite", testsuiteWithoutID.Type)
		require.NotEmpty(testsuiteWithoutID.Spec.ID)
		require.Equal("New TestSuite", testsuiteWithoutID.Spec.Name)
		require.Equal("a TestSuite", testsuiteWithoutID.Spec.Description)
		require.Len(testsuiteWithoutID.Spec.Steps, 2)
		require.Equal("./testsuite-step-1.yaml", testsuiteWithoutID.Spec.Steps[0])
		require.Equal("./testsuite-step-2.yaml", testsuiteWithoutID.Spec.Steps[1])
	})

	t.Run("should apply a legacy transaction", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to set up a legacy transaction
		// Then it should be applied with success
		newTestSuitePath := env.GetTestResourcePath(t, "legacy-transaction")

		result := tracetestcli.Exec(t, fmt.Sprintf("apply transaction --file %s", newTestSuitePath), tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		require.Contains(result.StdOut, "TestSuite")
		require.Contains(result.StdOut, "Qti5R3_VR")
		require.Contains(result.StdOut, "New Transaction")
		require.Contains(result.StdOut, "a Transaction")
		require.Contains(result.StdOut, "9wtAH2_Vg")
		require.Contains(result.StdOut, "ajksdkasjbd")

		// When I try to get the testsuite applied on the last step
		// Then it should return it
		result = tracetestcli.Exec(t, "get transaction --id Qti5R3_VR --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		require.Contains(result.StdOut, "TestSuite")
		require.Contains(result.StdOut, "Qti5R3_VR")
		require.Contains(result.StdOut, "New Transaction")
		require.Contains(result.StdOut, "a Transaction")
		require.Contains(result.StdOut, "9wtAH2_Vg")
		require.Contains(result.StdOut, "ajksdkasjbd")
	})
}
