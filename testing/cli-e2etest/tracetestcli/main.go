package tracetestcli

import (
	"strings"

	"github.com/kubeshop/tracetest/cli-e2etest/command"
)

const (
	tracetestCommand = "tracetest"
)

func Exec(tracetestSubCommand string) (string, int, error) {
	tracetestSubCommands := strings.Split(tracetestSubCommand, " ")

	return command.Exec(tracetestCommand, tracetestSubCommands...)
}
