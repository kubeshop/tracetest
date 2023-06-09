package trigger

const TriggerTypeTraceID TriggerType = "traceid"

type TraceIDRequest struct {
	ID string `expr_enabled:"true"`
}

type TraceIDResponse struct {
	ID string
}
