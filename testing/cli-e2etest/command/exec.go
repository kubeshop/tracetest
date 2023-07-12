package command

import (
	"fmt"
	"os/exec"
	"strings"
)

type ExecResult struct {
	CommandExecuted string
	StdOut          string
	StdErr          string
	ExitCode        int
}

func (r *ExecResult) String() string {
	return fmt.Sprintf("commandText: [%s] \nexitCode: %d \nstdout: [%s] \nstderr: [%s]", r.CommandExecuted, r.ExitCode, r.StdOut, r.StdErr)
}

func Exec(command string, args ...string) (*ExecResult, error) {
	checkIfCommandExists(command)

	fullCommand := fmt.Sprintf("%s %s", command, strings.Join(args, " "))

	cmd := exec.Command(command, args...)

	out, err := cmd.Output()
	output := string(out)

	if err != nil {
		return handleRunError(err, fullCommand, output, cmd)
	}

	exitCode := cmd.ProcessState.ExitCode()
	return &ExecResult{
		CommandExecuted: fullCommand,
		StdOut:          output,
		ExitCode:        exitCode,
	}, nil
}

func checkIfCommandExists(command string) {
	path, err := exec.LookPath(command)
	if err != nil {
		panic(fmt.Sprintf("error when checking if %s exists, error: %s", command, err.Error()))
	}

	if path == "" {
		panic(fmt.Sprintf("%s was not found on this machine", command))
	}
}

func handleRunError(err error, fullCommand string, output string, cmd *exec.Cmd) (*ExecResult, error) {
	exitError, castOk := err.(*exec.ExitError)

	if !castOk {
		return nil, fmt.Errorf("error when executing command: \n%s \nerror: %w", fullCommand, err)
	}

	commandResult := &ExecResult{
		CommandExecuted: fullCommand,
		StdOut:          output,
		StdErr:          string(exitError.Stderr),
		ExitCode:        exitError.ExitCode(),
	}

	return commandResult, nil
}
