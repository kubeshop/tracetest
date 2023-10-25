package config

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/testing/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestApplyConfig(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup isolated e2e environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to set up a new config
	// Then it should be applied with success
	configPath := env.GetTestResourcePath(t, "new-config")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply config --file %s", configPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	// When I try to get a config again
	// Then it should return the config applied on the last step, with analytics disabled
	result = tracetestcli.Exec(t, "get config --id current", tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	config := helpers.UnmarshalYAML[types.ConfigResource](t, result.StdOut)
	require.Equal("Config", config.Type)
	require.Equal("current", config.Spec.ID)
	require.Equal("Config", config.Spec.Name)
	require.False(config.Spec.AnalyticsEnabled) // disabling analytics
}
