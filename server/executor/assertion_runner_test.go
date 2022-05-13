package executor_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/kubeshop/tracetest/executor"
	"github.com/kubeshop/tracetest/openapi"
	"github.com/kubeshop/tracetest/testmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutorSuccessfulExecution(t *testing.T) {
	testCases := []struct {
		Name       string
		Tracefile  string
		ShouldPass bool
	}{
		{
			Name:       "pokeshop - import pokemon: should pass",
			Tracefile:  "../testmock/data/pokeshop_import_pokemon.json",
			ShouldPass: true,
		},
		{
			Name:       "pokeshop - import pokemon: should fail",
			Tracefile:  "../testmock/data/pokeshop_import_pokemon_failed_assertions.json",
			ShouldPass: false,
		},
	}

	postgresRepository, err := testmock.GetTestingDatabase("file://../migrations")
	require.NoError(t, err)

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			ctx := context.Background()

			test, result, err := loadTestFile(testCase.Tracefile)
			require.NoError(t, err)

			assertionExecutor := executor.NewAssertionRunner(postgresRepository)

			createdTestID, err := postgresRepository.CreateTest(ctx, &test)
			require.NoError(t, err)

			for _, assertion := range test.Assertions {
				postgresRepository.CreateAssertion(ctx, test.TestId, &assertion)
			}

			result.TestId = createdTestID
			err = postgresRepository.CreateResult(ctx, test.TestId, &result)
			require.NoError(t, err)

			testDefinition, err := executor.ConvertAssertionsIntoTestDefinition(test.Assertions)
			assert.NoError(t, err)

			assertionRequest := executor.AssertionRequest{
				TestDefinition: testDefinition,
				Result:         result,
			}

			assertionExecutor.Start(1)
			assertionExecutor.RunAssertions(assertionRequest)
			assertionExecutor.Stop()

			dbResult, err := postgresRepository.GetResult(ctx, result.ResultId)
			require.NoError(t, err)

			if testCase.ShouldPass {
				assert.Equal(t, executor.TestRunStateFinished, dbResult.State)
				for _, assertionResult := range dbResult.AssertionResult {
					for _, spanAssertionResult := range assertionResult.SpanAssertionResults {
						assert.True(t, spanAssertionResult.Passed)
					}
				}
				assert.True(t, dbResult.AssertionResultState)
			} else {
				assert.False(t, dbResult.AssertionResultState)
			}

		})
	}
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

	testFile.Test.TestId = uuid.NewString()
	testFile.Result.TestId = testFile.Test.TestId
	testFile.Result.ResultId = uuid.NewString()

	return testFile.Test, testFile.Result, nil
}
