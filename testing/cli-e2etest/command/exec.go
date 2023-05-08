package command

import (
	"fmt"
	"os/exec"
	"strings"
)

func Exec(command string, args ...string) (string, int, error) {
	checkIfCommandExists(command)

	fullCommand := fmt.Sprintf("%s %s", command, strings.Join(args, " "))

	cmd := exec.Command(command, args...)

	out, err := cmd.Output()
	output := string(out)

	if err != nil {
		return handleRunError(err, fullCommand, output, cmd)
	}

	exitCode := cmd.ProcessState.ExitCode()
	return output, exitCode, nil
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

func handleRunError(err error, fullCommand string, output string, cmd *exec.Cmd) (string, int, error) {
	exitError, castOk := err.(*exec.ExitError)

	if !castOk {
		return "", -1, fmt.Errorf("error when executing command '%s', stdin '%s', error %w", fullCommand, output, err)
	}

	errorOutput := string(exitError.Stderr)

	exitCode := exitError.ExitCode()
	return "", exitCode, fmt.Errorf("error when executing tracetest command '%s', exit code %d, stdin '%s', stderr '%s'", fullCommand, exitCode, output, errorOutput)
}
