package misc

import (
	"fmt"

	"github.com/kubeshop/tracetest/cli/global"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/spf13/cobra"
)

type Dashboard struct {
	args[any]
}

func (d Dashboard) Run(cmd *cobra.Command, args []string) (string, error) {
	if d.Setup.Config.IsEmpty() {
		return "", fmt.Errorf("missing Tracetest endpoint configuration")
	}

	err := utils.OpenBrowser(d.Setup.Config.URL())
	if err != nil {
		return "", fmt.Errorf("failed to open the dashboard url: %s", d.Setup.Config.URL())
	}

	return fmt.Sprintf("Opening \"%s\" in the default browser", d.Setup.Config.URL()), nil
}

func NewDashboard(root global.Root) Dashboard {
	defaults := NewDefaults("dashboard", root.Setup)

	dashboard := Dashboard{
		args: NewArgs[any](defaults, nil),
	}

	dashboard.Cmd = &cobra.Command{
		GroupID: global.GroupMisc.ID,
		Use:     "dashboard",
		Short:   "Opens the Tracetest Dashboard URL",
		Long:    "Opens the Tracetest Dashboard URL",
		PreRun:  defaults.PreRun,
		Run:     defaults.Run(dashboard.Run),
		PostRun: defaults.PostRun,
	}
	root.Cmd.AddCommand(dashboard.Cmd)

	return dashboard
}
