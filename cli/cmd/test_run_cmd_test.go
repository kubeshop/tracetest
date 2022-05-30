package cmd_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli/cmd/e2e"
	"github.com/stretchr/testify/assert"
)

func TestRunTestCmd(t *testing.T) {
	cli := e2e.NewCLI()

	output, err := cli.RunCommand("test", "run", "--config", "e2e/config.yml")
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
}
