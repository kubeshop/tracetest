package config

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func addGetConfigPreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// When I try to set up a config
	// Then it should be applied with success
	configPath := env.GetTestResourcePath(t, "new-config")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply config --file %s", configPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
}

func TestGetConfig(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("get with no config initialized", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And no config previously registered

		// When I try to get a config on yaml mode
		// Then it should print a YAML with the default config
		result := tracetestcli.Exec(t, "get config --id current --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		require.Equal(0, result.ExitCode)

		config := helpers.UnmarshalYAML[types.ConfigResource](t, result.StdOut)
		require.Equal("Config", config.Type)
		require.Equal("current", config.Spec.ID)
		require.Equal("Config", config.Spec.Name)
		require.True(config.Spec.AnalyticsEnabled)
	})

	addGetConfigPreReqs(t, env)

	t.Run("get with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a config already set

		// When I try to get a config on yaml mode
		// Then it should print a YAML
		result := tracetestcli.Exec(t, "get config --id current --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		require.Equal(0, result.ExitCode)

		config := helpers.UnmarshalYAML[types.ConfigResource](t, result.StdOut)
		require.Equal("Config", config.Type)
		require.Equal("current", config.Spec.ID)
		require.Equal("Config", config.Spec.Name)
		require.False(config.Spec.AnalyticsEnabled)
	})

	t.Run("get with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a config already set

		// When I try to get a config on json mode
		// Then it should print a json
		result := tracetestcli.Exec(t, "get config --id current --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		config := helpers.UnmarshalJSON[types.ConfigResource](t, result.StdOut)
		require.Equal("Config", config.Type)
		require.Equal("current", config.Spec.ID)
		require.Equal("Config", config.Spec.Name)
		require.False(config.Spec.AnalyticsEnabled)
	})

	t.Run("get with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a config already set

		// When I try to get a config on pretty mode
		// Then it should print a table with 4 lines printed: header, separator, a config item and empty line
		result := tracetestcli.Exec(t, "get config --id current --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		parsedTable := helpers.UnmarshalTable(t, result.StdOut)
		require.Len(parsedTable, 1)

		singleLine := parsedTable[0]

		require.Equal("current", singleLine["ID"])
		require.Equal("Config", singleLine["NAME"])
		require.Equal("false", singleLine["ANALYTICS ENABLED"])
	})
}
