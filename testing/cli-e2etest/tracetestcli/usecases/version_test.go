package usecases

import (
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestVersionCommand(t *testing.T) {
	// Given I am a Tracetest CLI user
	// When I try to check the tracetest version
	// Then I should receive a version string with sucess

	result, err := tracetestcli.Exec("version")
	require.NoError(t, err)
	require.Equal(t, 0, result.ExitCode)
	require.Greater(t, len(result.StdOut), 0)
}
