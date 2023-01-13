package actions

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"go.uber.org/zap"
)

type ApplyEnvironmentConfig struct {
	File string
}

type applyEnvironmentAction struct {
	logger *zap.Logger
	client *openapi.APIClient
}

var _ Action[ApplyEnvironmentConfig] = &applyEnvironmentAction{}

func NewApplyEnvironmentAction(logger *zap.Logger, client *openapi.APIClient) applyEnvironmentAction {
	return applyEnvironmentAction{
		logger: logger,
		client: client,
	}
}

func (a applyEnvironmentAction) Run(ctx context.Context, args ApplyEnvironmentConfig) error {
	if args.File == "" {
		return fmt.Errorf("you must specify a file to be applied")
	}

	a.logger.Debug(
		"applying environment",
		zap.String("file", args.File),
	)

	fileContent, err := file.Read(args.File)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	if fileContent.Definition().Type != "Environment" {
		return fmt.Errorf(`file must be of type "Environment"`)
	}

	fileContentString := fileContent.Contents()

	req := a.client.ApiApi.ExecuteDefinition(ctx)
	req = req.TextDefinition(openapi.TextDefinition{
		RunInformation: &openapi.RunInformation{},
		Content:        &fileContentString,
	})

	response, _, err := req.Execute()
	if err != nil {
		return fmt.Errorf("could not apply environment: %w", err)
	}

	if fileContent.HasID() || response.Id == nil {
		return nil
	}

	newFile, err := fileContent.SetID(*response.Id)
	if err != nil {
		return fmt.Errorf("could not set id into environment file: %w", err)
	}

	_, err = newFile.Write()
	return err
}
