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

	assert.Equal(t, pollingRequest.TestID, receivedPollingResponse.Data.TestID)
	assert.Equal(t, pollingRequest.RunID, receivedPollingResponse.Data.RunID)
	assert.Equal(t, pollingRequest.TraceID, receivedPollingResponse.Data.TraceID)

	require.Len(t, receivedPollingResponse.Data.Spans, len(pollingRequest.Spans))
	for i, span := range pollingRequest.Spans {
		assert.Equal(t, span.Id, receivedPollingResponse.Data.Spans[i].Id)
		assert.Equal(t, span.Name, receivedPollingResponse.Data.Spans[i].Name)
		assert.Equal(t, span.Kind, receivedPollingResponse.Data.Spans[i].Kind)
		assert.Equal(t, span.ParentId, receivedPollingResponse.Data.Spans[i].ParentId)
		assert.Equal(t, span.StartTime, receivedPollingResponse.Data.Spans[i].StartTime)
		assert.Equal(t, span.EndTime, receivedPollingResponse.Data.Spans[i].EndTime)
		for j := range span.Attributes {
			assert.Equal(t, span.Attributes[i].Key, receivedPollingResponse.Data.Spans[i].Attributes[j].Key)
			assert.Equal(t, span.Attributes[i].Value, receivedPollingResponse.Data.Spans[i].Attributes[j].Value)
		}
	}
}
