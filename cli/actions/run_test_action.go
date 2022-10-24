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
	a.logger.Debug(
		"Running test from definition",
		zap.String("definitionFile", args.DefinitionFile),
		zap.Bool("waitForResults", args.WaitForResult),
		zap.String("junit", args.JUnit),
	)
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
	f, err := file.Read(params.DefinitionFile)
	if err != nil {
		return err
	}

	defFile := f.Definition()
	if err = defFile.Validate(); err != nil {
		return fmt.Errorf("invalid definition file: %w", err)
	}

	return a.runDefinitionFile(ctx, f, params)
}

func (a runTestAction) runDefinitionFile(ctx context.Context, f file.File, params runDefParams) error {
	err := a.replaceEnvVariables(f)
	if err != nil {
		return err
	}

	body, resp, err := a.client.ApiApi.
		ExecuteDefinition(ctx).
		TextDefinition(openapi.TextDefinition{
			Content: openapi.PtrString(f.Contents()),
			RunInformation: &openapi.TestRunInformation{
				Metadata: params.Metadata,
			},
		}).
		Execute()

	if err != nil {
		return fmt.Errorf("could not execute definition: %w", err)
	}

	if resp.StatusCode == http.StatusCreated && !f.HasID() {
		f, err = f.SetID(body.GetId())
		if err != nil {
			return fmt.Errorf("could not update definition file: %w", err)
		}

		_, err = f.Write()
		if err != nil {
			return fmt.Errorf("could not update definition file: %w", err)
		}
	}

	runID := body.GetRunId()
	a.logger.Debug(
		"executed",
		zap.String("runID", runID),
		zap.String("runType", body.GetType()),
	)

	switch body.GetType() {
	case "Test":
		test, err := a.getTest(ctx, body.GetId())
		if err != nil {
			return fmt.Errorf("could not get test info: %w", err)
		}
		return a.testRun(ctx, test, runID, params)
	case "Transaction":
		panic("not implemented")
	}

	return fmt.Errorf(`unsuported run type "%s"`, body.GetType())
}

func (a runTestAction) replaceEnvVariables(f file.File) error {
	variableInjector := variable.NewInjector(variable.WithVariableProvider(
		variable.EnvironmentVariableProvider{},
	))

	err := variableInjector.Inject(&f)
	if err != nil {
		return err
	}

	return nil
}

func (a runTestAction) getTest(ctx context.Context, id string) (openapi.Test, error) {
	test, _, err := a.client.ApiApi.
		GetTest(ctx, id).
		Execute()
	if err != nil {
		return openapi.Test{}, fmt.Errorf("could not execute request: %w", err)
	}

	return *test, nil
}

func (a runTestAction) testRun(ctx context.Context, test openapi.Test, runID string, params runDefParams) error {
	a.logger.Debug("run test", zap.Bool("wait-for-results", params.WaitForResult))
	testID := test.GetId()
	testRun, err := a.getTestRun(ctx, testID, runID)
	if err != nil {
		return fmt.Errorf("could not run test: %w", err)
	}

	if params.WaitForResult {
		updatedTestRun, err := a.waitForResult(ctx, testID, testRun.GetId())
		if err != nil {
			return fmt.Errorf("could not wait for result: %w", err)
		}

		testRun = updatedTestRun

		if err := a.saveJUnitFile(ctx, testID, testRun.GetId(), params.JunitFile); err != nil {
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

func (a runTestAction) getTestRun(ctx context.Context, testID, runID string) (openapi.TestRun, error) {
	run, _, err := a.client.ApiApi.
		GetTestRun(ctx, testID, runID).
		Execute()
	if err != nil {
		return openapi.TestRun{}, fmt.Errorf("could not execute request: %w", err)
	}

	return *run, nil
}

func (a runTestAction) waitForResult(ctx context.Context, testID, testRunID string) (openapi.TestRun, error) {
	var (
		testRun   openapi.TestRun
		lastError error
		wg        sync.WaitGroup
	)
	wg.Add(1)
	ticker := time.NewTicker(1 * time.Second) // TODO: make this configurable
	go func() {
		for {
			select {
			case <-ticker.C:
				readyTestRun, err := a.isTestReady(ctx, testID, testRunID)
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

func (a runTestAction) isTestReady(ctx context.Context, testID, testRunId string) (*openapi.TestRun, error) {
	req := a.client.ApiApi.GetTestRun(ctx, testID, testRunId)
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
