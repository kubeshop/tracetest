package actions

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	cienvironment "github.com/cucumber/ci-environment/go"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/variable"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"go.uber.org/zap"
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

type runDefParams struct {
	DefinitionFile string
	WaitForResult  bool
	JunitFile      string
	Metadata       map[string]string
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

	metadata := a.getMetadata()
	a.logger.Debug("Running test from definition", zap.String("definitionFile", args.DefinitionFile))
	params := runDefParams{
		DefinitionFile: args.DefinitionFile,
		WaitForResult:  args.WaitForResult,
		JunitFile:      args.JUnit,
		Metadata:       metadata,
	}

	err := a.runDefinition(ctx, params)
	if err != nil {
		return fmt.Errorf("could not run definition: %w", err)
	}

	return nil
}

func (a runTestAction) runDefinition(ctx context.Context, params runDefParams) error {
	defFile, err := file.LoadDefinition(params.DefinitionFile)
	if err != nil {
		return err
	}

	err = defFile.Validate()
	if err != nil {
		return fmt.Errorf("invalid definition file: %w", err)
	}

	if test, err := defFile.Test(); err == nil {
		return a.runTestFile(ctx, test, params)
	}

	return fmt.Errorf("invalid file type")
}

func (a runTestAction) runTestFile(ctx context.Context, defTest yaml.Test, params runDefParams) error {
	a.logger.Debug("try to create test")

	test, exists, err := a.createTestFromDefinition(ctx, defTest)
	if err != nil {
		return fmt.Errorf("could not create test from definition: %w", err)
	}

	if exists {
		a.logger.Debug("test exists. Updating it")
		test, err = a.updateTestFromDefinition(ctx, defTest)
		if err != nil {
			return fmt.Errorf("could not update test using definition: %w", err)
		}
	}

	defTest.ID = *test.Id
	err = file.SetTestID(params.DefinitionFile, *test.Id)
	if err != nil {
		return fmt.Errorf("could not save test definition: %w", err)
	}

	testRun, err := a.runTest(ctx, defTest.ID, params.Metadata)
	if err != nil {
		return fmt.Errorf("could not run test: %w", err)
	}

	if params.WaitForResult {
		updatedTestRun, err := a.waitForResult(ctx, defTest.ID, testRun.GetId())
		if err != nil {
			return fmt.Errorf("could not wait for result: %w", err)
		}

		testRun = updatedTestRun

		if err := a.saveJUnitFile(ctx, defTest.ID, testRun.GetId(), params.JunitFile); err != nil {
			return fmt.Errorf("could not save junit file: %w", err)
		}
	}

	tro := formatters.TestRunOutput{
		HasResults: params.WaitForResult,
		Test:       test,
		Run:        testRun,
	}

	formatter := formatters.TestRun(a.config, true)
	formattedOutput := formatter.Format(tro)
	fmt.Print(formattedOutput)

	allPassed := tro.Run.Result.GetAllPassed()

	if params.WaitForResult && !allPassed {
		// It failed, so we have to return an error status
		os.Exit(1)
	}

	return nil
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

func (a runTestAction) createTestFromDefinition(ctx context.Context, definition yaml.Test) (_ openapi.Test, exists bool, _ error) {
	variableInjector := variable.NewInjector()
	variableInjector.Inject(&definition)

	yamlContentBytes, err := yaml.Encode(yaml.File{
		Type: yaml.FileTypeTest,
		Spec: definition,
	})

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

func (a runTestAction) updateTestFromDefinition(ctx context.Context, defTest yaml.Test) (openapi.Test, error) {
	variableInjector := variable.NewInjector()
	variableInjector.Inject(&defTest)

	yamlContentBytes, err := yaml.Encode(yaml.File{
		Type: yaml.FileTypeTest,
		Spec: defTest,
	})
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not marshal yaml: %w", err)
	}

	yamlContent := string(yamlContentBytes)
	textDefinition := openapi.TextDefinition{Content: &yamlContent}

	req := a.client.ApiApi.UpdateTestFromDefinition(ctx, defTest.ID)
	req = req.TextDefinition(textDefinition)

	a.logger.Debug("Sending request to update test", zap.ByteString("test", yamlContentBytes))
	openapiTest, _, err := a.client.ApiApi.UpdateTestFromDefinitionExecute(req)
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not execute request: %w", err)
	}

	return *openapiTest, nil
}

func (a runTestAction) runTest(ctx context.Context, testID string, metadata map[string]string) (openapi.TestRun, error) {
	req := a.client.ApiApi.RunTest(ctx, testID)
	req = req.TestRunInformation(openapi.TestRunInformation{
		Metadata: metadata,
	})
	run, _, err := a.client.ApiApi.RunTestExecute(req)
	if err != nil {
		return openapi.TestRun{}, fmt.Errorf("could not execute request: %w", err)
	}

	return *run, nil
}

func (a runTestAction) waitForResult(ctx context.Context, testId, testRunId string) (openapi.TestRun, error) {
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

func (a runTestAction) isTestReady(ctx context.Context, testId, testRunId string) (*openapi.TestRun, error) {
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

func (a runTestAction) getMetadata() map[string]string {
	ci := cienvironment.DetectCIEnvironment()
	if ci == nil {
		return map[string]string{}
	}

	metadata := map[string]string{
		"name":        ci.Name,
		"url":         ci.URL,
		"buildNumber": ci.BuildNumber,
	}

	if ci.Git != nil {
		metadata["branch"] = ci.Git.Branch
		metadata["tag"] = ci.Git.Tag
		metadata["revision"] = ci.Git.Revision
	}

	return metadata
}
