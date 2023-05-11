package flows_parameters

import "github.com/kubeshop/tracetest/cli/flows/actions/installer"

type Installer struct {
	Force             bool
	RunEnvironment    installer.RunEnvironmentType
	InstallationMode  installer.InstallationModeType
	KubernetesContext string
}

func NewInstaller() *Installer {
	return &Installer{}
}
