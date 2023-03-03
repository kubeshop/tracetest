package config

import "os"

func (c *Config) AnalyticsEnabled() bool {
	if os.Getenv("TRACETEST_DEV") != "" {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	return c.config.GA.Enabled
}
