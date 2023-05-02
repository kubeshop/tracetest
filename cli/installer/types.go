package installer

import (
	"errors"

	"github.com/spf13/pflag"
)

type RunEnvironmentType string

var _ pflag.Value = (*RunEnvironmentType)(nil)

const (
	DockerRunEnvironmentType     RunEnvironmentType = "docker"
	KubernetesRunEnvironmentType RunEnvironmentType = "kubernetes"
	NoneRunEnvironmentType       RunEnvironmentType = "none" // stands for "no option chosen"
)

func (e *RunEnvironmentType) String() string {
	return string(*e)
}

func (e *RunEnvironmentType) Set(v string) error {
	switch v {
	case "docker", "kubernetes":
		*e = RunEnvironmentType(v)
		return nil
	default:
		return errors.New(`must be "docker" or "kubernetes"`)
	}
}

func (e *RunEnvironmentType) Type() string {
	return "(docker|kubernetes)"
}

type InstallationModeType string

var _ pflag.Value = (*InstallationModeType)(nil)

const (
	WithDemoInstallationModeType    InstallationModeType = "with-demo"
	WithoutDemoInstallationModeType InstallationModeType = "just-tracetest"
	NotChosenInstallationModeType   InstallationModeType = "none" // stands for "no option chosen"
)

func (e *InstallationModeType) String() string {
	return string(*e)
}

func (e *InstallationModeType) Set(v string) error {
	switch v {
	case "with-demo", "just-tracetest":
		*e = InstallationModeType(v)
		return nil
	default:
		return errors.New(`must be "with-demo" or "just-tracetest"`)
	}
}

func (e *InstallationModeType) Type() string {
	return "(with-demo|just-tracetest)"
}
