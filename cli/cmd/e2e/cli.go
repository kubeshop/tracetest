package e2e

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/kubeshop/tracetest/cli/cmd"
)

type CLI struct {
	options []string
}

type Command struct {
	command string
	args    []string
	options []string
}

func NewCLI(options ...string) CLI {
	return CLI{options}
}

func (c CLI) NewCommand(command string, args ...string) Command {
	return Command{
		command: command,
		args:    args,
		options: c.options,
	}
}

func (c Command) Run(stdin ...string) (string, error) {
	executableArgs := make([]string, 0)
	executableArgs = append(executableArgs, "executable")
	executableArgs = append(executableArgs, c.command)
	executableArgs = append(executableArgs, c.options...)
	executableArgs = append(executableArgs, c.args...)

	os.Args = executableArgs

	oldStdout := os.Stdout
	stdoutReader, newStdout, _ := os.Pipe()
	os.Stdout = newStdout

	oldStdin := os.Stdin
	newStdin, err := ioutil.TempFile("", "stdin")
	if err != nil {
		return "", err
	}
	os.Stdin = newStdin

	err = c.WriteToStdin(stdin, newStdin)
	if err != nil {
		return "", fmt.Errorf("could not write to stdin: %w", err)
	}

	cmd.Execute()

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, stdoutReader)
		outC <- buf.String()
	}()

	newStdout.Close()
	newStdin.Close()
	os.Stdout = oldStdout
	os.Stdin = oldStdin
	out := <-outC

	return out, nil
}

func (c Command) WriteToStdin(stdinArgs []string, stdin *os.File) error {
	for _, input := range stdinArgs {
		line := fmt.Sprintf("%s\n", input)
		stdin.Write([]byte(line))
	}

	_, err := stdin.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("could not seek beginning of file")
	}

	return nil
}
