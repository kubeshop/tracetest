package trigger

import (
	"errors"
	"slices"
)

type TriggerType string

func (t TriggerType) IsTraceIDBased() bool {
	return slices.Contains(traceIDBasedTriggers, t)
}

func (t TriggerType) IsIntegration() bool {
	return slices.Contains(traceIDBasedIntegrationsTriggers, t)
}

type (
	Trigger struct {
		Type    TriggerType     `json:"type"`
		HTTP    *HTTPRequest    `json:"httpRequest,omitempty"`
		GRPC    *GRPCRequest    `json:"grpc,omitempty"`
		TraceID *TraceIDRequest `json:"traceid,omitempty"`
		Kafka   *KafkaRequest   `json:"kafka,omitempty"`
	}

	TriggerResult struct {
		Type    TriggerType      `json:"type"`
		HTTP    *HTTPResponse    `json:"http,omitempty"`
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

func (e *TriggerError) ToError() error {
	if e == nil || e.ErrorMessage == "" {
		return nil
	}

	return errors.New(e.ErrorMessage)
}
