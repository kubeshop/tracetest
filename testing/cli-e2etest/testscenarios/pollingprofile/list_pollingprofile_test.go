package pollingprofile

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func addListPollingProfilePreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// When I try to set up a new config
	// Then it should be applied with success
	pollingProfilePath := env.GetTestResourcePath(t, "new-pollingprofile")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply pollingprofile --file %s", pollingProfilePath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
}

func TestListPollingProfile(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("list with no polling profile initialized", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list polling profile on pretty mode and there is no polling profile previously registered
		// Then it should print an empty table
		// Then it should print a table with 4 lines printed: header, separator, the default polling profile item and empty line
		result := tracetestcli.Exec(t, "list pollingprofile --sortBy name --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "current")  // id
		require.Contains(result.StdOut, "Default")  // name
		require.Contains(result.StdOut, "periodic") // strategy

		lines := strings.Split(result.StdOut, "\n")
		require.Len(lines, 4)
	})

	addListPollingProfilePreReqs(t, env)

	t.Run("list with invalid sortBy field", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a polling profile created

		// When I try to list a polling profile by an invalid field
		// Then I should receive an error
		result := tracetestcli.Exec(t, "list pollingprofile --sortBy id --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 1)
		require.Contains(result.StdErr, "invalid sort field: id") // TODO: think on how to improve this error handling
	})

	t.Run("list with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a polling profile created

		// When I try to list polling profile again on yaml mode
		// Then it should print a YAML list with one item
		result := tracetestcli.Exec(t, "list pollingprofile --sortBy name --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		pollingProfileYAML := helpers.UnmarshalYAMLSequence[types.PollingProfileResource](t, result.StdOut)

		require.Len(pollingProfileYAML, 1)
		require.Equal("PollingProfile", pollingProfileYAML[0].Type)
		require.Equal("current", pollingProfileYAML[0].Spec.ID)
		require.Equal("current", pollingProfileYAML[0].Spec.Name)
		require.True(pollingProfileYAML[0].Spec.Default)
		require.Equal("periodic", pollingProfileYAML[0].Spec.Strategy)
		require.Equal("50s", pollingProfileYAML[0].Spec.Periodic.RetryDelay)
		require.Equal("10m", pollingProfileYAML[0].Spec.Periodic.Timeout)
	})

	t.Run("list with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a polling profile created

		// When I try to list polling profile again on json mode
		// Then it should print a JSON list with one item
		result := tracetestcli.Exec(t, "list pollingprofile --sortBy name --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		pollingProfileYAML := helpers.UnmarshalJSON[[]types.PollingProfileResource](t, result.StdOut)

		require.Len(pollingProfileYAML, 1)
		require.Equal("PollingProfile", pollingProfileYAML[0].Type)
		require.Equal("current", pollingProfileYAML[0].Spec.ID)
		require.Equal("current", pollingProfileYAML[0].Spec.Name)
		require.True(pollingProfileYAML[0].Spec.Default)
		require.Equal("periodic", pollingProfileYAML[0].Spec.Strategy)
		require.Equal("50s", pollingProfileYAML[0].Spec.Periodic.RetryDelay)
		require.Equal("10m", pollingProfileYAML[0].Spec.Periodic.Timeout)
	})

	t.Run("list with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a polling profile created

		// When I try to list polling profile again on pretty mode
		// Then it should print a table with 4 lines printed: header, separator, polling profile item and empty line
		result := tracetestcli.Exec(t, "list pollingprofile --sortBy name --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		parsedTable := helpers.UnmarshalTable(t, result.StdOut)
		require.Len(parsedTable, 1)

		singleLine := parsedTable[0]

		require.Equal("current", singleLine["ID"])
		require.Equal("current", singleLine["NAME"])
		require.Equal("periodic", singleLine["STRATEGY"])
	})
}
