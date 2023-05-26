package testscenarios

import (
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/environment"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestVersionCommand(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// setup test server environment
	env := environment.CreateAndStart(t)
	defer env.Close(t)

	// Given I am a Tracetest CLI user
	// When I try to check the tracetest version
	// Then I should receive a version string with success
	result := tracetestcli.Exec(t, "version")

	require.Equal(0, result.ExitCode)
	require.Greater(len(result.StdOut), 0)
}
