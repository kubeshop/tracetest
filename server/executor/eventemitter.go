package executor

import (
	"context"

	"github.com/kubeshop/tracetest/server/model"
)

type EventEmitter interface {
	Emit(ctx context.Context, event model.TestRunEvent) error
}

type publisher interface {
	Publish(eventID string, message any)
}

type internalEventEmitter struct {
	repository model.TestRunEventRepository
	publisher  publisher
}

func NewEventEmitter(repository model.TestRunEventRepository, publisher publisher) EventEmitter {
	return &internalEventEmitter{
		repository: repository,
		publisher:  publisher,
	}
}

func (em *internalEventEmitter) Emit(ctx context.Context, event model.TestRunEvent) error {
	err := em.repository.CreateTestRunEvent(ctx, event)
	if err != nil {
		return err
	}

	em.publisher.Publish(event.ResourceID(), event)

	return nil
}
