package config_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/config"
	"gotest.tools/v3/assert"
)

func TestServerConfig(t *testing.T) {
	t.Run("DefaultValues", func(t *testing.T) {
		t.Parallel()
		cfg, _ := config.New(nil)

		assert.Equal(t, "host=postgres user=postgres password=postgres port=5432 dbname=tracetest", cfg.PostgresConnString())

		assert.Equal(t, 11633, cfg.ServerPort())
		assert.Equal(t, "/", cfg.ServerPathPrefix())

		assert.DeepEqual(t, []string{}, cfg.ExperimentalFeatures())

		assert.Equal(t, false, cfg.InternalTelemetryEnabled())
		assert.Equal(t, "", cfg.InternalTelemetryOtelCollectorAddress())
	})
}
