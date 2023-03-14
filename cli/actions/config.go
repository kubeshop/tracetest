package actions

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/utils"
)

type configActions struct {
	resourceArgs
}

var _ ResourceActions = &configActions{}
var currentConfigID = "current"

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

	fileContent, err := file.ReadRaw(args.File)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	if fileContent.Definition().Type != "Config" {
		return fmt.Errorf(`file must be of type "Config"`)
	}

	url := fmt.Sprintf("%s/%s", config.resourceClient.BaseUrl, currentConfigID)
	request, err := config.resourceClient.NewRequest(url, http.MethodPut, fileContent.Contents())
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	resp, err := config.resourceClient.Client.Do(request)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not create config: %s", resp.Status)
	}

	_, err = fileContent.SaveChanges(utils.IOReadCloserToString(resp.Body))
	return err
}

func (config configActions) Get(ctx context.Context, ID string) error {
	configResponse, err := config.get(ctx)
	if err != nil {
		return err
	}

	fmt.Println(configResponse)
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

	file, err := file.NewFromRaw(filePath, []byte(configResponse))
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}

	_, err = file.WriteRaw()
	return err
}

func (config configActions) Delete(ctx context.Context, ID string) error {
	return ErrNotSupportedResourceAction
}

func (config configActions) get(ctx context.Context) (string, error) {
	url := fmt.Sprintf("%s/%s", config.resourceClient.BaseUrl, currentConfigID)
	request, err := config.resourceClient.NewRequest(url, http.MethodGet, "")
	if err != nil {
		return "", fmt.Errorf("could not create request: %w", err)
	}

	resp, err := config.resourceClient.Client.Do(request)
	if err != nil {
		return "", fmt.Errorf("could not send request: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		validationError := string(body)
		return "", fmt.Errorf("invalid config: %s", validationError)
	}

	if err != nil {
		return "", fmt.Errorf("could not get config: %w", err)
	}

	return utils.IOReadCloserToString(resp.Body), nil
}
