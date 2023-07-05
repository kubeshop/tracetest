package actions

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/openapi"
	"go.uber.org/zap"
)

type ListTestConfig struct{}

type listTestsAction struct {
	config config.Config
	logger *zap.Logger
	client *openapi.APIClient
}

func NewListTestsAction(config config.Config, logger *zap.Logger, client *openapi.APIClient) listTestsAction {
	return listTestsAction{config, logger, client}
}

func (a listTestsAction) Run(ctx context.Context, args ListTestConfig) error {
	tests, err := a.executeRequest(ctx)
	if err != nil {
		return err
	}

	formatter := formatters.TestsList(a.config)
	formattedOutput := formatter.Format(tests)
	fmt.Println(formattedOutput)

	return nil
}

func (a listTestsAction) executeRequest(ctx context.Context) ([]openapi.Test, error) {
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
