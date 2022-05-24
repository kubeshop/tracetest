package http

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/kubeshop/tracetest/analytics"
	"github.com/kubeshop/tracetest/assertions"
	"github.com/kubeshop/tracetest/assertions/comparator"
	"github.com/kubeshop/tracetest/assertions/selectors"
	"github.com/kubeshop/tracetest/executor"
	"github.com/kubeshop/tracetest/model"
	"github.com/kubeshop/tracetest/openapi"
	"github.com/kubeshop/tracetest/testdb"
	"github.com/kubeshop/tracetest/tracedb"
)

type controller struct {
	testDB          model.Repository
	runner          executor.Runner
	assertionRunner executor.AssertionRunner

	openapi openapiMapper
	model   modelMapper
}

func NewController(
	traceDB tracedb.TraceDB,
	testDB model.Repository,
	runner executor.Runner,
	assertionRunner executor.AssertionRunner,
) openapi.ApiApiServicer {
	return &controller{
		testDB:          testDB,
		runner:          runner,
		assertionRunner: assertionRunner,
		openapi:         openapiMapper{},
		model:           modelMapper{Comparators: comparator.DefaultRegistry()},
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
	test, err := c.testDB.CreateTest(ctx, c.model.Test(in))
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	analytics.CreateAndSendEvent("test_created_backend", "test")

	return openapi.Response(200, c.openapi.Test(test)), nil
}

func (c *controller) DeleteTest(ctx context.Context, testID string) (openapi.ImplResponse, error) {
	id, err := uuid.Parse(testID)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	test, err := c.testDB.GetLatestTestVersion(ctx, id)
	if err != nil {
		return handleDBError(err), err
	}

	err = c.testDB.DeleteTest(ctx, test)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	analytics.CreateAndSendEvent("test_deleted_backend", "test")

	return openapi.Response(204, nil), nil

}

func (c *controller) GetTest(ctx context.Context, testID string) (openapi.ImplResponse, error) {
	id, err := uuid.Parse(testID)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	test, err := c.testDB.GetLatestTestVersion(ctx, id)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.openapi.Test(test)), nil
}

func (c *controller) GetTestDefinition(ctx context.Context, testID string) (openapi.ImplResponse, error) {
	id, err := uuid.Parse(testID)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	test, err := c.testDB.GetLatestTestVersion(ctx, id)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.openapi.Definition(test.Definition)), nil
}

func (c *controller) GetTestResultSelectedSpans(ctx context.Context, _ string, runID string, selectorQuery string) (openapi.ImplResponse, error) {
	rid, err := uuid.Parse(runID)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	selector, err := selectors.New(selectorQuery)
	if err != nil {
		return handleDBError(err), err
	}

	run, err := c.testDB.GetRun(ctx, rid)
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

	return openapi.Response(http.StatusOK, selectedSpanIds), nil
}

func (c *controller) GetTestRun(ctx context.Context, _ string, runID string) (openapi.ImplResponse, error) {
	rid, err := uuid.Parse(runID)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	run, err := c.testDB.GetRun(ctx, rid)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.openapi.Run(&run)), nil
}

func (c *controller) GetTestRuns(ctx context.Context, testID string, take, skip int32) (openapi.ImplResponse, error) {
	if take == 0 {
		take = 20
	}

	id, err := uuid.Parse(testID)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	test, err := c.testDB.GetLatestTestVersion(ctx, id)
	if err != nil {
		return handleDBError(err), err
	}

	runs, err := c.testDB.GetTestRuns(ctx, test, take, skip)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.openapi.Runs(runs)), nil
}

func (c *controller) GetTests(ctx context.Context, take, skip int32) (openapi.ImplResponse, error) {
	if take == 0 {
		take = 20
	}

	tests, err := c.testDB.GetTests(ctx, take, skip)
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(200, c.openapi.Tests(tests)), nil
}

func (c *controller) RerunTestRun(ctx context.Context, testID string, runID string) (openapi.ImplResponse, error) {
	id, err := uuid.Parse(testID)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	test, err := c.testDB.GetLatestTestVersion(ctx, id)
	if err != nil {
		return handleDBError(err), err
	}

	rid, err := uuid.Parse(runID)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	run, err := c.testDB.GetRun(ctx, rid)
	if err != nil {
		return handleDBError(err), err
	}

	run.State = model.RunStateAwaitingTestResults
	err = c.testDB.UpdateRun(ctx, run)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	assertionRequest := executor.AssertionRequest{
		Test: test,
		Run:  run,
	}

	c.assertionRunner.RunAssertions(assertionRequest)

	return openapi.Response(http.StatusOK, c.openapi.Run(&run)), nil
}

func (c *controller) RunTest(ctx context.Context, testID string) (openapi.ImplResponse, error) {
	id, err := uuid.Parse(testID)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	test, err := c.testDB.GetLatestTestVersion(ctx, id)
	if err != nil {
		return handleDBError(err), err
	}

	run := c.runner.Run(test)

	analytics.CreateAndSendEvent("test_run_backend", "test")

	return openapi.Response(200, c.openapi.Run(&run)), nil
}

func (c *controller) SetTestDefinition(ctx context.Context, testID string, def openapi.TestDefinition) (openapi.ImplResponse, error) {
	id, err := uuid.Parse(testID)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	test, err := c.testDB.GetLatestTestVersion(ctx, id)
	if err != nil {
		return handleDBError(err), err
	}

	err = c.testDB.SetDefiniton(ctx, test, c.model.Definition(def))
	if err != nil {
		return handleDBError(err), err
	}

	return openapi.Response(204, nil), nil
}

func (c *controller) UpdateTest(ctx context.Context, testID string, in openapi.Test) (openapi.ImplResponse, error) {
	id, err := uuid.Parse(testID)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	test, err := c.testDB.GetLatestTestVersion(ctx, id)
	if err != nil {
		return handleDBError(err), err
	}

	updated := c.model.Test(in)
	updated.ID = test.ID
	updated.ReferenceRun = nil

	testHasChanged, err := c.testHasChanged(test, updated)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	if testHasChanged {
		updated.Version = test.Version + 1

		_, err = c.testDB.CreateTestVersion(ctx, updated)
		if err != nil {
			return handleDBError(err), err
		}
	}

	return openapi.Response(204, nil), nil
}

func (c *controller) DryRunAssertion(ctx context.Context, _, runID string, def openapi.TestDefinition) (openapi.ImplResponse, error) {
	rid, err := uuid.Parse(runID)
	if err != nil {
		return openapi.Response(http.StatusUnprocessableEntity, err.Error()), err
	}

	run, err := c.testDB.GetRun(ctx, rid)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, ""), nil
	}

	if run.Trace == nil {
		return openapi.Response(http.StatusUnprocessableEntity, fmt.Sprintf(`run "%s" has no trace associated`, runID)), nil
	}

	results, allPassed := assertions.Assert(c.model.Definition(def), *run.Trace)
	res := c.openapi.Result(&model.RunResults{
		AllPassed: allPassed,
		Results:   results,
	})

	return openapi.Response(200, res), nil
}

func (c *controller) testHasChanged(oldTest model.Test, newTest model.Test) (bool, error) {
	definitionHasChanged, err := c.testFieldHasChanged(oldTest.Definition, newTest.Definition)
	if err != nil {
		return false, err
	}

	serviceUnderTestHasChanged, err := c.testFieldHasChanged(oldTest.ServiceUnderTest, newTest.ServiceUnderTest)
	if err != nil {
		return false, err
	}

	nameHasChanged := oldTest.Name != newTest.Name
	descriptionHasChanged := oldTest.Description != newTest.Description

	return definitionHasChanged || serviceUnderTestHasChanged || nameHasChanged || descriptionHasChanged, nil
}

func (c controller) testFieldHasChanged(oldField interface{}, newField interface{}) (bool, error) {
	oldFieldJSON, err := json.Marshal(oldField)
	if err != nil {
		return false, err
	}

	newFieldJSON, err := json.Marshal(newField)
	if err != nil {
		return false, err
	}

	return string(oldFieldJSON) != string(newFieldJSON), nil
}
