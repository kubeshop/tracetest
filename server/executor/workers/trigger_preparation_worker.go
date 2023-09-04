package workers

import (
	"context"

	"github.com/kubeshop/tracetest/server/executor"
	triggerer "github.com/kubeshop/tracetest/server/executor/trigger"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/test"
)

type EventEmitter interface {
	Emit(ctx context.Context, event model.TestRunEvent) error
}

type TriggerPreparationWorker struct {
	runRepository test.RunRepository
	eventEmiter   EventEmitter
	triggers      *triggerer.Registry
}

func (w *TriggerPreparationWorker) ProcessItem(ctx context.Context, job executor.Job) {

}
