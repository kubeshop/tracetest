package model

import (
	"context"

	"github.com/kubeshop/tracetest/server/pkg/id"
)

type List[T any] struct {
	Items      []T
	TotalCount int
}

type TestRunEventRepository interface {
	CreateTestRunEvent(context.Context, TestRunEvent) error
	GetTestRunEvents(context.Context, id.ID, int) ([]TestRunEvent, error)
}

type Repository interface {
	TestRunEventRepository

	ServerID() (id string, isNew bool, _ error)
	Close() error
	Drop() error
}
