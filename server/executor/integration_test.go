package executor_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/app"
	"github.com/kubeshop/tracetest/model"
	"github.com/kubeshop/tracetest/openapi"
	"github.com/kubeshop/tracetest/testmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const appEndpoint = "http://localhost:8080"

func TestExecutorIntegration(t *testing.T) {
	demoApp, err := testmock.GetDemoApplicationInstance()
	require.NoError(t, err)
	defer demoApp.Stop()

	tracetestApp, err := testmock.GetTestingApp(demoApp)
	require.NoError(t, err)

	go tracetestApp.Start()

	time.Sleep(1 * time.Second)

	t.Run("HappyPath", func(t *testing.T) {
		happyPath(t, tracetestApp, demoApp)
	})

}

func happyPath(t *testing.T, app *app.App, demoApp *testmock.DemoApp) {
	testID, err := createImportPokemonTest(app, demoApp)
	assert.NoError(t, err)
	assert.NotEmpty(t, testID)

	createTestDefinition(testID)

	runID, err := runTest(app, testID)
	assert.NoError(t, err)
	assert.NotEmpty(t, runID)

	run := waitForRunState(app, testID, runID, string(model.RunStateFinished), 30*time.Second)
	require.NotNil(t, run)

	assert.Equal(t, model.RunStateFinished, run.State)
	assert.Greater(t, len(run.Result.Results), 0)
	assert.True(t, run.Result.AllPassed)

	count := 0
	for _, res := range run.Result.Results {
		for _, assertionResult := range res.Results {
			for _, spanRes := range assertionResult.SpanResults {
				assert.True(t, spanRes.Passed)
				count = count + 1
			}
		}
	}

	assert.Equal(t, 2, count)
}

func createImportPokemonTest(app *app.App, demoApp *testmock.DemoApp) (string, error) {
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
		return "", fmt.Errorf("could not convert body into json: %w", err)
	}
	bytesBuffer := bytes.NewBuffer(jsonBytes)
	url := fmt.Sprintf("%s/api/tests", appEndpoint)

	response, err := http.Post(url, "application/json", bytesBuffer)
	if err != nil {
		return "", fmt.Errorf("could not create test: %w", err)
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("could not create test. Expected status 200, got %d", response.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("could not read response body: %w", err)
	}

	test := openapi.Test{}
	err = json.Unmarshal(bodyBytes, &test)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return test.Id, nil
}

func createTestDefinition(testID string) {
	body := openapi.TestDefinition{
		Definitions: []openapi.TestDefinitionDefinitions{
			{
				Selector: `span[service.name="pokeshop" tracetest.span.type="http" name="POST /pokemon/import"]`,
				Assertions: []openapi.Assertion{
					{
						Attribute:  "http.status_code",
						Comparator: "=",
						Expected:   "200",
					},
					{
						Attribute:  "tracetest.span.duration",
						Comparator: "<",
						Expected:   "100",
					},
				},
			},
		},
	}
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		panic(fmt.Errorf("could not convert body into json: %w", err))
	}
	bytesBuffer := bytes.NewBuffer(jsonBytes)
	url := fmt.Sprintf("%s/api/tests/%s/definition", appEndpoint, testID)

	req, _ := http.NewRequest("PUT", url, bytesBuffer)
	req.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Errorf("could not send request to create assertion: %w", err))
	}

	if response.StatusCode != 204 {
		r, _ := ioutil.ReadAll(response.Body)
		fmt.Println("Response", string(r))
		panic(fmt.Errorf("could not create definition. Expected status 201, got %d", response.StatusCode))
	}
}

func runTest(app *app.App, testID string) (string, error) {
	url := fmt.Sprintf("%s/api/tests/%s/run", appEndpoint, testID)
	response, err := http.Post(url, "application/json", nil)
	if err != nil {
		return "", fmt.Errorf("could not send request to run test: %w", err)
	}

	if response.StatusCode != 200 {
		return "", fmt.Errorf("could not run test. Expected status 200, got %d", response.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("could not read response body: %w", err)
	}

	testRun := openapi.TestRun{}
	err = json.Unmarshal(bodyBytes, &testRun)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return testRun.Id, nil
}

func waitForRunState(app *app.App, testID, resultID, state string, timeout time.Duration) *openapi.TestRun {
	timeoutTicker := time.NewTicker(timeout)
	executionTicker := time.NewTicker(1 * time.Second)
	outputChannel := make(chan *openapi.TestRun, 1)
	go func() {
		for {
			select {
			case <-timeoutTicker.C:
				outputChannel <- nil
				return
			case <-executionTicker.C:
				testRun := getRunInState(app, testID, resultID, state)
				if testRun != nil {
					outputChannel <- testRun
					return
				}
			}
		}
	}()

	testRun := <-outputChannel
	return testRun
}

func getRunInState(app *app.App, testID, resultID, state string) *openapi.TestRun {
	run, err := getRun(app, testID, resultID)
	if err != nil {
		return nil
	}

	if run.State != state {
		return nil

	}

	return &run
}

func getRun(app *app.App, testID, runID string) (openapi.TestRun, error) {
	url := fmt.Sprintf("%s/api/tests/%s/run/%s", appEndpoint, testID, runID)
	response, err := http.Get(url)
	if err != nil {
		return openapi.TestRun{}, fmt.Errorf("could not send request to run test: %w", err)
	}

	if response.StatusCode != 200 {
		return openapi.TestRun{}, fmt.Errorf("could not run test. Expected status 200, got %d", response.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return openapi.TestRun{}, fmt.Errorf("could not read response body: %w", err)
	}

	run := openapi.TestRun{}
	err = json.Unmarshal(bodyBytes, &run)
	if err != nil {
		return openapi.TestRun{}, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return run, nil
}
