package trigger

const TriggerTypeTraceID TriggerType = "traceid"
const TriggerTypeCypress TriggerType = "cypress"
const TriggerTypePlaywright TriggerType = "playwright"
const TriggerTypeArtillery TriggerType = "artillery"
const TriggerTypeK6 TriggerType = "k6"

var traceIDBasedTriggers = []TriggerType{TriggerTypeTraceID, TriggerTypeCypress, TriggerTypePlaywright, TriggerTypeArtillery, TriggerTypeK6}

type TraceIDRequest struct {
	ID string `json:"id,omitempty" expr_enabled:"true"`
}

type TraceIDResponse struct {
	ID string
}
