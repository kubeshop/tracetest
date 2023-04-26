package actions

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/utils"
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

func (environment environmentsActions) Apply(ctx context.Context, fileContent file.File) (*file.File, error) {
	var envResource openapi.EnvironmentResource
	mapstructure.Decode(fileContent.Definition().Spec, &envResource.Spec)

	if envResource.Spec.Id == nil || *envResource.Spec.Id == "" {
		return environment.resourceClient.Create(ctx, fileContent)
	}

	return environment.resourceClient.Update(ctx, fileContent, *envResource.Spec.Id)
}

func (environment environmentsActions) List(ctx context.Context, listArgs utils.ListArgs) (*file.File, error) {
	return environment.resourceClient.List(ctx, listArgs)
}

func (environment environmentsActions) Delete(ctx context.Context, ID string) error {
	return environment.resourceClient.Delete(ctx, ID)
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

func (environment environmentsActions) FromFile(ctx context.Context, filePath string) (openapi.EnvironmentResource, error) {
	if !utils.StringReferencesFile(filePath) {
		return openapi.EnvironmentResource{}, fmt.Errorf(`env file "%s" does not exist`, filePath)
	}

	envVars, err := godotenv.Read(filePath)
	if err != nil {
		return openapi.EnvironmentResource{}, fmt.Errorf(`cannot read env file "%s": %w`, filePath, err)
	}

	values := make([]openapi.EnvironmentValue, 0, len(envVars))
	for k, v := range envVars {
		values = append(values, openapi.EnvironmentValue{
			Key:   openapi.PtrString(k),
			Value: openapi.PtrString(v),
		})
	}

	name := filepath.Base(filePath)

	env := openapi.Environment{
		Id:     &name,
		Name:   &name,
		Values: values,
	}
	envType := "Environment"

	resource := openapi.EnvironmentResource{
		Type: &envType,
		Spec: &env,
	}

	return resource, nil
}
