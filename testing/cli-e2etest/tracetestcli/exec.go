package tracetestcli

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/command"
	"github.com/stretchr/testify/require"
)

const (
	defaultTracetestCommand = "tracetest"
)

type ExecOption func(*executionState)

type executionState struct {
	cliConfigFile string
}

func Exec(t *testing.T, tracetestSubCommand string, options ...ExecOption) *command.ExecResult {
	state := &executionState{}
	for _, option := range options {
		option(state)
	}

	if state.cliConfigFile != "" {
		// append config at the start of the command
		tracetestSubCommand = fmt.Sprintf("--config %s %s", state.cliConfigFile, tracetestSubCommand)
	}

	tracetestCommand := getTracetestCommand()
	tracetestSubCommands := strings.Split(tracetestSubCommand, " ")

	result, err := command.Exec(tracetestCommand, tracetestSubCommands...)
	require.NoError(t, err)

	return result
}

func getTracetestCommand() string {
	tracetestCommand := os.Getenv("TRACETEST_COMMAND")

	if tracetestCommand == "" {
		return defaultTracetestCommand
	}

	return tracetestCommand
}

func WithCLIConfig(cliConfig string) ExecOption {
	return func(es *executionState) {
		es.cliConfigFile = cliConfig
	}
}
