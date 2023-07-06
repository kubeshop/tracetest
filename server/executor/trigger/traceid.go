package trigger

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"
)

func TRACEID() Triggerer {
	return &traceidTriggerer{}
}

type traceidTriggerer struct{}

func (t *traceidTriggerer) Trigger(ctx context.Context, test test.Test, opts *TriggerOptions) (Response, error) {
	response := Response{
		Result: trigger.TriggerResult{
			Type:    t.Type(),
			TraceID: &trigger.TraceIDResponse{ID: test.Trigger.TraceID.ID},
		},
	}

	return response, nil
}

func (t *traceidTriggerer) Type() trigger.TriggerType {
	return trigger.TriggerTypeTraceID
}

func (t *traceidTriggerer) Resolve(ctx context.Context, test test.Test, opts *TriggerOptions) (test.Test, error) {
	traceid := test.Trigger.TraceID
	if traceid == nil {
		return test, fmt.Errorf("no settings provided for TRACEID triggerer")
	}

	id, err := opts.Executor.ResolveStatement(WrapInQuotes(traceid.ID, "\""))
	if err != nil {
		return test, err
	}

	traceid.ID = id
	test.Trigger.TraceID = traceid

	return test, nil
}
