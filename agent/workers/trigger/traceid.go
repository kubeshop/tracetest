package trigger

import (
	"context"
)

func TRACEID() Triggerer {
	return &traceidTriggerer{}
}

type traceidTriggerer struct{}

func (t *traceidTriggerer) Trigger(ctx context.Context, triggerConfig Trigger, opts *Options) (Response, error) {
	response := Response{
		Result: TriggerResult{
			Type:    t.Type(),
			TraceID: &TraceIDResponse{ID: triggerConfig.TraceID.ID},
		},
	}

	return response, nil
}

func (t *traceidTriggerer) Type() TriggerType {
	return TriggerTypeTraceID
}

const TriggerTypeTraceID TriggerType = "traceid"
const TriggerTypeCypress TriggerType = "cypress"
const TriggerTypePlaywright TriggerType = "playwright"

var traceIDBasedTriggers = []TriggerType{TriggerTypeTraceID, TriggerTypeCypress, TriggerTypePlaywright}

type TraceIDRequest struct {
	ID string `json:"id,omitempty" expr_enabled:"true"`
}

type TraceIDResponse struct {
	ID string
}
