package executor_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/kubeshop/tracetest/server/tracing"
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
			Name:       "ImportPokemonSucess",
			Tracefile:  "../testmock/data/pokeshop_import_pokemon.json",
			ShouldPass: true,
		},
		{
			Name:       "ImportPokemonFail",
			Tracefile:  "../testmock/data/pokeshop_import_pokemon_failed_assertions.json",
			ShouldPass: false,
		},
	}

	repo, err := testmock.GetTestingDatabase("file://../migrations")
	require.NoError(t, err)

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			ctx := context.Background()

			test, run, err := loadTestFile(testCase.Tracefile)
			require.NoError(t, err)

			tracer, _ := tracing.NewTracer(ctx, config.Config{})
			assertionExecutor := executor.NewAssertionExecutor(tracer)
			assertionRunner := executor.NewAssertionRunner(executor.NewDBUpdater(repo), assertionExecutor)

			test, err = repo.CreateTest(ctx, test)
			require.NoError(t, err)

			err = repo.SetDefiniton(ctx, test, test.Spec)
			require.NoError(t, err)

			run, err = repo.CreateRun(ctx, test, run)
			require.NoError(t, err)

			assertionRequest := executor.AssertionRequest{
				Test: test,
				Run:  run,
			}

			assertionRunner.Start(1)
			assertionRunner.RunAssertions(ctx, assertionRequest)
			assertionRunner.Stop()

			dbResult, err := repo.GetRun(ctx, run.ID)
			require.NoError(t, err)

			require.NotNil(t, dbResult.Results)
			if testCase.ShouldPass {
				assert.Equal(t, model.RunStateFinished, dbResult.State)
				dbResult.Results.Results.Map(func(_ model.SpanQuery, results []model.AssertionResult) {
					for _, assertRes := range results {
						for _, spanAssertionRes := range assertRes.Results {
							assert.NoError(t, spanAssertionRes.CompareErr)
						}
					}
				})
				assert.True(t, dbResult.Results.AllPassed)
			} else {
				assert.False(t, dbResult.Results.AllPassed)
			}

		})
	}
}

type testFile struct {
	Test model.Test `json:"test"`
	Run  model.Run  `json:"run"`
}

func loadTestFile(filePath string) (model.Test, model.Run, error) {
	fileContent, err := os.Open(filePath)
	if err != nil {
		return model.Test{}, model.Run{}, fmt.Errorf("could not open test file: %w", err)
	}

	fileBytes, err := ioutil.ReadAll(fileContent)
	if err != nil {
		return model.Test{}, model.Run{}, fmt.Errorf("could not read test file: %w", err)
	}

	testFile := testFile{}
	err = json.Unmarshal(fileBytes, &testFile)
	if err != nil {
		return model.Test{}, model.Run{}, fmt.Errorf("could not parse test file: %w", err)
	}

	testFile.Test.ID = id.NewRandGenerator().UUID()
	testFile.Run.ID = id.NewRandGenerator().UUID()

	return testFile.Test, testFile.Run, nil
}
