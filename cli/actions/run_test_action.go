package actions

import (
	"context"
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

func NewRunTestAction(config config.Config, logger *zap.Logger, client *openapi.APIClient) runTestAction {
	return runTestAction{config, logger, client}
}

func (a runTestAction) Run(ctx context.Context, args RunTestConfig) error {
	if args.DefinitionFile == "" {
		return fmt.Errorf("you must specify a definition file to run a test")
	}

	a.logger.Debug("Running test from definition", zap.String("definitionFile", args.DefinitionFile))
	return a.runDefinition(ctx, args.DefinitionFile)
}

func (a runTestAction) runDefinition(ctx context.Context, definitionFile string) error {
	definition, err := file.LoadDefinition(definitionFile)
	if err != nil {
		return err
	}

	err = definition.Validate()
	if err != nil {
		return fmt.Errorf("invalid definition file: %w", err)
	}

	if definition.Id == "" {
		a.logger.Debug("test doesn't exist. Creating it")
		testID, err := a.createTestFromDefinition(ctx, definition)
		if err != nil {
			return fmt.Errorf("could not create test from definition: %w", err)
		}

		definition.Id = testID
		err = file.SaveDefinition(definitionFile, definition)
		if err != nil {
			return fmt.Errorf("could not save definition: %w", err)
		}
	}

	// TODO: update definition

	testRunId, err := a.runTest(ctx, definition.Id)
	if err != nil {
		return fmt.Errorf("could not run test: %w", err)
	}

	fmt.Println("testRunId", testRunId)

	return nil
}

func (a runTestAction) createTestFromDefinition(ctx context.Context, definition definition.Test) (string, error) {
	openapiTest, err := conversion.ConvertTestDefinitionIntoOpenAPIObject(definition)
	if err != nil {
		return "", err
	}

	req := a.client.ApiApi.CreateTest(ctx)
	req.Test(openapiTest)

	createdTest, _, err := a.client.ApiApi.CreateTestExecute(req)
	if err != nil {
		return "", fmt.Errorf("could not execute request: %w", err)
	}

	return *createdTest.Id, nil
}

func (a runTestAction) runTest(ctx context.Context, testID string) (string, error) {
	req := a.client.ApiApi.RunTest(ctx, testID)
	run, _, err := a.client.ApiApi.RunTestExecute(req)
	if err != nil {
		return "", fmt.Errorf("could not execute request: %w", err)
	}

	return *run.Id, nil
}
