package assertions_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/kubeshop/tracetest/assertions"
	"github.com/kubeshop/tracetest/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutor(t *testing.T) {
	test, result, err := loadTestFile("../test/data/pokeshop_import_pokemon.json")
	require.NoError(t, err)

	runAssertionsMessage := assertions.RunAssertionsMessage{
		Test:   test,
		Result: result,
	}

	inputChannel := make(chan assertions.RunAssertionsMessage, 1)
	outputChannel := make(chan assertions.RunAssertionsMessage, 1)

	executor := assertions.NewExecutor(inputChannel, outputChannel)

	go executor.Start()

	inputChannel <- runAssertionsMessage
	outputMessage := <-outputChannel

	assert.NotNil(t, outputMessage)
}

type testFile struct {
	Test   openapi.Test          `json:"test"`
	Result openapi.TestRunResult `json:"result"`
}

func loadTestFile(filePath string) (openapi.Test, openapi.TestRunResult, error) {
	fileContent, err := os.Open(filePath)
	if err != nil {
		return openapi.Test{}, openapi.TestRunResult{}, fmt.Errorf("could not open test file: %w", err)
	}

	fileBytes, err := ioutil.ReadAll(fileContent)
	if err != nil {
		return openapi.Test{}, openapi.TestRunResult{}, fmt.Errorf("could not read test file: %w", err)
	}

	testFile := testFile{}
	err = json.Unmarshal(fileBytes, &testFile)
	if err != nil {
		return openapi.Test{}, openapi.TestRunResult{}, fmt.Errorf("could not parse test file: %w", err)
	}

	return testFile.Test, testFile.Result, nil
}
