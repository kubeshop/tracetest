package trigger

import (
	"encoding/json"

	"github.com/fluidtruck/deepcopy"
	jsoniter "github.com/json-iterator/go"
)

type triggerJSONV3 struct {
	Type    TriggerType     `json:"type"`
	OldHTTP *HTTPRequest    `json:"http,omitempty"`
	HTTP    *HTTPRequest    `json:"httpRequest,omitempty"`
	GRPC    *GRPCRequest    `json:"grpc,omitempty"`
	TraceID *TraceIDRequest `json:"traceid,omitempty"`
}

func (v3 triggerJSONV3) valid() bool {
	// has a valid type and at least one not nil trigger type settings
	return v3.Type != "" &&
		(v3.HTTP != nil || v3.OldHTTP != nil ||
			v3.GRPC != nil ||
			v3.TraceID != nil)
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
		err := deepcopy.DeepCopy(v3, t)
		if err != nil {
			return err
		}

		if v3.OldHTTP != nil {
			t.HTTP = v3.OldHTTP
		}
	}

	return nil
}
