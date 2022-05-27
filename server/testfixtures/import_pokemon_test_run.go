package testfixtures

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kubeshop/tracetest/server/app"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
)

func init() {
	RegisterFixture(IMPORT_POKEMON_TEST_RUN, getImportPokemonTestRun)
}

func GetPokemonTestRun(options ...Option) (*openapi.TestRun, error) {
	return GetFixtureValue[*openapi.TestRun](IMPORT_POKEMON_TEST_RUN, options...)
}

func getImportPokemonTestRun(options FixtureOptions) (*openapi.TestRun, error) {
	tracetestApp, err := GetTracetestApp()
	if err != nil {
		return nil, fmt.Errorf("could not get tracetest app: %w", err)
	}

	importPokemonTest, err := GetPokemonTest()
	if err != nil {
		return nil, fmt.Errorf("could not get import pokemon test: %w", err)
	}

	runID, err := runTest(tracetestApp, importPokemonTest.Id)
	if err != nil {
		return nil, fmt.Errorf("could not run test: %w", err)
	}

	if runID == "" {
		return nil, fmt.Errorf("run id must not be empty")
	}

	run := waitForRunState(tracetestApp, importPokemonTest.Id, runID, string(model.RunStateFinished), 30*time.Second)
	if run == nil {
		return nil, fmt.Errorf("test run must not be nil")
	}

	return run, nil
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

func waitForRunState(app *app.App, testID, runID, state string, timeout time.Duration) *openapi.TestRun {
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
				testRun := getRunInState(app, testID, runID, state)
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
