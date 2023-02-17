package config

import (
	"time"
)

var poolingOptions = options{
	{
		key:          "poolingConfig.maxWaitTimeForTrace",
		defaultValue: "30s",
		description:  "pooling config: max wait time for trace",
		validate:     validateDuration("poolingConfig.maxWaitTimeForTrace"),
	},
	{
		key:          "poolingConfig.retryDelay",
		defaultValue: "1s",
		description:  "pooling config: interval between poll retries",
		validate:     validateDuration("poolingConfig.retryDelay"),
	},
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
