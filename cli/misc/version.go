package misc

import (
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/global"
	"github.com/spf13/cobra"
)

type Version struct {
	args[any]
}

func (d Version) Run(cmd *cobra.Command, args []string) (string, error) {
	analytics.Track("Version", "cmd", map[string]string{})

	return *d.Setup.VersionText, nil
}

func NewVersion(root global.Root) Version {
	defaults := NewDefaults("test", root.Setup)
	version := Version{
		args: NewArgs[any](defaults, nil),
	}

	version.Cmd = &cobra.Command{
		GroupID: global.GroupMisc.ID,
		Use:     "version",
		Short:   "Display this CLI tool version",
		Long:    "Display this CLI tool version",
		PreRun:  defaults.PreRun,
		Run:     defaults.Run(version.Run),
		PostRun: defaults.PostRun,
	}

	root.Cmd.AddCommand(version.Cmd)
	return version
}
