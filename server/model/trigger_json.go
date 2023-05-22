package model

import (
	"encoding/json"
	"fmt"

	"github.com/fluidtruck/deepcopy"
)

type triggerJSONV2 struct {
	Type    TriggerType     `json:"triggerType"`
	HTTP    *HTTPRequest    `json:"http,omitempty"`
	GRPC    *GRPCRequest    `json:"grpc,omitempty"`
	TraceID *TRACEIDRequest `json:"traceid,omitempty"`
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
	TraceID *TRACEIDRequest
}

func (v1 triggerJSONV1) valid() bool {
	// has a valid type and at least one not nil trigger type settings
	return v1.Type != "" &&
		(v1.HTTP != nil ||
			v1.GRPC != nil ||
			v1.TraceID != nil)
}

func (t Trigger) MarshalJSON() ([]byte, error) {
	jt := triggerJSONV2{}
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

	return fmt.Errorf("unexpected json format")
}
