package analyzer

import (
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/testing/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/testscenarios/types"
	"github.com/kubeshop/tracetest/testing/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestApplyAnalyzer(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup isolated e2e environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	cliConfig := env.GetCLIConfigPath(t)

	// Given I am a Tracetest CLI user
	// And I have my server recently created

	// When I try to set up a new analyzer
	// Then it should be applied with success
	configPath := env.GetTestResourcePath(t, "new-analyzer")

	result := tracetestcli.Exec(t, fmt.Sprintf("apply analyzer --file %s", configPath), tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	// When I try to get a analyzer again
	// Then it should return the analyzer applied on the last step, with analytics disabled
	result = tracetestcli.Exec(t, "get analyzer --id current", tracetestcli.WithCLIConfig(cliConfig))
	helpers.RequireExitCodeEqual(t, result, 0)

	analyzer := helpers.UnmarshalYAML[types.AnalyzerResource](t, result.StdOut)
	require.Equal("Analyzer", analyzer.Type)
	require.Equal("current", analyzer.Spec.ID)
	require.Equal("analyzer", analyzer.Spec.Name)
	require.True(analyzer.Spec.Enabled)
	require.Equal(analyzer.Spec.MinimumScore, 95)
	require.Len(analyzer.Spec.Plugins, 3)

	plugin1 := analyzer.Spec.Plugins[0]
	require.Len(plugin1.Rules, 4)
	require.Equal(plugin1.Rules[0].ID, "span-naming")
	for _, rule := range plugin1.Rules {
		require.Equal(rule.Weight, 25)
	}
}
