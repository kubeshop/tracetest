package trigger

type (
	TriggerType string

	Trigger struct {
		Type    TriggerType     `json:"type"`
		HTTP    *HTTPRequest    `json:"http,omitempty"`
		GRPC    *GRPCRequest    `json:"grpc,omitempty"`
		TraceID *TraceIDRequest `json:"traceid,omitempty"`
	}

	TriggerResult struct {
		Type    TriggerType      `json:"type"`
		HTTP    *HTTPResponse    `json:"http,omitempty"`
		GRPC    *GRPCResponse    `json:"grpc,omitempty"`
		TraceID *TraceIDResponse `json:"traceid,omitempty"`
	}
)
