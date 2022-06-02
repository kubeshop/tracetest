package cmd_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/kubeshop/tracetest/cli/cmd/e2e"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunTestCmd(t *testing.T) {
	cli := e2e.NewCLI()

	definitionFile := "test_run_cmd_test_definition.yml"
	err := copyFile("../testdata/definitions/valid_http_test_definition.yml", definitionFile)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := deleteFile(definitionFile)
		require.NoError(t, err)
	})

	output, err := cli.RunCommand("test", "run", "--config", "e2e/config.yml", "--definition", definitionFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, output)

	definition, err := file.LoadDefinition(definitionFile)
	require.NoError(t, err)
	assert.NotEmpty(t, definition.Id)

	type cliOutput struct {
		TestId    string `json:"testId"`
		TestRunId string `json:"testRunId"`
	}

	var outputObject cliOutput
	err = json.Unmarshal([]byte(output), &outputObject)
	require.NoError(t, err)

	assert.NotEmpty(t, outputObject.TestId)
	assert.NotEmpty(t, outputObject.TestRunId)
	assert.Equal(t, outputObject.TestId, definition.Id)
}

func copyFile(source string, destination string) error {
	input, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("could not read source file: %w", err)
	}

	err = ioutil.WriteFile(destination, input, 0644)
	if err != nil {
		return fmt.Errorf("could not write destination file: %w", err)
	}

	return nil
}

func deleteFile(target string) error {
	err := os.Remove(target)
	if err != nil {
		return fmt.Errorf("could not remove file: %w", err)
	}

	return nil
}
