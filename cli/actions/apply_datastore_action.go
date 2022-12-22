package actions

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"go.uber.org/zap"
)

type ApplyDataStoreConfig struct {
	File string
}

type applyDataStoreAction struct {
	logger *zap.Logger
	client *openapi.APIClient
}

var _ Action[ApplyDataStoreConfig] = &applyDataStoreAction{}

func NewApplyDataStoreAction(logger *zap.Logger, client *openapi.APIClient) applyDataStoreAction {
	return applyDataStoreAction{
		logger: logger,
		client: client,
	}
}

func (a applyDataStoreAction) Run(ctx context.Context, args ApplyDataStoreConfig) error {
	if args.File == "" {
		return fmt.Errorf("you must specify a file to be applied")
	}

	a.logger.Debug(
		"applying data store",
		zap.String("file", args.File),
	)

	fileContent, err := file.Read(args.File)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	if fileContent.Definition().Type != "DataStore" {
		return fmt.Errorf(`file must be of type "DataStore"`)
	}

	var dataStore openapi.DataStore
	deepcopy.DeepCopy(fileContent.Definition().Spec, &dataStore)

	if dataStore.Id == nil || *dataStore.Id == "" {
		err := a.createDataStore(ctx, fileContent, dataStore)
		if err != nil {
			return err
		}
	} else {
		err := a.updateDataStore(ctx, fileContent, dataStore)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a applyDataStoreAction) createDataStore(ctx context.Context, file file.File, dataStore openapi.DataStore) error {
	req := a.client.ApiApi.CreateDataStore(ctx)
	req = req.DataStore(dataStore)
	createdDataStore, resp, err := a.client.ApiApi.CreateDataStoreExecute(req)
	if resp.StatusCode == http.StatusUnprocessableEntity {
		// validation error
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		validationError := string(body)
		return fmt.Errorf("invalid data store: %s", validationError)
	}
	if err != nil {
		return fmt.Errorf("could not create data store: %w", err)
	}

	f, err := file.SetID(*createdDataStore.Id)
	if err != nil {
		return fmt.Errorf("could no set data store id: %w", err)
	}

	_, err = f.Write()
	if err != nil {
		return fmt.Errorf("could not write to data store file: %w", err)
	}

	return nil
}

func (a applyDataStoreAction) updateDataStore(ctx context.Context, file file.File, dataStore openapi.DataStore) error {
	req := a.client.ApiApi.UpdateDataStore(ctx, *dataStore.Id)
	req = req.DataStore(dataStore)
	resp, err := a.client.ApiApi.UpdateDataStoreExecute(req)
	if resp.StatusCode == http.StatusUnprocessableEntity {
		// validation error
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		validationError := string(body)
		return fmt.Errorf("invalid data store: %s", validationError)
	}
	if err != nil {
		return fmt.Errorf("could not update data store: %w", err)
	}

	return nil
}
