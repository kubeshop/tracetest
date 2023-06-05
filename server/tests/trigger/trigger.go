package trigger

type (
	TriggerType string

	Trigger struct {
		Type    TriggerType
		HTTP    *HTTPRequest
		GRPC    *GRPCRequest
		TraceID *TraceIDRequest
	}

	TriggerResult struct {
		Type    TriggerType
		HTTP    *HTTPResponse
		GRPC    *GRPCResponse
		TraceID *TraceIDResponse
	}
)
