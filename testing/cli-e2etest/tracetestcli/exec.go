package tracetestcli

import (
	"fmt"
	"os"
	"strings"

	"github.com/kubeshop/tracetest/cli-e2etest/command"
)

const (
	defaultTracetestCommand = "tracetest"
)

type ExecOption func(*executionState)

type executionState struct {
	cliConfigFile string
}

func Exec(tracetestSubCommand string, options ...ExecOption) (*command.ExecResult, error) {
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

	return command.Exec(tracetestCommand, tracetestSubCommands...)
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
