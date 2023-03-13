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
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
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

	fileContent, err := file.Read(args.File)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	if fileContent.Definition().Type != "PollingProfile" {
		return fmt.Errorf(`file must be of type "Config"`)
	}

	var pollingProfile openapi.PollingProfile
	deepcopy.DeepCopy(fileContent.Definition(), &pollingProfile)
	deepcopy.DeepCopy(fileContent.Definition().Spec, &pollingProfile.Spec)

	if pollingProfile.Spec.Id == "" {
		err := polling.create(ctx, fileContent, pollingProfile)
		return err
	} else {
		err := polling.update(ctx, fileContent, pollingProfile)
		return err
	}
}

func (polling pollingActions) create(ctx context.Context, file file.File, pollingProfile openapi.PollingProfile) error {
	request := polling.client.ResourceApiApi.CreatePollingProfile(ctx)
	request = request.PollingProfile(pollingProfile)
	createdPollingProfile, resp, err := polling.client.ResourceApiApi.CreatePollingProfileExecute(request)

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

	f, err := file.SetID(createdPollingProfile.Spec.Id)
	if err != nil {
		return fmt.Errorf("could no set polling profile id: %w", err)
	}

	_, err = f.Write()
	if err != nil {
		return fmt.Errorf("could not write to polling profile file: %w", err)
	}

	return nil
}

func (polling pollingActions) update(ctx context.Context, file file.File, pollingProfile openapi.PollingProfile) error {
	request := polling.client.ResourceApiApi.UpdatePollingProfile(ctx, pollingProfile.Spec.Id)
	request = request.PollingProfile(pollingProfile)
	_, resp, err := polling.client.ResourceApiApi.UpdatePollingProfileExecute(request)

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
		return fmt.Errorf("could not update polling profile: %w", err)
	}

	return nil
}

func (polling pollingActions) List(ctx context.Context, listArgs ListArgs) error {
	request := polling.client.ResourceApiApi.ListPollingProfiles(ctx).Skip(listArgs.Skip).Take(listArgs.Take).SortBy(listArgs.SortBy).SortDirection(listArgs.SortBy)
	pollingProfileList, _, err := polling.client.ResourceApiApi.ListPollingProfilesExecute(request)

	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	formatter := formatters.PollingProfileList(polling.config)
	output := formatter.Format(pollingProfileList.Items)
	fmt.Println(output)

	return nil
}

func (polling pollingActions) Export(ctx context.Context, ID string, filePath string) error {
	pollingProfileResponse, err := polling.get(ctx, ID)
	if err != nil {
		return err
	}

	yamlData, err := yaml.Marshal(&pollingProfileResponse)
	if err != nil {
		return fmt.Errorf("could not marshal polling profile: %w", err)
	}

	file, err := file.New(filePath, []byte(yamlData))
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}

	_, err = file.Write()
	return err
}

func (polling pollingActions) Delete(ctx context.Context, ID string) error {
	request := polling.client.ResourceApiApi.DeletePollingProfile(ctx, ID)
	_, err := polling.client.ResourceApiApi.DeletePollingProfileExecute(request)

	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}

	return nil
}

func (polling pollingActions) Get(ctx context.Context, ID string) error {
	pollingProfileResponse, err := polling.get(ctx, ID)
	if err != nil {
		return err
	}

	formatter := formatters.PollingProfileFormatter(polling.config)
	fmt.Println(formatter.Format(*pollingProfileResponse))

	return err
}

func (polling pollingActions) get(ctx context.Context, ID string) (*openapi.PollingProfile, error) {
	request := polling.client.ResourceApiApi.GetPollingProfile(ctx, ID)
	pollingProfileResponse, resp, err := polling.client.ResourceApiApi.GetPollingProfileExecute(request)

	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		validationError := string(body)
		return nil, fmt.Errorf("invalid polling profile: %s", validationError)
	}

	if err != nil {
		return nil, fmt.Errorf("could not get polling profile: %w", err)
	}

	return pollingProfileResponse, nil
}
