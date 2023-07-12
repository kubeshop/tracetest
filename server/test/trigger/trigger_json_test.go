package trigger_test

import (
	"encoding/json"
	"testing"

	"github.com/kubeshop/tracetest/server/test/trigger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTriggerFormatV1(t *testing.T) {
	v1 := struct {
		Type    trigger.TriggerType
		HTTP    *trigger.HTTPRequest
		GRPC    *trigger.GRPCRequest
		TraceID *trigger.TraceIDRequest
	}{
		Type: trigger.TriggerTypeHTTP,
		HTTP: &trigger.HTTPRequest{
			Method: trigger.HTTPMethodGET,
			URL:    "http://example.com/list",
		},
	}

	expected := trigger.Trigger{
		Type: trigger.TriggerTypeHTTP,
		HTTP: &trigger.HTTPRequest{
			Method: trigger.HTTPMethodGET,
			URL:    "http://example.com/list",
		},
	}

	v1Json, err := json.Marshal(v1)
	require.NoError(t, err)

	current := trigger.Trigger{}
	err = json.Unmarshal(v1Json, &current)
	require.NoError(t, err)

	assert.Equal(t, expected, current)
}

func TestTriggerFormatV2(t *testing.T) {
	v2 := struct {
		Type    trigger.TriggerType     `json:"triggerType"`
		HTTP    *trigger.HTTPRequest    `json:"http,omitempty"`
		GRPC    *trigger.GRPCRequest    `json:"grpc,omitempty"`
		TraceID *trigger.TraceIDRequest `json:"traceid,omitempty"`
	}{
		Type: trigger.TriggerTypeHTTP,
		HTTP: &trigger.HTTPRequest{
			Method: trigger.HTTPMethodGET,
			URL:    "http://example.com/list",
		},
	}

	expected := trigger.Trigger{
		Type: trigger.TriggerTypeHTTP,
		HTTP: &trigger.HTTPRequest{
			Method: trigger.HTTPMethodGET,
			URL:    "http://example.com/list",
		},
	}

	v2Json, err := json.Marshal(v2)
	require.NoError(t, err)

	current := trigger.Trigger{}
	err = json.Unmarshal(v2Json, &current)
	require.NoError(t, err)

	assert.Equal(t, expected, current)
}

func TestTriggerFormatV3(t *testing.T) {
	v3 := struct {
		Type    trigger.TriggerType     `json:"type"`
		HTTP    *trigger.HTTPRequest    `json:"httpRequest,omitempty"`
		GRPC    *trigger.GRPCRequest    `json:"grpc,omitempty"`
		TraceID *trigger.TraceIDRequest `json:"traceid,omitempty"`
	}{
		Type: trigger.TriggerTypeHTTP,
		HTTP: &trigger.HTTPRequest{
			Method: trigger.HTTPMethodGET,
			URL:    "http://example.com/list",
		},
	}

	expected := trigger.Trigger{
		Type: trigger.TriggerTypeHTTP,
		HTTP: &trigger.HTTPRequest{
			Method: trigger.HTTPMethodGET,
			URL:    "http://example.com/list",
		},
	}

	v3Json, err := json.Marshal(v3)
	require.NoError(t, err)

	current := trigger.Trigger{}
	err = json.Unmarshal(v3Json, &current)
	require.NoError(t, err)

	assert.Equal(t, expected, current)
}

func TestTriggerResultV1(t *testing.T) {
	v1 := struct {
		Type    trigger.TriggerType
		HTTP    *trigger.HTTPResponse
		GRPC    *trigger.GRPCResponse
		TraceID *trigger.TraceIDResponse
	}{
		Type: trigger.TriggerTypeHTTP,
		HTTP: &trigger.HTTPResponse{
			Status:     "200",
			StatusCode: 200,
			Headers: []trigger.HTTPHeader{
				{Key: "Content-Type", Value: "application/json"},
			},
			Body: "this is a body",
		},
	}

	expected := trigger.TriggerResult{
		Type: trigger.TriggerTypeHTTP,
		HTTP: &trigger.HTTPResponse{
			Status:     "200",
			StatusCode: 200,
			Headers: []trigger.HTTPHeader{
				{Key: "Content-Type", Value: "application/json"},
			},
			Body: "this is a body",
		},
	}

	v1Json, err := json.Marshal(v1)
	require.NoError(t, err)

	current := trigger.TriggerResult{}
	err = json.Unmarshal(v1Json, &current)
	require.NoError(t, err)

	assert.Equal(t, expected, current)
}

func TestTriggerResultV2(t *testing.T) {
	v2 := struct {
		Type    trigger.TriggerType      `json:"type,omitEmpty"`
		HTTP    *trigger.HTTPResponse    `json:"httpRequest,omitEmpty"`
		GRPC    *trigger.GRPCResponse    `json:"grpc,omitEmpty"`
		TraceID *trigger.TraceIDResponse `json:"traceid,omitEmpty"`
	}{
		Type: trigger.TriggerTypeHTTP,
		HTTP: &trigger.HTTPResponse{
			Status:     "200",
			StatusCode: 200,
			Headers: []trigger.HTTPHeader{
				{Key: "Content-Type", Value: "application/json"},
			},
			Body: "this is a body",
		},
	}

	expected := trigger.TriggerResult{
		Type: trigger.TriggerTypeHTTP,
		HTTP: &trigger.HTTPResponse{
			Status:     "200",
			StatusCode: 200,
			Headers: []trigger.HTTPHeader{
				{Key: "Content-Type", Value: "application/json"},
			},
			Body: "this is a body",
		},
	}

	v1Json, err := json.Marshal(v2)
	require.NoError(t, err)

	current := trigger.TriggerResult{}
	err = json.Unmarshal(v1Json, &current)
	require.NoError(t, err)

	assert.Equal(t, expected, current)
}
