package config_test

import (
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/logger"
	"gotest.tools/v3/assert"
)

func TestPoolingConfig(t *testing.T) {
	t.Run("DefaultValues", func(t *testing.T) {
		cfg, _ := config.New(nil, logger.Default())

		assert.Equal(t, 30*time.Second, cfg.PoolingMaxWaitTimeForTraceDuration())
		assert.Equal(t, 1*time.Second, cfg.PoolingRetryDelay())
	})

	t.Run("InvalidFormat", func(t *testing.T) {
		flags := []string{
			"--poolingConfig.maxWaitTimeForTrace", "not-valid",
			"--poolingConfig.retryDelay", "not-valid",
		}

		_, err := configWithFlagsE(t, flags)

		assert.ErrorContains(t, err, "invalid config value for 'poolingConfig.maxWaitTimeForTrace': invalid duration format 'not-valid'")
		assert.ErrorContains(t, err, "invalid config value for 'poolingConfig.retryDelay': invalid duration format 'not-valid'")
	})

	t.Run("Flags", func(t *testing.T) {
		flags := []string{
			"--poolingConfig.maxWaitTimeForTrace", "2m",
			"--poolingConfig.retryDelay", "10s",
		}

		cfg := configWithFlags(t, flags)

		assert.Equal(t, 2*time.Minute, cfg.PoolingMaxWaitTimeForTraceDuration())
		assert.Equal(t, 10*time.Second, cfg.PoolingRetryDelay())
	})

	t.Run("EnvVars", func(t *testing.T) {
		env := map[string]string{
			"TRACETEST_POOLINGCONFIG_MAXWAITTIMEFORTRACE": "2m",
			"TRACETEST_POOLINGCONFIG_RETRYDELAY":          "10s",
		}

		cfg := configWithEnv(t, env)

		assert.Equal(t, 2*time.Minute, cfg.PoolingMaxWaitTimeForTraceDuration())
		assert.Equal(t, 10*time.Second, cfg.PoolingRetryDelay())
	})
}
