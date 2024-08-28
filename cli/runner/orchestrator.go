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

	"github.com/davecgh/go-spew/spew"
	"github.com/kubeshop/tracetest/cli/metadata"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/kubeshop/tracetest/cli/variable"
	"github.com/kubeshop/tracetest/cli/varset"
	"go.uber.org/zap"
)

// RunOptions defines options for running a resource
// ID and DefinitionFile are mutually exclusive and the only required options
type RunOptions struct {
	// ID of the resource to run
	ID string

	// path to the file with resource definition
	// the file will be applied before running
	DefinitionFile string

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

	// ResourceType defines what is the type of resource that is being run. It's value
	// is filled automatically when the user define the type of resource that will be run
	// when they enter: tracetest run <resource-name> --id <id>
	ResourceType string
}

// RunResult holds the result of the run
// Resources
type RunResult struct {
	// The resource being run. If has been preprocessed, this needs to be the updated version
	Resource any

	// The result of the run. It can be anything the resource needs for validating and formatting the result
	Run any

	// If true, it means that the current run is ready to be presented to the user
	Finished bool

	// Whether the run has passed or not. Used to determine exit code
	Passed bool
}

// Runner defines interface for running a resource
type Runner interface {
	// Name of the runner. must match the resource name it handles
	Name() string

	// Apply the given file and return a resource. The resource can be of any type.
	// It will then be used by Run method
	Apply(context.Context, fileutil.File) (resource any, _ error)

	// GetByID gets the resource by ID. This method is used to get the resource when running from id
	GetByID(_ context.Context, id string) (resource any, _ error)

	// StartRun starts running the resource and return the result. This method should not wait for the test to finish
	StartRun(_ context.Context, resource any, _ openapi.RunInformation) (RunResult, error)

	// UpdateResult is regularly called by the orchestrator to check the status of the run
	UpdateResult(context.Context, RunResult) (RunResult, error)

	// JUnitResult returns the result of the run in JUnit format
	JUnitResult(context.Context, RunResult) (string, error)

	// Format the result of the run into a string
	FormatResult(_ RunResult, format string) string
}

func Orchestrator(
	logger *zap.Logger,
	openapiClient *openapi.APIClient,
	variableSets resourcemanager.Client,
	runnerRegistry Registry,
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
	runnerRegistry Registry
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

func (o orchestrator) Run(ctx context.Context, opts RunOptions, outputFormat string) (exitCode int, _ error) {
	o.logger.Debug(
		"Running test from definition",
		zap.String("definitionFile", opts.DefinitionFile),
		zap.String("ID", opts.ID),
		zap.String("varSetID", opts.VarsID),
		zap.Bool("skipResultsWait", opts.SkipResultWait),
		zap.String("junitOutputFile", opts.JUnitOuptutFile),
		zap.Strings("requiredGates", opts.RequiredGates),
	)

	variableSetFetcher := GetVariableSetFetcher(o.logger, o.variableSets)

	varsID, err := variableSetFetcher.Fetch(ctx, opts.VarsID)
	if err != nil {
		return ExitCodeGeneralError, fmt.Errorf("cannot resolve variable set id: %w", err)
	}

	resourceFetcher := GetResourceFetcher(o.logger, o.runnerRegistry)

	var resource any

	if opts.DefinitionFile != "" {
		resource, err = resourceFetcher.FetchWithDefinitionFile(ctx, opts.DefinitionFile)
	} else {
		resource, err = resourceFetcher.FetchWithID(ctx, opts.ID, opts.ResourceType)
	}
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

	var result RunResult
	var ev varset.VarSets

	// iterate until we have all env vars,
	// or the server returns an actual error
	for {
		runInfo := openapi.RunInformation{
			VariableSetId: &varsID,
			Variables:     ev.ToOpenapi(),
			Metadata:      metadata.GetMetadata(),
			RequiredGates: getRequiredGates(opts.RequiredGates),
		}

		result, err = runner.StartRun(ctx, resource, runInfo)
		if err == nil {
			break
		}
		if !errors.Is(err, varset.MissingVarsError{}) {
			// actual error, return
			return ExitCodeGeneralError, fmt.Errorf("cannot run test: %w", err)
		}

		// missing vars error
		ev = varset.AskForMissingVars([]varset.VarSet(err.(varset.MissingVarsError)))
		o.logger.Debug("filled variables", zap.Any("variables", ev))
	}

	if opts.SkipResultWait {
		fmt.Println(runner.FormatResult(result, outputFormat))
		return ExitCodeSuccess, nil
	}

	result, err = o.waitForResult(ctx, runner, result)
	if err != nil {
		return ExitCodeGeneralError, fmt.Errorf("cannot wait for test result: %w", err)
	}

	fmt.Println(runner.FormatResult(result, outputFormat))

	err = o.writeJUnitReport(ctx, runner, result, opts.JUnitOuptutFile)
	if err != nil {
		return ExitCodeGeneralError, fmt.Errorf("cannot write junit report: %w", err)
	}

	exitCode = ExitCodeSuccess
	if !result.Passed {
		exitCode = ExitCodeTestNotPassed
	}

	return exitCode, nil
}

func (o orchestrator) injectLocalEnvVars(ctx context.Context, df fileutil.File) (fileutil.File, error) {
	variableInjector := variable.NewInjector(variable.WithVariableProvider(
		variable.EnvironmentVariableProvider{},
	))

	injected, err := variableInjector.ReplaceInString(string(df.Contents()))
	if err != nil {
		return df, fmt.Errorf("cannot inject local variable set: %w", err)
	}

	df = fileutil.New(df.AbsPath(), []byte(injected))

	return df, nil
}

func (o orchestrator) waitForResult(ctx context.Context, r Runner, result RunResult) (RunResult, error) {
	var (
		updatedResult RunResult
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
		return RunResult{}, lastError
	}

	return updatedResult, nil
}

var ErrJUnitNotSupported = errors.New("junit report is not supported for this resource type")

func (a orchestrator) writeJUnitReport(ctx context.Context, r Runner, result RunResult, outputFile string) error {
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
		return varset.BuildMissingVarsError(body)
	}

	if ok, msg := attemptToParseStructuredError(body); ok {
		return fmt.Errorf("could not run resouce: %s", msg)
	}

	if reqErr != nil {
		return fmt.Errorf("could not run resouce: %w", reqErr)
	}

	return nil
}

func attemptToParseStructuredError(body []byte) (bool, string) {
	var parsed struct {
		Status int    `json:"status"`
		Detail string `json:"detail"`
	}

	err := jsonFormat.Unmarshal(body, &parsed)
	if err != nil || parsed.Status == 0 {
		return false, ""
	}

	msg := fmt.Sprintf("%s (code %d)", parsed.Detail, parsed.Status)

	return true, msg
}
