package cmd

import (
	"github.com/kubeshop/tracetest/cli/installer"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

var installerParams = &installerParameters{
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

var (
	AllowedRunEnvironments = []installer.RunEnvironmentType{
		installer.DockerRunEnvironmentType,
		installer.KubernetesRunEnvironmentType,
		installer.NoneRunEnvironmentType,
	}
	AllowedInstallationMode = []installer.InstallationModeType{
		installer.WithDemoInstallationModeType,
		installer.WithoutDemoInstallationModeType,
		installer.NotChosenInstallationModeType,
	}
)

type installerParameters struct {
	Force             bool
	RunEnvironment    installer.RunEnvironmentType
	InstallationMode  installer.InstallationModeType
	KubernetesContext string
}

func (p installerParameters) Validate(cmd *cobra.Command, args []string) []ParamError {
	errors := make([]ParamError, 0)

	if cmd.Flags().Lookup("run-environment").Changed && slices.Contains(AllowedRunEnvironments, p.RunEnvironment) {
		errors = append(errors, ParamError{
			Parameter: "run-environment",
			Message:   "run-environment must be one of 'none', 'docker' or 'kubernetes'",
		})
	}

	if cmd.Flags().Lookup("mode").Changed && slices.Contains(AllowedInstallationMode, p.InstallationMode) {
		errors = append(errors, ParamError{
			Parameter: "mode",
			Message:   "mode must be one of 'not-chosen', 'with-demo' or 'just-tracetest'",
		})
	}

	if cmd.Flags().Lookup("kubernetes-context").Changed && p.KubernetesContext == "" {
		errors = append(errors, ParamError{
			Parameter: "kubernetes-context",
			Message:   "kubernetes-context cannot be empty",
		})
	}

	return errors
}
