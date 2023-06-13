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

func addGetPollingProfilePreReqs(t *testing.T, env environment.Manager) {
	cliConfig := env.GetCLIConfigPath(t)

	// When I try to set up a polling profile
	// Then it should be applied with success
	pollingProfilePath := env.GetTestResourcePath(t, "new-pollingprofile")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply pollingprofile --file %s", pollingProfilePath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)
}

func TestGetPollingProfile(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	t.Run("get with no polling profile initialized", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And no polling profile previously registered

		// When I try to get a polling profile on yaml mode
		// Then it should print a YAML with the default polling profile
		result := tracetestcli.Exec(t, "get pollingprofile --id current --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		require.Equal(0, result.ExitCode)

		pollingProfile := helpers.UnmarshalYAML[types.PollingProfileResource](t, result.StdOut)
		require.Equal("PollingProfile", pollingProfile.Type)
		require.Equal("current", pollingProfile.Spec.ID)
		require.Equal("Default", pollingProfile.Spec.Name)
		require.True(pollingProfile.Spec.Default)
		require.Equal("periodic", pollingProfile.Spec.Strategy)
		require.Equal("5s", pollingProfile.Spec.Periodic.RetryDelay)
		require.Equal("10m", pollingProfile.Spec.Periodic.Timeout)
	})

	addGetPollingProfilePreReqs(t, env)

	t.Run("get with YAML format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a polling profile already set

		// When I try to get a polling profile on yaml mode
		// Then it should print a YAML
		result := tracetestcli.Exec(t, "get pollingprofile --id current --output yaml", tracetestcli.WithCLIConfig(cliConfig))
		require.Equal(0, result.ExitCode)

		pollingProfile := helpers.UnmarshalYAML[types.PollingProfileResource](t, result.StdOut)
		require.Equal("PollingProfile", pollingProfile.Type)
		require.Equal("current", pollingProfile.Spec.ID)
		require.Equal("current", pollingProfile.Spec.Name)
		require.True(pollingProfile.Spec.Default)
		require.Equal("periodic", pollingProfile.Spec.Strategy)
		require.Equal("50s", pollingProfile.Spec.Periodic.RetryDelay)
		require.Equal("10m", pollingProfile.Spec.Periodic.Timeout)
	})

	t.Run("get with JSON format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a polling profile already set

		// When I try to get a polling profile on json mode
		// Then it should print a json
		result := tracetestcli.Exec(t, "get pollingprofile --id current --output json", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)

		pollingProfile := helpers.UnmarshalJSON[types.PollingProfileResource](t, result.StdOut)
		require.Equal("PollingProfile", pollingProfile.Type)
		require.Equal("current", pollingProfile.Spec.ID)
		require.Equal("current", pollingProfile.Spec.Name)
		require.True(pollingProfile.Spec.Default)
		require.Equal("periodic", pollingProfile.Spec.Strategy)
		require.Equal("50s", pollingProfile.Spec.Periodic.RetryDelay)
		require.Equal("10m", pollingProfile.Spec.Periodic.Timeout)
	})

	t.Run("get with pretty format", func(t *testing.T) {
		// Given I am a Tracetest CLI user
		// And I have my server recently created
		// And I have a polling profile already set

		// When I try to get a polling profile on pretty mode
		// Then it should print a table with 4 lines printed: header, separator, a polling profile item and empty line
		result := tracetestcli.Exec(t, "get pollingprofile --id current --output pretty", tracetestcli.WithCLIConfig(cliConfig))
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Contains(result.StdOut, "current")
		require.Contains(result.StdOut, "periodic")

		lines := strings.Split(result.StdOut, "\n")
		require.Len(lines, 4)
	})
}
