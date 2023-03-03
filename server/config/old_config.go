package config

type (
	oldConfig struct {
		Server    serverConfig    `yaml:",omitempty" mapstructure:"server"`
		GA        googleAnalytics `yaml:"googleAnalytics,omitempty" mapstructure:"googleAnalytics"`
		Telemetry telemetry       `yaml:",omitempty" mapstructure:"telemetry"`
		Demo      demo            `yaml:",omitempty" mapstructure:"demo"`
	}
)
