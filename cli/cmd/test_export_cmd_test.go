package cmd_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli/cmd/e2e"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestExportCmd(t *testing.T) {
	cli := e2e.NewCLI()

	definitionFile := "definition_file.yml"
	t.Cleanup(func() {
		deleteFile(definitionFile)
	})

	testExportCommand := cli.NewCommand("test", "export", "--config", "e2e/config.yml", "--id", "a6074714-e986-4747-a399-b87748143884", "--output", definitionFile)
	output, err := testExportCommand.Run()
	assert.NoError(t, err)
	assert.Empty(t, output)

	definition, err := file.LoadDefinition(definitionFile)
	require.NoError(t, err)

	assert.NotEmpty(t, definition.Id)
	assert.NotEmpty(t, definition.Name)
	assert.NotEmpty(t, definition.Trigger.Type)
	assert.NotEmpty(t, definition.Trigger.HTTPRequest.URL)
}
