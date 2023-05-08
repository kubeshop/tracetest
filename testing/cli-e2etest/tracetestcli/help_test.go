package tracetestcli_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
	"github.com/stretchr/testify/require"
)

func TestHelpCommand(t *testing.T) {
	expectedExitCode := 0

	possibleCommands := []string{"help", "--help", "-h"}

	for _, helpCommand := range possibleCommands {
		t.Run(helpCommand, func(t *testing.T) {
			_, exitCode, err := tracetestcli.Exec(helpCommand)
			require.NoError(t, err)

			require.Equal(t, expectedExitCode, exitCode)
		})
	}
}
