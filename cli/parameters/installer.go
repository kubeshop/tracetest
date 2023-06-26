package parameters

import (
	"github.com/kubeshop/tracetest/cli/installer"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
)

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

type InstallerParams struct {
	Force             bool
	RunEnvironment    installer.RunEnvironmentType
	InstallationMode  installer.InstallationModeType
	KubernetesContext string
}

func (p *InstallerParams) Validate(cmd *cobra.Command, args []string) []paramError {
	errors := make([]paramError, 0)

	if cmd.Flags().Lookup("run-environment").Changed && slices.Contains(AllowedRunEnvironments, p.RunEnvironment) {
		errors = append(errors, paramError{
			Parameter: "run-environment",
			Message:   "run-environment must be one of 'none', 'docker' or 'kubernetes'",
		})
	}

	if cmd.Flags().Lookup("mode").Changed && slices.Contains(AllowedInstallationMode, p.InstallationMode) {
		errors = append(errors, paramError{
			Parameter: "mode",
			Message:   "mode must be one of 'not-chosen', 'with-demo' or 'just-tracetest'",
		})
	}

	if cmd.Flags().Lookup("kubernetes-context").Changed && p.KubernetesContext == "" {
		errors = append(errors, paramError{
			Parameter: "kubernetes-context",
			Message:   "kubernetes-context cannot be empty",
		})
	}

	return errors
}
