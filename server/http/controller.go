package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/kubeshop/tracetest/analytics"
	"github.com/kubeshop/tracetest/executor"
	"github.com/kubeshop/tracetest/openapi"
	"github.com/kubeshop/tracetest/testdb"
	"github.com/kubeshop/tracetest/tracedb"
)

type controller struct {
	traceDB tracedb.TraceDB
	testDB  testdb.Repository
	runner  executor.Runner
}

func NewController(traceDB tracedb.TraceDB, testDB testdb.Repository, runner executor.Runner) openapi.ApiApiServicer {
	return &controller{
		traceDB: traceDB,
		testDB:  testDB,
		runner:  runner,
	}
}

// CreateTest - Create new test
func (s *controller) CreateTest(ctx context.Context, test openapi.Test) (openapi.ImplResponse, error) {
	id, err := s.testDB.CreateTest(ctx, &test)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	analytics.CreateAndSendEvent("test_created", "test")

	test.TestId = id
	return openapi.Response(200, test), nil
}

// UpdateTest - Create new test
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

	analytics.CreateAndSendEvent("test_updated", "test")

	return openapi.Response(204, nil), nil
}

// GetTest - Get a test
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

// GetTests - Gets all tests
func (s *controller) GetTests(ctx context.Context) (openapi.ImplResponse, error) {
	tests, err := s.testDB.GetTests(ctx)
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

	analytics.CreateAndSendEvent("test_run", "test")

	return openapi.Response(200, result), nil
}

// GetTestResults -
func (s *controller) GetTestResults(ctx context.Context, id string) (openapi.ImplResponse, error) {
	res, err := s.testDB.GetResultsByTestID(ctx, id)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(http.StatusOK, res), nil

}

// GetTestResult -
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
	testResult.State = executor.TestRunStateFinished

	err = s.testDB.UpdateResult(ctx, testResult)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(http.StatusOK, *testResult), nil
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

	analytics.CreateAndSendEvent("assertion_created", "test")

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

	analytics.CreateAndSendEvent("assertion_updated", "test")

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

	analytics.CreateAndSendEvent("assertion_deleted", "test")

	return openapi.Response(http.StatusNoContent, nil), nil
}

func (s *controller) GetAssertions(ctx context.Context, testID string) (openapi.ImplResponse, error) {
	assertions, err := s.testDB.GetAssertionsByTestID(ctx, testID)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, err.Error()), err
	}

	return openapi.Response(http.StatusOK, assertions), nil
}
