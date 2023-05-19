package tracetestcli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/kubeshop/tracetest/cli-e2etest/command"
	"github.com/kubeshop/tracetest/cli-e2etest/config"
	"golang.org/x/exp/slices"

	"github.com/kubeshop/tracetest/cli/cmd"
	"github.com/stretchr/testify/require"
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

	tracetestCommand := config.GetConfigAsEnvVars().TracetestCommand
	tracetestSubCommands := strings.Split(tracetestSubCommand, " ")

	if config.GetConfigAsEnvVars().EnableCLIDebug {
		return runTracetestAsInternalCommand(t, tracetestCommand, tracetestSubCommands)
	}

	result, err := command.Exec(tracetestCommand, tracetestSubCommands...)
	require.NoError(t, err)

	return result
}

func WithCLIConfig(cliConfig string) ExecOption {
	return func(es *executionState) {
		es.cliConfigFile = cliConfig
	}
}

func runTracetestAsInternalCommand(t *testing.T, tracetestCommand string, tracetestSubCommands []string) *command.ExecResult {
	// keep backup of the real stdout
	stdoutBackup := os.Stdout
	stdoutRead, stdoutWriter, _ := os.Pipe()
	os.Stdout = stdoutWriter

	// keep backup of the real stderr
	stderrBackup := os.Stderr
	stderrRead, stderrWriter, _ := os.Pipe()
	os.Stderr = stderrWriter

	argsBackup := os.Args
	os.Args = slices.Insert(tracetestSubCommands, 0, tracetestCommand)

	exitCode := 0
	cmd.RegisterCLIExitInterceptor(func(i int) {
		exitCode = i
	})

	cmd.Execute()

	os.Args = argsBackup

	stdoutChannel := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, stdoutRead)
		stdoutChannel <- buf.String()
	}()

	// back to normal state
	stdoutWriter.Close()
	os.Stdout = stdoutBackup // restoring the real stdout

	stderrChannel := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, stderrRead)
		stderrChannel <- buf.String()
	}()

	// back to normal state
	stderrWriter.Close()
	os.Stderr = stderrBackup // restoring the real stderr

	// TODO: need to intercept exitCode

	return &command.ExecResult{
		CommandExecuted: fmt.Sprintf("%s %s", tracetestCommand, strings.Join(tracetestSubCommands, " ")),
		StdOut:          <-stdoutChannel,
		StdErr:          <-stderrChannel,
		ExitCode:        exitCode,
	}
}
