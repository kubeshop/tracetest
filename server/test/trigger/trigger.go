package trigger

type (
	TriggerType string

	Trigger struct {
		Type    TriggerType     `json:"type"`
		HTTP    *HTTPRequest    `json:"httpRequest,omitempty"`
		GRPC    *GRPCRequest    `json:"grpc,omitempty"`
		TraceID *TraceIDRequest `json:"traceid,omitempty"`
	}

	TriggerResult struct {
		Type    TriggerType      `json:"type"`
		HTTP    *HTTPResponse    `json:"httpRequest,omitempty"`
		GRPC    *GRPCResponse    `json:"grpc,omitempty"`
		TraceID *TraceIDResponse `json:"traceid,omitempty"`
	}
)
