package actions

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"github.com/mitchellh/mapstructure"
)

type environmentsActions struct {
	resourceArgs
}

var _ ResourceActions = &environmentsActions{}

func NewEnvironmentsActions(options ...ResourceArgsOption) environmentsActions {
	args := NewResourceArgs(options...)

	return environmentsActions{
		resourceArgs: args,
	}
}

func (environmentsActions) FileType() yaml.FileType {
	return yaml.FileTypeEnvironment
}

func (environmentsActions) Name() string {
	return "environment"
}

func (action environmentsActions) GetID(file *file.File) (string, error) {
	resource, err := action.formatter.ToStruct(file)
	if err != nil {
		return "", err
	}

	return resource.(openapi.EnvironmentResource).Spec.GetId(), nil
}

func (environment environmentsActions) Apply(ctx context.Context, fileContent file.File) (result *file.File, err error) {
	envResource := openapi.EnvironmentResource{
		Spec: &openapi.Environment{},
	}

	mapstructure.Decode(fileContent.Definition().Spec, &envResource.Spec)

	return environment.resourceClient.Upsert(ctx, fileContent)
}

func (environment environmentsActions) List(ctx context.Context, listArgs utils.ListArgs) (*file.File, error) {
	return environment.resourceClient.List(ctx, listArgs)
}

func (environment environmentsActions) Delete(ctx context.Context, ID string) (string, error) {
	return "Environment successfully deleted", environment.resourceClient.Delete(ctx, ID)
}

func (environment environmentsActions) Get(ctx context.Context, ID string) (*file.File, error) {
	return environment.resourceClient.Get(ctx, ID)
}

func (environment environmentsActions) ApplyResource(ctx context.Context, resource openapi.EnvironmentResource) (*file.File, error) {
	content, err := resource.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("could not marshal environment: %w", err)
	}

	file, err := file.NewFromRaw("env.yaml", content)
	if err != nil {
		return nil, fmt.Errorf("could not create environment file: %w", err)
	}

	return environment.Apply(ctx, file)
}

func (action environmentsActions) FromFile(ctx context.Context, filePath string) (openapi.EnvironmentResource, error) {
	if !utils.StringReferencesFile(filePath) {
		return openapi.EnvironmentResource{}, fmt.Errorf(`env file "%s" does not exist`, filePath)
	}

	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return openapi.EnvironmentResource{}, fmt.Errorf(`cannot read env file "%s": %w`, filePath, err)
	}

	model, err := yaml.Decode(fileContent)
	if err != nil {
		return openapi.EnvironmentResource{}, fmt.Errorf(`cannot parse env file "%s": %w`, filePath, err)
	}

	envModel := model.Spec.(environment.Environment)

	values := make([]openapi.EnvironmentValue, 0, len(envModel.Values))
	for _, value := range envModel.Values {
		v := value
		values = append(values, openapi.EnvironmentValue{
			Key:   &v.Key,
			Value: &v.Value,
		})
	}

	environmentResource := openapi.EnvironmentResource{
		Type: (*string)(&model.Type),
		Spec: &openapi.Environment{
			Id:          openapi.PtrString(envModel.ID.String()),
			Name:        &envModel.Name,
			Description: &envModel.Description,
			Values:      values,
		},
	}

	return environmentResource, nil
}
