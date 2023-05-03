package cmd

import (
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/installer"
	"github.com/spf13/cobra"
)

var (
	force             = false
	runEnvironment    = installer.NoneRunEnvironmentType
	installationMode  = installer.NotChosenInstallationModeType
	kubernetesContext = ""
)

var serverInstallCmd = &cobra.Command{
	Use:    "install",
	Short:  "Install a new Tracetest server",
	Long:   "Install a new Tracetest server",
	PreRun: setupCommand(SkipConfigValidation()),
	Run: func(cmd *cobra.Command, args []string) {
		installer.Force = force
		installer.RunEnvironment = runEnvironment
		installer.InstallationMode = installationMode
		installer.KubernetesContext = kubernetesContext

		analytics.Track("Server Install", "cmd", map[string]string{})
		installer.Start()
	},
	PostRun: teardownCommand,
}

func init() {
	serverInstallCmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing files")
	serverInstallCmd.Flags().StringVar(&kubernetesContext, "kubernetes-context", "", "Kubernetes context used to install Tracetest. It will be only used if 'run-environment' is set as 'kubernetes'.")

	// these commands will not have shorthand parameters to avoid colision with existing ones in other commands
	serverInstallCmd.Flags().Var(&installationMode, "mode", "Indicate the type of demo environment to be installed with Tracetest. It can be 'with-demo' or 'just-tracetest'.")
	serverInstallCmd.Flags().Var(&runEnvironment, "run-environment", "Type of environment were Tracetest will be installed. It can be 'docker' or 'kubernetes'.")

	serverCmd.AddCommand(serverInstallCmd)
}
