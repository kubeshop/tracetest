package e2e_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli/e2e"
	"github.com/stretchr/testify/assert"
)

func TestTestListCmd(t *testing.T) {
	cli := e2e.NewCLI()

	output, err := cli.RunCommand("test", "list")
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
}
