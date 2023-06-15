package actions

import (
	"context"

	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/kubeshop/tracetest/server/model/yaml"
)

type analyzerActions struct {
	resourceArgs
}

var _ ResourceActions = &analyzerActions{}

func NewAnalyzerActions(options ...ResourceArgsOption) analyzerActions {
	args := NewResourceArgs(options...)

	return analyzerActions{
		resourceArgs: args,
	}
}

func (analyzerActions) FileType() yaml.FileType {
	return yaml.FileTypeAnalyzer
}

func (analyzerActions) Name() string {
	return "analyzer"
}

func (analyzer analyzerActions) GetID(file *file.File) (string, error) {
	resource, err := analyzer.formatter.ToStruct(file)
	if err != nil {
		return "", err
	}

	return *resource.(openapi.LinterResource).Spec.Id, nil
}

func (analyzer analyzerActions) Apply(ctx context.Context, fileContent file.File) (result *file.File, err error) {
	result, err = analyzer.resourceClient.Update(ctx, fileContent, currentConfigID)
	return result, err
}

func (analyzer analyzerActions) Get(ctx context.Context, ID string) (*file.File, error) {
	return analyzer.resourceClient.Get(ctx, currentConfigID)
}

func (analyzer analyzerActions) List(ctx context.Context, listArgs utils.ListArgs) (*file.File, error) {
	return analyzer.resourceClient.List(ctx, listArgs)
}

func (analyzer analyzerActions) Delete(ctx context.Context, ID string) (string, error) {
	return "", ErrNotSupportedResourceAction
}
