package actions

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

type pollingActions struct {
	resourceArgs
}

var _ ResourceActions = &pollingActions{}

func NewPollingActions(options ...ResourceArgsOption) pollingActions {
	args := NewResourceArgs(options...)

	return pollingActions{
		resourceArgs: args,
	}
}

func (polling pollingActions) Apply(ctx context.Context, args ApplyArgs) error {
	if args.File == "" {
		return fmt.Errorf("you must specify a file to be applied")
	}

	polling.logger.Debug(
		"applying analytics config",
		zap.String("file", args.File),
	)

	fileContent, err := file.ReadRaw(args.File)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	if fileContent.Definition().Type != "PollingProfile" {
		return fmt.Errorf(`file must be of type "Config"`)
	}

	var pollingProfile openapi.PollingProfile
	mapstructure.Decode(fileContent.Definition().Spec, &pollingProfile.Spec)

	if pollingProfile.Spec.Id == "" {
		return polling.create(ctx, fileContent)
	}

	return polling.update(ctx, fileContent, pollingProfile.Spec.Id)
}

func (polling pollingActions) create(ctx context.Context, file file.File) error {
	request, err := polling.resourceClient.GetRequest(polling.resourceClient.BaseUrl, http.MethodPost, file.Contents())
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	resp, err := polling.resourceClient.Client.Do(request)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	if resp.StatusCode == http.StatusUnprocessableEntity {
		// validation error
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		validationError := string(body)
		return fmt.Errorf("invalid polling profile: %s", validationError)
	}
	if err != nil {
		return fmt.Errorf("could not create polling profile: %w", err)
	}

	_, err = file.SaveChanges(utils.IOReadCloserToString(resp.Body))
	return err
}

func (polling pollingActions) update(ctx context.Context, file file.File, ID string) error {
	url := fmt.Sprintf("%s/%s", polling.resourceClient.BaseUrl, ID)
	request, err := polling.resourceClient.GetRequest(url, http.MethodPut, file.Contents())
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	resp, err := polling.resourceClient.Client.Do(request)

	if err != nil {
		return fmt.Errorf("could not update polling profile: %w", err)
	}
	if resp.StatusCode == http.StatusUnprocessableEntity {
		// validation error
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		validationError := string(body)
		return fmt.Errorf("invalid polling profile: %s", validationError)
	}

	_, err = file.SaveChanges(utils.IOReadCloserToString(resp.Body))
	return err
}

func (polling pollingActions) List(ctx context.Context, listArgs ListArgs) error {
	url := fmt.Sprintf("%s?skip=%d&take=%d&sortBy=%s&sortDirection=%s", polling.resourceClient.BaseUrl, listArgs.Skip, listArgs.Take, listArgs.SortBy, listArgs.SortDirection)
	request, err := polling.resourceClient.GetRequest(url, http.MethodGet, "")
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	resp, err := polling.resourceClient.Client.Do(request)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	fmt.Println(utils.IOReadCloserToString(resp.Body))
	return nil
}

func (polling pollingActions) Export(ctx context.Context, ID string, filePath string) error {
	pollingProfile, err := polling.get(ctx, ID)
	if err != nil {
		return err
	}

	file, err := file.NewFromRaw(filePath, []byte(pollingProfile))
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}

	_, err = file.WriteRaw()
	return err
}

func (polling pollingActions) Delete(ctx context.Context, ID string) error {
	url := fmt.Sprintf("%s/%s", polling.resourceClient.BaseUrl, ID)
	request, err := polling.resourceClient.GetRequest(url, http.MethodDelete, "")
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	_, err = polling.resourceClient.Client.Do(request)
	return err
}

func (polling pollingActions) Get(ctx context.Context, ID string) error {
	pollingProfileResponse, err := polling.get(ctx, ID)
	if err != nil {
		return err
	}

	fmt.Println(pollingProfileResponse)
	return err
}

func (polling pollingActions) get(ctx context.Context, ID string) (string, error) {
	request, err := polling.resourceClient.GetRequest(fmt.Sprintf("%s/%s", polling.resourceClient.BaseUrl, ID), http.MethodGet, "")
	if err != nil {
		return "", fmt.Errorf("could not create request: %w", err)
	}

	resp, err := polling.resourceClient.Client.Do(request)
	if err != nil {
		return "", fmt.Errorf("could not get polling profile: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		validationError := string(body)
		return "", fmt.Errorf("invalid polling profile: %s", validationError)
	}

	return utils.IOReadCloserToString(resp.Body), nil
}
