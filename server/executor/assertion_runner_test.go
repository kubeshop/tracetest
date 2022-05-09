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
	"github.com/kubeshop/tracetest/test"
	"github.com/kubeshop/tracetest/testdb"
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
			Tracefile:  "../test/data/pokeshop_import_pokemon.json",
			ShouldPass: true,
		},
		{
			Name:       "pokeshop - import pokemon: should fail",
			Tracefile:  "../test/data/pokeshop_import_pokemon_failed_assertions.json",
			ShouldPass: false,
		},
	}

	sqlDb, err := test.GetTestingDatabase()
	require.NoError(t, err)

	postgresRepository, err := testdb.Postgres(
		testdb.WithDB(sqlDb),
		testdb.WithMigrations("file://../migrations"),
	)
	require.NoError(t, err)

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			ctx := context.Background()
			t.Parallel()

			test, result, err := loadTestFile(testCase.Tracefile)
			require.NoError(t, err)

			inputChannel := make(chan openapi.TestRunResult)
			assertionExecutor := executor.NewAssertionRunner(postgresRepository, postgresRepository, inputChannel)

			_, err = postgresRepository.CreateTest(ctx, &test)
			require.NoError(t, err)

			err = postgresRepository.CreateResult(ctx, test.TestId, &result)
			require.NoError(t, err)

			assertionExecutor.Start(1)
			assertionExecutor.RunAssertions(result)
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
	testFile.Result.ResultId = uuid.NewString()

	return testFile.Test, testFile.Result, nil
}
