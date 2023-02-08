package config

import (
	"os"
	"sync"
)

type Config struct {
	config *config
	mu     sync.Mutex
}

func New() *Config {
	return &Config{
		config: &config{},
	}
}

func (c *Config) AnalyticsEnabled() bool {
	if os.Getenv("TRACETEST_DEV") != "" {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	return c.config.GA.Enabled
}
