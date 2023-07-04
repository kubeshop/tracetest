package analyzer

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func addGetAnalyzerPreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// When I try to set up a analyzer
	// Then it should be applied with success
	configPath := env.GetTestResourcePath(t, "new-analyzer")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply analyzer --file %s", configPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
}

func TestGetAnalyzer(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("get with no analyzer initialized", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And no analyzer previously registered

		// When I try to get a analyzer on yaml mode
		// Then it should print a YAML with the default analyzer
		result := tracetestcli.Exec(t, "get analyzer --id current --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		require.Equal(0, result.ExitCode)

		analyzer := helpers.UnmarshalYAML[types.AnalyzerResource](t, result.StdOut)
		require.Equal("Analyzer", analyzer.Type)
		require.Equal("current", analyzer.Spec.Id)
		require.Equal("analyzer", analyzer.Spec.Name)
		require.True(analyzer.Spec.Enabled)
		require.Equal(analyzer.Spec.MinimumScore, 0)
		require.Len(analyzer.Spec.Plugins, 3)
	})

	addGetAnalyzerPreReqs(t, env)

	t.Run("get with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a config already set

		// When I try to get a config on yaml mode
		// Then it should print a YAML
		result := tracetestcli.Exec(t, "get analyzer --id current --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		require.Equal(0, result.ExitCode)

		analyzer := helpers.UnmarshalYAML[types.AnalyzerResource](t, result.StdOut)
		require.Equal("Analyzer", analyzer.Type)
		require.Equal("current", analyzer.Spec.Id)
		require.Equal("analyzer", analyzer.Spec.Name)
		require.True(analyzer.Spec.Enabled)
		require.Equal(analyzer.Spec.MinimumScore, 95)
		require.Len(analyzer.Spec.Plugins, 3)
	})

	t.Run("get with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a analyzer already set

		// When I try to get a analyzer on json mode
		// Then it should print a json
		result := tracetestcli.Exec(t, "get analyzer --id current --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		analyzer := helpers.UnmarshalJSON[types.AnalyzerResource](t, result.StdOut)
		require.Equal("Analyzer", analyzer.Type)
		require.Equal("current", analyzer.Spec.Id)
		require.Equal("analyzer", analyzer.Spec.Name)
		require.True(analyzer.Spec.Enabled)
		require.Equal(analyzer.Spec.MinimumScore, 95)
		require.Len(analyzer.Spec.Plugins, 3)
	})

	t.Run("get with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a analyzer already set

		// When I try to get a analyzer on pretty mode
		// Then it should print a table with 4 lines printed: header, separator, a analyzer item and empty line
		result := tracetestcli.Exec(t, "get analyzer --id current --output pretty", tracetestcli.WithCLIConfig(cliConfig))
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
