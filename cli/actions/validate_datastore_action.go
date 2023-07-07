package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"github.com/kubeshop/tracetest/cli/pkg/termutil"

	"github.com/goccy/go-yaml"
)

const (
	PASSED_STEP_ICON = "✔"
	FAILED_STEP_ICON = "✘"
)

type validateDataStoreAction struct {
	openAPIClient *openapi.APIClient
}

func NewValidateDataStoreAction(openAPIClient *openapi.APIClient) *validateDataStoreAction {
	return &validateDataStoreAction{
		openAPIClient: openAPIClient,
	}
}

func (a *validateDataStoreAction) ValidateDatastore(dataStoreFile string) (string, error) {
	f, err := fileutil.Read(dataStoreFile)
	if err != nil {
		return "", fmt.Errorf("cannot read datastore file %s: %w", dataStoreFile, err)
	}

	dataStore, err := a.mapDataStoreFileToOpenAPI(f)
	if err != nil {
		return "", fmt.Errorf("cannot map datastore to validate connection: %w", err)
	}

	ctx := context.Background()
	request := a.openAPIClient.ApiApi.TestConnection(ctx).DataStore(*dataStore)

	result, response, err := request.Execute()
	if err != nil && response.StatusCode != http.StatusUnprocessableEntity {
		return "", fmt.Errorf("cannot test connection with datastore file: %w", err)
	}
	if response.StatusCode == http.StatusUnprocessableEntity {
		// it could be a valid response, try to update result
		result, err = a.tryToRetrieveResultFromResponse(response)
		if err != nil {
			return "", err // it is something else, propagate the error
		}
	}

	return a.mapOpenAPIResponseToCLIOutput(result)
}

func (a *validateDataStoreAction) tryToRetrieveResultFromResponse(response *http.Response) (*openapi.ConnectionResult, error) {
	result := &openapi.ConnectionResult{}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read test connection result: %w", err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("could not parse test connection result: %w", err)
	}

	if !result.HasPortCheck() {
		// it is a validation error, return the same error
		return nil, fmt.Errorf("cannot test connection with datastore file: %w", err)
	}

	return result, nil
}

func (a *validateDataStoreAction) mapDataStoreFileToOpenAPI(f fileutil.File) (*openapi.DataStore, error) {
	var dataStoreResource openapi.DataStoreResource
	err := yaml.Unmarshal(f.Contents(), &dataStoreResource)
	if err != nil {
		return nil, fmt.Errorf("could not convert datastore file content to API format: %w", err)
	}

	return dataStoreResource.Spec, nil
}

func (a *validateDataStoreAction) mapOpenAPIResponseToCLIOutput(response *openapi.ConnectionResult) (string, error) {
	result := []string{
		a.printMessage("Port checking:", response.PortCheck),
		a.printMessage("Connectivity:", response.Connectivity),
		a.printMessage("Authentication:", response.Authentication),
		a.printMessage("Fetch traces:", response.FetchTraces),
	}

	return strings.Join(result, ""), nil
}

func (a *validateDataStoreAction) stringPointerToString(stringPointer *string) string {
	if stringPointer == nil {
		return ""
	}

	return *stringPointer
}

func (a *validateDataStoreAction) printMessage(topic string, step *openapi.ConnectionTestStep) string {
	if step == nil || step.Message == nil {
		return ""
	}

	passed := (step != nil && step.Passed != nil && *step.Passed)

	icon := PASSED_STEP_ICON
	paintedText := termutil.GetGreenText(topic)
	finalMessage := a.stringPointerToString(step.Message)

	if !passed {
		icon = FAILED_STEP_ICON
		paintedText = termutil.GetRedText(topic)
		finalMessage = fmt.Sprintf("%s - Error: %s", a.stringPointerToString(step.Message), a.stringPointerToString(step.Error))
	}

	lineBreak := fmt.Sprintln("")
	return fmt.Sprintf("%s %s %s%s", icon, paintedText, finalMessage, lineBreak)
}
