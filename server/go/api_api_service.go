/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.0.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/kubeshop/tracetest/server/go/tracedb"
	"go.opentelemetry.io/otel/trace"
)

var ErrNotFound = errors.New("record not found")

//go:generate mockgen -package=mocks -destination=mocks/testdb.go . TestDB
type TestDB interface {
	CreateTest(ctx context.Context, test *Test) (string, error)
	UpdateTest(ctx context.Context, test *Test) error
	GetTests(ctx context.Context) ([]Test, error)
	GetTest(ctx context.Context, id string) (*Test, error)

	CreateResult(ctx context.Context, testID string, res *TestRunResult) error
	UpdateResult(ctx context.Context, res *TestRunResult) error
	GetResult(ctx context.Context, id string) (*TestRunResult, error)
	GetResultsByTestID(ctx context.Context, testid string) ([]TestRunResult, error)

	CreateAssertion(ctx context.Context, testid string, assertion *Assertion) (string, error)
	GetAssertion(ctx context.Context, id string) (*Assertion, error)
	GetAssertionsByTestID(ctx context.Context, testID string) ([]Assertion, error)
}

//go:generate mockgen -package=mocks -destination=mocks/executor.go . TestExecutor
type TestExecutor interface {
	Execute(test *Test, tid trace.TraceID, sid trace.SpanID) (*TestRunResult, error)
}

// ApiApiService is a service that implements the logic for the ApiApiServicer
// This service should implement the business logic for every endpoint for the ApiApi API.
// Include any external packages or services that will be required by this service.
type ApiApiService struct {
	traceDB tracedb.TraceDB
	testDB  TestDB
	runner  Runner
}

// NewApiApiService creates a default api service
func NewApiApiService(traceDB tracedb.TraceDB, testDB TestDB, runner Runner) ApiApiServicer {
	return &ApiApiService{
		traceDB: traceDB,
		testDB:  testDB,
		runner:  runner,
	}
}

// CreateTest - Create new test
func (s *ApiApiService) CreateTest(ctx context.Context, test Test) (ImplResponse, error) {
	id, err := s.testDB.CreateTest(ctx, &test)
	if err != nil {
		return Response(http.StatusInternalServerError, err.Error()), err
	}

	test.TestId = id
	return Response(200, test), nil
}

// UpdateTest - Create new test
func (s *ApiApiService) UpdateTest(ctx context.Context, testid string, updated Test) (ImplResponse, error) {
	test, err := s.testDB.GetTest(ctx, testid)
	if err != nil {
		switch {
		case errors.Is(ErrNotFound, err):
			return Response(http.StatusNotFound, err.Error()), err
		default:
			return Response(http.StatusInternalServerError, err.Error()), err
		}
	}

	updated.TestId = test.TestId

	err = s.testDB.UpdateTest(ctx, &updated)
	if err != nil {
		return Response(http.StatusInternalServerError, err.Error()), err
	}

	return Response(204, nil), nil
}

// GetTest - Get a test
func (s *ApiApiService) GetTest(ctx context.Context, testid string) (ImplResponse, error) {
	test, err := s.testDB.GetTest(ctx, testid)
	if err != nil {
		switch {
		case errors.Is(ErrNotFound, err):
			return Response(http.StatusNotFound, err.Error()), err
		default:
			return Response(http.StatusInternalServerError, err.Error()), err
		}
	}

	if test.ReferenceTestRunResult.TraceId != "" {
		// TOODO: implement this code
		// test.ReferenceTestRunResult = getResultByTraceID(test.ReferenceTestRunResult.TraceId)
	}

	return Response(200, test), nil
}

// GetTests - Gets all tests
func (s *ApiApiService) GetTests(ctx context.Context) (ImplResponse, error) {
	tests, err := s.testDB.GetTests(ctx)
	if err != nil {
		return Response(http.StatusInternalServerError, err.Error()), err
	}

	return Response(200, tests), nil
}

func (s *ApiApiService) RunTest(ctx context.Context, testid string) (ImplResponse, error) {
	test, err := s.testDB.GetTest(ctx, testid)
	if err != nil {
		switch {
		case errors.Is(ErrNotFound, err):
			return Response(http.StatusNotFound, err.Error()), err
		default:
			return Response(http.StatusInternalServerError, err.Error()), err
		}
	}

	id := s.runner.Run(ctx, *test)

	return Response(200, TestRun{
		TestRunId: id,
	}), nil
}

// GetTestResults -
func (s *ApiApiService) GetTestResults(ctx context.Context, id string) (ImplResponse, error) {
	res, err := s.testDB.GetResultsByTestID(ctx, id)
	if err != nil {
		return Response(http.StatusInternalServerError, err.Error()), err
	}

	return Response(http.StatusOK, res), nil

}

// GetTestResult -
func (s *ApiApiService) GetTestResult(ctx context.Context, testid string, id string) (ImplResponse, error) {
	res, err := s.testDB.GetResult(ctx, id)
	if err != nil {
		return Response(http.StatusInternalServerError, err.Error()), err
	}
	tr, err := s.traceDB.GetTraceByID(ctx, res.TraceId)
	if err == nil {
		sid, err := trace.SpanIDFromHex(res.SpanId)
		if err != nil {
			return Response(http.StatusInternalServerError, err.Error()), err
		}
		tid, err := trace.TraceIDFromHex(res.TraceId)
		if err != nil {
			return Response(http.StatusInternalServerError, err.Error()), err
		}
		ttr := FixParent(tr, string(tid[:]), string(sid[:]), res.Response)
		res.Trace = mapTrace(ttr)
	}
	return Response(http.StatusOK, *res), nil
}

func (s *ApiApiService) UpdateTestResult(ctx context.Context, testid string, id string, testRunResult TestAssertionResult) (ImplResponse, error) {
	testResult, err := s.testDB.GetResult(ctx, id)
	if err != nil {
		return Response(http.StatusInternalServerError, err.Error()), err
	}

	testResult.AssertionResultState = testRunResult.AssertionResultState
	testResult.AssertionResult = testRunResult.AssertionResult

	err = s.testDB.UpdateResult(ctx, testResult)
	if err != nil {
		return Response(http.StatusInternalServerError, err.Error()), err
	}

	return Response(http.StatusOK, *testResult), nil
}

func (s *ApiApiService) CreateAssertion(ctx context.Context, testID string, assertion Assertion) (ImplResponse, error) {
	test, err := s.testDB.GetTest(ctx, testID)
	if err != nil {
		switch {
		case errors.Is(ErrNotFound, err):
			return Response(http.StatusNotFound, err.Error()), err
		default:
			return Response(http.StatusInternalServerError, err.Error()), err
		}
	}

	id, err := s.testDB.CreateAssertion(ctx, testID, &assertion)
	if err != nil {
		return Response(http.StatusInternalServerError, err.Error()), err
	}
	assertion.AssertionId = id

	// Mark reference result as empty after test is updated,
	// so that next test run will update the reference result.
	test.ReferenceTestRunResult.ResultId = ""
	if err = s.testDB.UpdateTest(ctx, test); err != nil {
		return Response(http.StatusInternalServerError, err.Error()), err
	}

	return Response(http.StatusOK, assertion), nil
}

func (s *ApiApiService) GetAssertions(ctx context.Context, testID string) (ImplResponse, error) {
	assertions, err := s.testDB.GetAssertionsByTestID(ctx, testID)
	if err != nil {
		return Response(http.StatusInternalServerError, err.Error()), err
	}

	return Response(http.StatusOK, assertions), nil
}
