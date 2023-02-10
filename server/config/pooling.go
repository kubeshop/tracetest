package config

import (
	"time"

	"github.com/spf13/viper"
)

func init() {
	defaultSetters = append(defaultSetters, poolingDefaultSetter)
}

func poolingDefaultSetter(vp *viper.Viper) {
	vp.SetDefault("poolingConfig.maxWaitTimeForTrace", "30s")
	vp.SetDefault("poolingConfig.retryDelay", "1s")
}

func (c *Config) SetPoolingConfig(pc PoolingConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.vp.Set("poolingConfig.maxWaitTimeForTrace", pc.MaxWaitTimeForTrace)
	c.vp.Set("poolingConfig.retryDelay", pc.RetryDelay)

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
