package executor

import (
	"context"
	"log"

	"github.com/kubeshop/tracetest/server/model/events"
)

type tracePollerStartWorker struct {
	eventEmitter EventEmitter
	outputQueue  Enqueuer
}

func NewTracePollerStartWorker(
	eventEmitter EventEmitter,
) *tracePollerStartWorker {
	return &tracePollerStartWorker{
		eventEmitter: eventEmitter,
	}
}

func (w *tracePollerStartWorker) SetOutputQueue(queue Enqueuer) {
	w.outputQueue = queue
}

func (w *tracePollerStartWorker) ProcessItem(ctx context.Context, job Job) {
	if w.isFirstRequest(job) {
		err := w.eventEmitter.Emit(ctx, events.TraceFetchingStart(job.Test.ID, job.Run.ID))
		if err != nil {
			log.Printf("[TracePoller] Test %s Run %d: fail to emit TracePollingStart event: %s", job.Test.ID, job.Run.ID, err.Error())
		}
	}

	log.Println("[TracePoller] processJob", job.EnqueueCount())

	// TODO: check if there is more "pre-processing" things to do here

	w.outputQueue.Enqueue(ctx, job)
}

func (w *tracePollerStartWorker) isFirstRequest(job Job) bool {
	return job.EnqueueCount() == 0
}
