package e2e

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/kubeshop/tracetest/cli/cmd"
)

type CLI struct {
	options []string
}

func NewCLI(options ...string) CLI {
	return CLI{options}
}

func (c CLI) RunCommand(command string, args ...string) (string, error) {
	executableArgs := make([]string, 0)
	executableArgs = append(executableArgs, "executable")
	executableArgs = append(executableArgs, command)
	executableArgs = append(executableArgs, c.options...)
	executableArgs = append(executableArgs, args...)

	os.Args = executableArgs
	fmt.Println(executableArgs)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	cmd.Execute()

	os.Stdout = old
	out := <-outC
	w.Close()

	return out, nil
}
