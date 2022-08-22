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

type ExportTestConfig struct {
	TestId     string
	OutputFile string
}

type exportTestAction struct {
	config config.Config
	logger *zap.Logger
	client *openapi.APIClient
}

var _ Action[ExportTestConfig] = &exportTestAction{}

func NewExportTestAction(config config.Config, logger *zap.Logger, client *openapi.APIClient) exportTestAction {
	return exportTestAction{
		config: config,
		logger: logger,
		client: client,
	}
}

func (a exportTestAction) Run(ctx context.Context, args ExportTestConfig) error {
	if args.OutputFile == "" {
		return fmt.Errorf("output file must be provided")
	}

	if args.TestId == "" {
		return fmt.Errorf("test id must be provided")
	}

	a.logger.Debug("exporting test", zap.String("testID", args.TestId), zap.String("outputFile", args.OutputFile))
	definition, err := a.getDefinitionFromServer(ctx, args.TestId)
	if err != nil {
		return fmt.Errorf("could not get definition from server: %w", err)
	}

	err = file.SaveDefinition(args.OutputFile, definition)
	if err != nil {
		return fmt.Errorf("could not save exported definition into file: %w", err)
	}

	return nil
}

func (a exportTestAction) getDefinitionFromServer(ctx context.Context, testID string) (definition.Test, error) {
	openapiTest, err := a.getTestFromServer(ctx, testID)
	if err != nil {
		return definition.Test{}, fmt.Errorf("could not get test from server: %w", err)
	}

	spec, err := conversion.ConvertOpenAPITestIntoSpecObject(openapiTest)
	if err != nil {
		return definition.Test{}, fmt.Errorf("could not convert openapi object into a defintion object: %w", err)
	}

	return spec, nil
}

func (a exportTestAction) getTestFromServer(ctx context.Context, testID string) (openapi.Test, error) {
	req := a.client.ApiApi.GetTest(ctx, testID)
	openapiTest, _, err := a.client.ApiApi.GetTestExecute(req)
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not execute getTest request: %w", err)
	}

	return *openapiTest, nil
}
