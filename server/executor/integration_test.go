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
	"github.com/kubeshop/tracetest/executor"
	"github.com/kubeshop/tracetest/openapi"
	"github.com/kubeshop/tracetest/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const appEndpoint = "http://localhost:8080"

func TestExecutorIntegration(t *testing.T) {
	demoApp, err := test.GetDemoApplicationInstance()
	require.NoError(t, err)
	defer demoApp.Stop()

	app, err := test.GetTestingApp(demoApp)
	require.NoError(t, err)

	go app.Start()

	time.Sleep(1 * time.Second)

	testID, err := createImportPokemonTest(app, demoApp)
	assert.NoError(t, err)
	assert.NotEmpty(t, testID)

	err = createTestAssertions(testID)
	assert.NoError(t, err)

	resultID, err := runTest(app, testID)
	assert.NoError(t, err)
	assert.NotEmpty(t, resultID)

	testRunResult := waitForResultState(app, testID, resultID, executor.TestRunStateFinished, 30*time.Second)
	assert.NotNil(t, testRunResult, "Test result should not be nil")
	assert.Greater(t, len(testRunResult.Trace.ResourceSpans), 0, "Number of affected spans should be greater than 0")
	assert.Equal(t, executor.TestRunStateFinished, testRunResult.State, "Result should be in FINISHED state")

	spanAssertions := make([]openapi.SpanAssertionResult, 0)
	for _, assertionResult := range testRunResult.AssertionResult {
		spanAssertions = append(spanAssertions, assertionResult.SpanAssertionResults...)
	}

	assert.Equal(t, len(spanAssertions), 2, "Test should contain 2 assertions")
	for _, spanAssertion := range spanAssertions {
		assert.True(t, spanAssertion.Passed, "All assertions should pass")
	}
}

func createImportPokemonTest(app *app.App, demoApp *test.DemoApp) (string, error) {
	body := openapi.Test{
		Name:        "Import Pokemon",
		Description: "Import a pokemon into the api",
		ServiceUnderTest: openapi.TestServiceUnderTest{
			Request: openapi.HttpRequest{
				Url:    fmt.Sprintf("http://%s/pokemon/import", demoApp.Endpoint()),
				Method: "POST",
				Headers: []openapi.HttpResponseHeaders{
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

	return test.TestId, nil
}

func createTestAssertions(testID string) error {
	body := openapi.Assertion{
		Selector: `span[service.name="pokeshop" tracetest.span.type="http" name="POST /pokemon/import"]`,
		SpanAssertions: []openapi.SpanAssertion{
			{
				PropertyName:    "http.status_code",
				Operator:        "EQUALS",
				ComparisonValue: "200",
			},
			{
				PropertyName:    "tracetest.span.duration",
				Operator:        "LESSTHAN",
				ComparisonValue: "100",
			},
		},
	}
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("could not convert body into json: %w", err)
	}
	bytesBuffer := bytes.NewBuffer(jsonBytes)
	url := fmt.Sprintf("%s/api/tests/%s/assertions", appEndpoint, testID)

	response, err := http.Post(url, "application/json", bytesBuffer)
	if err != nil {
		return fmt.Errorf("could not send request to create assertion: %w", err)
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("could not create assertion. Expected status 200, got %d", response.StatusCode)
	}

	return nil
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

	testRunResult := openapi.TestRunResult{}
	err = json.Unmarshal(bodyBytes, &testRunResult)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return testRunResult.ResultId, nil
}

func waitForResultState(app *app.App, testID, resultID, state string, timeout time.Duration) *openapi.TestRunResult {
	timeoutTicker := time.NewTicker(timeout)
	executionTicker := time.NewTicker(1 * time.Second)
	outputChannel := make(chan *openapi.TestRunResult, 1)
	go func() {
		for {
			select {
			case <-timeoutTicker.C:
				outputChannel <- nil
				return
			case <-executionTicker.C:
				testRunResult := getTestRunResultInState(app, testID, resultID, state)
				if testRunResult != nil {
					outputChannel <- testRunResult
					return
				}
			}
		}
	}()

	testRunResult := <-outputChannel
	return testRunResult
}

func getTestRunResultInState(app *app.App, testID, resultID, state string) *openapi.TestRunResult {
	result, err := getTestResult(app, testID, resultID)
	if err != nil {
		return nil
	}

	if result.State == state {
		return result
	}

	return nil
}

func getTestResult(app *app.App, testID, resultID string) (*openapi.TestRunResult, error) {
	url := fmt.Sprintf("%s/api/tests/%s/results/%s", appEndpoint, testID, resultID)
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not send request to run test: %w", err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("could not run test. Expected status 200, got %d", response.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	testRunResult := openapi.TestRunResult{}
	err = json.Unmarshal(bodyBytes, &testRunResult)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal response body: %w", err)
	}

	return &testRunResult, nil
}
