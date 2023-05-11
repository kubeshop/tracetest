package flows

import (
	"github.com/kubeshop/tracetest/cli/global"
	"github.com/spf13/cobra"
)

type Server struct {
	args[any]
}

func NewServer(root global.Root) Server {
	defaults := NewDefaults("server", root.Setup)
	server := Server{
		args: NewArgs[any](defaults, nil),
	}

	server.Cmd = &cobra.Command{
		GroupID: global.GroupConfig.ID,
		Use:     "server",
		Short:   "Manage your tracetest server",
		Long:    "Manage your tracetest server",
		PreRun:  defaults.PreRun,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		PostRun: defaults.PostRun,
	}

	root.Cmd.AddCommand(server.Cmd)
	return server
}

var server Server

func GetServer() Server {
	if server.Cmd == nil {
		panic("server not initialized")
	}

	return server
}
