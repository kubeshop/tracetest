package config

import (
	"context"
	"os"
)

func (cfg *Config) AnalyticsEnabled() bool {
	if os.Getenv("TRACETEST_DEV") != "" {
		return false
	}

	cr := cfg.resources.config.Current(context.TODO())

	return cr.AnalyticsEnabled
}
