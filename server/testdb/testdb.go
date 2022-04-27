package testdb

import (
	"context"
	"errors"

	"github.com/kubeshop/tracetest/openapi"
)

var ErrNotFound = errors.New("record not found")

type Tests interface {
	CreateTest(ctx context.Context, test *openapi.Test) (string, error)
	UpdateTest(ctx context.Context, test *openapi.Test) error
	GetTests(ctx context.Context) ([]openapi.Test, error)
	GetTest(ctx context.Context, id string) (*openapi.Test, error)
}

type Results interface {
	CreateResult(ctx context.Context, testID string, res *openapi.TestRunResult) error
	UpdateResult(ctx context.Context, res *openapi.TestRunResult) error
	GetResult(ctx context.Context, id string) (*openapi.TestRunResult, error)
	GetResultsByTestID(ctx context.Context, testid string) ([]openapi.TestRunResult, error)
	GetResultByTraceID(ctx context.Context, testid, traceid string) (openapi.TestRunResult, error)
}

type Assertions interface {
	CreateAssertion(ctx context.Context, testid string, assertion *openapi.Assertion) (string, error)
	UpdateAssertion(ctx context.Context, testID, assertionID string, assertion openapi.Assertion) error
	DeleteAssertion(ctx context.Context, testID, assertionID string) error
	GetAssertion(ctx context.Context, id string) (*openapi.Assertion, error)
	GetAssertionsByTestID(ctx context.Context, testID string) ([]openapi.Assertion, error)
}

type Database interface {
	Tests
	Results
	Assertions

	Drop() error
}
