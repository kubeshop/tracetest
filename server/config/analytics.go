package config

import "os"

type googleAnalytics struct {
	Enabled bool `yaml:",omitempty" mapstructure:"enabled"`
}

func (c *Config) Analyticsnabled() bool {
	if os.Getenv("TRACETEST_DEV") != "" {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	return c.config.GA.Enabled
}
