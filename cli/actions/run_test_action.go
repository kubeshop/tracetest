package actions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/conversion"
	"github.com/kubeshop/tracetest/cli/definition"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"go.uber.org/zap"
)

type RunTestConfig struct {
	DefinitionFile string
}

type runTestAction struct {
	config config.Config
	logger *zap.Logger
	client *openapi.APIClient
}

var _ Action[RunTestConfig] = &runTestAction{}

type runTestOutput struct {
	TestID    string `json:"testId"`
	RunID     string `json:"testRunId"`
	RunWebURL string `json:"testRunWebUrl"`
}

func NewRunTestAction(config config.Config, logger *zap.Logger, client *openapi.APIClient) runTestAction {
	return runTestAction{config, logger, client}
}

func (a runTestAction) Run(ctx context.Context, args RunTestConfig) error {
	if args.DefinitionFile == "" {
		return fmt.Errorf("you must specify a definition file to run a test")
	}

	a.logger.Debug("Running test from definition", zap.String("definitionFile", args.DefinitionFile))
	output, err := a.runDefinition(ctx, args.DefinitionFile)

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

func (a runTestAction) runDefinition(ctx context.Context, definitionFile string) (runTestOutput, error) {
	definition, err := file.LoadDefinition(definitionFile)
	if err != nil {
		return runTestOutput{}, err
	}

	err = definition.Validate()
	if err != nil {
		return runTestOutput{}, fmt.Errorf("invalid definition file: %w", err)
	}

	if definition.Id == "" {
		a.logger.Debug("test doesn't exist. Creating it")
		testID, err := a.createTestFromDefinition(ctx, definition)
		if err != nil {
			return runTestOutput{}, fmt.Errorf("could not create test from definition: %w", err)
		}

		definition.Id = testID
		err = file.SaveDefinition(definitionFile, definition)
		if err != nil {
			return runTestOutput{}, fmt.Errorf("could not save definition: %w", err)
		}
	} else {
		a.logger.Debug("test exists. Updating it")
		err = a.updateTestFromDefinition(ctx, definition)
		if err != nil {
			return runTestOutput{}, fmt.Errorf("could not update test using definition: %w", err)
		}
	}

	testRun, err := a.runTest(ctx, definition.Id)
	if err != nil {
		return runTestOutput{}, fmt.Errorf("could not run test: %w", err)
	}

	return runTestOutput{
		TestID:    definition.Id,
		RunID:     *testRun.Id,
		RunWebURL: fmt.Sprintf("%s://%s/test/%s/run/%s", a.config.Scheme, a.config.Endpoint, definition.Id, *testRun.Id),
	}, nil
}

func (a runTestAction) createTestFromDefinition(ctx context.Context, definition definition.Test) (string, error) {
	openapiTest, err := conversion.ConvertTestDefinitionIntoOpenAPIObject(definition)
	if err != nil {
		return "", err
	}

	req := a.client.ApiApi.CreateTest(ctx)
	req = req.Test(openapiTest)

	testBytes, err := json.Marshal(openapiTest)
	if err != nil {
		return "", fmt.Errorf("could not marshal test: %w", err)
	}

	a.logger.Debug("Sending request to create test", zap.ByteString("test", testBytes))
	createdTest, _, err := a.client.ApiApi.CreateTestExecute(req)
	if err != nil {
		return "", fmt.Errorf("could not execute request: %w", err)
	}

	return *createdTest.Id, nil
}

func (a runTestAction) updateTestFromDefinition(ctx context.Context, definition definition.Test) error {
	openapiTest, err := conversion.ConvertTestDefinitionIntoOpenAPIObject(definition)
	if err != nil {
		return err
	}

	req := a.client.ApiApi.UpdateTest(ctx, definition.Id)
	req = req.Test(openapiTest)

	testBytes, err := json.Marshal(openapiTest)
	if err != nil {
		return fmt.Errorf("could not marshal test: %w", err)
	}

	a.logger.Debug("Sending request to update test", zap.ByteString("test", testBytes))
	_, err = a.client.ApiApi.UpdateTestExecute(req)
	if err != nil {
		return fmt.Errorf("could not execute request: %w", err)
	}

	return nil
}

func (a runTestAction) runTest(ctx context.Context, testID string) (openapi.TestRun, error) {
	req := a.client.ApiApi.RunTest(ctx, testID)
	run, _, err := a.client.ApiApi.RunTestExecute(req)
	if err != nil {
		return openapi.TestRun{}, fmt.Errorf("could not execute request: %w", err)
	}

	return *run, nil
}
