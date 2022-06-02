package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/conversion"
	"github.com/kubeshop/tracetest/cli/definition"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"go.uber.org/zap"
)

type RunTestConfig struct {
	DefinitionFile string
	WaitForResult  bool
}

type runTestAction struct {
	config config.Config
	logger *zap.Logger
	client *openapi.APIClient
}

var _ Action[RunTestConfig] = &runTestAction{}

type runTestOutput struct {
	Test      openapi.Test    `json:"test"`
	Run       openapi.TestRun `json:"testRun"`
	RunWebURL string          `json:"testRunWebUrl"`
}

func NewRunTestAction(config config.Config, logger *zap.Logger, client *openapi.APIClient) runTestAction {
	return runTestAction{config, logger, client}
}

func (a runTestAction) Run(ctx context.Context, args RunTestConfig) error {
	if args.DefinitionFile == "" {
		return fmt.Errorf("you must specify a definition file to run a test")
	}

	a.logger.Debug("Running test from definition", zap.String("definitionFile", args.DefinitionFile))
	output, err := a.runDefinition(ctx, args.DefinitionFile, args.WaitForResult)

	if err != nil {
		return fmt.Errorf("could not run definition: %w", err)
	}

	bytes, err := json.Marshal(output)
	if err != nil {
		return fmt.Errorf("could not marshal output json: %w", err)
	}

	fmt.Print(string(bytes))
	return nil
}

func (a runTestAction) runDefinition(ctx context.Context, definitionFile string, waitForResult bool) (runTestOutput, error) {
	definition, err := file.LoadDefinition(definitionFile)
	if err != nil {
		return runTestOutput{}, err
	}

	err = definition.Validate()
	if err != nil {
		return runTestOutput{}, fmt.Errorf("invalid definition file: %w", err)
	}

	var test openapi.Test

	if definition.Id == "" {
		a.logger.Debug("test doesn't exist. Creating it")
		test, err = a.createTestFromDefinition(ctx, definition)
		if err != nil {
			return runTestOutput{}, fmt.Errorf("could not create test from definition: %w", err)
		}

		definition.Id = *test.Id
		err = file.SaveDefinition(definitionFile, definition)
		if err != nil {
			return runTestOutput{}, fmt.Errorf("could not save definition: %w", err)
		}
	} else {
		a.logger.Debug("test exists. Updating it")
		test, err = a.updateTestFromDefinition(ctx, definition)
		if err != nil {
			return runTestOutput{}, fmt.Errorf("could not update test using definition: %w", err)
		}
	}

	testRun, err := a.runTest(ctx, definition.Id)
	if err != nil {
		return runTestOutput{}, fmt.Errorf("could not run test: %w", err)
	}

	if waitForResult {
		updatedTestRun, err := a.waitForResult(ctx, definition.Id, *testRun.Id)
		if err != nil {
			return runTestOutput{}, fmt.Errorf("could not wait for result: %w", err)
		}

		testRun = updatedTestRun
	}

	return runTestOutput{
		Test:      test,
		Run:       testRun,
		RunWebURL: fmt.Sprintf("%s://%s/test/%s/run/%s", a.config.Scheme, a.config.Endpoint, definition.Id, *testRun.Id),
	}, nil
}

func (a runTestAction) createTestFromDefinition(ctx context.Context, definition definition.Test) (openapi.Test, error) {
	openapiTest, err := conversion.ConvertTestDefinitionIntoOpenAPIObject(definition)
	if err != nil {
		return openapi.Test{}, err
	}

	req := a.client.ApiApi.CreateTest(ctx)
	req = req.Test(openapiTest)

	testBytes, err := json.Marshal(openapiTest)
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not marshal test: %w", err)
	}

	a.logger.Debug("Sending request to create test", zap.ByteString("test", testBytes))
	createdTest, _, err := a.client.ApiApi.CreateTestExecute(req)
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not execute request: %w", err)
	}

	return *createdTest, nil
}

func (a runTestAction) updateTestFromDefinition(ctx context.Context, definition definition.Test) (openapi.Test, error) {
	openapiTest, err := conversion.ConvertTestDefinitionIntoOpenAPIObject(definition)
	if err != nil {
		return openapi.Test{}, err
	}

	req := a.client.ApiApi.UpdateTest(ctx, definition.Id)
	req = req.Test(openapiTest)

	testBytes, err := json.Marshal(openapiTest)
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not marshal test: %w", err)
	}

	a.logger.Debug("Sending request to update test", zap.ByteString("test", testBytes))
	_, err = a.client.ApiApi.UpdateTestExecute(req)
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not execute request: %w", err)
	}

	return openapiTest, nil
}

func (a runTestAction) runTest(ctx context.Context, testID string) (openapi.TestRun, error) {
	req := a.client.ApiApi.RunTest(ctx, testID)
	run, _, err := a.client.ApiApi.RunTestExecute(req)
	if err != nil {
		return openapi.TestRun{}, fmt.Errorf("could not execute request: %w", err)
	}

	return *run, nil
}

func (a runTestAction) waitForResult(ctx context.Context, testId string, testRunId string) (openapi.TestRun, error) {
	var testRun openapi.TestRun
	var lastError error
	var wg sync.WaitGroup
	wg.Add(1)
	ticker := time.NewTicker(1 * time.Second) // TODO: make this configurable
	go func() {
		for {
			select {
			case <-ticker.C:
				readyTestRun, err := a.isTestReady(ctx, testId, testRunId)
				if err != nil {
					lastError = err
					wg.Done()
					return
				}

				if readyTestRun != nil {
					testRun = *readyTestRun
					wg.Done()
					return
				}
			}
		}
	}()
	wg.Wait()

	if lastError != nil {
		return openapi.TestRun{}, lastError
	}

	return testRun, nil
}

func (a runTestAction) isTestReady(ctx context.Context, testId string, testRunId string) (*openapi.TestRun, error) {
	req := a.client.ApiApi.GetTestRun(ctx, testId, testRunId)
	run, _, err := a.client.ApiApi.GetTestRunExecute(req)
	if err != nil {
		return &openapi.TestRun{}, fmt.Errorf("could not execute GetTestRun request: %w", err)
	}

	if *run.State == "FAILED" || *run.State == "FINISHED" {
		return run, nil
	}

	return nil, nil
}
