package actions

import (
	"context"

	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"github.com/mitchellh/mapstructure"
)

type dataStoreActions struct {
	resourceArgs
}

var _ ResourceActions = &dataStoreActions{}

func NewDataStoreActions(options ...ResourceArgsOption) *dataStoreActions {
	args := NewResourceArgs(options...)

	return &dataStoreActions{
		resourceArgs: args,
	}
}

func (d *dataStoreActions) FileType() yaml.FileType {
	return yaml.FileTypeDataStore
}

func (d *dataStoreActions) Name() string {
	return "datastore"
}

func (d *dataStoreActions) Apply(ctx context.Context, fileContent file.File) (*file.File, error) {
	var dataStore openapi.DataStore
	mapstructure.Decode(fileContent.Definition().Spec, &dataStore)

	return d.resourceClient.Update(ctx, fileContent, currentConfigID)
}

func (d *dataStoreActions) List(ctx context.Context, args utils.ListArgs) (*file.File, error) {
	return nil, ErrNotSupportedResourceAction
}

func (d *dataStoreActions) Get(ctx context.Context, id string) (*file.File, error) {
	return d.resourceClient.Get(ctx, currentConfigID)
}

func (d *dataStoreActions) Delete(ctx context.Context, id string) error {
	return d.resourceClient.Delete(ctx, currentConfigID)
}
