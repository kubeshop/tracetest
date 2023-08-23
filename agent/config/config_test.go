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
	assert.Equal(t, "https://cloud.tracetest.io", cfg.ServerURL)
}

func TestConfigWithEnvs(t *testing.T) {
	t.Cleanup(func() {
		os.Unsetenv("TRACETEST_AGENT_NAME")
		os.Unsetenv("TRACETEST_API_KEY")
		os.Unsetenv("TRACETEST_DEV_MODE")
		os.Unsetenv("TRACETEST_SERVER_URL")
	})

	os.Setenv("TRACETEST_AGENT_NAME", "my-agent-name")
	os.Setenv("TRACETEST_API_KEY", "my-agent-api-key")
	os.Setenv("TRACETEST_DEV_MODE", "true")
	os.Setenv("TRACETEST_SERVER_URL", "https://custom.server.com")

	cfg, err := config.LoadConfig()

	require.NoError(t, err)

	assert.Equal(t, "my-agent-api-key", cfg.APIKey)
	assert.Equal(t, "my-agent-name", cfg.Name)
	assert.Equal(t, "https://custom.server.com", cfg.ServerURL)
}
