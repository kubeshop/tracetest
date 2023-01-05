package actions

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"go.uber.org/zap"
)

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

	if args.ID == "" {
		return fmt.Errorf("id is required. Use --id to specify it")
	}

	req := a.client.ApiApi.GetDataStoreDefinitionFile(ctx, args.ID)
	dataStoreContent, _, err := a.client.ApiApi.GetDataStoreDefinitionFileExecute(req)
	if err != nil {
		return fmt.Errorf("could not get data store by its id: %w", err)
	}

	file, err := file.New(args.OutputFile, []byte(dataStoreContent))
	if err != nil {
		return fmt.Errorf("could not process definition from server: %w", err)
	}

	_, err = file.Write()
	if err != nil {
		return fmt.Errorf("could not save exported definition into file: %w", err)
	}

	return nil
}
