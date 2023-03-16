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

type demoActions struct {
	resourceArgs
}

var _ ResourceActions = &demoActions{}

func NewDemoActions(options ...ResourceArgsOption) demoActions {
	args := NewResourceArgs(options...)

	return demoActions{
		resourceArgs: args,
	}
}

func (demo demoActions) Apply(ctx context.Context, args ApplyArgs) error {
	if args.File == "" {
		return fmt.Errorf("you must specify a file to be applied")
	}

	demo.logger.Debug(
		"applying demo",
		zap.String("file", args.File),
	)

	fileContent, err := file.ReadRaw(args.File)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	if fileContent.Definition().Type != "Demo" {
		return fmt.Errorf(`file must be of type "Demo"`)
	}

	var demoResource openapi.Demo
	mapstructure.Decode(fileContent.Definition().Spec, &demoResource.Spec)

	if demoResource.Spec.Id == nil || *demoResource.Spec.Id == "" {
		return demo.create(ctx, fileContent)
	}

	return demo.update(ctx, fileContent, *demoResource.Spec.Id)
}

func (demo demoActions) create(ctx context.Context, file file.File) error {
	request, err := demo.resourceClient.NewRequest(demo.resourceClient.BaseUrl, http.MethodPost, file.Contents())
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	resp, err := demo.resourceClient.Client.Do(request)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnprocessableEntity {
		// validation error
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("could not read validation error: %w", err)
		}

		validationError := string(body)
		return fmt.Errorf("invalid demo profile: %s", validationError)
	}

	_, err = file.SaveChanges(utils.IOReadCloserToString(resp.Body))
	return err
}

func (demo demoActions) update(ctx context.Context, file file.File, ID string) error {
	url := fmt.Sprintf("%s/%s", demo.resourceClient.BaseUrl, ID)
	request, err := demo.resourceClient.NewRequest(url, http.MethodPut, file.Contents())
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	resp, err := demo.resourceClient.Client.Do(request)
	if err != nil {
		return fmt.Errorf("could not update demo profile: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnprocessableEntity {
		// validation error
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("could not send request: %w", err)
		}

		validationError := string(body)
		return fmt.Errorf("invalid demo profile: %s", validationError)
	}

	_, err = file.SaveChanges(utils.IOReadCloserToString(resp.Body))
	return err
}

func (demo demoActions) List(ctx context.Context, listArgs ListArgs) error {
	url := fmt.Sprintf("%s?skip=%d&take=%d&sortBy=%s&sortDirection=%s", demo.resourceClient.BaseUrl, listArgs.Skip, listArgs.Take, listArgs.SortBy, listArgs.SortDirection)
	request, err := demo.resourceClient.NewRequest(url, http.MethodGet, "")
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	resp, err := demo.resourceClient.Client.Do(request)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	defer resp.Body.Close()
	fmt.Println(utils.IOReadCloserToString(resp.Body))
	return nil
}

func (demo demoActions) Export(ctx context.Context, ID string, filePath string) error {
	demoProfile, err := demo.get(ctx, ID)
	if err != nil {
		return err
	}

	file, err := file.NewFromRaw(filePath, []byte(demoProfile))
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}

	_, err = file.WriteRaw()
	return err
}

func (demo demoActions) Delete(ctx context.Context, ID string) error {
	url := fmt.Sprintf("%s/%s", demo.resourceClient.BaseUrl, ID)
	request, err := demo.resourceClient.NewRequest(url, http.MethodDelete, "")
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	_, err = demo.resourceClient.Client.Do(request)
	return err
}

func (demo demoActions) Get(ctx context.Context, ID string) error {
	demoProfileResponse, err := demo.get(ctx, ID)
	if err != nil {
		return err
	}

	fmt.Println(demoProfileResponse)
	return err
}

func (demo demoActions) get(ctx context.Context, ID string) (string, error) {
	request, err := demo.resourceClient.NewRequest(fmt.Sprintf("%s/%s", demo.resourceClient.BaseUrl, ID), http.MethodGet, "")
	if err != nil {
		return "", fmt.Errorf("could not create request: %w", err)
	}

	resp, err := demo.resourceClient.Client.Do(request)
	if err != nil {
		return "", fmt.Errorf("could not get demo profile: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		validationError := string(body)
		return "", fmt.Errorf("invalid demo profile: %s", validationError)
	}

	return utils.IOReadCloserToString(resp.Body), nil
}
