package testscenarios

import (
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/helpers"
	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestHelpCommand(t *testing.T) {
	// instantiate require with testing helper
	require := require.New(t)

	// Given I am a Tracetest CLI user
	// When I try to get help with the commands "tracetest help", "tracetest --help" or "tracetest -h"
	// Then I should receive a message with sucess

	possibleCommands := []string{"help", "--help", "-h"}

	for _, helpCommand := range possibleCommands {
		result := tracetestcli.Exec(t, helpCommand)
		helpers.RequireExitCodeEqual(t, result, 0)
		require.Greater(len(result.StdOut), 0)
	}
}
