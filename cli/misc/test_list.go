package misc

import (
	"context"

	misc_actions "github.com/kubeshop/tracetest/cli/misc/actions"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type TestList struct {
	args[any]
}

func (tl TestList) Run(cmd *cobra.Command, args []string) (string, error) {
	ctx := context.Background()
	tl.Setup.Logger.Debug("Retrieving list of tests", zap.String("endpoint", tl.Setup.Config.Endpoint))
	listTestsAction := misc_actions.NewListTestsAction(*tl.Setup.Config, tl.Setup.Logger, tl.Setup.Client)

	actionArgs := misc_actions.ListTestConfig{}
	err := listTestsAction.Run(ctx, actionArgs)
	return "", err
}

func NewTestList(root Test) TestList {
	defaults := NewDefaults("test export", root.Setup.Setup)

	testList := TestList{
		args: NewArgs[any](defaults, nil),
	}

	testList.Cmd = &cobra.Command{
		Use:     "list",
		Short:   "List all tests",
		Long:    "List all tests",
		PreRun:  defaults.PreRun,
		Run:     defaults.Run(testList.Run),
		PostRun: defaults.PostRun,
	}
	root.Cmd.AddCommand(testList.Cmd)

	return testList
}
