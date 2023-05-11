package misc

import (
	"github.com/kubeshop/tracetest/cli/global"
	"github.com/spf13/cobra"
)

type Test struct {
	args[any]
}

func NewTest(root global.Root) Test {
	defaults := NewDefaults("test", root.Setup)

	test := Test{
		args: NewArgs[any](defaults, nil),
	}

	test.Cmd = &cobra.Command{
		GroupID: global.GroupTests.ID,
		Use:     "test",
		Short:   "Manage your tracetest tests",
		Long:    "Manage your tracetest tests",
		PreRun:  defaults.PreRun,
		Run: defaults.Run(func(cmd *cobra.Command, args []string) (string, error) {
			cmd.Help()

			return "", nil
		}),
		PostRun: defaults.PostRun,
	}

	root.Cmd.AddCommand(test.Cmd)

	return test
}
