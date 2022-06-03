package cmd_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli/cmd/e2e"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigureCmd(t *testing.T) {
	cli := e2e.NewCLI()

	configureCommand := cli.NewCommand("configure")

	t.Cleanup(func() {
		deleteFile("./config.yml")
	})

	output, err := configureCommand.Run("https://localhost:1234")
	assert.NoError(t, err)
	assert.NotEmpty(t, output)

	tracetestConfig, err := config.LoadConfig("./config.yml")
	require.NoError(t, err)

	assert.Equal(t, "https", tracetestConfig.Scheme)
	assert.Equal(t, "localhost:1234", tracetestConfig.Endpoint)
	assert.Nil(t, tracetestConfig.ServerPath)

}
