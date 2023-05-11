package resources_actions

import (
	"context"

	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"github.com/mitchellh/mapstructure"
)

type demoActions struct {
	resourceArgs
}

var _ ResourceActionsInterface = &demoActions{}

func NewDemoActions(options ...ResourceArgsOption) demoActions {
	args := NewResourceArgs(options...)

	return demoActions{
		resourceArgs: args,
	}
}

func (demoActions) FileType() yaml.FileType {
	return yaml.FileTypeDemo
}

func (demoActions) Name() string {
	return "demo"
}

func (demo demoActions) GetID(file *file.File) (string, error) {
	resource, err := demo.formatter.ToStruct(file)
	if err != nil {
		return "", err
	}

	return *resource.(openapi.Demo).Spec.Id, nil
}

func (demo demoActions) Apply(ctx context.Context, fileContent file.File) (result *file.File, err error) {
	var demoResource openapi.Demo
	mapstructure.Decode(fileContent.Definition().Spec, &demoResource.Spec)

	if demoResource.Spec.Id == nil || *demoResource.Spec.Id == "" {
		result, err = demo.resourceClient.Create(ctx, fileContent)
		return result, err
	}

	result, err = demo.resourceClient.Update(ctx, fileContent, *demoResource.Spec.Id)
	return result, err
}

func (demo demoActions) List(ctx context.Context, listArgs utils.ListArgs) (*file.File, error) {
	return demo.resourceClient.List(ctx, listArgs)
}

func (demo demoActions) Delete(ctx context.Context, ID string) (string, error) {
	err := demo.resourceClient.Delete(ctx, ID)
	return "Demo successfully deleted", err
}

func (demo demoActions) Get(ctx context.Context, ID string) (*file.File, error) {
	return demo.resourceClient.Get(ctx, ID)
}
