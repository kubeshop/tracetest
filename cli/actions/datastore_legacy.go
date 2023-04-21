package actions

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
	"go.uber.org/zap"
)

type ApplyDataStoreConfig struct {
	File string
}

type applyDataStoreAction struct{}

var _ Action[ApplyDataStoreConfig] = &applyDataStoreAction{}

func NewApplyDataStoreAction(logger *zap.Logger, client *openapi.APIClient) applyDataStoreAction {
	return applyDataStoreAction{}
}

func (a applyDataStoreAction) Run(ctx context.Context, args ApplyDataStoreConfig) error {
	if args.File == "" {
		return fmt.Errorf("you must specify a file to be applied")
	}

	return nil
}

type ExportDataStoreConfig struct {
	OutputFile string
	ID         string
}

type exportDataStoreAction struct {
	logger *zap.Logger
	client *openapi.APIClient
}

var _ Action[ExportDataStoreConfig] = &exportDataStoreAction{}

func NewExportDataStoreAction(logger *zap.Logger, client *openapi.APIClient) exportDataStoreAction {
	return exportDataStoreAction{
		logger: logger,
		client: client,
	}
}

func (a exportDataStoreAction) Run(ctx context.Context, args ExportDataStoreConfig) error {
	if args.OutputFile == "" {
		return fmt.Errorf("output file is required. Use --file to specify it")
	}

	return nil
}

type ListDataStoreConfig struct{}

type listDataStoreAction struct {
	logger *zap.Logger
	client *openapi.APIClient
}

var _ Action[ListDataStoreConfig] = &listDataStoreAction{}

func NewListDataStoreAction(config config.Config, logger *zap.Logger, client *openapi.APIClient) listDataStoreAction {
	return listDataStoreAction{
		logger: logger,
		client: client,
	}
}

func (a listDataStoreAction) Run(ctx context.Context, args ListDataStoreConfig) error {

	return nil
}
