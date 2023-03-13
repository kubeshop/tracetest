package provisioning

import (
	"github.com/kubeshop/tracetest/server/config/configresource"
)

type config struct {
	Type             string `mapstructure:"type"`
	AnalyticsEnabled bool   `mapstructure:"analyticsEnabled"`
}

func (cfg config) model() configresource.Config {
	m := configresource.Config{
		AnalyticsEnabled: cfg.AnalyticsEnabled,
	}

	return m
}
