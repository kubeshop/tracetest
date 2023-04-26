package actions

import (
	"context"

	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"github.com/mitchellh/mapstructure"
)

type pollingActions struct {
	resourceArgs
}

var _ ResourceActions = &pollingActions{}

func NewPollingActions(options ...ResourceArgsOption) pollingActions {
	args := NewResourceArgs(options...)

	return pollingActions{
		resourceArgs: args,
	}
}

func (pollingActions) FileType() yaml.FileType {
	return yaml.FileTypePollingProfile
}

func (pollingActions) Name() string {
	return "pollingprofile"
}

func (polling pollingActions) Apply(ctx context.Context, fileContent file.File) (*file.File, error) {
	var pollingProfile openapi.PollingProfile
	mapstructure.Decode(fileContent.Definition().Spec, &pollingProfile.Spec)

	return polling.resourceClient.Update(ctx, fileContent, currentConfigID)
}

func (polling pollingActions) List(ctx context.Context, listArgs utils.ListArgs) (*file.File, error) {
	return nil, ErrNotSupportedResourceAction
}

func (polling pollingActions) Delete(ctx context.Context, ID string) error {
	return ErrNotSupportedResourceAction
}

func (polling pollingActions) Get(ctx context.Context, ID string) (*file.File, error) {
	return polling.resourceClient.Get(ctx, currentConfigID)
}
