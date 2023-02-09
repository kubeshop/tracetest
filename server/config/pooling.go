package config

import "time"

func (c *Config) PoolingConfig() PoolingConfig {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.config.PoolingConfig
}

func (c *Config) SetPoolingConfig(pc PoolingConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.config.PoolingConfig = pc
}

func (c *Config) PoolingRetryDelay() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	delay, err := time.ParseDuration(c.config.PoolingConfig.RetryDelay)
	if err != nil {
		return 1 * time.Second
	}

	return delay
}

func (c *Config) MaxWaitTimeForTraceDuration() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	maxWaitTimeForTrace, err := time.ParseDuration(c.config.PoolingConfig.MaxWaitTimeForTrace)
	if err != nil {
		// use a default value
		maxWaitTimeForTrace = 30 * time.Second
	}
	return maxWaitTimeForTrace
}
