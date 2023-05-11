package flows

import (
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/flows/actions/installer"
	flows_parameters "github.com/kubeshop/tracetest/cli/flows/parameters"
	"github.com/spf13/cobra"
)

type ServerInstall struct {
	args[*flows_parameters.Installer]
}

func NewServerInstall(root Server) ServerInstall {
	parameters := flows_parameters.NewInstaller()
	defaults := NewDefaults("Server Install", root.Setup)

	serverInstall := ServerInstall{
		args: NewArgs(defaults, parameters),
	}

	serverInstall.Cmd = &cobra.Command{
		Use:    "install",
		Short:  "Install a new Tracetest server",
		Long:   "Install a new Tracetest server",
		PreRun: defaults.PreRun,
		Run: func(cmd *cobra.Command, args []string) {

			installer.Force = serverInstall.Parameters.Force
			installer.RunEnvironment = serverInstall.Parameters.RunEnvironment
			installer.InstallationMode = serverInstall.Parameters.InstallationMode
			installer.KubernetesContext = serverInstall.Parameters.KubernetesContext

			analytics.Track("Server Install", "cmd", map[string]string{})
			installer.Start()
		},
		PostRun: defaults.PostRun,
	}

	serverInstall.Cmd.Flags().BoolVarP(&serverInstall.Parameters.Force, "force", "f", false, "Overwrite existing files")
	serverInstall.Cmd.Flags().StringVar(&serverInstall.Parameters.KubernetesContext, "kubernetes-context", "", "Kubernetes context used to install Tracetest. It will be only used if 'run-environment' is set as 'kubernetes'.")

	// these commands will not have shorthand parameters to avoid colision with existing ones in other commands
	serverInstall.Cmd.Flags().Var(&serverInstall.Parameters.InstallationMode, "mode", "Indicate the type of demo environment to be installed with Tracetest. It can be 'with-demo' or 'just-tracetest'.")
	serverInstall.Cmd.Flags().Var(&serverInstall.Parameters.RunEnvironment, "run-environment", "Type of environment were Tracetest will be installed. It can be 'docker' or 'kubernetes'.")

	root.Cmd.AddCommand(serverInstall.Cmd)

	return serverInstall
}
