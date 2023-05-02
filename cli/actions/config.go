package actions

import (
	"context"

	"github.com/kubeshop/tracetest/cli/file"
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

func (config configActions) Apply(ctx context.Context, fileContent file.File) (*file.File, error) {
	return config.resourceClient.Update(ctx, fileContent, currentConfigID)
}

func (config configActions) Get(ctx context.Context, ID string) (*file.File, error) {
	return config.resourceClient.Get(ctx, currentConfigID)
}

func (config configActions) List(ctx context.Context, listArgs utils.ListArgs) (*file.File, error) {
	return nil, ErrNotSupportedResourceAction
}

func (config configActions) Delete(ctx context.Context, ID string) error {
	return ErrNotSupportedResourceAction
}
