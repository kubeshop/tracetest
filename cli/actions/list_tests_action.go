package actions

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
	"go.uber.org/zap"
)

type listTestsAction struct {
	config config.Config
	logger *zap.Logger
	client *openapi.APIClient
}

var _ Action = &listTestsAction{}

func NewListTestsAction(config config.Config, logger *zap.Logger, client *openapi.APIClient) listTestsAction {
	return listTestsAction{config, logger, client}
}

func (a listTestsAction) Run(ctx context.Context, args []string) error {
	tests, err := a.executeRequest(ctx)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(tests)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))
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
