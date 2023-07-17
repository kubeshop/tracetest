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

type transactionFormatter interface {
	Format(output formatters.TransactionRunOutput, format formatters.Output) string
}

type transactionRunner struct {
	client        resourcemanager.Client
	openapiClient *openapi.APIClient
	formatter     transactionFormatter
}

func TransactionRunner(
	client resourcemanager.Client,
	openapiClient *openapi.APIClient,
	formatter transactionFormatter,
) Runner {
	return transactionRunner{
		client:        client,
		openapiClient: openapiClient,
		formatter:     formatter,
	}
}

func (r transactionRunner) Name() string {
	return "transaction"
}

func (r transactionRunner) GetByID(_ context.Context, id string) (resource any, _ error) {
	jsonTransaction, err := r.client.Get(context.Background(), id, jsonFormat)
	if err != nil {
		return nil, fmt.Errorf("cannot get transaction '%s': %w", id, err)
	}

	var transaction openapi.TransactionResource
	err = jsonFormat.Unmarshal([]byte(jsonTransaction), &transaction)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal transaction definition file: %w", err)
	}

	return transaction, nil
}

func (r transactionRunner) Apply(ctx context.Context, df fileutil.File) (resource any, _ error) {
	updated, err := r.client.Apply(ctx, df, yamlFormat)
	if err != nil {
		return nil, fmt.Errorf("could not read transaction file: %w", err)
	}

	var parsed openapi.TransactionResource
	err = yamlFormat.Unmarshal([]byte(updated), &parsed)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal transaction definition file: %w", err)
	}

	return parsed, nil
}

func (r transactionRunner) StartRun(ctx context.Context, resource any, runInfo openapi.RunInformation) (RunResult, error) {
	tran := resource.(openapi.TransactionResource)
	run, resp, err := r.openapiClient.ApiApi.
		RunTransaction(ctx, tran.Spec.GetId()).
		RunInformation(runInfo).
		Execute()

	err = HandleRunError(resp, err)
	if err != nil {
		return RunResult{}, err
	}

	full, err := r.client.Get(ctx, tran.Spec.GetId(), jsonFormat)
	if err != nil {
		return RunResult{}, fmt.Errorf("cannot get full transaction '%s': %w", tran.Spec.GetId(), err)
	}
	err = jsonFormat.Unmarshal([]byte(full), &tran)
	if err != nil {
		return RunResult{}, fmt.Errorf("cannot get full transaction '%s': %w", tran.Spec.GetId(), err)
	}

	return RunResult{
		Resource: tran,
		Run:      *run,
	}, nil
}

func (r transactionRunner) UpdateResult(ctx context.Context, result RunResult) (RunResult, error) {
	transaction := result.Resource.(openapi.TransactionResource)
	run := result.Run.(openapi.TransactionRun)
	runID, err := strconv.Atoi(run.GetId())
	if err != nil {
		return RunResult{}, fmt.Errorf("invalid transaction run id format: %w", err)
	}

	updated, _, err := r.openapiClient.ApiApi.
		GetTransactionRun(ctx, transaction.Spec.GetId(), int32(runID)).
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
		Resource: transaction,
		Run:      *updated,
		Finished: isStateFinished(updated.GetState()),
		Passed:   passed,
	}, nil
}

func (r transactionRunner) JUnitResult(ctx context.Context, result RunResult) (string, error) {
	return "", ErrJUnitNotSupported
}

func (r transactionRunner) FormatResult(result RunResult, format string) string {
	transaction := result.Resource.(openapi.TransactionResource)
	run := result.Run.(openapi.TransactionRun)

	tro := formatters.TransactionRunOutput{
		HasResults:  result.Finished,
		Transaction: transaction.GetSpec(),
		Run:         run,
	}

	return r.formatter.Format(tro, formatters.Output(format))
}
