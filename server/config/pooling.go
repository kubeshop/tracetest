package config

import (
	"time"
)

var poolingOptions = options{
	{"poolingConfig.maxWaitTimeForTrace", "30s", "pooling config: max wait time for trace", validateDuration("poolingConfig.maxWaitTimeForTrace")},
	{"poolingConfig.retryDelay", "1s", "pooling config: interval between poll retries", validateDuration("poolingConfig.retryDelay")},
}

func init() {
	configOptions = append(configOptions, poolingOptions...)
}

func (c *Config) PoolingMaxWaitTimeForTraceDuration() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.vp.GetDuration("poolingConfig.maxWaitTimeForTrace")
}

func (c *Config) PoolingRetryDelay() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.vp.GetDuration("poolingConfig.retryDelay")
}
