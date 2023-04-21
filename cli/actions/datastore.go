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

type dataStoreActions struct {
	resourceArgs
}

var _ ResourceActions = &dataStoreActions{}

func NewDataStoreActions(options ...ResourceArgsOption) *dataStoreActions {
	args := NewResourceArgs(options...)

	return &dataStoreActions{
		resourceArgs: args,
	}
}

func (d *dataStoreActions) Name() string {
	return "datastore"
}

func (d *dataStoreActions) Apply(ctx context.Context, args ApplyArgs) error {
	if args.File == "" {
		return fmt.Errorf("you must specify a file to be applied")
	}

	d.logger.Debug(
		"applying analytics config",
		zap.String("file", args.File),
	)

	fileContent, err := file.ReadRaw(args.File)
	if err != nil {
		return fmt.Errorf("could not read file: %w", err)
	}

	if fileContent.Definition().Type != "DataStore" {
		return fmt.Errorf(`file must be of type "DataStore"`)
	}

	var dataStore openapi.DataStore
	mapstructure.Decode(fileContent.Definition().Spec, &dataStore)

	return d.update(ctx, fileContent, currentConfigID)
}

func (d *dataStoreActions) List(ctx context.Context, args ListArgs) error {
	return ErrNotSupportedResourceAction
}

func (d *dataStoreActions) Get(ctx context.Context, id string) error {
	dataStoreResponse, err := d.get(ctx)
	if err != nil {
		return err
	}

	fmt.Println(dataStoreResponse)
	return nil
}

func (d *dataStoreActions) Export(ctx context.Context, id string, filePath string) error {
	dataStore, err := d.get(ctx)
	if err != nil {
		return err
	}

	file, err := file.NewFromRaw(filePath, []byte(dataStore))
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}

	_, err = file.WriteRaw()
	return err
}

func (d *dataStoreActions) Delete(ctx context.Context, id string) error {
	return nil
}

func (d *dataStoreActions) get(ctx context.Context) (string, error) {
	request, err := d.resourceClient.NewRequest(fmt.Sprintf("%s/%s", d.resourceClient.BaseUrl, currentConfigID), http.MethodGet, "")
	if err != nil {
		return "", fmt.Errorf("could not create request: %w", err)
	}

	resp, err := d.resourceClient.Client.Do(request)
	if err != nil {
		return "", fmt.Errorf("could not get data store: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		validationError := string(body)
		return "", fmt.Errorf("invalid data store: %s", validationError)
	}

	return utils.IOReadCloserToString(resp.Body), nil
}

func (d *dataStoreActions) update(ctx context.Context, file file.File, ID string) error {
	url := fmt.Sprintf("%s/%s", d.resourceClient.BaseUrl, ID)
	request, err := d.resourceClient.NewRequest(url, http.MethodPut, file.Contents())
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	resp, err := d.resourceClient.Client.Do(request)
	if err != nil {
		return fmt.Errorf("could not update data store: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("datastore id doesn't exist on server. Remove it from the definition file and try again")
	}

	if resp.StatusCode == http.StatusUnprocessableEntity {
		// validation error
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("could not send request: %w", err)
		}

		validationError := string(body)
		return fmt.Errorf("invalid datastore: %s", validationError)
	}

	_, err = file.SaveChanges(utils.IOReadCloserToString(resp.Body))
	return err
}
