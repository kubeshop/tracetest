package cmd_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/kubeshop/tracetest/cli/cmd/e2e"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type cliOutput struct {
	Test    openapi.Test    `json:"test"`
	TestRun openapi.TestRun `json:"testRun"`
}

func TestRunTestCmd(t *testing.T) {
	cli := e2e.NewCLI()

	definitionFile := "test_run_cmd_test_definition.yml"
	err := copyFile("../testdata/definitions/valid_http_test_definition.yml", definitionFile)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := deleteFile(definitionFile)
		require.NoError(t, err)
	})

	testRunCommand := cli.NewCommand("test", "run", "--config", "e2e/config.yml", "--definition", definitionFile)
	output, err := testRunCommand.Run()
	assert.NoError(t, err)
	assert.NotEmpty(t, output)

	definition, err := file.LoadDefinition(definitionFile)
	require.NoError(t, err)
	assert.NotEmpty(t, definition.Id)

	var outputObject cliOutput
	err = json.Unmarshal([]byte(output), &outputObject)
	require.NoError(t, err)

	assert.NotEmpty(t, outputObject.Test.Id)
	assert.NotEmpty(t, outputObject.TestRun.Id)
	assert.Equal(t, *outputObject.Test.Id, definition.Id)
}

func TestRunTestCmdWhenEditingTest(t *testing.T) {
	cli := e2e.NewCLI()

	definitionFile := "test_run_cmd_test_definition.yml"
	err := copyFile("../testdata/definitions/valid_http_test_definition.yml", definitionFile)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := deleteFile(definitionFile)
		require.NoError(t, err)
	})

	testRunCommand := cli.NewCommand("test", "run", "--config", "e2e/config.yml", "--definition", definitionFile)
	output, err := testRunCommand.Run()
	assert.NoError(t, err)
	assert.NotEmpty(t, output)

	var outputObject cliOutput
	err = json.Unmarshal([]byte(output), &outputObject)
	require.NoError(t, err)

	updateCmdOutput, err := testRunCommand.Run()
	assert.NoError(t, err)
	assert.NotEmpty(t, output)

	var updateOutputObject cliOutput
	err = json.Unmarshal([]byte(updateCmdOutput), &updateOutputObject)
	require.NoError(t, err)

	// Assert a new test wasn't created
	assert.Equal(t, *outputObject.Test.Id, *updateOutputObject.Test.Id)
}

func TestRunTestJUnitCmdValidation(t *testing.T) {
	cli := e2e.NewCLI()

	definitionFile := "test_run_cmd_test_definition.yml"
	err := copyFile("../testdata/definitions/valid_http_test_definition.yml", definitionFile)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := deleteFile(definitionFile)
		require.NoError(t, err)
	})

	testRunCommand := cli.NewCommand("test", "run", "--config", "e2e/config.yml", "--definition", definitionFile, "--junit", "junit_output.xml")
	output, err := testRunCommand.Run()
	assert.NoError(t, err)
	assert.NotEmpty(t, output)
	assert.Contains(t, output, "--junit option requires --wait-for-result")
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
