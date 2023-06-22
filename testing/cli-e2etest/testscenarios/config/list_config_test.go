package config

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

func addListConfigPreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// When I try to set up a new config
	// Then it should be applied with success
	configPath := env.GetTestResourcePath(t, "new-config")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply config --file %s", configPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
}

func TestListConfig(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("list with no config initialized", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created

		// When I try to list config on pretty mode and there is no config previously registered
		// Then it should print an empty table
		// Then it should print a table with 4 lines printed: header, separator, the default config item and empty line
		result := tracetestcli.Exec(t, "list config --sortBy name --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "current")

		lines := strings.Split(result.StdOut, "\n")
		require.Len(lines, 4)
	})

	addListConfigPreReqs(t, env)

	t.Run("list with invalid sortBy field", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a config created

		// When I try to list a config by an invalid field
		// Then I should receive an error
		result := tracetestcli.Exec(t, "list config --sortBy id --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 1)
		require.Contains(result.StdErr, "invalid sort field: id") // TODO: think on how to improve this error handling
	})

	t.Run("list with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a config created

		// When I try to list config again on yaml mode
		// Then it should print a YAML list with one item
		result := tracetestcli.Exec(t, "list config --sortBy name --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		configYAML := helpers.UnmarshalYAMLSequence[types.ConfigResource](t, result.StdOut)

		require.Len(configYAML, 1)
		require.Equal("Config", configYAML[0].Type)
		require.Equal("current", configYAML[0].Spec.ID)
		require.Equal("Config", configYAML[0].Spec.Name)
		require.False(configYAML[0].Spec.AnalyticsEnabled)
	})

	t.Run("list with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a config created

		// When I try to list config again on json mode
		// Then it should print a JSON list with one item
		result := tracetestcli.Exec(t, "list config --sortBy name --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		configYAML := helpers.UnmarshalJSON[[]types.ConfigResource](t, result.StdOut)

		require.Len(configYAML, 1)
		require.Equal("Config", configYAML[0].Type)
		require.Equal("current", configYAML[0].Spec.ID)
		require.Equal("Config", configYAML[0].Spec.Name)
		require.False(configYAML[0].Spec.AnalyticsEnabled)
	})

	t.Run("list with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I already have a config created

		// When I try to list config again on pretty mode
		// Then it should print a table with 4 lines printed: header, separator, config item and empty line
		result := tracetestcli.Exec(t, "list config --sortBy name --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		parsedTable := helpers.UnmarshalTable(t, result.StdOut)
		require.Len(parsedTable, 1)

		singleLine := parsedTable[0]

		require.Equal("current", singleLine["ID"])
		require.Equal("Config", singleLine["NAME"])
		require.Equal("false", singleLine["ANALYTICS ENABLED"])
	})
}
