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

func (polling pollingActions) GetID(file *file.File) (string, error) {
	resource, err := polling.formatter.ToStruct(file)
	if err != nil {
		return "", err
	}

	return resource.(openapi.PollingProfile).Spec.Id, nil
}

func (polling pollingActions) Apply(ctx context.Context, fileContent file.File) (result *file.File, err error) {
	var pollingProfile openapi.PollingProfile
	mapstructure.Decode(fileContent.Definition().Spec, &pollingProfile.Spec)

	result, err = polling.resourceClient.Update(ctx, fileContent, currentConfigID)
	return result, err
}

func (polling pollingActions) List(ctx context.Context, listArgs utils.ListArgs) (*file.File, error) {
	return polling.resourceClient.List(ctx, listArgs)
}

func (polling pollingActions) Delete(ctx context.Context, ID string) (string, error) {
	return "PollingProfile successfully reset to default", ErrNotSupportedResourceAction
}

func (polling pollingActions) Get(ctx context.Context, ID string) (*file.File, error) {
	return polling.resourceClient.Get(ctx, currentConfigID)
}
