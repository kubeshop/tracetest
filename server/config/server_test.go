package config_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/config"
	"gotest.tools/v3/assert"
)

func TestServerConfig(t *testing.T) {
	t.Run("DefaultValues", func(t *testing.T) {
		cfg, _ := config.New()

		assert.Equal(t, "postgres://postgres:postgres@postgres:5432/tracetest?sslmode=disable", cfg.PostgresConnString())

		assert.Equal(t, 11633, cfg.ServerPort())
		assert.Equal(t, "", cfg.ServerPathPrefix())

		assert.DeepEqual(t, []string{}, cfg.ExperimentalFeatures())

		assert.Equal(t, false, cfg.InternalTelemetryEnabled())
		assert.Equal(t, "", cfg.InternalTelemetryOtelCollectorAddress())
	})

	t.Run("Flags", func(t *testing.T) {
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
			"--testPipelines.triggerExecute.enabled", "false",
			"--testPipelines.traceFetch.enabled", "false",
		}

		cfg := configWithFlags(t, flags)

		assert.Equal(t, "postgres://user:passwd@localhost:1234/other_dbname?custom=params", cfg.PostgresConnString())

		assert.Equal(t, 4321, cfg.ServerPort())
		assert.Equal(t, "/prefix", cfg.ServerPathPrefix())

		assert.DeepEqual(t, []string{"a", "b"}, cfg.ExperimentalFeatures())

		assert.Equal(t, true, cfg.InternalTelemetryEnabled())
		assert.Equal(t, "otel-collector.tracetest", cfg.InternalTelemetryOtelCollectorAddress())

		assert.Equal(t, false, cfg.TestPipelineTriggerExecutionEnabled())
		assert.Equal(t, false, cfg.TestPipelineTraceFetchingEnabled())
	})

	t.Run("EnvVars", func(t *testing.T) {
		env := map[string]string{
			"TRACETEST_POSTGRES_DBNAME":                         "other_dbname",
			"TRACETEST_POSTGRES_HOST":                           "localhost",
			"TRACETEST_POSTGRES_USER":                           "user",
			"TRACETEST_POSTGRES_PASSWORD":                       "passwd",
			"TRACETEST_POSTGRES_PORT":                           "1234",
			"TRACETEST_POSTGRES_PARAMS":                         "custom=params",
			"TRACETEST_SERVER_HTTPPORT":                         "4321",
			"TRACETEST_SERVER_PATHPREFIX":                       "/prefix",
			"TRACETEST_EXPERIMENTALFEATURES":                    "a b",
			"TRACETEST_INTERNALTELEMETRY_ENABLED":               "true",
			"TRACETEST_INTERNALTELEMETRY_OTELCOLLECTORENDPOINT": "otel-collector.tracetest",
			"TRACETEST_TESTPIPELINES_TRIGGEREXECUTE_ENABLED":    "false",
			"TRACETEST_TESTPIPELINES_TRACEFETCH_ENABLED":        "false",
		}

		cfg := configWithEnv(t, env)

		assert.Equal(t, "postgres://user:passwd@localhost:1234/other_dbname?custom=params", cfg.PostgresConnString())

		assert.Equal(t, 4321, cfg.ServerPort())
		assert.Equal(t, "/prefix", cfg.ServerPathPrefix())

		assert.DeepEqual(t, []string{"a", "b"}, cfg.ExperimentalFeatures())

		assert.Equal(t, true, cfg.InternalTelemetryEnabled())
		assert.Equal(t, "otel-collector.tracetest", cfg.InternalTelemetryOtelCollectorAddress())

		assert.Equal(t, false, cfg.TestPipelineTriggerExecutionEnabled())
		assert.Equal(t, false, cfg.TestPipelineTraceFetchingEnabled())
	})
}
