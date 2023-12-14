package trigger

const TriggerTypeTraceID TriggerType = "traceid"
const TriggerTypeCypress TriggerType = "cypress"

var traceIDBasedTriggers = []TriggerType{TriggerTypeTraceID, TriggerTypeCypress}

type TraceIDRequest struct {
	ID string `json:"id,omitempty" expr_enabled:"true"`
}

type TraceIDResponse struct {
	ID string
}
