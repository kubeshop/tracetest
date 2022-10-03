package cmd_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli/cmd/e2e"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigureWithAnalytics(t *testing.T) {
	cli := e2e.NewCLI()

	configureCommand := cli.NewCommand("configure", "--endpoint", "https://localhost:1234", "--analytics")

	t.Cleanup(func() {
		deleteFile("./config.yml")
	})

	_, err := configureCommand.Run()
	assert.NoError(t, err)

	tracetestConfig, err := config.LoadConfig("./config.yml")
	require.NoError(t, err)

	assert.Equal(t, "https", tracetestConfig.Scheme)
	assert.Equal(t, "localhost:1234", tracetestConfig.Endpoint)
	assert.True(t, tracetestConfig.AnalyticsEnabled)
	assert.Nil(t, tracetestConfig.ServerPath)
}

func TestConfigureWithoutAnalytics(t *testing.T) {
	cli := e2e.NewCLI()

	configureCommand := cli.NewCommand("configure", "--endpoint", "https://localhost:1234", "--analytics=false")

	t.Cleanup(func() {
		deleteFile("./config.yml")
	})

	_, err := configureCommand.Run()
	assert.NoError(t, err)

	tracetestConfig, err := config.LoadConfig("./config.yml")
	require.NoError(t, err)

	assert.Equal(t, "https", tracetestConfig.Scheme)
	assert.Equal(t, "localhost:1234", tracetestConfig.Endpoint)
	assert.False(t, tracetestConfig.AnalyticsEnabled)
	assert.Nil(t, tracetestConfig.ServerPath)
}
