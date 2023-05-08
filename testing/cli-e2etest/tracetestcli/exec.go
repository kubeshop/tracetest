package tracetestcli

import (
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli-e2etest/command"
)

const (
	tracetestCommand = "tracetest"
)

type CLIExecCommand func(string, ...ExecOption) (string, int, error)

func Exec(tracetestSubCommand string, options ...ExecOption) (string, int, error) {
	state := composeExecutionState(options...)

	if state.cliConfigFile != "" {
		// append config at the start of the command
		tracetestSubCommand = fmt.Sprintf("-c %s %s", state.cliConfigFile, tracetestSubCommand)
	}

	tracetestSubCommands := strings.Split(tracetestSubCommand, " ")

	return command.Exec(tracetestCommand, tracetestSubCommands...)
}

func GetComposedExecFunc(command string) CLIExecCommand {
	return func(subCommand string, opts ...ExecOption) (string, int, error) {
		fullSubCommand := fmt.Sprintf("%s %s", command, subCommand)
		return Exec(fullSubCommand, opts...)
	}
}
