package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	cienvironment "github.com/cucumber/ci-environment/go"
	"github.com/davecgh/go-spew/spew"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/kubeshop/tracetest/cli/runner"
	"go.uber.org/zap"
)

// RunOptions defines options for running a resource
// IDs and DefinitionFiles are mutually exclusive and the only required options
type RunOptions struct {
	// path to the file with resource definition
	// the file will be applied before running
	DefinitionFiles []string

	// varsID or path to the file with environment definition
	VarsID string

	// By default the runner will wait for the result of the run
	// if this option is true, the wait will be skipped
	SkipResultWait bool

	// Optional path to the file where the result of the run will be saved
	// in JUnit xml format
	JUnitOuptutFile string

	// Overrides the default required gates for the resource
	RequiredGates []string
}

func MultiFileOrchestrator(
	logger *zap.Logger,
	openapiClient *openapi.APIClient,
	variableSets resourcemanager.Client,
	runnerRegistry runner.Registry,
) orchestrator {
	return orchestrator{
		logger:         logger,
		openapiClient:  openapiClient,
		variableSets:   variableSets,
		runnerRegistry: runnerRegistry,
	}
}

type orchestrator struct {
	logger *zap.Logger

	openapiClient  *openapi.APIClient
	variableSets   resourcemanager.Client
	runnerRegistry runner.Registry
}

var (
	yamlFormat = resourcemanager.Formats.Get(resourcemanager.FormatYAML)
	jsonFormat = resourcemanager.Formats.Get(resourcemanager.FormatJSON)
)

const (
	ExitCodeSuccess       = 0
	ExitCodeGeneralError  = 1
	ExitCodeTestNotPassed = 2
)

func (o orchestrator) Run(ctx context.Context, opts RunOptions) (exitCode int, _ error) {
	o.logger.Debug(
		"Running tests from definition",
		zap.Strings("definitionFiles", opts.DefinitionFiles),
		zap.String("varSetID", opts.VarsID),
		zap.Bool("skipResultsWait", opts.SkipResultWait),
		zap.String("junitOutputFile", opts.JUnitOuptutFile),
		zap.Strings("requiredGates", opts.RequiredGates),
	)

	variableSetFetcher := runner.GetVariableSetFetcher(o.logger, o.variableSets)

	varsID, err := variableSetFetcher.Fetch(ctx, opts.VarsID)
	if err != nil {
		return ExitCodeGeneralError, fmt.Errorf("cannot resolve variable set id: %w", err)
	}
	o.logger.Debug("env resolved", zap.String("ID", varsID))

	var resources []any
	hasDefinitionFilesDefined := opts.DefinitionFiles != nil && len(opts.DefinitionFiles) > 0

	resourceFetcher := runner.GetResourceFetcher(o.logger, o.runnerRegistry)

	if !hasDefinitionFilesDefined {
		return ExitCodeGeneralError, fmt.Errorf("you must define at least two files to use the multifile orchestrator")
	}

	resources = make([]any, 0, len(opts.DefinitionFiles))

	// call run group creation

	for _, definitionFile := range opts.DefinitionFiles {
		resource, err := resourceFetcher.FetchWithDefinitionFile(ctx, definitionFile)
		if err != nil {
			return ExitCodeGeneralError, err
		}

		resourceType, err := resourcemanager.GetResourceType(resource)
		if err != nil {
			return ExitCodeGeneralError, fmt.Errorf("cannot extract type from resource: %w", err)
		}

		runner, err := o.runnerRegistry.Get(resourceType)
		if err != nil {
			return ExitCodeGeneralError, fmt.Errorf("cannot find runner for resource type %s: %w", resourceType, err)
		}

		runInfo := openapi.RunInformation{
			VariableSetId: &varsID,
			Variables:     ev.toOpenapi(),
			Metadata:      getMetadata(),
			RequiredGates: getRequiredGates(opts.RequiredGates),
		}

		runner.StartRun(ctx, resource, runInfo)
		resources = append(resources, resource)
	}

	var result RunResult
	var ev varSets

	// iterate until we have all env vars,
	// or the server returns an actual error
	for {

		result, err = r.StartRun(ctx, resource, runInfo)
		if err == nil {
			break
		}
		if !errors.Is(err, missingVarsError{}) {
			// actual error, return
			return ExitCodeGeneralError, fmt.Errorf("cannot run test: %w", err)
		}

		// missing vars error
		ev = askForMissingVars([]varSet(err.(missingVarsError)))
		o.logger.Debug("filled variables", zap.Any("variables", ev))
	}

	if opts.SkipResultWait {
		fmt.Println(r.FormatResult(result, outputFormat))
		return ExitCodeSuccess, nil
	}

	result, err = o.waitForResult(ctx, r, result)
	if err != nil {
		return ExitCodeGeneralError, fmt.Errorf("cannot wait for test result: %w", err)
	}

	fmt.Println(r.FormatResult(result, outputFormat))

	err = o.writeJUnitReport(ctx, r, result, opts.JUnitOuptutFile)
	if err != nil {
		return ExitCodeGeneralError, fmt.Errorf("cannot write junit report: %w", err)
	}

	exitCode = ExitCodeSuccess
	if !result.Passed {
		exitCode = ExitCodeTestNotPassed
	}

	return exitCode, nil
}

func (o orchestrator) waitForResult(ctx context.Context, r runner.Runner, result runner.RunResult) (runner.RunResult, error) {
	var (
		updatedResult runner.RunResult
		lastError     error
		wg            sync.WaitGroup
	)

	wg.Add(1)
	ticker := time.NewTicker(1 * time.Second) // TODO: change to websockets
	go func() {
		for range ticker.C {
			updated, err := r.UpdateResult(ctx, result)
			o.logger.Debug("updated result", zap.String("result", spew.Sdump(updated)))
			if err != nil {
				o.logger.Debug("UpdateResult failed", zap.Error(err))
				lastError = err
				wg.Done()
				return
			}

			if updated.Finished {
				o.logger.Debug("result is finished")
				updatedResult = updated
				wg.Done()
				return
			}
			o.logger.Debug("still waiting")
		}
	}()
	wg.Wait()

	if lastError != nil {
		return runner.RunResult{}, lastError
	}

	return updatedResult, nil
}

var ErrJUnitNotSupported = errors.New("junit report is not supported for this resource type")

func (a orchestrator) writeJUnitReport(ctx context.Context, r runner.Runner, result runner.RunResult, outputFile string) error {
	if outputFile == "" {
		a.logger.Debug("no junit output file specified")
		return nil
	}

	a.logger.Debug("saving junit report", zap.String("outputFile", outputFile))

	report, err := r.JUnitResult(ctx, result)
	if err != nil {
		return err
	}

	a.logger.Debug("junit report", zap.String("report", report))
	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("could not create junit output file: %w", err)
	}

	_, err = f.WriteString(report)

	return err
}

var source = "cli"

func getMetadata() map[string]string {
	ci := cienvironment.DetectCIEnvironment()
	if ci == nil {
		return map[string]string{
			"source": source,
		}
	}

	metadata := map[string]string{
		"name":        ci.Name,
		"url":         ci.URL,
		"buildNumber": ci.BuildNumber,
		"source":      source,
	}

	if ci.Git != nil {
		metadata["branch"] = ci.Git.Branch
		metadata["tag"] = ci.Git.Tag
		metadata["revision"] = ci.Git.Revision
	}

	return metadata
}

func getRequiredGates(gates []string) []openapi.SupportedGates {
	if len(gates) == 0 {
		return nil
	}
	requiredGates := make([]openapi.SupportedGates, 0, len(gates))

	for _, g := range gates {
		requiredGates = append(requiredGates, openapi.SupportedGates(g))
	}

	return requiredGates
}

// HandleRunError handles errors returned by the server when running a test.
// It normalizes the handling of general errors, like 404,
// but more importantly, it processes the missing environment variables error
// so the orchestrator can request them from the user.
func HandleRunError(resp *http.Response, reqErr error) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read response body: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("resource not found in server")
	}

	if resp.StatusCode == http.StatusUnprocessableEntity {
		return buildMissingVarsError(body)
	}

	if reqErr != nil {
		return fmt.Errorf("could not run test suite: %w", reqErr)
	}

	return nil
}
