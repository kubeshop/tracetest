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
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/metadata"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
	"github.com/kubeshop/tracetest/cli/runner"
	"github.com/kubeshop/tracetest/cli/varset"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"go.uber.org/zap"
)

// RunOptions defines options for running a resource
// IDs and DefinitionFiles are mutually exclusive and the only required options
type RunOptions struct {
	// if IDs is used it needs to have the ResourceName defined
	IDs          []string
	ResourceName string

	// path to the file with resource definition
	// the file will be applied before running
	DefinitionFiles []string

	// varsID or path to the file with environment definition
	VarsID string

	// runGroupID as string to define it for the entire run execution
	RunGroupID string

	// By default the runner will wait for the result of the run
	// if this option is true, the wait will be skipped
	SkipResultWait bool

	// Optional path to the file where the result of the run will be saved
	// in JUnit xml format
	JUnitOuptutFile string

	// Overrides the default required gates for the resource
	RequiredGates []string
}

type MultipleRunFormatter interface {
	Format(output formatters.MultipleRunOutput[runner.RunResult], format formatters.Output) string
}

func MultiFileOrchestrator(
	logger *zap.Logger,
	openapiClient *openapi.APIClient,
	variableSets resourcemanager.Client,
	runnerRegistry runner.Registry,
	multipleRunFormatter MultipleRunFormatter,
) orchestrator {
	return orchestrator{
		logger:               logger,
		openapiClient:        openapiClient,
		variableSets:         variableSets,
		runnerRegistry:       runnerRegistry,
		multipleRunFormatter: multipleRunFormatter,
	}
}

type orchestrator struct {
	logger *zap.Logger

	openapiClient        *openapi.APIClient
	variableSets         resourcemanager.Client
	runnerRegistry       runner.Registry
	multipleRunFormatter MultipleRunFormatter
}

const (
	ExitCodeSuccess       = 0
	ExitCodeGeneralError  = 1
	ExitCodeTestNotPassed = 2
)

func (o orchestrator) Run(ctx context.Context, opts RunOptions, outputFormat string) (exitCode int, _ error) {
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

	vars := varset.VarSets{}
	resourceFetcher := runner.GetResourceFetcher(o.logger, o.runnerRegistry)

	runGroupID := opts.RunGroupID
	if runGroupID == "" {
		runGroupID = id.GenerateID().String()
	}

	var resources []any
	var runsResults []runner.RunResult

	if len(opts.DefinitionFiles) > 0 {
		resources, runsResults, err = o.runByFiles(ctx, opts, resourceFetcher, &vars, varsID, runGroupID)
		if err != nil {
			return ExitCodeGeneralError, fmt.Errorf("cannot run files: %w", err)
		}
	} else {
		resources, runsResults, err = o.runByIDs(ctx, opts, resourceFetcher, &vars, varsID, runGroupID)
		if err != nil {
			return ExitCodeGeneralError, fmt.Errorf("cannot run by id: %w", err)
		}
	}

	runnerGetter := func(resource any) (formatters.Runner[runner.RunResult], error) {
		resourceType, err := resourcemanager.GetResourceType(resource)
		if err != nil {
			return nil, fmt.Errorf("cannot extract type from resource: %w", err)
		}

		runner, err := o.runnerRegistry.Get(resourceType)
		if err != nil {
			return nil, fmt.Errorf("cannot find runner for resource type %s: %w", resourceType, err)
		}

		return runner, nil
	}

	// 3. if skip wait, print results and exit
	if opts.SkipResultWait {
		output := formatters.MultipleRunOutput[runner.RunResult]{
			Runs:         runsResults,
			Resources:    resources,
			RunGroup:     openapi.RunGroup{Id: runGroupID},
			RunnerGetter: runnerGetter,
			HasResults:   false,
		}

		summary := o.multipleRunFormatter.Format(output, formatters.Output(outputFormat))
		fmt.Println(summary)
		return ExitCodeSuccess, nil
	}

	// 4. wait for the run group
	runGroup, err := o.waitForRunGroup(ctx, runGroupID)
	if err != nil {
		return ExitCodeGeneralError, fmt.Errorf("cannot wait for test result: %w", err)
	}

	// 5. update runs and print results
	for i, result := range runsResults {
		resource := resources[i]

		resourceType, err := resourcemanager.GetResourceType(resource)
		if err != nil {
			return ExitCodeGeneralError, fmt.Errorf("cannot extract type from resource: %w", err)
		}

		runner, err := o.runnerRegistry.Get(resourceType)
		if err != nil {
			return ExitCodeGeneralError, fmt.Errorf("cannot find runner for resource type %s: %w", resourceType, err)
		}

		// TODO: I think we can just pull the test runs from the group instead of updating each of them
		updated, err := runner.UpdateResult(ctx, result)
		if err != nil {
			return ExitCodeGeneralError, fmt.Errorf("cannot update test result: %w", err)
		}
		runsResults[i] = updated

		err = o.writeJUnitReport(ctx, runner, result, opts.JUnitOuptutFile)
		if err != nil {
			return ExitCodeGeneralError, fmt.Errorf("cannot write junit report: %w", err)
		}
	}

	output := formatters.MultipleRunOutput[runner.RunResult]{
		Runs:         runsResults,
		Resources:    resources,
		RunGroup:     runGroup,
		RunnerGetter: runnerGetter,
		HasResults:   true,
	}

	summary := o.multipleRunFormatter.Format(output, formatters.Output(outputFormat))
	fmt.Println(summary)

	exitCode = ExitCodeSuccess
	if runGroup.GetStatus() == "failed" {
		exitCode = ExitCodeTestNotPassed
	}

	return exitCode, nil
}

func (o orchestrator) getDefinitionFiles(file []string) ([]string, error) {
	files := make([]string, 0)

	for _, f := range file {
		files = append(files, fileutil.ReadDirFileNames(f)...)
	}

	return files, nil
}

func (o orchestrator) runByFiles(ctx context.Context, opts RunOptions, resourceFetcher runner.ResourceFetcher, vars *varset.VarSets, varsID string, runGroupID string) ([]any, []runner.RunResult, error) {
	resources := make([]any, 0)
	runsResults := make([]runner.RunResult, 0)
	var mainErr error

	hasDefinitionFilesDefined := opts.DefinitionFiles != nil && len(opts.DefinitionFiles) > 0
	if !hasDefinitionFilesDefined {
		return resources, runsResults, fmt.Errorf("no definition files defined")
	}

	definitionFiles, err := o.getDefinitionFiles(opts.DefinitionFiles)
	if err != nil {
		return resources, runsResults, fmt.Errorf("cannot get definition files: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(definitionFiles))
	for _, definitionFile := range definitionFiles {
		go func(def string) {
			defer wg.Done()
			resource, err := resourceFetcher.FetchWithDefinitionFile(ctx, def)
			if err != nil {
				mainErr = fmt.Errorf("cannot fetch resource from definition file: %w", err)
				return
			}
			result, resource, err := o.createRun(ctx, resource, vars, opts.RequiredGates, varsID, runGroupID)
			if err != nil {
				mainErr = fmt.Errorf("cannot run test: %w", err)
				return
			}

			runsResults = append(runsResults, result)
			resources = append(resources, resource)
		}(definitionFile)
	}

	wg.Wait()
	return resources, runsResults, mainErr
}

func (o orchestrator) runByIDs(ctx context.Context, opts RunOptions, resourceFetcher runner.ResourceFetcher, vars *varset.VarSets, varsID string, runGroupID string) ([]any, []runner.RunResult, error) {
	resources := make([]any, 0)
	runsResults := make([]runner.RunResult, 0)

	for _, id := range opts.IDs {
		resource, err := resourceFetcher.FetchWithID(ctx, opts.ResourceName, id)
		if err != nil {
			return resources, runsResults, err
		}

		result, resource, err := o.createRun(ctx, resource, vars, opts.RequiredGates, varsID, runGroupID)
		if err != nil {
			return resources, runsResults, fmt.Errorf("cannot run test: %w", err)
		}

		runsResults = append(runsResults, result)
		resources = append(resources, resource)
	}

	return resources, runsResults, nil
}

func (o orchestrator) createRun(ctx context.Context, resource any, vars *varset.VarSets, requiredGates []string, varsID, runGroupID string) (runner.RunResult, any, error) {
	resourceType, err := resourcemanager.GetResourceType(resource)
	if err != nil {
		return runner.RunResult{}, nil, fmt.Errorf("cannot extract type from resource: %w", err)
	}

	r, err := o.runnerRegistry.Get(resourceType)
	if err != nil {
		return runner.RunResult{}, nil, fmt.Errorf("cannot find runner for resource type %s: %w", resourceType, err)
	}

	runInfo := openapi.RunInformation{
		VariableSetId: &varsID,
		Variables:     vars.ToOpenapi(),
		Metadata:      metadata.GetMetadata(),
		RequiredGates: getRequiredGates(requiredGates),
		RunGroupId:    &runGroupID,
	}

	// 2. validate missing vars
	for {
		result, err := r.StartRun(ctx, resource, runInfo)
		if err == nil {
			return result, resource, nil
		}
		if !errors.Is(err, varset.MissingVarsError{}) {
			// actual error, return
			return result, resource, fmt.Errorf("cannot run test: %w", err)
		}

		// missing vars error
		newVars := varset.AskForMissingVars([]varset.VarSet(err.(varset.MissingVarsError)))
		vars = &newVars
		o.logger.Debug("filled variables", zap.Any("variables", vars))
	}
}

func (o orchestrator) waitForRunGroup(ctx context.Context, runGroupID string) (openapi.RunGroup, error) {
	var (
		updatedResult openapi.RunGroup
		lastError     error
		wg            sync.WaitGroup
	)

	wg.Add(1)
	ticker := time.NewTicker(1 * time.Second) // TODO: change to websockets
	go func() {
		for range ticker.C {
			req := o.openapiClient.ApiApi.GetRunGroup(ctx, runGroupID)
			runGroup, _, err := req.Execute()

			// updatedResult = runGroup
			o.logger.Debug("updated run group", zap.String("result", spew.Sdump(runGroup)))
			if err != nil {
				o.logger.Debug("UpdateResult failed", zap.Error(err))
				lastError = err
				wg.Done()
				return
			}

			if runGroup.GetStatus() == "succeed" || runGroup.GetStatus() == "failed" {
				o.logger.Debug("result is finished")
				updatedResult = *runGroup
				wg.Done()
				return
			}
			o.logger.Debug("still waiting")
		}
	}()
	wg.Wait()

	if lastError != nil {
		return openapi.RunGroup{}, lastError
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

	if reqErr != nil {
		return fmt.Errorf("could not run test suite: %w", reqErr)
	}

	return nil
}
