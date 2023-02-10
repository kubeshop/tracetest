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

	t.Run("Flags", func(t *testing.T) {
		t.Parallel()

		flags := []string{
			"--postgres.dbname", "other_dbname",
			"--postgres.host", "localhost",
			"--postgres.user", "user",
			"--postgres.password", "passwd",
			"--postgres.port", "1234",
			"--postgres.params", "custom=params",
			"--server.httpPort", "4321",
			"--server.pathPrefix", "/prefix",
			"--experimentalFeatures", "a",
			"--experimentalFeatures", "b",
			"--internalTelemetry.enabled", "true",
			"--internalTelemetry.otelCollectorEndpoint", "otel-collector.tracetest",
		}

		cfg := configWithFlags(t, flags)

		assert.Equal(t, "host=localhost user=user password=passwd port=1234 dbname=other_dbname custom=params", cfg.PostgresConnString())

		assert.Equal(t, 4321, cfg.ServerPort())
		assert.Equal(t, "/prefix", cfg.ServerPathPrefix())

		assert.DeepEqual(t, []string{"a", "b"}, cfg.ExperimentalFeatures())

		assert.Equal(t, true, cfg.InternalTelemetryEnabled())
		assert.Equal(t, "otel-collector.tracetest", cfg.InternalTelemetryOtelCollectorAddress())
	})
}
