package runner

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/pkg/fileutil"
	"github.com/kubeshop/tracetest/cli/pkg/resourcemanager"
)

type testFormatter interface {
	Format(output formatters.TestRunOutput, format formatters.Output) string
}

type testRunner struct {
	client        resourcemanager.Client
	openapiClient *openapi.APIClient
	formatter     testFormatter
}

func TestRunner(
	client resourcemanager.Client,
	openapiClient *openapi.APIClient,
	formatter testFormatter,
) Runner {
	return testRunner{
		client:        client,
		openapiClient: openapiClient,
		formatter:     formatter,
	}
}

func (r testRunner) Name() string {
	return "test"
}

func (r testRunner) Apply(ctx context.Context, df fileutil.File) (resource any, _ error) {
	updated, err := r.client.Apply(ctx, df, yamlFormat)
	if err != nil {
		return nil, fmt.Errorf("could not read test file: %w", err)
	}

	var parsed openapi.TestResource
	err = yamlFormat.Unmarshal([]byte(updated), &parsed)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal test definition file: %w", err)
	}

	return parsed, nil
}

func (r testRunner) StartRun(ctx context.Context, resource any, runInfo openapi.RunInformation) (RunResult, error) {
	test := resource.(openapi.TestResource)
	run, resp, err := r.openapiClient.ApiApi.
		RunTest(ctx, test.Spec.GetId()).
		RunInformation(runInfo).
		Execute()

	err = HandleRunError(resp, err)
	if err != nil {
		return RunResult{}, err
	}

	full, err := r.client.Get(ctx, test.Spec.GetId(), jsonFormat)
	if err != nil {
		return RunResult{}, fmt.Errorf("cannot get full test '%s': %w", test.Spec.GetId(), err)
	}
	err = jsonFormat.Unmarshal([]byte(full), &test)
	if err != nil {
		return RunResult{}, fmt.Errorf("cannot get full test '%s': %w", test.Spec.GetId(), err)
	}

	return RunResult{
		Resource: test,
		Run:      *run,
	}, nil
}

func (r testRunner) UpdateResult(ctx context.Context, result RunResult) (RunResult, error) {
	test := result.Resource.(openapi.TestResource)
	run := result.Run.(openapi.TestRun)
	runID, err := strconv.Atoi(run.GetId())
	if err != nil {
		return RunResult{}, fmt.Errorf("invalid test run id format: %w", err)
	}

	updated, _, err := r.openapiClient.ApiApi.
		GetTestRun(ctx, test.Spec.GetId(), int32(runID)).
		Execute()

	if err != nil {
		return RunResult{}, err
	}

	passed := !isStateFailed(updated.GetState()) && updated.Result.GetAllPassed()

	return RunResult{
		Resource: test,
		Run:      *updated,
		Finished: isStateFinished(updated.GetState()),
		Passed:   passed,
	}, nil
}

func (r testRunner) FormatResult(result RunResult, format string) string {
	test := result.Resource.(openapi.TestResource)
	run := result.Run.(openapi.TestRun)

	tro := formatters.TestRunOutput{
		HasResults: result.Finished,
		IsFailed:   isStateFailed(run.GetState()),
		Test:       test.GetSpec(),
		Run:        run,
	}

	return r.formatter.Format(tro, formatters.Output(format))
}

func isStateFinished(state string) bool {
	return isStateFailed(state) || state == "FINISHED"
}

func isStateFailed(state string) bool {
	return state == "TRIGGER_FAILED" ||
		state == "TRACE_FAILED" ||
		state == "ASSERTION_FAILED" ||
		state == "ANALYZING_ERROR" ||
		state == "FAILED" // this one is for backwards compatibility
}
