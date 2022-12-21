package actions

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/openapi"
	"go.uber.org/zap"
)

type ListDataStoreConfig struct{}

type listDataStoreAction struct {
	config config.Config
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
	req := a.client.ApiApi.GetDataStores(ctx)
	dataStores, _, err := a.client.ApiApi.GetDataStoresExecute(req)
	if err != nil {
		return fmt.Errorf("could not list data stores: %w", err)
	}

	formatter := formatters.DataStoreList(a.config)
	output := formatter.Format(dataStores)
	fmt.Println(output)

	return nil
}
