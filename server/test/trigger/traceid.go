package trigger

import "golang.org/x/exp/slices"

const TriggerTypeTraceID TriggerType = "traceid"
const TriggerTypeCypress TriggerType = "cypress"

var traceIDBasedTriggers = []TriggerType{TriggerTypeTraceID, TriggerTypeCypress}

func IsTraceIDBasedTrigger(triggerType TriggerType) bool {
	return slices.Contains(traceIDBasedTriggers, triggerType)
}

type TraceIDRequest struct {
	ID string `json:"id,omitempty" expr_enabled:"true"`
}

type TraceIDResponse struct {
	ID string
}
