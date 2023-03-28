package event

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
)

type Emitter interface {
	Emit(ctx context.Context, event model.TestRunEvent) error
}

type publisher interface {
	Publish(eventID string, message any)
}

type internalEmitter struct {
	repository model.TestRunEventRepository
	publisher  publisher
}

func NewEmitter(repository model.TestRunEventRepository, publisher publisher) Emitter {
	return &internalEmitter{
		repository: repository,
		publisher:  publisher,
	}
}

func (em *internalEmitter) Emit(ctx context.Context, event model.TestRunEvent) error {
	err := em.repository.CreateTestRunEvent(ctx, event)
	if err != nil {
		return err
	}

	eventIDAsString := fmt.Sprintf("%d", event.ID)
	em.publisher.Publish(eventIDAsString, event)

	return nil
}
