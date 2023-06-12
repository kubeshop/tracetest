package actions

import (
	"context"

	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/kubeshop/tracetest/server/model/yaml"
)

type configActions struct {
	resourceArgs
}

var _ ResourceActions = &configActions{}
var currentConfigID = "current"

func NewConfigActions(options ...ResourceArgsOption) configActions {
	args := NewResourceArgs(options...)

	return configActions{
		resourceArgs: args,
	}
}

func (configActions) FileType() yaml.FileType {
	return yaml.FileTypeConfig
}

func (configActions) Name() string {
	return "config"
}

func (config configActions) GetID(file *file.File) (string, error) {
	resource, err := config.formatter.ToStruct(file)
	if err != nil {
		return "", err
	}

	return *resource.(openapi.ConfigurationResource).Spec.Id, nil
}

func (config configActions) Apply(ctx context.Context, fileContent file.File) (result *file.File, err error) {
	result, err = config.resourceClient.Update(ctx, fileContent, currentConfigID)
	return result, err
}

func (config configActions) Get(ctx context.Context, ID string) (*file.File, error) {
	return config.resourceClient.Get(ctx, currentConfigID)
}

func (config configActions) List(ctx context.Context, listArgs utils.ListArgs) (*file.File, error) {
	return config.resourceClient.List(ctx, listArgs)
}

func (config configActions) Delete(ctx context.Context, ID string) (string, error) {
	return "", ErrNotSupportedResourceAction // we don't have support to delete config today
}
