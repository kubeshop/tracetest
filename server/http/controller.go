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
	"github.com/kubeshop/tracetest/server/model/yaml/yamlconvert"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/tests"
	"github.com/kubeshop/tracetest/server/tracedb"
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

	testDB            model.Repository
	transactions      transactionsRepository
	environmentGetter environmentGetter
}

type transactionsRepository interface {
	SetID(tests.Transaction, id.ID) tests.Transaction
	IDExists(ctx context.Context, id id.ID) (bool, error)
	ListAugmented(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]tests.Transaction, error)
	GetAugmented(context.Context, id.ID) (tests.Transaction, error)
	Count(ctx context.Context, query string) (int, error)
	Get(context.Context, id.ID) (tests.Transaction, error)
	GetVersion(context.Context, id.ID, int) (tests.Transaction, error)
	Create(context.Context, tests.Transaction) (tests.Transaction, error)
	Update(context.Context, tests.Transaction) (tests.Transaction, error)

	GetTransactionRun(ctx context.Context, transactionID id.ID, runID int) (tests.TransactionRun, error)
	GetTransactionsRuns(ctx context.Context, transactionID id.ID, take, skip int32) ([]tests.TransactionRun, error)
	DeleteTransactionRun(ctx context.Context, tr tests.TransactionRun) error
}

type runner interface {
	StopTest(testID id.ID, runID int)
	RunTest(ctx context.Context, test model.Test, rm model.RunMetadata, env environment.Environment) model.Run
	RunTransaction(ctx context.Context, tr tests.Transaction, rm model.RunMetadata, env environment.Environment) tests.TransactionRun
	RunAssertions(ctx context.Context, request executor.AssertionRequest)
}

type environmentGetter interface {
	Get(context.Context, id.ID) (environment.Environment, error)
}

func NewController(
	testDB model.Repository,
	transactions transactionsRepository,
	newTraceDBFn func(ds datastore.DataStore) (tracedb.TraceDB, error),
	runner runner,
	mappers mappings.Mappings,
	envGetter environmentGetter,
	triggerRegistry *trigger.Registry,
	tracer trace.Tracer,
	version string,
) openapi.ApiApiServicer {
	return &controller{
		tracer:            tracer,
		testDB:            testDB,
		transactions:      transactions,
		environmentGetter: envGetter,
		runner:            runner,
		newTraceDBFn:      newTraceDBFn,
		mappers:           mappers,
		triggerRegistry:   triggerRegistry,
		version:           version,
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

func (c *controller) CreateTest(ctx context.Context, in openapi.Test) (openapi.ImplResponse, error) {
	test, err := c.mappers.In.Test(in)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), nil
	}

	return c.doCreateTest(ctx, test)
}

var errTestExists = errors.New("test already exists")

func (c *controller) doCreateTest(ctx context.Context, test model.Test) (openapi.ImplResponse, error) {
	// if they try to create a test with preset ID, we need to make sure that ID doesn't exist already
	if test.HasID() {
		exists, err := c.testDB.TestIDExists(ctx, test.ID)

		if err != nil {
			return handleDBError(err), err
		}

		if exists {
			err := fmt.Errorf(`cannot create test with ID "%s: %w`, test.ID, errTestExists)
			r := map[string]string{
				"error": err.Error(),
			}
			return openapi.Response(http.StatusBadRequest, r), err
		}
	}

	test, err := c.testDB.CreateTest(ctx, test)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(200, c.mappers.Out.Test(test)), nil
}

func (c *controller) DeleteTest(ctx context.Context, testID string) (openapi.ImplResponse, error) {
	test, err := c.testDB.GetLatestTestVersion(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	err = c.testDB.DeleteTest(ctx, test)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(204, nil), nil
}

func (c *controller) GetTest(ctx context.Context, testID string) (openapi.ImplResponse, error) {
	test, err := c.testDB.GetLatestTestVersion(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.mappers.Out.Test(test)), nil
}

func (c *controller) GetTestSpecs(ctx context.Context, testID string) (openapi.ImplResponse, error) {
	test, err := c.testDB.GetLatestTestVersion(ctx, id.ID(testID))
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

	run, err := c.testDB.GetRun(ctx, id.ID(testID), int(runID))
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
		Selector: c.mappers.Out.Selector(model.SpanQuery(selectorQuery)),
		SpanIds:  selectedSpanIds,
	}

	return openapi.Response(http.StatusOK, res), nil
}

func (c *controller) GetTestRun(ctx context.Context, testID string, runID int32) (openapi.ImplResponse, error) {
	run, err := c.testDB.GetRun(ctx, id.ID(testID), int(runID))
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
	run, err := c.testDB.GetRun(ctx, id.ID(testID), int(runID))
	if err != nil {
		return handleDBError(err), err
	}

	err = c.testDB.DeleteRun(ctx, run)
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

	test, err := c.testDB.GetLatestTestVersion(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	runs, err := c.testDB.GetTestRuns(ctx, test, take, skip)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, paginated[openapi.TestRun]{
		items: c.mappers.Out.Runs(runs.Items),
		count: runs.TotalCount,
	}), nil
}

func (c *controller) GetTests(ctx context.Context, take, skip int32, query string, sortBy string, sortDirection string) (openapi.ImplResponse, error) {
	if take == 0 {
		take = 20
	}

	tests, err := c.testDB.GetTests(ctx, take, skip, query, sortBy, sortDirection)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, paginated[openapi.Test]{
		items: c.mappers.Out.Tests(tests.Items),
		count: tests.TotalCount,
	}), nil
}

func (c *controller) RerunTestRun(ctx context.Context, testID string, runID int32) (openapi.ImplResponse, error) {
	test, err := c.testDB.GetLatestTestVersion(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	run, err := c.testDB.GetRun(ctx, id.ID(testID), int(runID))
	if err != nil {
		return handleDBError(err), err
	}

	newTestRun, err := c.testDB.CreateRun(ctx, test, run.Copy())
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	newTestRun = newTestRun.SuccessfullyPolledTraces(run.Trace)
	err = c.testDB.UpdateRun(ctx, newTestRun)
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
	test, err := c.testDB.GetLatestTestVersion(ctx, id.ID(testID))
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

	missingVariablesError, err := validation.ValidateMissingVariables(ctx, c.testDB, test, environment)
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

func (c *controller) UpdateTest(ctx context.Context, testID string, in openapi.Test) (openapi.ImplResponse, error) {
	updated, err := c.mappers.In.Test(in)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), nil
	}

	return c.doUpdateTest(ctx, id.ID(testID), updated)
}

func (c *controller) doUpdateTest(ctx context.Context, testID id.ID, updated model.Test) (openapi.ImplResponse, error) {
	test, err := c.testDB.GetLatestTestVersion(ctx, testID)
	if err != nil {
		return handleDBError(err), err
	}

	updated.Version = test.Version
	updated.ID = test.ID

	_, err = c.testDB.UpdateTest(ctx, updated)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(204, nil), nil
}

func (c *controller) DryRunAssertion(ctx context.Context, testID string, runID int32, def openapi.TestSpecs) (openapi.ImplResponse, error) {
	run, err := c.testDB.GetRun(ctx, id.ID(testID), int(runID))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, ""), nil
	}

	if run.Trace == nil {
		return openapi.Response(http.StatusUnprocessableEntity, fmt.Sprintf(`run "%d" has no trace associated`, runID)), nil
	}

	definition, err := c.mappers.In.Definition(def.Specs)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), nil
	}

	ds := []expression.DataStore{expression.EnvironmentDataStore{
		Values: run.Environment.Values,
	}}

	assertionExecutor := executor.NewAssertionExecutor(c.tracer)

	results, allPassed := assertionExecutor.Assert(ctx, definition, *run.Trace, ds)
	res := c.mappers.Out.Result(&model.RunResults{
		AllPassed: allPassed,
		Results:   results,
	})

	return openapi.Response(200, res), nil
}

func (c *controller) GetRunResultJUnit(ctx context.Context, testID string, runID int32) (openapi.ImplResponse, error) {
	run, err := c.testDB.GetRun(ctx, id.ID(testID), int(runID))
	if err != nil {
		return handleDBError(err), err
	}

	test, err := c.testDB.GetTestVersion(ctx, id.ID(testID), run.TestVersion)
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
	test, err := c.testDB.GetTestVersion(ctx, id.ID(testID), int(version))
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.mappers.Out.Test(test)), nil
}

func (c controller) GetTestVersionDefinitionFile(ctx context.Context, testID string, version int32) (openapi.ImplResponse, error) {
	test, err := c.testDB.GetTestVersion(ctx, id.ID(testID), int(version))
	if err != nil {
		return handleDBError(err), err
	}

	enc, err := yaml.Encode(yamlconvert.Test(test))
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	return openapi.Response(200, enc), nil
}

func (c controller) ExportTestRun(ctx context.Context, testID string, runID int32) (openapi.ImplResponse, error) {
	run, err := c.testDB.GetRun(ctx, id.ID(testID), int(runID))
	if err != nil {
		return handleDBError(err), err
	}

	test, err := c.testDB.GetTestVersion(ctx, id.ID(testID), run.TestVersion)
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

	createdTest, err := c.testDB.CreateTest(ctx, test)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	createdRun, err := c.testDB.CreateRun(ctx, createdTest, *run)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	createdRun.State = run.State

	err = c.testDB.UpdateRun(ctx, createdRun)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	response := openapi.ExportedTestInformation{
		Test: c.mappers.Out.Test(createdTest),
		Run:  c.mappers.Out.Run(&createdRun),
	}

	return openapi.Response(http.StatusOK, response), nil
}

func (c *controller) UpsertDefinition(ctx context.Context, testDefinition openapi.TextDefinition) (openapi.ImplResponse, error) {
	def, err := yaml.Decode([]byte(testDefinition.Content))
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	if test, err := def.Test(); err == nil {
		return c.upsertTest(ctx, test.Model())
	}

	if transaction, err := def.Transaction(); err == nil {
		return c.upsertTransaction(ctx, transaction.Model())
	}

	return openapi.Response(http.StatusUnprocessableEntity, nil), nil
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

func (c *controller) executeTransaction(ctx context.Context, transaction tests.Transaction, runInfo openapi.RunInformation) (openapi.ImplResponse, error) {
	// create or update transaction
	resp, err := c.doCreateTransaction(ctx, transaction)
	if err != nil {
		if errors.Is(err, errTransactionExists) {
			resp, err = c.doUpdateTransaction(ctx, transaction.ID, transaction)
			if err != nil {
				return resp, err
			}
		} else {
			return resp, err
		}
	} else {
		// the transaction was created, make sure we have the correct ID
		transaction = resp.Body.(tests.Transaction)
	}

	transactionID := transaction.ID

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

func (c *controller) upsertTransaction(ctx context.Context, transaction tests.Transaction) (openapi.ImplResponse, error) {
	resp, err := c.doCreateTransaction(ctx, transaction)
	var status int
	if err != nil {
		if errors.Is(err, errTransactionExists) {
			resp, err := c.doUpdateTransaction(ctx, transaction.ID, transaction)
			if err != nil {
				return resp, err
			}
			status = http.StatusOK
		} else {
			return resp, err
		}
	} else {
		status = http.StatusCreated
		transaction.ID = id.ID(resp.Body.(openapi.Transaction).Id)
	}

	return openapi.ImplResponse{
		Code: status,
		Body: openapi.UpsertDefinitionResponse{
			Id:   transaction.ID.String(),
			Type: yaml.FileTypeTransaction.String(),
		},
	}, nil
}

var errTransactionExists = errors.New("transaction already exists")

func (c *controller) doCreateTransaction(ctx context.Context, transaction tests.Transaction) (openapi.ImplResponse, error) {
	// if they try to create a transaction with preset ID, we need to make sure that ID doesn't exists already
	if transaction.HasID() {
		exists, err := c.transactions.IDExists(ctx, transaction.ID)
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
		transaction = c.transactions.SetID(transaction, id.GenerateID())
	}

	transaction, err := c.transactions.Create(ctx, transaction)
	if err != nil {
		return handleDBError(err), err
	}
	transaction, err = c.transactions.GetAugmented(ctx, transaction.ID)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, transaction), nil
}

func (c *controller) GetTransactionVersionDefinitionFile(ctx context.Context, transactionId string, version int32) (openapi.ImplResponse, error) {
	transaction, err := c.transactions.GetVersion(ctx, id.ID(transactionId), int(version))
	if err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), err
	}

	return openapi.Response(200, transaction), nil
}

func (c *controller) doUpdateTransaction(ctx context.Context, transactionID id.ID, updated tests.Transaction) (openapi.ImplResponse, error) {
	transaction, err := c.transactions.Get(ctx, transactionID)
	if err != nil {
		return handleDBError(err), err
	}

	updated.Version = transaction.Version
	updated.ID = transaction.ID

	_, err = c.transactions.Update(ctx, updated)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(204, nil), nil
}

func metadata(in *map[string]string) model.RunMetadata {
	if in == nil {
		return nil
	}

	return model.RunMetadata(*in)
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

func (c *controller) executeTest(ctx context.Context, test model.Test, runInfo openapi.RunInformation) (openapi.ImplResponse, error) {
	resp, err := c.upsertTest(ctx, test)
	if err != nil {
		return resp, err
	}
	testID := id.ID(resp.Body.(openapi.UpsertDefinitionResponse).Id)
	// test ready, execute it
	resp, err = c.RunTest(ctx, testID.String(), runInfo)
	if resp.Code != http.StatusOK || err != nil {
		return resp, err
	}

	res := openapi.ExecuteDefinitionResponse{
		Id:    testID.String(),
		RunId: resp.Body.(openapi.TestRun).Id,
		Type:  yaml.FileTypeTest.String(),
	}
	return openapi.Response(200, res), nil
}

func (c *controller) upsertTest(ctx context.Context, test model.Test) (openapi.ImplResponse, error) {
	resp, err := c.doCreateTest(ctx, test)
	var status int
	if err != nil {
		if errors.Is(err, errTestExists) {
			resp, err := c.doUpdateTest(ctx, test.ID, test)
			if err != nil {
				return resp, err
			}
			status = http.StatusOK
		} else {
			return resp, err
		}
	} else {
		status = http.StatusCreated
		test.ID = id.ID(resp.Body.(openapi.Test).Id)
	}

	return openapi.ImplResponse{
		Code: status,
		Body: openapi.UpsertDefinitionResponse{
			Id:   test.ID.String(),
			Type: yaml.FileTypeTest.String(),
		},
	}, nil
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

		run, err := c.testDB.GetRun(ctx, id.ID(context.TestId), runId)
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
	transaction, err := c.transactions.GetVersion(ctx, id.ID(tID), int(version))

	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(http.StatusOK, transaction), nil
}

// RunTransaction implements openapi.ApiApiServicer
func (c *controller) RunTransaction(ctx context.Context, transactionID string, runInformation openapi.RunInformation) (openapi.ImplResponse, error) {
	transaction, err := c.transactions.GetAugmented(ctx, id.ID(transactionID))
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

	missingVariablesError, err := validation.ValidateMissingVariablesFromTransaction(ctx, c.testDB, transaction, environment)
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
	run, err := c.transactions.GetTransactionRun(ctx, id.ID(transactionId), int(runId))
	if err != nil {
		return handleDBError(err), err
	}

	openapiRun := c.mappers.Out.TransactionRun(run)
	return openapi.Response(http.StatusOK, openapiRun), nil
}

func (c *controller) GetTransactionRuns(ctx context.Context, transactionId string, take, skip int32) (openapi.ImplResponse, error) {
	runs, err := c.transactions.GetTransactionsRuns(ctx, id.ID(transactionId), take, skip)
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
	run, err := c.transactions.GetTransactionRun(ctx, id.ID(transactionId), int(runId))
	if err != nil {
		return handleDBError(err), err
	}

	err = c.transactions.DeleteTransactionRun(ctx, run)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

func (c *controller) GetResources(ctx context.Context, take, skip int32, query, sortBy, sortDirection string) (openapi.ImplResponse, error) {
	// TODO: this is endpoint is a hack to unblock the team quickly.
	// This is not production ready because it might take too long to respond if there are numerous
	// transactions and tests.

	if take == 0 {
		take = 20
	}

	newTake := take + skip

	transactions, err := c.transactions.ListAugmented(ctx, int(newTake), 0, query, sortBy, sortDirection)
	if err != nil {
		return handleDBError(err), err
	}

	transactionCount, err := c.transactions.Count(ctx, query)
	if err != nil {
		return handleDBError(err), err
	}

	getTestsResponse, err := c.GetTests(ctx, newTake, 0, query, sortBy, sortDirection)
	if err != nil {
		return getTestsResponse, err
	}

	testPaginatedResponse := getTestsResponse.Body.(paginated[openapi.Test])
	tests := testPaginatedResponse.items

	totalResources := transactionCount + testPaginatedResponse.count

	items := takeResources(transactions, tests, take, skip)

	paginatedResponse := paginated[openapi.Resource]{
		items: items,
		count: totalResources,
	}

	return openapi.Response(http.StatusOK, paginatedResponse), nil
}

func takeResources(transactions []tests.Transaction, tests []openapi.Test, take, skip int32) []openapi.Resource {
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

		if transaction.CreatedAt.After(test.CreatedAt) {
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
