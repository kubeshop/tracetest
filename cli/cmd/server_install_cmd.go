package cmd

import (
	"github.com/kubeshop/tracetest/cli/analytics"
	"github.com/kubeshop/tracetest/cli/installer"
	"github.com/kubeshop/tracetest/cli/parameters"
	"github.com/spf13/cobra"
)

var installerParams = &parameters.InstallerParams{
	Force:             false,
	RunEnvironment:    installer.NoneRunEnvironmentType,
	InstallationMode:  installer.NotChosenInstallationModeType,
	KubernetesContext: "",
}

var serverInstallCmd = &cobra.Command{
	Use:    "install",
	Short:  "Install a new Tracetest server",
	Long:   "Install a new Tracetest server",
	PreRun: setupCommand(SkipConfigValidation(), SkipVersionMismatchCheck()),
	Run: func(_ *cobra.Command, _ []string) {
		installer.Force = installerParams.Force
		installer.RunEnvironment = installerParams.RunEnvironment
		installer.InstallationMode = installerParams.InstallationMode
		installer.KubernetesContext = installerParams.KubernetesContext

		analytics.Track("Server Install", "cmd", map[string]string{})
		installer.Start()
	},
	PostRun: teardownCommand,
}

func init() {
	serverInstallCmd.Flags().BoolVarP(&installerParams.Force, "force", "f", false, "Overwrite existing files")
	serverInstallCmd.Flags().StringVar(&installerParams.KubernetesContext, "kubernetes-context", "", "Kubernetes context used to install Tracetest. It will be only used if 'run-environment' is set as 'kubernetes'.")

	// these commands will not have shorthand parameters to avoid colision with existing ones in other commands
	serverInstallCmd.Flags().Var(&installerParams.InstallationMode, "mode", "Indicate the type of demo environment to be installed with Tracetest. It can be 'with-demo' or 'just-tracetest'.")
	serverInstallCmd.Flags().Var(&installerParams.RunEnvironment, "run-environment", "Type of environment were Tracetest will be installed. It can be 'docker' or 'kubernetes'.")

	serverCmd.AddCommand(serverInstallCmd)
}
