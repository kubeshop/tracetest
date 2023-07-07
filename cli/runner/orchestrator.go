package runner

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	cienvironment "github.com/cucumber/ci-environment/go"
	"github.com/davecgh/go-spew/spew"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// RunOptions defines options for running a resource
// ID and DefinitionFile are mutually exclusive and the only required options
type RunOptions struct {
	// ID of the resource to run
	ID string

	// path to the file with resource definition
	// the file will be applied before running
	DefinitionFile string

	// environmentID or path to the file with environment definition
	EnvID string

	// By default the runner will wait for the result of the run
	// if this option is true, the wait will be skipped
	SkipResultWait bool

	// Optional path to the file where the result of the run will be saved
	// in JUnit xml format
	JUnitOuptutFile string
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
	// Apply the given file and return a resource. The resource can be of any type.
	// It will then be used by Run method
	Apply(context.Context, fileutil.File) (resource any, _ error)

	// Run the resource and return the result. This method should not wait for the test to finish
	StartRun(_ context.Context, resource any, _ openapi.RunInformation) (RunResult, error)
	UpdateResult(context.Context, RunResult) (RunResult, error)
	FormatResult(RunResult) (string, error)
}

func Orchestrator(
	logger *zap.Logger,
	openapiClient *openapi.APIClient,
	environments resourcemanager.Client,
) orchestrator {
	return orchestrator{
		logger:        logger,
		openapiClient: openapiClient,
		environments:  environments,
	}
}

type orchestrator struct {
	logger *zap.Logger

	openapiClient *openapi.APIClient
	environments  resourcemanager.Client
}

var yamlFormat resourcemanager.Format

func init() {
	var err error
	yamlFormat, err = resourcemanager.Formats.Get(resourcemanager.FormatYAML)
	if err != nil {
		panic(fmt.Errorf("could not get yaml format: %w", err))
	}
}

const (
	ExitCodeSuccess        = 0
	ExitCodeGeneralError   = 1
	ExitCodeTestNotPasswed = 2
)

func (o orchestrator) Run(ctx context.Context, r Runner, opts RunOptions) (exitCode int, _ error) {

	o.logger.Debug(
		"Running test from definition",
		zap.String("definitionFile", opts.DefinitionFile),
		zap.String("ID", opts.ID),
		zap.String("envID", opts.EnvID),
		zap.Bool("skipResultsWait", opts.SkipResultWait),
		zap.String("junitOutputFile", opts.JUnitOuptutFile),
	)

	f, err := fileutil.Read(opts.DefinitionFile)
	if err != nil {
		return ExitCodeGeneralError, fmt.Errorf("cannot read definition file %s: %w", opts.DefinitionFile, err)
	}
	df := f
	o.logger.Debug("Definition file read", zap.String("absolutePath", df.AbsPath()))

	envID, err := o.resolveEnvID(ctx, opts.EnvID)
	if err != nil {
		return ExitCodeGeneralError, fmt.Errorf("cannot resolve environment id: %w", err)
	}
	o.logger.Debug("env resolved", zap.String("ID", envID))

	resource, err := r.Apply(ctx, df)
	if err != nil {
		return ExitCodeGeneralError, fmt.Errorf("cannot apply definition file: %w", err)
	}

	o.logger.Debug("Definition file applied", zap.String("updated", string(df.Contents())))

	var result RunResult
	var ev envVars

	// iterate until we have all env vars,
	// or the server returns an actual error
	for {
		runInfo := openapi.RunInformation{
			EnvironmentId: &envID,
			Variables:     ev.toOpenapi(),
			Metadata:      getMetadata(),
		}

		result, err = r.StartRun(ctx, resource, runInfo)
		if err == nil {
			break
		}
		missingEnvVarsErr, ok := err.(missingEnvVarsError)
		if !ok {
			// actual error, return
			return ExitCodeGeneralError, fmt.Errorf("cannot run test: %w", err)
		}

		// missing vars error
		ev = askForMissingVars([]envVar(missingEnvVarsErr))
		o.logger.Debug("filled variables", zap.Any("variables", ev))
	}

	if opts.SkipResultWait {
		fmt.Println(r.FormatResult(result))
		return ExitCodeSuccess, nil
	}

	result, err = o.waitForResult(ctx, r, result)
	if err != nil {
		return ExitCodeGeneralError, fmt.Errorf("cannot wait for test result: %w", err)
	}

	fmt.Println(r.FormatResult(result))

	err = o.writeJUnitReport(ctx, result, opts.JUnitOuptutFile)
	if err != nil {
		return ExitCodeGeneralError, fmt.Errorf("cannot write junit report: %w", err)
	}

	exitCode = ExitCodeSuccess
	if !result.Passed {
		exitCode = ExitCodeTestNotPasswed
	}

	return exitCode, nil
}

func (o orchestrator) resolveEnvID(ctx context.Context, envID string) (string, error) {
	if !fileutil.IsFilePath(envID) {
		o.logger.Debug("envID is not a file path", zap.String("envID", envID))
		return envID, nil
	}

	f, err := fileutil.Read(envID)
	if err != nil {
		return "", fmt.Errorf("cannot read environment file %s: %w", envID, err)
	}

	o.logger.Debug("envID is a file path", zap.String("filePath", envID), zap.Any("file", f))
	updatedEnv, err := o.environments.Apply(ctx, f, yamlFormat)
	if err != nil {
		return "", fmt.Errorf("could not read environment file: %w", err)
	}

	var env openapi.EnvironmentResource
	err = yaml.Unmarshal([]byte(updatedEnv), &env)
	if err != nil {
		o.logger.Error("error parsing json", zap.String("content", updatedEnv), zap.Error(err))
		return "", fmt.Errorf("could not unmarshal environment json: %w", err)
	}

	return env.Spec.GetId(), nil
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

func (a orchestrator) writeJUnitReport(ctx context.Context, result RunResult, outputFile string) error {
	if outputFile == "" {
		return nil
	}

	test := result.Resource.(openapi.TestResource)
	run := result.Run.(openapi.TestRun)
	runID, err := strconv.Atoi(run.GetId())
	if err != nil {
		return fmt.Errorf("invalid run id format: %w", err)
	}

	junit, _, err := a.openapiClient.ApiApi.
		GetRunResultJUnit(
			ctx,
			test.Spec.GetId(),
			int32(runID),
		).
		Execute()
	if err != nil {
		return fmt.Errorf("could not execute request: %w", err)
	}

	f, err := os.Create(junit)
	if err != nil {
		return fmt.Errorf("could not create junit output file: %w", err)
	}

	_, err = f.WriteString(outputFile)

	return err
}

func getMetadata() map[string]string {
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
