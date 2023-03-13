package actions

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/openapi"
	"gopkg.in/yaml.v3"
)

type configActions struct {
	resourceArgs
}

var _ ResourceActions = &configActions{}

func NewConfigActions(options ...ResourceArgsOption) configActions {
	args := NewResourceArgs(options...)

	return configActions{
		resourceArgs: args,
	}
}

func (config configActions) Apply(ctx context.Context, args ApplyArgs) error {
	if args.File == "" {
		return fmt.Errorf("you must specify a file to be applied")
	}

	fileContent, err := file.Read(args.File)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	if fileContent.Definition().Type != "Config" {
		return fmt.Errorf(`file must be of type "Config"`)
	}

	var configurationResource openapi.ConfigurationResource
	deepcopy.DeepCopy(fileContent.Definition(), &configurationResource)
	deepcopy.DeepCopy(fileContent.Definition().Spec, &configurationResource.Spec)

	request := config.client.ResourceApiApi.UpdateConfiguration(ctx, "current")
	request = request.ConfigurationResource(configurationResource)

	_, res, err := config.client.ResourceApiApi.UpdateConfigurationExecute(request)

	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("could not create config: %s", res.Status)
	}

	return err
}

func (config configActions) Get(ctx context.Context, ID string) error {
	configResponse, err := config.get(ctx)
	if err != nil {
		return err
	}

	formatter := formatters.ConfigFormatter(config.config)
	fmt.Println(formatter.Format(*configResponse))

	return err
}

func (config configActions) List(ctx context.Context, listArgs ListArgs) error {
	return ErrNotSupportedResourceAction
}

func (config configActions) Export(ctx context.Context, ID string, filePath string) error {
	configResponse, err := config.get(ctx)
	if err != nil {
		return err
	}

	yamlData, err := yaml.Marshal(&configResponse)
	if err != nil {
		return fmt.Errorf("could not marshal config: %w", err)
	}

	file, err := file.New(filePath, []byte(yamlData))
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}

	_, err = file.Write()
	return err
}

func (config configActions) Delete(ctx context.Context, ID string) error {
	return ErrNotSupportedResourceAction
}

func (config configActions) get(ctx context.Context) (*openapi.ConfigurationResource, error) {
	request := config.client.ResourceApiApi.GetConfiguration(ctx, "current")
	configResponse, resp, err := config.client.ResourceApiApi.GetConfigurationExecute(request)

	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		validationError := string(body)
		return nil, fmt.Errorf("invalid config: %s", validationError)
	}

	if err != nil {
		return nil, fmt.Errorf("could not get config: %w", err)
	}

	return configResponse, nil
}
