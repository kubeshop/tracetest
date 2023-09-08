package trigger

import "errors"

type (
	TriggerType string

	Trigger struct {
		Type    TriggerType     `json:"type"`
		HTTP    *HTTPRequest    `json:"httpRequest,omitempty"`
		GRPC    *GRPCRequest    `json:"grpc,omitempty"`
		TraceID *TraceIDRequest `json:"traceid,omitempty"`
		Kafka   *KafkaRequest   `json:"kafka,omitempty"`
	}

	TriggerResult struct {
		Type    TriggerType      `json:"type"`
		HTTP    *HTTPResponse    `json:"httpRequest,omitempty"`
		GRPC    *GRPCResponse    `json:"grpc,omitempty"`
		TraceID *TraceIDResponse `json:"traceid,omitempty"`
		Kafka   *KafkaResponse   `json:"kafka,omitempty"`
		Error   *TriggerError    `json:"error,omitempty"`
	}

	TriggerError struct {
		ConnectionError    bool   `json:"connectionError"`
		RunningOnContainer bool   `json:"runningOnContainer"`
		TargetsLocalhost   bool   `json:"targetsLocalhost"`
		ErrorMessage       string `json:"message"`
	}
)

func (e TriggerError) Error() error {
	return errors.New(e.ErrorMessage)
}
