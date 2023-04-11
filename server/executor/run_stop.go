package executor

import (
	"context"
	"log"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/events"
	"github.com/kubeshop/tracetest/server/subscription"
)

type StopRequest struct {
	TestID id.ID
	RunID  int
}

func (sr StopRequest) ResourceID() string {
	runID := (model.Run{ID: sr.RunID, TestID: sr.TestID}).ResourceID()
	return runID + "/stop"
}

func (r persistentRunner) listenForStopRequests(ctx context.Context, cancelCtx context.CancelFunc, run model.Run) {
	var sfn subscription.Subscriber

	sfn = subscription.NewSubscriberFunction(func(m subscription.Message) error {
		stopRequest, ok := m.Content.(StopRequest)
		if !ok {
			return nil
		}

		ctx, _ := r.tracer.Start(ctx, "User Requested Stop Run")
		// refresh data from DB to make sure we have the altest possible data before updating
		run, err := r.runs.GetRun(ctx, stopRequest.TestID, stopRequest.RunID)
		if err != nil {
			log.Printf("[TracePoller] Test %s Run %d: fail to emit TracePollingStart event: %s \n", stopRequest.TestID, stopRequest.RunID, err.Error())
			return err
		}

		run = run.Stopped()
		r.handleDBError(run, r.updater.Update(ctx, run))

		evt := events.TraceStoppedInfo(stopRequest.TestID, stopRequest.RunID)
		err = r.eventEmitter.Emit(ctx, evt)
		if err != nil {
			log.Printf("[TracePoller] Test %s Run %d: fail to emit TraceStoppedInfo event: %s \n", stopRequest.TestID, stopRequest.RunID, err.Error())
			return err
		}

		cancelCtx()

		r.subscriptionManager.Unsubscribe(stopRequest.ResourceID(), sfn.ID())

		return nil
	})

	r.subscriptionManager.Subscribe((StopRequest{run.TestID, run.ID}).ResourceID(), sfn)
}
