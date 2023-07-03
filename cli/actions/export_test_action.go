package actions

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"go.uber.org/zap"
)

type ExportTestConfig struct {
	TestId     string
	OutputFile string
	Version    int32
}

type exportTestAction struct {
	config config.Config
	logger *zap.Logger
	client *openapi.APIClient
}

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
	definition, err := a.getDefinitionFromServer(ctx, args)
	if err != nil {
		return fmt.Errorf("could not get definition from server: %w", err)
	}

	f, err := file.New(args.OutputFile, []byte(definition))
	if err != nil {
		return fmt.Errorf("could not process definition from server: %w", err)
	}

	_, err = f.Write()
	if err != nil {
		return fmt.Errorf("could not save exported definition into file: %w", err)
	}

	return nil
}

func (a exportTestAction) getDefinitionFromServer(ctx context.Context, args ExportTestConfig) (string, error) {
	expectedVersion := args.Version
	if expectedVersion < 0 {
		test, err := a.getTestFromServer(ctx, args.TestId)
		if err != nil {
			return "", fmt.Errorf("could not get test: %w", err)
		}

		expectedVersion = *test.Version
	}

	request := a.client.ApiApi.GetTestVersionDefinitionFile(ctx, args.TestId, expectedVersion)
	definitionString, _, err := a.client.ApiApi.GetTestVersionDefinitionFileExecute(request)
	if err != nil {
		return "", fmt.Errorf("could not get test definition: %w", err)
	}

	return definitionString, nil
}

func (a exportTestAction) getTestFromServer(ctx context.Context, testID string) (openapi.Test, error) {
	req := a.client.ApiApi.GetTest(ctx, testID)
	openapiTest, _, err := a.client.ApiApi.GetTestExecute(req)
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not execute getTest request: %w", err)
	}

	return *openapiTest, nil
}
