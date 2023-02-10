package config

import (
	"time"
)

var poolingOptions = options{
	{"poolingConfig.maxWaitTimeForTrace", "30s", "pooling config: max wait time for trace"},
	{"poolingConfig.retryDelay", "1s", "pooling config: interval between poll retries"},
}

func init() {
	configOptions = append(configOptions, poolingOptions...)
}

func (c *Config) SetPoolingConfig(pc PoolingConfig) {
	c.Set("poolingConfig.maxWaitTimeForTrace", pc.MaxWaitTimeForTrace)
	c.Set("poolingConfig.retryDelay", pc.RetryDelay)

}

func (c *Config) PoolingMaxWaitTimeForTraceDuration() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	maxWaitTimeForTrace, err := time.ParseDuration(c.vp.GetString("poolingConfig.maxWaitTimeForTrace"))
	if err != nil {
		// use a default value
		maxWaitTimeForTrace = 30 * time.Second
	}
	return maxWaitTimeForTrace
}

func (c *Config) PoolingRetryDelay() time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	delay, err := time.ParseDuration(c.vp.GetString("poolingConfig.retryDelay"))
	if err != nil {
		return 1 * time.Second
	}

	return delay
}
