package actions

import (
	"context"
	"fmt"

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

func (polling pollingActions) Apply(ctx context.Context, fileContent file.File) error {
	var pollingProfile openapi.PollingProfile
	mapstructure.Decode(fileContent.Definition().Spec, &pollingProfile.Spec)

	return polling.resourceClient.Update(ctx, fileContent, currentConfigID)
}

func (polling pollingActions) List(ctx context.Context, listArgs utils.ListArgs) (string, error) {
	return "", ErrNotSupportedResourceAction
}

func (polling pollingActions) Export(ctx context.Context, ID string, filePath string) error {
	pollingProfile, err := polling.resourceClient.Get(ctx, currentConfigID)
	if err != nil {
		return err
	}

	file, err := file.NewFromRaw(filePath, []byte(pollingProfile))
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}

	_, err = file.WriteRaw()
	return err
}

func (polling pollingActions) Delete(ctx context.Context, ID string) error {
	return ErrNotSupportedResourceAction
}

func (polling pollingActions) Get(ctx context.Context, ID string) (string, error) {
	return polling.resourceClient.Get(ctx, currentConfigID)
}
