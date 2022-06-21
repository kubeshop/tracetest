package testfixtures

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kubeshop/tracetest/server/app"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/testmock"
)

const appEndpoint = "http://localhost:8080"

func init() {
	RegisterFixture(IMPORT_POKEMON_TEST, getImportPokemonTest)
}

func GetPokemonTest(options ...Option) (*openapi.Test, error) {
	return GetFixtureValue[*openapi.Test](IMPORT_POKEMON_TEST, options...)
}

func getImportPokemonTest(options FixtureOptions) (*openapi.Test, error) {
	pokeshopApp, err := GetPokeshopApp()
	if err != nil {
		return nil, fmt.Errorf("could not get pokeshop app: %w", err)
	}

	tracetestApp, err := GetTracetestApp()
	if err != nil {
		return nil, fmt.Errorf("could not get tracetest app: %w", err)
	}

	test, err := createImportPokemonTest(tracetestApp, pokeshopApp)
	if err != nil {
		return nil, fmt.Errorf("could not create import pokemon test: %w", err)
	}

	if test.Id == "" {
		return nil, fmt.Errorf("testID cannot be empty")
	}

	err = createTestDefinition(test.Id)
	if err != nil {
		return nil, fmt.Errorf("could not create test definition: %w", err)
	}

	updatedTest, err := getTest(test.Id)
	if err != nil {
		return nil, fmt.Errorf("could not get test: %w", err)
	}

	return &updatedTest, nil
}

func createImportPokemonTest(app *app.App, demoApp *testmock.DemoApp) (openapi.Test, error) {
	body := openapi.Test{
		Name:        "Import Pokemon",
		Description: "Import a pokemon into the api",
		ServiceUnderTest: openapi.TestServiceUnderTest{
			Request: openapi.HttpRequest{
				Url:    fmt.Sprintf("http://%s/pokemon/import", demoApp.Endpoint()),
				Method: "POST",
				Headers: []openapi.HttpHeader{
					{
						Key:   "Content-Type",
						Value: "application/json",
					},
				},
				Body: `{ "id": 52 }`,
			},
		},
	}
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not convert body into json: %w", err)
	}
	bytesBuffer := bytes.NewBuffer(jsonBytes)
	url := fmt.Sprintf("%s/api/tests", appEndpoint)

	response, err := http.Post(url, "application/json", bytesBuffer)
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not create test: %w", err)
	}

	if response.StatusCode != 200 {
		return openapi.Test{}, fmt.Errorf("could not create test. Expected status 200, got %d", response.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not read response body: %w", err)
	}

	test := openapi.Test{}
	err = json.Unmarshal(bodyBytes, &test)
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return test, nil
}

func createTestDefinition(testID string) error {
	body := openapi.TestDefinition{
		Definitions: []openapi.TestDefinitionDefinitions{
			{
				Selector: openapi.Selector{
					Query: `span[service.name="pokeshop" tracetest.span.type="http" name="POST /pokemon/import"]`,
				},
				Assertions: []openapi.Assertion{
					{
						Attribute:  "http.status_code",
						Comparator: "=",
						Expected:   "200",
					},
					{
						Attribute:  "tracetest.span.duration",
						Comparator: "<",
						Expected:   "200",
					},
				},
			},
		},
	}
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("could not convert body into json: %w", err)
	}
	bytesBuffer := bytes.NewBuffer(jsonBytes)
	url := fmt.Sprintf("%s/api/tests/%s/definition", appEndpoint, testID)

	req, _ := http.NewRequest("PUT", url, bytesBuffer)
	req.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not send request to create assertion: %w", err)
	}

	if response.StatusCode != 204 {
		return fmt.Errorf("could not create definition. Expected status 201, got %d", response.StatusCode)
	}

	return nil
}

func getTest(testID string) (openapi.Test, error) {
	url := fmt.Sprintf("%s/api/tests/%s", appEndpoint, testID)

	response, err := http.Get(url)
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not send request to create assertion: %w", err)
	}

	if response.StatusCode != 200 {
		return openapi.Test{}, fmt.Errorf("could not create definition. Expected status 201, got %d", response.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not read response body: %w", err)
	}

	test := openapi.Test{}
	err = json.Unmarshal(bodyBytes, &test)
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return test, nil
}
