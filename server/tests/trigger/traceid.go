package trigger

const TriggerTypeTRACEID TriggerType = "traceid"

type TraceIDRequest struct {
	ID string `expr_enabled:"true"`
}

type TraceIDResponse struct {
	ID string
}
