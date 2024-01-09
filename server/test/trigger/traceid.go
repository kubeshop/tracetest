package trigger

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
