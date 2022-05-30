package actions

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/config"
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
		return fmt.Errorf("You must specify a definition file to run a test")
	}

	a.logger.Debug("Running test from definition", zap.String("definitionFile", args.DefinitionFile))
	return a.runDefinition(ctx, args.DefinitionFile)
}

func (a runTestAction) runDefinition(ctx context.Context, definitionFile string) error {
	definition, err := file.LoadDefinition(definitionFile)
	if err != nil {
		return err
	}

	if definition.Id == "" {
		a.logger.Debug("test doesn't exist. Creating it")
		testID, err := a.createTestFromDefinition(ctx, definition)
	}

	return nil
}

func (a runTestAction) createTestFromDefinition(ctx context.Context, definition definition.Test) (string, error) {
	headers := make([]openapi.HTTPHeader, 0, len(definition.Trigger.HTTPRequest.Headers))
	for _, header := range definition.Trigger.HTTPRequest.Headers {
		headers = append(headers, openapi.HTTPHeader{
			Key:   &header.Key,
			Value: &header.Value,
		})
	}
	testModel := openapi.Test{
		Name:        &definition.Name,
		Description: &definition.Description,
		ServiceUnderTest: &openapi.TestServiceUnderTest{
			Request: &openapi.HTTPRequest{
				Url:     &definition.Trigger.HTTPRequest.URL,
				Method:  &definition.Trigger.HTTPRequest.Method,
				Headers: headers,
				Body:    &definition.Trigger.HTTPRequest.Body.Raw,
				Auth: &openapi.HTTPAuth{
					Type:   &definition.Trigger.HTTPRequest.Authentication.Type,
					ApiKey: &openapi.HTTPAuthApiKey{},
				},
			},
		},
	}
}

func (a runTestAction) executeRequest(ctx context.Context) ([]openapi.Test, error) {
	request := a.client.ApiApi.GetTests(ctx)
	tests, response, err := a.client.ApiApi.GetTestsExecute(request)
	if err != nil {
		return []openapi.Test{}, fmt.Errorf("could not get tests: %w", err)
	}

	if response.StatusCode != 200 {
		return []openapi.Test{}, fmt.Errorf("get tests request failed. Expected 200, got %d", response.StatusCode)
	}

	if tests == nil {
		return []openapi.Test{}, nil
	}

	return tests, nil
}
