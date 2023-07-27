package http

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kubeshop/tracetest/server/assertions/selectors"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/http/mappings"
	"github.com/kubeshop/tracetest/server/http/validation"
	"github.com/kubeshop/tracetest/server/junit"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/transaction"
	"go.opentelemetry.io/otel/trace"
)

type controller struct {
	tracer trace.Tracer

	testRunner        testRunner
	transactionRunner transactionRunner

	testRunEvents            model.TestRunEventRepository
	testRunRepository        test.RunRepository
	testRepository           testsRepository
	transactionRepository    transactionsRepository
	transactionRunRepository transactionRunRepository

	environmentGetter environmentGetter
	newTraceDBFn      func(ds datastore.DataStore) (tracedb.TraceDB, error)
	mappers           mappings.Mappings
	version           string
}

type transactionsRepository interface {
	GetAugmented(context.Context, id.ID) (transaction.Transaction, error)
	ListAugmented(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]transaction.Transaction, error)
	Count(ctx context.Context, query string) (int, error)
	GetVersion(context.Context, id.ID, int) (transaction.Transaction, error)
}

type testsRepository interface {
	Get(context.Context, id.ID) (test.Test, error)
	GetAugmented(context.Context, id.ID) (test.Test, error)
	ListAugmented(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]test.Test, error)
	Count(ctx context.Context, query string) (int, error)
	GetVersion(context.Context, id.ID, int) (test.Test, error)
	Create(context.Context, test.Test) (test.Test, error)
}

type transactionRunRepository interface {
	GetTransactionRun(ctx context.Context, transactionID id.ID, runID int) (transaction.TransactionRun, error)
	GetTransactionsRuns(ctx context.Context, transactionID id.ID, take, skip int32) ([]transaction.TransactionRun, error)
	DeleteTransactionRun(ctx context.Context, tr transaction.TransactionRun) error
}

type testRunner interface {
	StopTest(_ context.Context, testID id.ID, runID int)
	Run(context.Context, test.Test, test.RunMetadata, environment.Environment, *[]testrunner.RequiredGate) test.Run
	Rerun(_ context.Context, _ test.Test, runID int) test.Run
}

type transactionRunner interface {
	Run(context.Context, transaction.Transaction, test.RunMetadata, environment.Environment, *[]testrunner.RequiredGate) transaction.TransactionRun
}

type environmentGetter interface {
	Get(context.Context, id.ID) (environment.Environment, error)
}

func NewController(
	tracer trace.Tracer,

	testRunner testRunner,
	transactionRunner transactionRunner,

	testRunEvents model.TestRunEventRepository,
	transactionRepository transactionsRepository,
	transactionRunRepository transactionRunRepository,
	testRepository testsRepository,
	testRunRepository test.RunRepository,
	environmentGetter environmentGetter,

	newTraceDBFn func(ds datastore.DataStore) (tracedb.TraceDB, error),
	mappers mappings.Mappings,
	version string,
) openapi.ApiApiServicer {
	return &controller{
		testRunEvents:            testRunEvents,
		transactionRepository:    transactionRepository,
		transactionRunRepository: transactionRunRepository,
		testRepository:           testRepository,
		testRunRepository:        testRunRepository,
		environmentGetter:        environmentGetter,

		testRunner:        testRunner,
		transactionRunner: transactionRunner,

		tracer:       tracer,
		newTraceDBFn: newTraceDBFn,
		mappers:      mappers,
		version:      version,
	}
}

func handleDBError(err error) openapi.ImplResponse {
	switch {
	case errors.Is(sql.ErrNoRows, err):
		return openapi.Response(http.StatusNotFound, err.Error())
	default:
		return openapi.Response(http.StatusInternalServerError, err.Error())
	}
}

func (c *controller) GetTestSpecs(ctx context.Context, testID string) (openapi.ImplResponse, error) {
	test, err := c.testRepository.Get(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.mappers.Out.Specs(test.Specs)), nil
}

func (c *controller) GetTestResultSelectedSpans(ctx context.Context, testID string, runID int32, selectorQuery string) (openapi.ImplResponse, error) {
	selector, err := selectors.New(selectorQuery)

	if err != nil {
		return handleDBError(err), err
	}

	run, err := c.testRunRepository.GetRun(ctx, id.ID(testID), int(runID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, ""), nil
	}

	if run.Trace == nil {
		return openapi.Response(http.StatusInternalServerError, "trace not available"), nil
	}

	selectedSpans := selector.Filter(*run.Trace)
	selectedSpanIds := make([]string, len(selectedSpans))

	for i, span := range selectedSpans {
		selectedSpanIds[i] = hex.EncodeToString(span.ID[:])
	}

	res := openapi.SelectedSpansResult{
		Selector: c.mappers.Out.Selector(test.SpanQuery(selectorQuery)),
		SpanIds:  selectedSpanIds,
	}

	return openapi.Response(http.StatusOK, res), nil
}

func (c *controller) GetTestRun(ctx context.Context, testID string, runID int32) (openapi.ImplResponse, error) {
	run, err := c.testRunRepository.GetRun(ctx, id.ID(testID), int(runID))
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.mappers.Out.Run(&run)), nil
}

func (c *controller) GetTestRunEvents(ctx context.Context, testID string, runID int32) (openapi.ImplResponse, error) {
	events, err := c.testRunEvents.GetTestRunEvents(ctx, id.ID(testID), int(runID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), err
	}

	return openapi.Response(http.StatusOK, c.mappers.Out.TestRunEvents(events)), nil
}

func (c *controller) DeleteTestRun(ctx context.Context, testID string, runID int32) (openapi.ImplResponse, error) {
	run, err := c.testRunRepository.GetRun(ctx, id.ID(testID), int(runID))
	if err != nil {
		return handleDBError(err), err
	}

	err = c.testRunRepository.DeleteRun(ctx, run)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(204, nil), nil
}

type paginated[T any] struct {
	items []T
	count int
}

func (c *controller) GetTestRuns(ctx context.Context, testID string, take, skip int32) (openapi.ImplResponse, error) {
	if take == 0 {
		take = 20
	}

	test, err := c.testRepository.Get(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	runs, err := c.testRunRepository.GetTestRuns(ctx, test, take, skip)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, paginated[openapi.TestRun]{
		items: c.mappers.Out.Runs(runs),
		count: len(runs), // TODO: find a way of returning the proper number
	}), nil
}

func (c *controller) RerunTestRun(ctx context.Context, testID string, runID int32) (openapi.ImplResponse, error) {
	testObj, err := c.testRepository.GetAugmented(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	newTestRun := c.testRunner.Rerun(ctx, testObj, int(runID))

	return openapi.Response(http.StatusOK, c.mappers.Out.Run(&newTestRun)), nil
}

func (c *controller) RunTest(ctx context.Context, testID string, runInfo openapi.RunInformation) (openapi.ImplResponse, error) {
	test, err := c.testRepository.GetAugmented(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	variablesEnv := c.mappers.In.Environment(openapi.Environment{
		Values: runInfo.Variables,
	})

	environment, err := c.getEnvironment(ctx, runInfo.EnvironmentId, variablesEnv)
	if err != nil {
		return handleDBError(err), err
	}

	missingVariablesError, err := validation.ValidateMissingVariables(ctx, c.testRepository, c.testRunRepository, test, environment)
	if err != nil {
		if errors.Is(err, validation.ErrMissingVariables) {
			return openapi.Response(http.StatusUnprocessableEntity, missingVariablesError), nil
		}

		return handleDBError(err), err
	}

	requiredGates := c.mappers.In.RequiredGates(runInfo.RequiredGates)

	run := c.testRunner.Run(ctx, test, metadata(runInfo.Metadata), environment, requiredGates)

	return openapi.Response(200, c.mappers.Out.Run(&run)), nil
}

func (c *controller) StopTestRun(ctx context.Context, testID string, runID int32) (openapi.ImplResponse, error) {
	c.testRunner.StopTest(ctx, id.ID(testID), int(runID))

	return openapi.Response(http.StatusOK, map[string]string{"result": "success"}), nil
}

func (c *controller) DryRunAssertion(ctx context.Context, testID string, runID int32, def openapi.TestSpecs) (openapi.ImplResponse, error) {
	run, err := c.testRunRepository.GetRun(ctx, id.ID(testID), int(runID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, ""), nil
	}

	if run.Trace == nil {
		return openapi.Response(http.StatusUnprocessableEntity, fmt.Sprintf(`run "%d" has no trace associated`, runID)), nil
	}

	definition := c.mappers.In.Definition(def.Specs)

	ds := []expression.DataStore{expression.EnvironmentDataStore{
		Values: run.Environment.Values,
	}}

	assertionExecutor := executor.NewAssertionExecutor(c.tracer)

	results, allPassed := assertionExecutor.Assert(ctx, definition, *run.Trace, ds)
	res := c.mappers.Out.Result(&test.RunResults{
		AllPassed: allPassed,
		Results:   results,
	})

	return openapi.Response(200, res), nil
}

func (c *controller) GetRunResultJUnit(ctx context.Context, testID string, runID int32) (openapi.ImplResponse, error) {
	run, err := c.testRunRepository.GetRun(ctx, id.ID(testID), int(runID))
	if err != nil {
		return handleDBError(err), err
	}

	test, err := c.testRepository.GetVersion(ctx, id.ID(testID), run.TestVersion)
	if err != nil {
		return handleDBError(err), err
	}

	res, err := junit.FromRunResult(test, run)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(200, res), nil
}

func (c controller) GetTestVersion(ctx context.Context, testID string, version int32) (openapi.ImplResponse, error) {
	test, err := c.testRepository.GetVersion(ctx, id.ID(testID), int(version))
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.mappers.Out.Test(test)), nil
}

func (c controller) ExportTestRun(ctx context.Context, testID string, runID int32) (openapi.ImplResponse, error) {
	run, err := c.testRunRepository.GetRun(ctx, id.ID(testID), int(runID))
	if err != nil {
		return handleDBError(err), err
	}

	test, err := c.testRepository.GetVersion(ctx, id.ID(testID), run.TestVersion)
	if err != nil {
		return handleDBError(err), err
	}

	response := openapi.ExportedTestInformation{
		Test: c.mappers.Out.Test(test),
		Run:  c.mappers.Out.Run(&run),
	}

	return openapi.Response(http.StatusOK, response), nil
}

func (c controller) ImportTestRun(ctx context.Context, exportedTest openapi.ExportedTestInformation) (openapi.ImplResponse, error) {
	test, err := c.mappers.In.Test(exportedTest.Test)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), nil
	}

	run, err := c.mappers.In.Run(exportedTest.Run)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), nil
	}

	createdTest, err := c.testRepository.Create(ctx, test)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	createdRun, err := c.testRunRepository.CreateRun(ctx, createdTest, *run)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	createdRun.State = run.State

	err = c.testRunRepository.UpdateRun(ctx, createdRun)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	response := openapi.ExportedTestInformation{
		Test: c.mappers.Out.Test(createdTest),
		Run:  c.mappers.Out.Run(&createdRun),
	}

	return openapi.Response(http.StatusOK, response), nil
}

func metadata(in *map[string]string) test.RunMetadata {
	if in == nil {
		return nil
	}

	return test.RunMetadata(*in)
}

func (c *controller) getEnvironment(ctx context.Context, environmentID string, variablesEnv environment.Environment) (environment.Environment, error) {
	if environmentID == "" {
		return variablesEnv, nil
	}

	environment, err := c.environmentGetter.Get(ctx, id.ID(environmentID))

	if err != nil {
		return variablesEnv, err
	}

	return environment.Merge(variablesEnv), nil

}

// Expressions
func (c *controller) ExpressionResolve(ctx context.Context, in openapi.ResolveRequestInfo) (openapi.ImplResponse, error) {
	dsList, err := c.buildDataStores(ctx, in)

	if err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), err
	}

	resolvedValues := make([]string, len(dsList))
	for i, ds := range dsList {
		parsed, err := expression.NewExecutor(ds...).ResolveStatement(in.Expression)

		if err != nil {
			return openapi.Response(http.StatusBadRequest, err.Error()), err
		}

		resolvedValues[i] = parsed
	}

	return openapi.Response(200, openapi.ResolveResponseInfo{ResolvedValues: resolvedValues}), nil
}

func (c *controller) buildDataStores(ctx context.Context, info openapi.ResolveRequestInfo) ([][]expression.DataStore, error) {
	context := info.Context

	ds := []expression.DataStore{}

	if context.EnvironmentId != "" {
		environment, err := c.environmentGetter.Get(ctx, id.ID(context.EnvironmentId))

		if err != nil {
			return [][]expression.DataStore{}, err
		}

		ds = append([]expression.DataStore{expression.EnvironmentDataStore{
			Values: environment.Values,
		}}, ds...)
	}

	if context.RunId != "" && context.TestId != "" {
		runId, err := strconv.Atoi(context.RunId)

		if err != nil {
			return [][]expression.DataStore{}, err
		}

		run, err := c.testRunRepository.GetRun(ctx, id.ID(context.TestId), runId)
		if err != nil {
			return [][]expression.DataStore{}, err
		}

		if context.SpanId != "" {
			spanId, err := trace.SpanIDFromHex(context.SpanId)

			if err != nil {
				return [][]expression.DataStore{}, err
			}

			span := run.Trace.Flat[spanId]

			ds = append([]expression.DataStore{expression.AttributeDataStore{
				Span: *span,
			}}, ds...)
		}

		selector, err := selectors.New(context.Selector)
		if err != nil {
			return [][]expression.DataStore{}, err
		}

		spans := selector.Filter(*run.Trace)
		ds = append([]expression.DataStore{expression.MetaAttributesDataStore{
			SelectedSpans: spans,
		}}, ds...)

		if context.SpanId == "" {
			dsList := make([][]expression.DataStore, len(spans))
			for i, span := range spans {

				dsList[i] = append([]expression.DataStore{expression.AttributeDataStore{
					Span: span,
				}}, ds...)
			}

			return dsList, nil
		}
	}

	return [][]expression.DataStore{ds}, nil
}

func (c *controller) GetTransactionVersion(ctx context.Context, tID string, version int32) (openapi.ImplResponse, error) {
	transaction, err := c.transactionRepository.GetVersion(ctx, id.ID(tID), int(version))

	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(http.StatusOK, transaction), nil
}

// RunTransaction implements openapi.ApiApiServicer
func (c *controller) RunTransaction(ctx context.Context, transactionID string, runInfo openapi.RunInformation) (openapi.ImplResponse, error) {
	transaction, err := c.transactionRepository.GetAugmented(ctx, id.ID(transactionID))
	if err != nil {
		return handleDBError(err), err
	}

	metadata := metadata(runInfo.Metadata)
	variablesEnv := c.mappers.In.Environment(openapi.Environment{
		Values: runInfo.Variables,
	})

	environment, err := c.getEnvironment(ctx, runInfo.EnvironmentId, variablesEnv)
	if err != nil {
		return handleDBError(err), err
	}

	missingVariablesError, err := validation.ValidateMissingVariablesFromTransaction(ctx, c.testRepository, c.testRunRepository, transaction, environment)
	if err != nil {
		if err == validation.ErrMissingVariables {
			return openapi.Response(http.StatusUnprocessableEntity, missingVariablesError), nil
		}

		return handleDBError(err), err
	}

	requiredGates := c.mappers.In.RequiredGates(runInfo.RequiredGates)

	run := c.transactionRunner.Run(ctx, transaction, metadata, environment, requiredGates)

	return openapi.Response(http.StatusOK, c.mappers.Out.TransactionRun(run)), nil
}

func (c *controller) GetTransactionRun(ctx context.Context, transactionId string, runId int32) (openapi.ImplResponse, error) {
	run, err := c.transactionRunRepository.GetTransactionRun(ctx, id.ID(transactionId), int(runId))
	if err != nil {
		return handleDBError(err), err
	}

	openapiRun := c.mappers.Out.TransactionRun(run)
	return openapi.Response(http.StatusOK, openapiRun), nil
}

func (c *controller) GetTransactionRuns(ctx context.Context, transactionId string, take, skip int32) (openapi.ImplResponse, error) {
	runs, err := c.transactionRunRepository.GetTransactionsRuns(ctx, id.ID(transactionId), take, skip)
	if err != nil {
		return handleDBError(err), err
	}

	openapiRuns := make([]openapi.TransactionRun, 0, len(runs))
	for _, run := range runs {
		openapiRuns = append(openapiRuns, c.mappers.Out.TransactionRun(run))
	}

	return openapi.Response(http.StatusOK, openapiRuns), nil
}

func (c *controller) DeleteTransactionRun(ctx context.Context, transactionId string, runId int32) (openapi.ImplResponse, error) {
	run, err := c.transactionRunRepository.GetTransactionRun(ctx, id.ID(transactionId), int(runId))
	if err != nil {
		return handleDBError(err), err
	}

	err = c.transactionRunRepository.DeleteTransactionRun(ctx, run)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

func (c *controller) GetResources(ctx context.Context, take, skip int32, query, sortBy, sortDirection string) (openapi.ImplResponse, error) {
	// TODO: this is endpoint is a hack to unblock the team quickly.
	// This is not production ready because it might take too long to respond if there are numerous
	// transactions and transaction.

	if take == 0 {
		take = 20
	}

	newTake := take + skip

	transactions, err := c.transactionRepository.ListAugmented(ctx, int(newTake), 0, query, sortBy, sortDirection)
	if err != nil {
		return handleDBError(err), err
	}

	transactionCount, err := c.transactionRepository.Count(ctx, query)
	if err != nil {
		return handleDBError(err), err
	}

	tests, err := c.testRepository.ListAugmented(ctx, int(newTake), 0, query, sortBy, sortDirection)
	if err != nil {
		return handleDBError(err), err
	}

	testCount, err := c.testRepository.Count(ctx, query)
	if err != nil {
		return handleDBError(err), err
	}

	totalResources := transactionCount + testCount

	items := takeResources(transactions, tests, take, skip)

	paginatedResponse := paginated[openapi.Resource]{
		items: items,
		count: totalResources,
	}

	return openapi.Response(http.StatusOK, paginatedResponse), nil
}

func takeResources(transactions []transaction.Transaction, tests []test.Test, take, skip int32) []openapi.Resource {
	numItems := len(transactions) + len(tests)
	items := make([]openapi.Resource, numItems)
	maxNumItems := len(transactions) + len(tests)
	currentNumberItens := 0

	var i, j int
	for currentNumberItens < int(numItems) && currentNumberItens < maxNumItems {
		if i >= len(transactions) {
			test := tests[j]
			testInterface := any(test)
			items[currentNumberItens] = openapi.Resource{Type: "test", Item: &testInterface}
			j++
			currentNumberItens++
			continue
		}

		if j >= len(tests) {
			transaction := transactions[i]
			transactionInterface := any(transaction)
			items[currentNumberItens] = openapi.Resource{Type: "transaction", Item: &transactionInterface}
			i++
			currentNumberItens++
			continue
		}

		transaction := transactions[i]
		test := tests[j]
		transactionInterface := any(transaction)
		testInterface := any(test)

		if transaction.CreatedAt.After(*test.CreatedAt) {
			items[currentNumberItens] = openapi.Resource{Type: "transaction", Item: &transactionInterface}
			i++
		} else {
			items[currentNumberItens] = openapi.Resource{Type: "test", Item: &testInterface}
			j++
		}

		currentNumberItens++
	}

	upperLimit := int(skip + take)
	if upperLimit > currentNumberItens {
		upperLimit = currentNumberItens
	}

	return items[skip:upperLimit]
}

// TestConnection implements openapi.ApiApiService
func (c *controller) TestConnection(ctx context.Context, dataStore openapi.DataStore) (openapi.ImplResponse, error) {
	ds := c.mappers.In.DataStore(dataStore)

	if err := ds.Validate(); err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), err
	}

	tdb, err := c.newTraceDBFn(ds)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), err
	}

	testResult := model.ConnectionResult{}
	statusCode := http.StatusOK

	if testableTraceDB, ok := tdb.(tracedb.TestableTraceDB); ok {
		testResult = testableTraceDB.TestConnection(ctx)
		statusCode = http.StatusOK
		if !testResult.HasSucceed() {
			statusCode = http.StatusUnprocessableEntity
		}
	}

	return openapi.Response(statusCode, c.mappers.Out.ConnectionTestResult(testResult)), nil
}

func (c *controller) GetVersion(ctx context.Context) (openapi.ImplResponse, error) {
	version := openapi.Version{
		Version: c.version,
	}

	return openapi.Response(http.StatusOK, version), nil
}
