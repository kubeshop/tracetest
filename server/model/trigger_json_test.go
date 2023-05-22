package model_test

import (
	"encoding/json"
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTriggerFormatV1(t *testing.T) {
	v1 := struct {
		Type    model.TriggerType
		HTTP    *model.HTTPRequest
		GRPC    *model.GRPCRequest
		TraceID *model.TRACEIDRequest
	}{
		Type: model.TriggerTypeHTTP,
		HTTP: &model.HTTPRequest{
			Method: model.HTTPMethodGET,
			URL:    "http://example.com/list",
		},
	}

	expected := model.Trigger{
		Type: model.TriggerTypeHTTP,
		HTTP: &model.HTTPRequest{
			Method: model.HTTPMethodGET,
			URL:    "http://example.com/list",
		},
	}

	v1Json, err := json.Marshal(v1)
	require.NoError(t, err)

	current := model.Trigger{}
	err = json.Unmarshal(v1Json, &current)
	require.NoError(t, err)

	assert.Equal(t, expected, current)
}

func TestTriggerFormatV2(t *testing.T) {
	v2 := struct {
		Type    model.TriggerType     `json:"triggerType"`
		HTTP    *model.HTTPRequest    `json:"http,omitempty"`
		GRPC    *model.GRPCRequest    `json:"grpc,omitempty"`
		TraceID *model.TRACEIDRequest `json:"traceid,omitempty"`
	}{
		Type: model.TriggerTypeHTTP,
		HTTP: &model.HTTPRequest{
			Method: model.HTTPMethodGET,
			URL:    "http://example.com/list",
		},
	}

	expected := model.Trigger{
		Type: model.TriggerTypeHTTP,
		HTTP: &model.HTTPRequest{
			Method: model.HTTPMethodGET,
			URL:    "http://example.com/list",
		},
	}

	v2Json, err := json.Marshal(v2)
	require.NoError(t, err)

	current := model.Trigger{}
	err = json.Unmarshal(v2Json, &current)
	require.NoError(t, err)

	assert.Equal(t, expected, current)
}

func TestTriggerFormatInvalid(t *testing.T) {
	v2 := struct {
		Type    model.TriggerType     `json:"invalid"`
		HTTP    *model.HTTPRequest    `json:"http,omitempty"`
		GRPC    *model.GRPCRequest    `json:"grpc,omitempty"`
		TraceID *model.TRACEIDRequest `json:"traceid,omitempty"`
	}{
		Type: model.TriggerTypeHTTP,
	}

	v2Json, err := json.Marshal(v2)
	require.NoError(t, err)

	current := model.Trigger{}
	err = json.Unmarshal(v2Json, &current)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "unexpected json format")

}
