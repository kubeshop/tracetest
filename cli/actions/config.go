package actions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fluidtruck/deepcopy"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"go.uber.org/zap"
)

type configActions = resourceArgs

var _ ResourceActions = &configActions{}

func NewConfigActions(options ...resourceArgsOption) configActions {
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

	config.logger.Debug(
		"applying analytics config",
		zap.String("file", args.File),
	)

	fileContent, err := file.Read(args.File)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	if fileContent.Definition().Type != "Config" {
		return fmt.Errorf(`file must be of type "Config"`)
	}

	var configurationResource openapi.ConfigurationResource
	deepcopy.DeepCopy(fileContent.Definition().Spec, &configurationResource)

	request := config.client.ResourceApiApi.UpdateConfiguration(ctx, "current")
	request = request.ConfigurationResource(configurationResource)

	_, res, err := config.client.ResourceApiApi.UpdateConfigurationExecute(request)

	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("could not create config: %s", res.Status)
	}

	return err
}

func (config configActions) List(ctx context.Context) error {
	return nil
}

func (config configActions) Export(ctx context.Context, ID string) error {
	return nil
}

func (config configActions) Delete(ctx context.Context, ID string) error {
	return nil
}
