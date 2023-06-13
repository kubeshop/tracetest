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

func addGetAnalyzerPreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// When I try to set up a config
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
		// And no config previously registered

		// When I try to get a config on yaml mode
		// Then it should print a YAML with the default config
		result := tracetestcli.Exec(t, "get analyzer --id current --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		require.Equal(0, result.ExitCode)

		config := helpers.UnmarshalYAML[types.AnalyzerResource](t, result.StdOut)
		require.Equal("Analyzer", config.Type)
		require.Equal("current", config.Spec.Id)
		require.Equal("analyzer", config.Spec.Name)
		require.True(config.Spec.Enabled)
		require.Equal(config.Spec.MinimumScore, 0)
		require.Len(config.Spec.Plugins, 3)
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

		config := helpers.UnmarshalYAML[types.AnalyzerResource](t, result.StdOut)
		require.Equal("Analyzer", config.Type)
		require.Equal("current", config.Spec.Id)
		require.Equal("analyzer", config.Spec.Name)
		require.True(config.Spec.Enabled)
		require.Equal(config.Spec.MinimumScore, 95)
		require.Len(config.Spec.Plugins, 3)
	})

	t.Run("get with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a config already set

		// When I try to get a config on json mode
		// Then it should print a json
		result := tracetestcli.Exec(t, "get analyzer --id current --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		config := helpers.UnmarshalJSON[types.AnalyzerResource](t, result.StdOut)
		require.Equal("Analyzer", config.Type)
		require.Equal("current", config.Spec.Id)
		require.Equal("analyzer", config.Spec.Name)
		require.True(config.Spec.Enabled)
		require.Equal(config.Spec.MinimumScore, 95)
		require.Len(config.Spec.Plugins, 3)
	})

	t.Run("get with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a config already set

		// When I try to get a config on pretty mode
		// Then it should print a table with 4 lines printed: header, separator, a config item and empty line
		result := tracetestcli.Exec(t, "get analyzer --id current --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "current")

		lines := strings.Split(result.StdOut, "\n")
		require.Len(lines, 4)
	})
}
