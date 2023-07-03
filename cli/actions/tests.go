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

type testActions struct {
	resourceArgs
}

var _ ResourceActions = &testActions{}

func NewTestActions(options ...ResourceArgsOption) testActions {
	args := NewResourceArgs(options...)

	return testActions{
		resourceArgs: args,
	}
}

func (testActions) FileType() yaml.FileType {
	return yaml.FileTypeEnvironment
}

func (testActions) Name() string {
	return "test"
}

func (action testActions) GetID(file *file.File) (string, error) {
	resource, err := action.formatter.ToStruct(file)
	if err != nil {
		return "", err
	}

	return resource.(openapi.TestResource).Spec.GetId(), nil
}

func (test testActions) Apply(ctx context.Context, fileContent file.File) (result *file.File, err error) {
	envResource := openapi.TestResource{
		Spec: &openapi.Test{},
	}

	mapstructure.Decode(fileContent.Definition().Spec, &envResource.Spec)

	return test.resourceClient.Upsert(ctx, fileContent)
}

func (test testActions) List(ctx context.Context, listArgs utils.ListArgs) (*file.File, error) {
	return test.resourceClient.List(ctx, listArgs)
}

func (test testActions) Delete(ctx context.Context, ID string) (string, error) {
	return "Test successfully deleted", test.resourceClient.Delete(ctx, ID)
}

func (test testActions) Get(ctx context.Context, ID string) (*file.File, error) {
	return test.resourceClient.Get(ctx, ID)
}

func (test testActions) ApplyResource(ctx context.Context, resource openapi.TestResource) (*file.File, error) {
	content, err := resource.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("could not marshal test: %w", err)
	}

	file, err := file.NewFromRaw("test.yaml", content)
	if err != nil {
		return nil, fmt.Errorf("could not create test file: %w", err)
	}

	return test.Apply(ctx, file)
}

// func (action testActions) FromFile(ctx context.Context, filePath string) (openapi.TestResource, error) {
// 	if !utils.StringReferencesFile(filePath) {
// 		return openapi.TestResource{}, fmt.Errorf(`env file "%s" does not exist`, filePath)
// 	}

// 	fileContent, err := ioutil.ReadFile(filePath)
// 	if err != nil {
// 		return openapi.TestResource{}, fmt.Errorf(`cannot read env file "%s": %w`, filePath, err)
// 	}

// 	model, err := yaml.Decode(fileContent)
// 	if err != nil {
// 		return openapi.TestResource{}, fmt.Errorf(`cannot parse env file "%s": %w`, filePath, err)
// 	}

// 	envModel := model.Spec.(test.)

// 	values := make([]openapi.EnvironmentValue, 0, len(envModel.Values))
// 	for _, value := range envModel.Values {
// 		v := value
// 		values = append(values, openapi.EnvironmentValue{
// 			Key:   &v.Key,
// 			Value: &v.Value,
// 		})
// 	}

// 	TestResource := openapi.TestResource{
// 		Type: (*string)(&model.Type),
// 		Spec: &openapi.Environment{
// 			Id:          openapi.PtrString(envModel.ID.String()),
// 			Name:        &envModel.Name,
// 			Description: &envModel.Description,
// 			Values:      values,
// 		},
// 	}

// 	return TestResource, nil
// }
