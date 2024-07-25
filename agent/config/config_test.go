package config_test

import (
	"os"
	"testing"

	"github.com/kubeshop/tracetest/agent/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigDefaults(t *testing.T) {
	cfg, err := config.LoadConfig()

	require.NoError(t, err)

	hostname, _ := os.Hostname()

	assert.Equal(t, "", cfg.APIKey)
	assert.Equal(t, hostname, cfg.Name)
	assert.Empty(t, cfg.EnvironmentID)
	assert.Equal(t, "https://app.tracetest.io", cfg.ServerURL)
	assert.Equal(t, 4317, cfg.OTLPServer.GRPCPort)
	assert.Equal(t, 4318, cfg.OTLPServer.HTTPPort)
}

func TestConfigWithEnvs(t *testing.T) {
	t.Cleanup(func() {
		os.Unsetenv("TRACETEST_AGENT_NAME")
		os.Unsetenv("TRACETEST_API_KEY")
		os.Unsetenv("TRACETEST_DEV")
		os.Unsetenv("TRACETEST_SERVER_URL")
		os.Unsetenv("TRACETEST_OTLP_SERVER_GRPC_PORT")
		os.Unsetenv("TRACETEST_OTLP_SERVER_HTTP_PORT")
		os.Unsetenv("TRACETEST_ENVIRONMENT_ID")
	})

	os.Setenv("TRACETEST_AGENT_NAME", "my-agent-name")
	os.Setenv("TRACETEST_API_KEY", "my-agent-api-key")
	os.Setenv("TRACETEST_ENVIRONMENT_ID", "123456")
	os.Setenv("TRACETEST_DEV", "true")
	os.Setenv("TRACETEST_SERVER_URL", "https://custom.server.com")
	os.Setenv("TRACETEST_OTLP_SERVER_GRPC_PORT", "1234")
	os.Setenv("TRACETEST_OTLP_SERVER_HTTP_PORT", "1235")

	cfg, err := config.LoadConfig()

	require.NoError(t, err)

	assert.Equal(t, "my-agent-api-key", cfg.APIKey)
	assert.Equal(t, "my-agent-name", cfg.Name)
	assert.Equal(t, "123456", cfg.EnvironmentID)
	assert.Equal(t, "https://custom.server.com", cfg.ServerURL)
	assert.Equal(t, 1234, cfg.OTLPServer.GRPCPort)
	assert.Equal(t, 1235, cfg.OTLPServer.HTTPPort)
}
