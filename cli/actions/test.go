package actions

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/kubeshop/tracetest/server/model/yaml"
)

type testActions struct {
	resourceArgs

	openapiClient *openapi.APIClient
}

var _ ResourceActions = &testActions{}

func NewTestActions(openapiClient *openapi.APIClient, options ...ResourceArgsOption) testActions {
	args := NewResourceArgs(options...)

	return testActions{
		resourceArgs:  args,
		openapiClient: openapiClient,
	}
}

// Apply implements ResourceActions
func (actions testActions) Apply(ctx context.Context, inputFile file.File) (*file.File, error) {
	return nil, fmt.Errorf("not implemented")
}

// Delete implements ResourceActions
func (actions testActions) Delete(ctx context.Context, id string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

// FileType implements ResourceActions
func (actions testActions) FileType() yaml.FileType {
	return yaml.FileTypeTest
}

// Get implements ResourceActions
func (actions testActions) Get(ctx context.Context, id string) (*file.File, error) {
	return nil, fmt.Errorf("not implemented")
}

// GetID implements ResourceActions
func (actions testActions) GetID(file *file.File) (string, error) {
	test, err := actions.formatter.ToStruct(file)
	if err != nil {
		return "", fmt.Errorf("could not convert file into struct: %w", err)
	}

	return test.(openapi.TestResource).Spec.GetId(), nil
}

// List implements ResourceActions
func (actions testActions) List(ctx context.Context, args utils.ListArgs) (*file.File, error) {
	ctx = context.WithValue(ctx, "X-Tracetest-Augmented", true)
	return actions.resourceClient.List(ctx, args)
}

// Name implements ResourceActions
func (actions testActions) Name() string {
	return "test"
}
