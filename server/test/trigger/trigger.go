package trigger

import "encoding/json"

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

// Compatibility with older versions
type triggerFormat struct {
	Type    TriggerType     `json:"type"`
	HTTP    *HTTPRequest    `json:"http,omitempty"`
	GRPC    *GRPCRequest    `json:"grpc,omitempty"`
	TraceID *TraceIDRequest `json:"traceid,omitempty"`
}

type oldTriggerFormat struct {
	Type    TriggerType     `json:"triggerType"`
	HTTP    *HTTPRequest    `json:"http,omitempty"`
	GRPC    *GRPCRequest    `json:"grpc,omitempty"`
	TraceID *TraceIDRequest `json:"traceid,omitempty"`
}

func (t *Trigger) UnmarshalJSON(data []byte) error {
	var trigger triggerFormat
	err := json.Unmarshal(data, &trigger)
	if err != nil || trigger.Type == "" {
		// probably an old format
		var oldFormat oldTriggerFormat
		err = json.Unmarshal(data, &oldFormat)
		if err != nil {
			return err
		}

		trigger.Type = oldFormat.Type
		trigger.HTTP = oldFormat.HTTP
		trigger.GRPC = oldFormat.GRPC
		trigger.TraceID = oldFormat.TraceID
	}

	t.Type = trigger.Type
	t.HTTP = trigger.HTTP
	t.GRPC = trigger.GRPC
	t.TraceID = trigger.TraceID

	return nil
}
