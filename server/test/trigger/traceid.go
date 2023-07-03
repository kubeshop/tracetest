package trigger

const TriggerTypeTraceID TriggerType = "traceid"

type TraceIDRequest struct {
	ID string `json:"id,omitempty" expr_enabled:"true"`
}

type TraceIDResponse struct {
	ID string
}
