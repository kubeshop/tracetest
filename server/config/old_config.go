package config

type (
	oldConfig struct {
		Server    serverConfig `yaml:",omitempty" mapstructure:"server"`
		Telemetry telemetry    `yaml:",omitempty" mapstructure:"telemetry"`
		Demo      demo         `yaml:",omitempty" mapstructure:"demo"`
	}
)
