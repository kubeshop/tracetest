package analyzer

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

func addListAnalyzerPreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// When I try to set up a new analyzer
	// Then it should be applied with success
	configPath := env.GetTestResourcePath(t, "new-analyzer")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply analyzer --file %s", configPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
}

func TestListAnalyzer(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("list with no analyzer initialized", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list analyzer on pretty mode and there is no analyzer previously registered
		// Then it should print an empty table
		// Then it should print a table with 4 lines printed: header, separator, the default analyzer item and empty line
		result := tracetestcli.Exec(t, "list analyzer --sortBy name --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "current")

		lines := strings.Split(result.StdOut, "\n")
		require.Len(lines, 4)
	})

	addListAnalyzerPreReqs(t, env)

	t.Run("list with invalid sortBy field", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a analyzer created

		// When I try to list a analyzer by an invalid field
		// Then I should receive an error
		result := tracetestcli.Exec(t, "list analyzer --sortBy id --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 1)
		require.Contains(result.StdErr, "invalid sort field: id") // TODO: think on how to improve this error handling
	})

	t.Run("list with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a analyzer created

		// When I try to list analyzer again on yaml mode
		// Then it should print a YAML list with one item
		result := tracetestcli.Exec(t, "list analyzer --sortBy name --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		analyzerYAML := helpers.UnmarshalYAMLSequence[types.AnalyzerResource](t, result.StdOut)

		require.Len(analyzerYAML, 1)
		require.Equal("Analyzer", analyzerYAML[0].Type)
		require.Equal("current", analyzerYAML[0].Spec.Id)
		require.Equal("analyzer", analyzerYAML[0].Spec.Name)
		require.True(analyzerYAML[0].Spec.Enabled)
		require.Equal(analyzerYAML[0].Spec.MinimumScore, 95)
		require.Len(analyzerYAML[0].Spec.Plugins, 3)
	})

	t.Run("list with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a analyzer created

		// When I try to list analyzer again on json mode
		// Then it should print a JSON list with one item
		result := tracetestcli.Exec(t, "list analyzer --sortBy name --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		analyzerList := helpers.UnmarshalJSON[types.ResourceList[types.AnalyzerResource]](t, result.StdOut)
		require.Len(analyzerList.Items, 1)
		require.Equal(len(analyzerList.Items), analyzerList.Count)

		require.Equal("Analyzer", analyzerList.Items[0].Type)
		require.Equal("current", analyzerList.Items[0].Spec.Id)
		require.Equal("analyzer", analyzerList.Items[0].Spec.Name)
		require.True(analyzerList.Items[0].Spec.Enabled)
		require.Equal(analyzerList.Items[0].Spec.MinimumScore, 95)
		require.Len(analyzerList.Items[0].Spec.Plugins, 3)
	})

	t.Run("list with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a analyzer created

		// When I try to list analyzer again on pretty mode
		// Then it should print a table with 4 lines printed: header, separator, analyzer item and empty line
		result := tracetestcli.Exec(t, "list analyzer --sortBy name --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		parsedTable := helpers.UnmarshalTable(t, result.StdOut)
		require.Len(parsedTable, 1)

		singleLine := parsedTable[0]

		require.Equal("current", singleLine["ID"])
		require.Equal("analyzer", singleLine["NAME"])
		require.Equal("true", singleLine["ENABLED"])
		require.Equal("95", singleLine["MINIMUM SCORE"])
	})
}
