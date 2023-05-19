package global_decorators

import (
	global_types "github.com/kubeshop/tracetest/cli/global/types"
	"github.com/spf13/cobra"
)

type Decorator func(command global_types.Command) global_types.Command
type CobraFn func(cmd *cobra.Command, args []string)

func Decorate[T global_types.Command](command global_types.Command, decorators ...Decorator) T {
	for _, decorator := range decorators {
		command = decorator(command)
	}

	return command.(T)
}

func NoopRun(cmd *cobra.Command, args []string) {}
