package http

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kubeshop/tracetest/server/assertions"
	"github.com/kubeshop/tracetest/server/assertions/selectors"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/expression"
	"github.com/kubeshop/tracetest/server/http/mappings"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/junit"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/yaml"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/testdb"
	"go.opentelemetry.io/otel/trace"
)

var IDGen = id.NewRandGenerator()

type controller struct {
	testDB          model.Repository
	runner          executor.Runner
	assertionRunner executor.AssertionRunner
	mappers         mappings.Mappings
}

func NewController(
	testDB model.Repository,
	runner executor.Runner,
	assertionRunner executor.AssertionRunner,
	mappers mappings.Mappings,
) openapi.ApiApiServicer {
	return &controller{
		testDB:          testDB,
		runner:          runner,
		assertionRunner: assertionRunner,
		mappers:         mappers,
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
	// if they try to create a test with preset ID, we need to make sure that ID doesn't exists already
	if test.HasID() {
		exists, err := c.testDB.IDExists(ctx, test.ID)

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

func (c *controller) GetTestResultSelectedSpans(ctx context.Context, testID, runID, selectorQuery string) (openapi.ImplResponse, error) {
	selector, err := selectors.New(selectorQuery)

	if err != nil {
		return handleDBError(err), err
	}

	rid, err := strconv.Atoi(runID)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, fmt.Errorf("%s is not a number", runID)), err
	}
	run, err := c.testDB.GetRun(ctx, id.ID(testID), rid)
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

func (c *controller) GetTestRun(ctx context.Context, testID, runID string) (openapi.ImplResponse, error) {
	rid, err := strconv.Atoi(runID)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, fmt.Errorf("%s is not a number", runID)), err
	}
	run, err := c.testDB.GetRun(ctx, id.ID(testID), rid)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.mappers.Out.Run(&run)), nil
}

func (c *controller) DeleteTestRun(ctx context.Context, testID, runID string) (openapi.ImplResponse, error) {
	rid, err := strconv.Atoi(runID)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, fmt.Errorf("%s is not a number", runID)), err
	}
	run, err := c.testDB.GetRun(ctx, id.ID(testID), rid)
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

func (c *controller) RerunTestRun(ctx context.Context, testID, runID string) (openapi.ImplResponse, error) {
	test, err := c.testDB.GetLatestTestVersion(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	rid, err := strconv.Atoi(runID)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, fmt.Errorf("%s is not a number", runID)), err
	}
	run, err := c.testDB.GetRun(ctx, id.ID(testID), rid)
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

	c.assertionRunner.RunAssertions(ctx, assertionRequest)

	return openapi.Response(http.StatusOK, c.mappers.Out.Run(&newTestRun)), nil
}

func (c *controller) RunTest(ctx context.Context, testID string, runInformation openapi.TestRunInformation) (openapi.ImplResponse, error) {
	test, err := c.testDB.GetLatestTestVersion(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	metadata := metadata(runInformation.Metadata)
	environment, err := environment(ctx, c.testDB, runInformation.EnvironmentId)

	if err != nil {
		return handleDBError(err), err
	}

	run := c.runner.Run(ctx, test, metadata, environment)

	return openapi.Response(200, c.mappers.Out.Run(&run)), nil
}

func (c *controller) SetTestOutputs(ctx context.Context, testID string, testOutputs []openapi.TestOutput) (openapi.ImplResponse, error) {
	test, err := c.testDB.GetLatestTestVersion(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	newOutputs, err := c.mappers.In.Outputs(testOutputs)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), nil
	}

	newTest, err := model.BumpVersionIfOutputsChanged(test, newOutputs)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	newTest.Outputs = newOutputs

	newTest, err = c.testDB.UpdateTest(ctx, newTest)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(204, nil), nil
}

func (c *controller) SetTestSpecs(ctx context.Context, testID string, def openapi.TestSpecs) (openapi.ImplResponse, error) {
	if err := c.mappers.In.ValidateDefinition(def); err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	test, err := c.testDB.GetLatestTestVersion(ctx, id.ID(testID))
	if err != nil {
		return handleDBError(err), err
	}

	newDefinition, err := c.mappers.In.Definition(def)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), nil
	}

	newTest, err := model.BumpVersionIfDefinitionChanged(test, newDefinition)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	newTest.Specs = newDefinition

	newTest, err = c.testDB.UpdateTest(ctx, newTest)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(204, nil), nil
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

func (c *controller) DryRunAssertion(ctx context.Context, testID, runID string, def openapi.TestSpecs) (openapi.ImplResponse, error) {
	rid, err := strconv.Atoi(runID)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, fmt.Errorf("%s is not a number", runID)), err
	}
	run, err := c.testDB.GetRun(ctx, id.ID(testID), rid)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, ""), nil
	}

	if run.Trace == nil {
		return openapi.Response(http.StatusUnprocessableEntity, fmt.Sprintf(`run "%s" has no trace associated`, runID)), nil
	}

	definition, err := c.mappers.In.Definition(def)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), nil
	}

	ds := []expression.DataStore{expression.EnvironmentDataStore{
		Values: run.Environment.Values,
	}}

	results, allPassed := assertions.Assert(definition, *run.Trace, ds)
	res := c.mappers.Out.Result(&model.RunResults{
		AllPassed: allPassed,
		Results:   results,
	})

	return openapi.Response(200, res), nil
}

func (c *controller) GetRunResultJUnit(ctx context.Context, testID, runID string) (openapi.ImplResponse, error) {
	rid, err := strconv.Atoi(runID)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, fmt.Errorf("%s is not a number", runID)), err
	}
	run, err := c.testDB.GetRun(ctx, id.ID(testID), rid)
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

	enc, err := yaml.Encode(yaml.File{
		Type: yaml.FileTypeTest,
		Spec: test,
	})
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	return openapi.Response(200, enc), nil
}

func (c controller) ExportTestRun(ctx context.Context, testID, runID string) (openapi.ImplResponse, error) {
	rid, err := strconv.Atoi(runID)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, fmt.Errorf("%s is not a number", runID)), err
	}
	run, err := c.testDB.GetRun(ctx, id.ID(testID), rid)
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

func (c *controller) ExecuteDefinition(ctx context.Context, testDefinition openapi.TextDefinition) (openapi.ImplResponse, error) {
	def, err := yaml.Decode([]byte(testDefinition.Content))
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	if test, err := def.Test(); err == nil {
		return c.executeTest(ctx, test.Model(), testDefinition.RunInformation)
	}

	if environment, err := def.Environment(); err == nil {
		return c.createEnvFromDefinition(ctx, environment)
	}

	return openapi.Response(http.StatusUnprocessableEntity, nil), nil
}

func (c *controller) createEnvFromDefinition(ctx context.Context, def yaml.Environment) (openapi.ImplResponse, error) {
	environment := def.Model()

	if environment.ID != "" {
		_, err := c.testDB.UpdateEnvironment(ctx, environment)

		if err != nil {
			return handleDBError(err), err
		}
	} else {
		env, err := c.testDB.CreateEnvironment(ctx, environment)
		environment.ID = env.ID

		if err != nil {
			return handleDBError(err), err
		}
	}

	res := openapi.ExecuteDefinitionResponse{
		Id: environment.ID,
	}

	return openapi.Response(200, res), nil
}

func metadata(in *map[string]string) model.RunMetadata {
	if in == nil {
		return nil
	}

	return model.RunMetadata(*in)
}

func environment(ctx context.Context, testDB model.Repository, environmentId string) (model.Environment, error) {
	if environmentId != "" {
		environment, err := testDB.GetEnvironment(ctx, environmentId)

		if err != nil {
			return model.Environment{}, err
		}

		return environment, nil
	}

	return model.Environment{}, nil
}

func (c *controller) executeTest(ctx context.Context, test model.Test, runInfo openapi.TestRunInformation) (openapi.ImplResponse, error) {
	// create or update test
	testID := test.ID
	resp, err := c.doCreateTest(ctx, test)
	if err != nil {
		if errors.Is(err, errTestExists) {
			resp, err := c.doUpdateTest(ctx, test.ID, test)
			if err != nil {
				return resp, err
			}

		} else {
			return resp, err
		}
	} else {
		testID = id.ID(resp.Body.(openapi.Test).Id)
	}

	// test ready, execute it
	resp, err = c.RunTest(ctx, testID.String(), runInfo)
	if err != nil {
		return resp, err
	}

	res := openapi.ExecuteDefinitionResponse{
		Id:    testID.String(),
		RunId: resp.Body.(openapi.TestRun).Id,
		Type:  string(yaml.FileTypeTest),
	}
	return openapi.Response(200, res), nil
}

// Environments

func (c *controller) CreateEnvironment(ctx context.Context, in openapi.Environment) (openapi.ImplResponse, error) {
	environment := c.mappers.In.Environment(in)

	if environment.ID != "" {
		exists, err := c.testDB.EnvironmentIDExists(ctx, environment.ID)

		if err != nil {
			return handleDBError(err), err
		}

		if exists {
			err := fmt.Errorf(`cannot create environment with ID "%s: %w`, environment.ID, errTestExists)
			r := map[string]string{
				"error": err.Error(),
			}
			return openapi.Response(http.StatusBadRequest, r), err
		}
	}

	environment, err := c.testDB.CreateEnvironment(ctx, environment)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(200, c.mappers.Out.Environment(environment)), nil
}

func (c *controller) DeleteEnvironment(ctx context.Context, environmentId string) (openapi.ImplResponse, error) {
	environment, err := c.testDB.GetEnvironment(ctx, environmentId)
	if err != nil {
		return handleDBError(err), err
	}

	err = c.testDB.DeleteEnvironment(ctx, environment)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

func (c *controller) GetEnvironment(ctx context.Context, environmentId string) (openapi.ImplResponse, error) {
	environment, err := c.testDB.GetEnvironment(ctx, environmentId)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.mappers.Out.Environment(environment)), nil
}

func (c *controller) GetEnvironments(ctx context.Context, take, skip int32, query string, sortBy string, sortDirection string) (openapi.ImplResponse, error) {
	if take == 0 {
		take = 20
	}

	environments, err := c.testDB.GetEnvironments(ctx, take, skip, query, sortBy, sortDirection)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, paginated[openapi.Environment]{
		items: c.mappers.Out.Environments(environments.Items),
		count: environments.TotalCount,
	}), nil
}

func (c *controller) UpdateEnvironment(ctx context.Context, environmentId string, in openapi.Environment) (openapi.ImplResponse, error) {
	updated := c.mappers.In.Environment(in)

	environment, err := c.testDB.GetEnvironment(ctx, environmentId)
	if err != nil {
		return handleDBError(err), err
	}

	updated.ID = environment.ID

	_, err = c.testDB.UpdateEnvironment(ctx, updated)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(204, nil), nil
}

// expressions
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
		environment, err := c.testDB.GetEnvironment(ctx, context.EnvironmentId)

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

func (c *controller) CreateTransaction(ctx context.Context, transaction openapi.Transaction) (openapi.ImplResponse, error) {
	transactionModel, err := c.mappers.In.Transaction(ctx, transaction)
	if err != nil {
		return handleDBError(err), err
	}

	createdTransaction, err := c.testDB.CreateTransaction(ctx, transactionModel)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(http.StatusOK, c.mappers.Out.Transaction(createdTransaction)), nil
}

func (c *controller) DeleteTransaction(ctx context.Context, tID string) (openapi.ImplResponse, error) {
	transaction, err := c.testDB.GetLatestTransactionVersion(ctx, id.ID(tID))
	if err != nil {
		return handleDBError(err), err
	}

	err = c.testDB.DeleteTransaction(ctx, transaction)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

func (c *controller) GetTransaction(ctx context.Context, tID string) (openapi.ImplResponse, error) {
	transaction, err := c.testDB.GetLatestTransactionVersion(ctx, id.ID(tID))
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(http.StatusOK, c.mappers.Out.Transaction(transaction)), nil
}

func (c *controller) GetTransactions(ctx context.Context, take, skip int32, query, sortBy, sortDirection string) (openapi.ImplResponse, error) {
	if take == 0 {
		take = 20
	}

	transactions, err := c.testDB.GetTransactions(ctx, take, skip, query, sortBy, sortDirection)
	if err != nil {
		return handleDBError(err), err
	}

	apiTransactions := make([]openapi.Transaction, len(transactions.Items))
	for i, transaction := range transactions.Items {
		apiTransactions[i] = c.mappers.Out.Transaction(transaction)
	}

	return openapi.Response(http.StatusOK, paginated[openapi.Transaction]{
		items: apiTransactions,
		count: transactions.TotalCount,
	}), nil
}

func (c *controller) UpdateTransaction(ctx context.Context, tID string, transaction openapi.Transaction) (openapi.ImplResponse, error) {
	modelTransaction, err := c.mappers.In.Transaction(ctx, transaction)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, err.Error()), err
	}

	modelTransaction.ID = id.ID(tID)

	_, err = c.testDB.UpdateTransaction(ctx, modelTransaction)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(http.StatusNoContent, nil), nil
}

// GetResources implements openapi.ApiApiServicer
func (c *controller) GetResources(ctx context.Context, take, skip int32, query, sortBy, sortDirection string) (openapi.ImplResponse, error) {
	// TODO: this is endpoint is a hack to unblock the team quickly.
	// This is not production ready because it might take too long to respond if there are numerous
	// transactions and tests.

	if take == 0 {
		take = 20
	}

	newTake := take + skip

	getTransactionsReponse, err := c.GetTransactions(ctx, newTake, 0, query, sortBy, sortDirection)
	if err != nil {
		return getTransactionsReponse, err
	}

	getTestsResponse, err := c.GetTests(ctx, newTake, 0, query, sortBy, sortDirection)
	if err != nil {
		return getTestsResponse, err
	}

	transactionPaginatedResponse := getTransactionsReponse.Body.(paginated[openapi.Transaction])
	transactions := transactionPaginatedResponse.items
	testPaginatedResponse := getTestsResponse.Body.(paginated[openapi.Test])
	tests := testPaginatedResponse.items

	totalResources := transactionPaginatedResponse.count + testPaginatedResponse.count

	items := takeResources(transactions, tests, take, skip)

	paginatedResponse := paginated[openapi.Resource]{
		items: items,
		count: totalResources,
	}

	return openapi.Response(http.StatusOK, paginatedResponse), nil
}

func takeResources(transactions []openapi.Transaction, tests []openapi.Test, take, skip int32) []openapi.Resource {
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
			i++
		}

		currentNumberItens++
	}

	upperLimit := int(skip + take)
	if upperLimit > currentNumberItens {
		upperLimit = currentNumberItens
	}

	return items[skip:upperLimit]
}
