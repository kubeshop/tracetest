package client_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/client/mocks"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendTrace(t *testing.T) {
	server := mocks.NewGrpcServer()
	defer server.Stop()

	client, err := client.Connect(context.Background(), server.Addr())
	require.NoError(t, err)

	err = client.Start(context.Background())
	require.NoError(t, err)

	pollingRequest := &proto.PollingResponse{
		TestID:  "test",
		RunID:   1,
		TraceID: "trace-id",
		Spans: []*proto.Span{
			{
				Name:      "GET /",
				Id:        "id",
				ParentId:  "parent-id",
				Kind:      "internal",
				StartTime: 0,
				EndTime:   45,
				Attributes: []*proto.KeyValuePair{
					{Key: "http.status", Value: "200"},
				},
			},
		},
	}

	err = client.SendTrace(context.Background(), pollingRequest)
	require.NoError(t, err)

	receivedPollingResponse := server.GetLastPollingResponse()

	assert.Equal(t, pollingRequest.TestID, receivedPollingResponse.TestID)
	assert.Equal(t, pollingRequest.RunID, receivedPollingResponse.RunID)
	assert.Equal(t, pollingRequest.TraceID, receivedPollingResponse.TraceID)

	require.Len(t, receivedPollingResponse.Spans, len(pollingRequest.Spans))
	for i, span := range pollingRequest.Spans {
		assert.Equal(t, span.Id, receivedPollingResponse.Spans[i].Id)
		assert.Equal(t, span.Name, receivedPollingResponse.Spans[i].Name)
		assert.Equal(t, span.Kind, receivedPollingResponse.Spans[i].Kind)
		assert.Equal(t, span.ParentId, receivedPollingResponse.Spans[i].ParentId)
		assert.Equal(t, span.StartTime, receivedPollingResponse.Spans[i].StartTime)
		assert.Equal(t, span.EndTime, receivedPollingResponse.Spans[i].EndTime)
		for j := range span.Attributes {
			assert.Equal(t, span.Attributes[i].Key, receivedPollingResponse.Spans[i].Attributes[j].Key)
			assert.Equal(t, span.Attributes[i].Value, receivedPollingResponse.Spans[i].Attributes[j].Value)
		}
	}
}
