package get

import (
	"fmt"

	"github.com/kubeshop/tracetest/cli-e2etest/tracetestcli"
)

func Exec(tracetestSubCommand string) (string, int, error) {
	subCommand := fmt.Sprintf("get %s", tracetestSubCommand)

	return tracetestcli.Exec(subCommand)
}
