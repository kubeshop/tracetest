package trigger

import (
	"encoding/json"

	"github.com/fluidtruck/deepcopy"
	jsoniter "github.com/json-iterator/go"
)

type triggerJSONV3 struct {
	Type    TriggerType     `json:"type"`
	HTTP    *HTTPRequest    `json:"httpRequest,omitempty"`
	GRPC    *GRPCRequest    `json:"grpc,omitempty"`
	TraceID *TraceIDRequest `json:"traceid,omitempty"`
	Kafka   *KafkaRequest   `json:"kafka,omitempty"`
}

func (v3 triggerJSONV3) valid() bool {
	// has a valid type and at least one not nil trigger type settings
	return v3.Type != "" &&
		(v3.HTTP != nil ||
			v3.GRPC != nil ||
			v3.TraceID != nil ||
			v3.Kafka != nil)
}

type triggerJSONV2 struct {
	Type    TriggerType     `json:"triggerType"`
	HTTP    *HTTPRequest    `json:"http,omitempty"`
	GRPC    *GRPCRequest    `json:"grpc,omitempty"`
	TraceID *TraceIDRequest `json:"traceid,omitempty"`
}

func (v2 triggerJSONV2) valid() bool {
	// has a valid type and at least one not nil trigger type settings
	return v2.Type != "" &&
		(v2.HTTP != nil ||
			v2.GRPC != nil ||
			v2.TraceID != nil)
}

type triggerJSONV1 struct {
	Type    TriggerType
	HTTP    *HTTPRequest
	GRPC    *GRPCRequest
	TraceID *TraceIDResponse
}

func (v1 triggerJSONV1) valid() bool {
	// has a valid type and at least one not nil trigger type settings
	return v1.Type != "" &&
		(v1.HTTP != nil ||
			v1.GRPC != nil ||
			v1.TraceID != nil)
}

func (t Trigger) MarshalJSON() ([]byte, error) {
	jt := triggerJSONV3{}
	err := deepcopy.DeepCopy(t, &jt)
	if err != nil {
		return nil, err
	}

	return json.Marshal(jt)
}

func (t *Trigger) UnmarshalJSON(data []byte) error {
	var err error
	// start with older versions and move up to the latest
	v1 := triggerJSONV1{}

	// DO NOT USE encoding/json here. the match is case insensitive, and can lead to unexpected results.
	// see https://stackoverflow.com/a/49006601
	var json = jsoniter.Config{
		EscapeHTML:    true,
		CaseSensitive: true,
	}.Froze()

	err = json.Unmarshal(data, &v1)
	if err != nil {
		return err
	}
	if v1.valid() {
		return deepcopy.DeepCopy(v1, t)
	}

	// v2
	v2 := triggerJSONV2{}
	err = json.Unmarshal(data, &v2)
	if err != nil {
		return err
	}
	if v2.valid() {
		return deepcopy.DeepCopy(v2, t)
	}

	// v3
	v3 := triggerJSONV3{}
	err = json.Unmarshal(data, &v3)
	if err != nil {
		return err
	}
	if v3.valid() {
		return deepcopy.DeepCopy(v3, t)
	}

	return nil
}

type triggerResultV1 struct {
	Type    TriggerType      `json:"type"`
	HTTP    *HTTPResponse    `json:"http,omitempty"`
	GRPC    *GRPCResponse    `json:"grpc,omitempty"`
	TraceID *TraceIDResponse `json:"traceid,omitempty"`
}

type triggerResultV2 struct {
	Type    TriggerType      `json:"type"`
	HTTP    *HTTPResponse    `json:"httpRequest,omitempty"`
	GRPC    *GRPCResponse    `json:"grpc,omitempty"`
	TraceID *TraceIDResponse `json:"traceid,omitempty"`
	Kafka   *KafkaResponse   `json:"kafka,omitempty"`
}

func (tr *triggerResultV2) valid() bool {
	return tr.HTTP != nil || tr.GRPC != nil || tr.TraceID != nil || tr.Kafka != nil
}

func (t *TriggerResult) UnmarshalJSON(data []byte) error {
	v2 := triggerResultV2{}
	json.Unmarshal(data, &v2)

	if v2.valid() {
		t.Type = v2.Type
		t.HTTP = v2.HTTP
		t.GRPC = v2.GRPC
		t.TraceID = v2.TraceID
		t.Kafka = v2.Kafka

		return nil
	}

	// Fallback to v1 in last case
	// TriggerResult might be empty at the time of the Unmarshal, so it's ok to not validate it
	v1 := triggerResultV1{}
	json.Unmarshal(data, &v1)

	t.Type = v1.Type
	t.HTTP = v1.HTTP
	t.GRPC = v1.GRPC
	t.TraceID = v1.TraceID

	return nil
}
