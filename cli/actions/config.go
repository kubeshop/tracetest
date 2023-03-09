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
)

type configActions = resourceArgs

var _ ResourceActions = &configActions{}

func NewConfigActions(options ...ResourceArgsOption) configActions {
	cfgActions := configActions{}

	for _, option := range options {
		option(&cfgActions)
	}

	return cfgActions
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
	request := config.client.ResourceApiApi.GetConfiguration(ctx, "current")
	configResponse, resp, err := config.client.ResourceApiApi.GetConfigurationExecute(request)

	if err != nil && resp == nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		validationError := string(body)
		return fmt.Errorf("invalid data store: %s", validationError)
	}

	if err != nil {
		return fmt.Errorf("could not create data store: %w", err)
	}

	formatter := formatters.ConfigFormatter(config.config)
	fmt.Println(formatter.Format(*configResponse))

	return err
}

func (config configActions) List(ctx context.Context) error {
	return ErrNotSupportedResourceAction
}

func (config configActions) Export(ctx context.Context, ID string) error {
	return ErrNotSupportedResourceAction
}

func (config configActions) Delete(ctx context.Context, ID string) error {
	return ErrNotSupportedResourceAction
}
