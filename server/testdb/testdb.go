package testdb

import (
	"context"
	"errors"

	"github.com/kubeshop/tracetest/openapi"
)

var ErrNotFound = errors.New("record not found")

type TestRepository interface {
	CreateTest(ctx context.Context, test *openapi.Test) (string, error)
	UpdateTest(ctx context.Context, test *openapi.Test) error
	DeleteTest(ctx context.Context, test *openapi.Test) error
	GetTests(ctx context.Context, take, skip int32) ([]openapi.Test, error)
	GetTest(ctx context.Context, id string) (*openapi.Test, error)
}

type RunRepository interface {
	CreateResult(ctx context.Context, testID string, res *openapi.TestRun) error
	UpdateResult(ctx context.Context, res *openapi.TestRun) error
	GetResult(ctx context.Context, id string) (*openapi.TestRun, error)
	GetResultsByTestID(ctx context.Context, testid string, take, skip int32) ([]openapi.TestRun, error)
	GetResultByTraceID(ctx context.Context, testid, traceid string) (openapi.TestRun, error)
}

type AssertionRepository interface {
	CreateAssertion(ctx context.Context, testid string, assertion *openapi.Assertion) (string, error)
	UpdateAssertion(ctx context.Context, testID, assertionID string, assertion openapi.Assertion) error
	DeleteAssertion(ctx context.Context, testID, assertionID string) error
	GetAssertion(ctx context.Context, id string) (*openapi.Assertion, error)
	GetAssertionsByTestID(ctx context.Context, testID string) ([]openapi.Assertion, error)
}

type Repository interface {
	TestRepository
	RunRepository
	AssertionRepository

	Drop() error
}
