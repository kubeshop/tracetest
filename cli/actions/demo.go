package actions

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

var _ ResourceActions = &demoActions{}

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

func (demo demoActions) Apply(ctx context.Context, fileContent file.File) error {
	var demoResource openapi.Demo
	mapstructure.Decode(fileContent.Definition().Spec, &demoResource.Spec)

	if demoResource.Spec.Id == nil || *demoResource.Spec.Id == "" {
		return demo.resourceClient.Create(ctx, fileContent)
	}

	return demo.resourceClient.Update(ctx, fileContent, *demoResource.Spec.Id)
}

func (demo demoActions) List(ctx context.Context, listArgs utils.ListArgs) (string, error) {
	return demo.resourceClient.List(ctx, listArgs)
}

func (demo demoActions) Delete(ctx context.Context, ID string) error {
	return demo.resourceClient.Delete(ctx, ID)
}

func (demo demoActions) Get(ctx context.Context, ID string) (string, error) {
	return demo.resourceClient.Get(ctx, ID)
}
