package global_types

import (
	"github.com/kubeshop/tracetest/cli/global"
	"github.com/spf13/cobra"
)

type Command interface {
	Cmd() *cobra.Command
}

type CommandFactory func(root global.Root) any
