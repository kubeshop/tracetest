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

type testSuiteFormatter interface {
	Format(output formatters.TestSuiteRunOutput, format formatters.Output) string
}

type testSuiteRunner struct {
	client        resourcemanager.Client
	openapiClient *openapi.APIClient
	formatter     testSuiteFormatter
}

func TestSuiteRunner(
	client resourcemanager.Client,
	openapiClient *openapi.APIClient,
	formatter testSuiteFormatter,
) Runner {
	return testSuiteRunner{
		client:        client,
		openapiClient: openapiClient,
		formatter:     formatter,
	}
}

func (r testSuiteRunner) Name() string {
	return "testsuite"
}

func (r testSuiteRunner) GetByID(_ context.Context, id string) (resource any, _ error) {
	jsonTestSuite, err := r.client.Get(context.Background(), id, jsonFormat)
	if err != nil {
		return nil, fmt.Errorf("cannot get test suite '%s': %w", id, err)
	}

	var testSuite openapi.TestSuiteResource
	err = jsonFormat.Unmarshal([]byte(jsonTestSuite), &testSuite)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal test suite definition file: %w", err)
	}

	return testSuite, nil
}

func (r testSuiteRunner) Apply(ctx context.Context, df fileutil.File) (resource any, _ error) {
	updated, err := r.client.Apply(ctx, df, yamlFormat)
	if err != nil {
		return nil, fmt.Errorf("could not read test suite file: %w", err)
	}

	var parsed openapi.TestSuiteResource
	err = yamlFormat.Unmarshal([]byte(updated), &parsed)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal test suite definition file: %w", err)
	}

	return parsed, nil
}

func (r testSuiteRunner) StartRun(ctx context.Context, resource any, runInfo openapi.RunInformation) (RunResult, error) {
	tran := resource.(openapi.TestSuiteResource)
	run, resp, err := r.openapiClient.ApiApi.
		RunTestSuite(ctx, tran.Spec.GetId()).
		RunInformation(runInfo).
		Execute()

	err = HandleRunError(resp, err)
	if err != nil {
		return RunResult{}, err
	}

	full, err := r.client.Get(ctx, tran.Spec.GetId(), jsonFormat)
	if err != nil {
		return RunResult{}, fmt.Errorf("cannot get full test suite '%s': %w", tran.Spec.GetId(), err)
	}
	err = jsonFormat.Unmarshal([]byte(full), &tran)
	if err != nil {
		return RunResult{}, fmt.Errorf("cannot get full test suite '%s': %w", tran.Spec.GetId(), err)
	}

	return RunResult{
		Resource: tran,
		Run:      *run,
	}, nil
}

func (r testSuiteRunner) UpdateResult(ctx context.Context, result RunResult) (RunResult, error) {
	testSuite := result.Resource.(openapi.TestSuiteResource)
	run := result.Run.(openapi.TestSuiteRun)
	runID, err := strconv.Atoi(run.GetId())
	if err != nil {
		return RunResult{}, fmt.Errorf("invalid test suite run id format: %w", err)
	}

	updated, _, err := r.openapiClient.ApiApi.
		GetTestSuiteRun(ctx, testSuite.Spec.GetId(), int32(runID)).
		Execute()

	if err != nil {
		return RunResult{}, err
	}

	allPassed := true
	for _, s := range updated.GetSteps() {
		if !s.Result.GetAllPassed() {
			allPassed = false
			break
		}
	}

	passed := !isStateFailed(updated.GetState()) && allPassed && updated.GetAllStepsRequiredGatesPassed()

	return RunResult{
		Resource: testSuite,
		Run:      *updated,
		Finished: isStateFinished(updated.GetState()),
		Passed:   passed,
	}, nil
}

func (r testSuiteRunner) JUnitResult(ctx context.Context, result RunResult) (string, error) {
	return "", ErrJUnitNotSupported
}

func (r testSuiteRunner) FormatResult(result RunResult, format string) string {
	testSuite := result.Resource.(openapi.TestSuiteResource)
	run := result.Run.(openapi.TestSuiteRun)

	tro := formatters.TestSuiteRunOutput{
		HasResults: result.Finished,
		TestSuite:  testSuite.GetSpec(),
		Run:        run,
	}

	return r.formatter.Format(tro, formatters.Output(format))
}
