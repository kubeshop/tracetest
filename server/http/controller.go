package http

import (
	"context"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/kubeshop/tracetest/analytics"
	"github.com/kubeshop/tracetest/assertions/selectors"
	"github.com/kubeshop/tracetest/executor"
	"github.com/kubeshop/tracetest/openapi"
	"github.com/kubeshop/tracetest/testdb"
	"github.com/kubeshop/tracetest/tracedb"
	"github.com/kubeshop/tracetest/traces"
)

type controller struct {
	traceDB         tracedb.TraceDB
	testDB          testdb.Repository
	runner          executor.Runner
	assertionRunner executor.AssertionRunner
}

func NewController(
	traceDB tracedb.TraceDB,
	testDB testdb.Repository,
	runner executor.Runner,
	assertionRunner executor.AssertionRunner,
) openapi.ApiApiServicer {
	return &controller{
		traceDB:         traceDB,
		testDB:          testDB,
		runner:          runner,
		assertionRunner: assertionRunner,
	}
}

func (s *controller) CreateTest(ctx context.Context, test openapi.Test) (openapi.ImplResponse, error) {
	id, err := s.testDB.CreateTest(ctx, &test)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	analytics.CreateAndSendEvent("test_created_backend", "test")

	test.TestId = id
	return openapi.Response(200, test), nil
}

func (s *controller) UpdateTest(ctx context.Context, testid string, updated openapi.Test) (openapi.ImplResponse, error) {
	test, err := s.testDB.GetTest(ctx, testid)
	if err != nil {
		switch {
		case errors.Is(testdb.ErrNotFound, err):
			return openapi.Response(http.StatusNotFound, err.Error()), err
		default:
			return openapi.Response(http.StatusInternalServerError, err.Error()), err
		}
	}

	updated.TestId = test.TestId

	err = s.testDB.UpdateTest(ctx, &updated)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	analytics.CreateAndSendEvent("test_updated_backend", "test")

	return openapi.Response(204, nil), nil
}

func (s *controller) DeleteTest(ctx context.Context, testid string) (openapi.ImplResponse, error) {
	test, err := s.testDB.GetTest(ctx, testid)
	if err != nil {
		switch {
		case errors.Is(testdb.ErrNotFound, err):
			return openapi.Response(http.StatusNotFound, err.Error()), err
		default:
			return openapi.Response(http.StatusInternalServerError, err.Error()), err
		}
	}

	err = s.testDB.DeleteTest(ctx, test)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	analytics.CreateAndSendEvent("test_deleted_backend", "test")

	return openapi.Response(204, nil), nil
}

func (s *controller) GetTest(ctx context.Context, testid string) (openapi.ImplResponse, error) {
	test, err := s.testDB.GetTest(ctx, testid)
	if err != nil {
		switch {
		case errors.Is(testdb.ErrNotFound, err):
			return openapi.Response(http.StatusNotFound, err.Error()), err
		default:
			return openapi.Response(http.StatusInternalServerError, err.Error()), err
		}
	}

	if test.ReferenceTestRunResult.TraceId != "" {
		res, err := s.testDB.GetResultByTraceID(ctx, test.TestId, test.ReferenceTestRunResult.TraceId)
		if err != nil {
			return openapi.Response(http.StatusInternalServerError, err.Error()), err
		}
		test.ReferenceTestRunResult = res
	}

	return openapi.Response(200, test), nil
}

func (s *controller) GetTests(ctx context.Context, take, skip int32) (openapi.ImplResponse, error) {
	if take == 0 {
		take = 20
	}

	tests, err := s.testDB.GetTests(ctx, take, skip)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(200, tests), nil
}

func (s *controller) RunTest(ctx context.Context, testid string) (openapi.ImplResponse, error) {
	test, err := s.testDB.GetTest(ctx, testid)
	if err != nil {
		switch {
		case errors.Is(testdb.ErrNotFound, err):
			return openapi.Response(http.StatusNotFound, err.Error()), err
		default:
			return openapi.Response(http.StatusInternalServerError, err.Error()), err
		}
	}

	result := s.runner.Run(*test)

	analytics.CreateAndSendEvent("test_run_backend", "test")

	return openapi.Response(200, result), nil
}

func (s *controller) GetTestResults(ctx context.Context, id string, take, skip int32) (openapi.ImplResponse, error) {
	if take == 0 {
		take = 20
	}

	res, err := s.testDB.GetResultsByTestID(ctx, id, take, skip)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(http.StatusOK, res), nil

}

func (s *controller) GetTestResult(ctx context.Context, testid string, id string) (openapi.ImplResponse, error) {
	res, err := s.testDB.GetResult(ctx, id)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}
	return openapi.Response(http.StatusOK, *res), nil
}

func (s *controller) UpdateTestResult(ctx context.Context, testid string, id string, testRunResult openapi.TestAssertionResult) (openapi.ImplResponse, error) {
	testResult, err := s.testDB.GetResult(ctx, id)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	if testRunResult.AssertionResult == nil {
		testRunResult.AssertionResult = []openapi.AssertionResult{}
	}

	for i := range testRunResult.AssertionResult {
		if testRunResult.AssertionResult[i].SpanAssertionResults == nil {
			testRunResult.AssertionResult[i].SpanAssertionResults = []openapi.SpanAssertionResult{}
		}
	}

	testResult.AssertionResultState = testRunResult.AssertionResultState
	testResult.AssertionResult = testRunResult.AssertionResult
	if isExpectingResultStateUpdate(testResult) {
		testResult.State = executor.TestRunStateFinished
	}

	err = s.testDB.UpdateResult(ctx, testResult)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(http.StatusOK, *testResult), nil
}

func isExpectingResultStateUpdate(r *openapi.TestRunResult) bool {
	return r.State == executor.TestRunStateAwaitingTestResults

}

func (s *controller) CreateAssertion(ctx context.Context, testID string, assertion openapi.Assertion) (openapi.ImplResponse, error) {
	test, err := s.testDB.GetTest(ctx, testID)
	if err != nil {
		switch {
		case errors.Is(testdb.ErrNotFound, err):
			return openapi.Response(http.StatusNotFound, err.Error()), err
		default:
			return openapi.Response(http.StatusInternalServerError, err.Error()), err
		}
	}

	id, err := s.testDB.CreateAssertion(ctx, testID, &assertion)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}
	assertion.AssertionId = id

	// Mark reference result as empty after test is updated,
	// so that next test run will update the reference result.
	test.ReferenceTestRunResult.ResultId = ""
	if err = s.testDB.UpdateTest(ctx, test); err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	analytics.CreateAndSendEvent("assertion-created-backend", "test")

	return openapi.Response(http.StatusOK, assertion), nil
}

func (s *controller) UpdateAssertion(ctx context.Context, testID string, assertionID string, updated openapi.Assertion) (openapi.ImplResponse, error) {
	_, err := s.testDB.GetAssertion(ctx, assertionID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	updated.AssertionId = assertionID

	err = s.testDB.UpdateAssertion(ctx, testID, assertionID, updated)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	analytics.CreateAndSendEvent("assertion_updated_backend", "test")

	return openapi.Response(http.StatusNoContent, nil), nil
}

func (s *controller) DeleteAssertion(ctx context.Context, testID string, assertionID string) (openapi.ImplResponse, error) {
	_, err := s.testDB.GetAssertion(ctx, assertionID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	err = s.testDB.DeleteAssertion(ctx, testID, assertionID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	analytics.CreateAndSendEvent("assertion_deleted_backend", "test")

	return openapi.Response(http.StatusNoContent, nil), nil
}

func (s *controller) GetAssertions(ctx context.Context, testID string) (openapi.ImplResponse, error) {
	assertions, err := s.testDB.GetAssertionsByTestID(ctx, testID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(http.StatusOK, assertions), nil
}

func (s *controller) GetTestResultSelectedSpans(ctx context.Context, testID string, resultID string, selectorQuery string) (openapi.ImplResponse, error) {
	selector, err := selectors.New(selectorQuery)
	if err != nil {
		return openapi.Response(http.StatusBadRequest, "invalid selector query"), nil
	}

	result, err := s.testDB.GetResult(ctx, resultID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, ""), nil
	}

	trace, err := traces.FromOtel(result.Trace)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, ""), nil
	}

	selectedSpans := selector.Filter(trace)
	selectedSpanIds := make([]string, len(selectedSpans))

	for i, span := range selectedSpans {
		selectedSpanIds[i] = hex.EncodeToString(span.ID[:])
	}

	return openapi.Response(http.StatusOK, selectedSpanIds), nil
}

func (s *controller) RerunTestResult(ctx context.Context, testID string, resultID string) (openapi.ImplResponse, error) {
	test, err := s.testDB.GetTest(ctx, testID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	result, err := s.testDB.GetResult(ctx, resultID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	result.State = executor.TestRunStateAwaitingTestResults
	err = s.testDB.UpdateResult(ctx, result)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	testDefinition, err := executor.ConvertAssertionsIntoTestDefinition(test.Assertions)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	assertionRequest := executor.AssertionRequest{
		TestDefinition: testDefinition,
		Result:         *result,
	}

	s.assertionRunner.RunAssertions(assertionRequest)

	return openapi.Response(http.StatusOK, result), nil
}
