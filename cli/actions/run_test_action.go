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

		fmt.Println(testID, err)
	}

	return nil
}

func (a runTestAction) createTestFromDefinition(ctx context.Context, definition definition.Test) (string, error) {
	openapiTest := conversion.ConvertTestDefinitionIntoOpenAPIObject(definition)
	fmt.Println(openapiTest)
	return "", nil
}
