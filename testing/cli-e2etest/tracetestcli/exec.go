package tracetestcli

import (
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli-e2etest/command"
)

const (
	tracetestCommand = "tracetest"
)

type CLIExecCommand func(string, ...ExecOption) (*command.ExecResult, error)

func Exec(tracetestSubCommand string, options ...ExecOption) (*command.ExecResult, error) {
	state := composeExecutionState(options...)

	if state.cliConfigFile != "" {
		// append config at the start of the command
		tracetestSubCommand = fmt.Sprintf("-c %s %s", state.cliConfigFile, tracetestSubCommand)
	}

	tracetestSubCommands := strings.Split(tracetestSubCommand, " ")

	return command.Exec(tracetestCommand, tracetestSubCommands...)
}

func GetComposedExecFunc(cliCommand string) CLIExecCommand {
	return func(subCommand string, opts ...ExecOption) (*command.ExecResult, error) {
		fullSubCommand := fmt.Sprintf("%s %s", cliCommand, subCommand)
		return Exec(fullSubCommand, opts...)
	}
}
