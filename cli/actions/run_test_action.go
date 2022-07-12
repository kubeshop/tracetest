package actions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/definition"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/variable"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type RunTestConfig struct {
	DefinitionFile string
	WaitForResult  bool
	JUnit          string
}

type runTestAction struct {
	config config.Config
	logger *zap.Logger
	client *openapi.APIClient
}

var _ Action[RunTestConfig] = &runTestAction{}

type runTestOutput struct {
	Test      openapi.Test    `json:"test"`
	Run       openapi.TestRun `json:"testRun"`
	RunWebURL string          `json:"testRunWebUrl"`
}

func NewRunTestAction(config config.Config, logger *zap.Logger, client *openapi.APIClient) runTestAction {
	return runTestAction{config, logger, client}
}

func (a runTestAction) Run(ctx context.Context, args RunTestConfig) error {
	if args.DefinitionFile == "" {
		return fmt.Errorf("you must specify a definition file to run a test")
	}

	if args.JUnit != "" && !args.WaitForResult {
		return fmt.Errorf("--junit option requires --wait-for-result")
	}

	a.logger.Debug("Running test from definition", zap.String("definitionFile", args.DefinitionFile))
	output, err := a.runDefinition(ctx, args.DefinitionFile, args.WaitForResult, args.JUnit)

	if err != nil {
		return fmt.Errorf("could not run definition: %w", err)
	}

	bytes, err := json.Marshal(output)
	if err != nil {
		return fmt.Errorf("could not marshal output json: %w", err)
	}

	fmt.Print(string(bytes))
	return nil
}

func (a runTestAction) runDefinition(ctx context.Context, definitionFile string, waitForResult bool, junit string) (runTestOutput, error) {
	definition, err := file.LoadDefinition(definitionFile)
	if err != nil {
		return runTestOutput{}, err
	}

	err = definition.Validate()
	if err != nil {
		return runTestOutput{}, fmt.Errorf("invalid definition file: %w", err)
	}

	var test openapi.Test

	a.logger.Debug("try to create test")

	test, exists, err := a.createTestFromDefinition(ctx, definition)
	if err != nil {
		return runTestOutput{}, fmt.Errorf("could not create test from definition: %w", err)
	}

	if exists {
		a.logger.Debug("test exists. Updating it")
		test, err = a.updateTestFromDefinition(ctx, definition)
		if err != nil {
			return runTestOutput{}, fmt.Errorf("could not update test using definition: %w", err)
		}
	}

	definition.Id = *test.Id
	err = file.SaveDefinition(definitionFile, definition)
	if err != nil {
		return runTestOutput{}, fmt.Errorf("could not save test definition: %w", err)
	}

	testRun, err := a.runTest(ctx, definition.Id)
	if err != nil {
		return runTestOutput{}, fmt.Errorf("could not run test: %w", err)
	}

	if waitForResult {
		updatedTestRun, err := a.waitForResult(ctx, definition.Id, *testRun.Id)
		if err != nil {
			return runTestOutput{}, fmt.Errorf("could not wait for result: %w", err)
		}

		testRun = updatedTestRun

		if err := a.saveJUnitFile(ctx, definition.Id, *testRun.Id, junit); err != nil {
			return runTestOutput{}, fmt.Errorf("could not save junit file: %w", err)
		}
	}

	return runTestOutput{
		Test:      test,
		Run:       testRun,
		RunWebURL: fmt.Sprintf("%s://%s/test/%s/run/%s", a.config.Scheme, a.config.Endpoint, definition.Id, *testRun.Id),
	}, nil
}

func (a runTestAction) saveJUnitFile(ctx context.Context, testId, testRunId, outputFile string) error {
	if outputFile == "" {
		return nil
	}

	req := a.client.ApiApi.GetRunResultJUnit(ctx, testId, testRunId)
	junit, _, err := a.client.ApiApi.GetRunResultJUnitExecute(req)
	if err != nil {
		return fmt.Errorf("could not execute request: %w", err)
	}

	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("could not create junit output file: %w", err)
	}

	_, err = f.WriteString(junit)

	return err

}

func (a runTestAction) createTestFromDefinition(ctx context.Context, definition definition.Test) (_ openapi.Test, exists bool, _ error) {
	variableInjector := variable.NewInjector()
	variableInjector.Inject(&definition)

	yamlContentBytes, err := yaml.Marshal(definition)
	if err != nil {
		return openapi.Test{}, false, fmt.Errorf("could not marshal yaml: %w", err)
	}

	yamlContent := string(yamlContentBytes)

	textDefinition := openapi.TextDefinition{Content: &yamlContent}
	req := a.client.ApiApi.CreateTestFromDefinition(ctx)
	req = req.TextDefinition(textDefinition)

	a.logger.Debug("Sending request to create test", zap.ByteString("test", yamlContentBytes))
	createdTest, resp, err := a.client.ApiApi.CreateTestFromDefinitionExecute(req)

	if resp != nil && resp.StatusCode == http.StatusBadRequest {
		// trying to create a test with already exsiting ID
		return openapi.Test{}, true, nil
	}

	if err != nil {
		return openapi.Test{}, false, fmt.Errorf("could not execute request: %w", err)
	}

	return *createdTest, false, nil
}

func (a runTestAction) updateTestFromDefinition(ctx context.Context, definition definition.Test) (openapi.Test, error) {
	variableInjector := variable.NewInjector()
	variableInjector.Inject(&definition)

	yamlContentBytes, err := yaml.Marshal(definition)
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not marshal yaml: %w", err)
	}

	yamlContent := string(yamlContentBytes)
	textDefinition := openapi.TextDefinition{Content: &yamlContent}

	req := a.client.ApiApi.UpdateTestFromDefinition(ctx)
	req = req.TextDefinition(textDefinition)

	a.logger.Debug("Sending request to update test", zap.ByteString("test", yamlContentBytes))
	openapiTest, _, err := a.client.ApiApi.UpdateTestFromDefinitionExecute(req)
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not execute request: %w", err)
	}

	return *openapiTest, nil
}

func (a runTestAction) runTest(ctx context.Context, testID string) (openapi.TestRun, error) {
	req := a.client.ApiApi.RunTest(ctx, testID)
	run, _, err := a.client.ApiApi.RunTestExecute(req)
	if err != nil {
		return openapi.TestRun{}, fmt.Errorf("could not execute request: %w", err)
	}

	return *run, nil
}

func (a runTestAction) waitForResult(ctx context.Context, testId string, testRunId string) (openapi.TestRun, error) {
	var testRun openapi.TestRun
	var lastError error
	var wg sync.WaitGroup
	wg.Add(1)
	ticker := time.NewTicker(1 * time.Second) // TODO: make this configurable
	go func() {
		for {
			select {
			case <-ticker.C:
				readyTestRun, err := a.isTestReady(ctx, testId, testRunId)
				if err != nil {
					lastError = err
					wg.Done()
					return
				}

				if readyTestRun != nil {
					testRun = *readyTestRun
					wg.Done()
					return
				}
			}
		}
	}()
	wg.Wait()

	if lastError != nil {
		return openapi.TestRun{}, lastError
	}

	return testRun, nil
}

func (a runTestAction) isTestReady(ctx context.Context, testId string, testRunId string) (*openapi.TestRun, error) {
	req := a.client.ApiApi.GetTestRun(ctx, testId, testRunId)
	run, _, err := a.client.ApiApi.GetTestRunExecute(req)
	if err != nil {
		return &openapi.TestRun{}, fmt.Errorf("could not execute GetTestRun request: %w", err)
	}

	if *run.State == "FAILED" || *run.State == "FINISHED" {
		return run, nil
	}

	return nil, nil
}
