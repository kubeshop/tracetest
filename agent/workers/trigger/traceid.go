package trigger

import (
	"context"

	"github.com/kubeshop/tracetest/server/test/trigger"
)

func TRACEID() Triggerer {
	return &traceidTriggerer{}
}

type traceidTriggerer struct{}

func (t *traceidTriggerer) Trigger(ctx context.Context, triggerConfig trigger.Trigger, opts Options) (Response, error) {
	response := Response{
		Result: trigger.TriggerResult{
			Type:    t.Type(),
			TraceID: &trigger.TraceIDResponse{ID: triggerConfig.TraceID.ID},
		},
	}

	return response, nil
}

func (t *traceidTriggerer) Type() trigger.TriggerType {
	return trigger.TriggerTypeTraceID
}
