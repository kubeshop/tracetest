package cmd_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli/cmd/e2e"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/stretchr/testify/assert"
)

func TestTestListCmd(t *testing.T) {
	cli := e2e.NewCLI()

	testListCommand := cli.NewCommand("test", "list", "--config", "e2e/config.yml")
	output, err := testListCommand.Run()
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
	e2e.IsJsonWithFormat(t, output, []openapi.Test{})
}
