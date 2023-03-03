package config_test

import (
	"os"
	"testing"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/stretchr/testify/assert"
)

func TestAnalyticsEnabledConfig(t *testing.T) {
	cleanEnv := func() {
		// make sure this env is empty
		os.Setenv("TRACETEST_DEV", "")
	}
	t.Run("DefaultValues", func(t *testing.T) {
		cleanEnv()

		cfg, _ := config.New()

		assert.False(t, cfg.AnalyticsEnabled())
	})

	t.Run("File", func(t *testing.T) {
		cleanEnv()

		cfg := configFromFile(t, "./testdata/analytics.yaml")

		assert.True(t, cfg.AnalyticsEnabled())
	})

	t.Run("EnvOverride", func(t *testing.T) {

		cfg := configFromFile(t, "./testdata/analytics.yaml")

		os.Setenv("TRACETEST_DEV", "yes")
		assert.False(t, cfg.AnalyticsEnabled())
	})
}
