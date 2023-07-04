package http

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kubeshop/tracetest/server/assertions/selectors"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/http/mappings"
	"github.com/kubeshop/tracetest/server/http/validation"
	"github.com/kubeshop/tracetest/server/junit"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/transaction"
	"go.opentelemetry.io/otel/trace"
)

var IDGen = id.NewRandGenerator()

type controller struct {
	tracer          trace.Tracer
	runner          runner
	newTraceDBFn    func(ds datastore.DataStore) (tracedb.TraceDB, error)
	mappers         mappings.Mappings
	triggerRegistry *trigger.Registry
	version         string

	testDB                   model.Repository
	transactionRepository    transactionsRepository
	transactionRunRepository transactionRunRepository
	testRepository           test.Repository
	testRunRepository        test.RunRepository
	environmentGetter        environmentGetter
}

type transactionsRepository interface {
	SetID(transaction.Transaction, id.ID) transaction.Transaction
	IDExists(ctx context.Context, id id.ID) (bool, error)
	ListAugmented(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]transaction.Transaction, error)
	GetAugmented(context.Context, id.ID) (transaction.Transaction, error)
	Count(ctx context.Context, query string) (int, error)
	Get(context.Context, id.ID) (transaction.Transaction, error)
	GetVersion(context.Context, id.ID, int) (transaction.Transaction, error)
	Create(context.Context, transaction.Transaction) (transaction.Transaction, error)
	Update(context.Context, transaction.Transaction) (transaction.Transaction, error)
}

type transactionRunRepository interface {
	GetTransactionRun(ctx context.Context, transactionID id.ID, runID int) (transaction.TransactionRun, error)
	GetTransactionsRuns(ctx context.Context, transactionID id.ID, take, skip int32) ([]transaction.TransactionRun, error)
	DeleteTransactionRun(ctx context.Context, tr transaction.TransactionRun) error
}

type runner interface {
	StopTest(testID id.ID, runID int)
	RunTest(ctx context.Context, test test.Test, rm test.RunMetadata, env environment.Environment) test.Run
	RunTransaction(ctx context.Context, tr transaction.Transaction, rm test.RunMetadata, env environment.Environment) transaction.TransactionRun
	RunAssertions(ctx context.Context, request executor.AssertionRequest)
}

type environmentGetter interface {
	Get(context.Context, id.ID) (environment.Environment, error)
}

func NewController(
	testDB model.Repository,
	transactionRepository transactionsRepository,
	transactionRunRepository transactionRunRepository,
	testRepository test.Repository,
	testRunRepository test.RunRepository,
	newTraceDBFn func(ds datastore.DataStore) (tracedb.TraceDB, error),
	runner runner,
	mappers mappings.Mappings,
	envGetter environmentGetter,
	triggerRegistry *trigger.Registry,
	tracer trace.Tracer,
	version string,
) openapi.ApiApiServicer {
	return &controller{
		tracer:                   tracer,
		testDB:                   testDB,
		transactionRepository:    transactionRepository,
		transactionRunRepository: transactionRunRepository,
		testRepository:           testRepository,
		testRunRepository:        testRunRepository,
		environmentGetter:        envGetter,
		runner:                   runner,
		newTraceDBFn:             newTraceDBFn,
		mappers:                  mappers,
		triggerRegistry:          triggerRegistry,
		version:                  version,
	}
}

func handleDBError(err error) openapi.ImplResponse {
	switch {
	case errors.Is(testdb.ErrNotFound, err):
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
	events, err := c.testDB.GetTestRunEvents(ctx, id.ID(testID), int(runID))
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
	test, err := c.testRepository.Get(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	run, err := c.testRunRepository.GetRun(ctx, id.ID(testID), int(runID))
	if err != nil {
		return handleDBError(err), err
	}

	newTestRun, err := c.testRunRepository.CreateRun(ctx, test, run.Copy())
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	newTestRun = newTestRun.SuccessfullyPolledTraces(run.Trace)
	err = c.testRunRepository.UpdateRun(ctx, newTestRun)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	assertionRequest := executor.AssertionRequest{
		Test: test,
		Run:  newTestRun,
	}

	c.runner.RunAssertions(ctx, assertionRequest)

	return openapi.Response(http.StatusOK, c.mappers.Out.Run(&newTestRun)), nil
}

func (c *controller) RunTest(ctx context.Context, testID string, runInformation openapi.RunInformation) (openapi.ImplResponse, error) {
	test, err := c.testRepository.Get(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	metadata := metadata(runInformation.Metadata)
	variablesEnv := c.mappers.In.Environment(openapi.Environment{
		Values: runInformation.Variables,
	})

	environment, err := getEnvironment(ctx, c.environmentGetter, runInformation.EnvironmentId, variablesEnv)
	if err != nil {
		return handleDBError(err), err
	}

	missingVariablesError, err := validation.ValidateMissingVariables(ctx, c.testRepository, c.testRunRepository, test, environment)
	if err != nil {
		if err == validation.ErrMissingVariables {
			return openapi.Response(http.StatusUnprocessableEntity, missingVariablesError), nil
		}

		return handleDBError(err), err
	}

	run := c.runner.RunTest(ctx, test, metadata, environment)

	return openapi.Response(200, c.mappers.Out.Run(&run)), nil
}

func (c *controller) StopTestRun(_ context.Context, testID string, runID int32) (openapi.ImplResponse, error) {
	c.runner.StopTest(id.ID(testID), int(runID))

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

func (c *controller) ExecuteDefinition(ctx context.Context, testDefinition openapi.TextDefinition) (openapi.ImplResponse, error) {
	def, err := yaml.Decode([]byte(testDefinition.Content))
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	if test, err := def.Test(); err == nil {
		return c.executeTest(ctx, test.Model(), testDefinition.RunInformation)
	}

	if transaction, err := def.Transaction(); err == nil {
		return c.executeTransaction(ctx, transaction.Model(), testDefinition.RunInformation)
	}

	return openapi.Response(http.StatusUnprocessableEntity, nil), nil
}

func (c *controller) executeTransaction(ctx context.Context, tran transaction.Transaction, runInfo openapi.RunInformation) (openapi.ImplResponse, error) {
	// create or update transaction
	resp, err := c.doCreateTransaction(ctx, tran)
	if err != nil {
		if errors.Is(err, errTransactionExists) {
			resp, err = c.doUpdateTransaction(ctx, tran.ID, tran)
			if err != nil {
				return resp, err
			}
		} else {
			return resp, err
		}
	} else {
		// the transaction was created, make sure we have the correct ID
		tran = resp.Body.(transaction.Transaction)
	}

	transactionID := tran.ID

	// transaction ready, execute it
	resp, err = c.RunTransaction(ctx, transactionID.String(), runInfo)
	if resp.Code != http.StatusOK || err != nil {
		return resp, err
	}

	res := openapi.ExecuteDefinitionResponse{
		Id:    transactionID.String(),
		RunId: resp.Body.(openapi.TransactionRun).Id,
		Type:  yaml.FileTypeTransaction.String(),
	}
	return openapi.Response(200, res), nil
}

var errTransactionExists = errors.New("transaction already exists")

func (c *controller) doCreateTransaction(ctx context.Context, transaction transaction.Transaction) (openapi.ImplResponse, error) {
	// if they try to create a transaction with preset ID, we need to make sure that ID doesn't exists already
	if transaction.HasID() {
		exists, err := c.transactionRepository.IDExists(ctx, transaction.ID)
		if err != nil {
			return handleDBError(err), err
		}

		if exists {
			err := fmt.Errorf(`cannot create transaction with ID "%s: %w`, transaction.ID, errTransactionExists)
			r := map[string]string{
				"error": err.Error(),
			}
			return openapi.Response(http.StatusBadRequest, r), err
		}
	} else {
		transaction = c.transactionRepository.SetID(transaction, id.GenerateID())
	}

	transaction, err := c.transactionRepository.Create(ctx, transaction)
	if err != nil {
		return handleDBError(err), err
	}
	transaction, err = c.transactionRepository.GetAugmented(ctx, transaction.ID)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, transaction), nil
}

func (c *controller) GetTransactionVersionDefinitionFile(ctx context.Context, transactionId string, version int32) (openapi.ImplResponse, error) {
	transaction, err := c.transactionRepository.GetVersion(ctx, id.ID(transactionId), int(version))
	if err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), err
	}

	return openapi.Response(200, transaction), nil
}

func (c *controller) doUpdateTransaction(ctx context.Context, transactionID id.ID, updated transaction.Transaction) (openapi.ImplResponse, error) {
	transaction, err := c.transactionRepository.Get(ctx, transactionID)
	if err != nil {
		return handleDBError(err), err
	}

	updated.Version = transaction.Version
	updated.ID = transaction.ID

	_, err = c.transactionRepository.Update(ctx, updated)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(204, nil), nil
}

func metadata(in *map[string]string) test.RunMetadata {
	if in == nil {
		return nil
	}

	return test.RunMetadata(*in)
}

func getEnvironment(ctx context.Context, environmentRepository environmentGetter, environmentId string, variablesEnv environment.Environment) (environment.Environment, error) {
	if environmentId != "" {
		environment, err := environmentRepository.Get(ctx, id.ID(environmentId))

		if err != nil {
			return variablesEnv, err
		}

		return environment.Merge(variablesEnv), nil
	}

	return variablesEnv, nil
}

func (c *controller) executeTest(ctx context.Context, test test.Test, runInfo openapi.RunInformation) (openapi.ImplResponse, error) {
	createdTest, err := c.testRepository.Create(ctx, test)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	// test ready, execute it
	resp, err := c.RunTest(ctx, createdTest.ID.String(), runInfo)
	if resp.Code != http.StatusOK || err != nil {
		return resp, err
	}

	res := openapi.ExecuteDefinitionResponse{
		Id:    createdTest.ID.String(),
		RunId: resp.Body.(openapi.TestRun).Id,
		Type:  yaml.FileTypeTest.String(),
	}
	return openapi.Response(200, res), nil
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
func (c *controller) RunTransaction(ctx context.Context, transactionID string, runInformation openapi.RunInformation) (openapi.ImplResponse, error) {
	transaction, err := c.transactionRepository.GetAugmented(ctx, id.ID(transactionID))
	if err != nil {
		return handleDBError(err), err
	}

	metadata := metadata(runInformation.Metadata)
	variablesEnv := c.mappers.In.Environment(openapi.Environment{
		Values: runInformation.Variables,
	})
	environment, err := getEnvironment(ctx, c.environmentGetter, runInformation.EnvironmentId, variablesEnv)

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

	run := c.runner.RunTransaction(ctx, transaction, metadata, environment)

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
