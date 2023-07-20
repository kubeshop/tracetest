package client_test

import (
	"context"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/client/mocks"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTriggerWorkflow(t *testing.T) {
	server := mocks.NewGrpcServer()

	client, err := client.Connect(context.Background(), server.Addr())
	require.NoError(t, err)

	var receivedTrigger *proto.TriggerRequest
	client.OnTriggerRequest(func(tr *proto.TriggerRequest) error {
		receivedTrigger = tr
		return nil
	})

	err = client.Start(context.Background())
	require.NoError(t, err)

	triggerRequest := &proto.TriggerRequest{
		TestID: "my test",
		RunID:  8,
		Trigger: &proto.Trigger{
			Type: "http",
			Http: &proto.HttpRequest{
				Method: "GET",
				Url:    "http://localhost:11633/api/tests",
				Headers: []*proto.HttpHeader{
					{Key: "Content-Type", Value: "application/json"},
				},
			},
		},
	}

	server.SendTriggerRequest(triggerRequest)

	// ensures there's enough time for networking between server and client
	time.Sleep(1 * time.Second)

	assert.NotNil(t, receivedTrigger)
	assert.Equal(t, triggerRequest.TestID, receivedTrigger.TestID)
	assert.Equal(t, triggerRequest.RunID, receivedTrigger.RunID)
	assert.Equal(t, triggerRequest.Trigger.Type, receivedTrigger.Trigger.Type)
	assert.Equal(t, triggerRequest.Trigger.Http.Method, receivedTrigger.Trigger.Http.Method)
	assert.Equal(t, triggerRequest.Trigger.Http.Url, receivedTrigger.Trigger.Http.Url)

	require.Len(t, receivedTrigger.Trigger.Http.Headers, 1)
	assert.Equal(t, triggerRequest.Trigger.Http.Headers[0].Key, receivedTrigger.Trigger.Http.Headers[0].Key)
	assert.Equal(t, triggerRequest.Trigger.Http.Headers[0].Value, receivedTrigger.Trigger.Http.Headers[0].Value)

	t.Run("Check if next trigger request gets accepted as well", func(t *testing.T) {
		anotherTriggerRequest := &proto.TriggerRequest{
			TestID: "other test",
			RunID:  10,
			Trigger: &proto.Trigger{
				Type: "http",
				Http: &proto.HttpRequest{
					Method: "POST",
					Url:    "http://localhost:11633/api/tests",
					Headers: []*proto.HttpHeader{
						{Key: "Content-Type", Value: "application/json"},
					},
				},
			},
		}

		server.SendTriggerRequest(anotherTriggerRequest)

		// ensures there's enough time for networking between server and client
		time.Sleep(1 * time.Second)

		assert.NotNil(t, receivedTrigger)
		assert.Equal(t, anotherTriggerRequest.TestID, receivedTrigger.TestID)
		assert.Equal(t, anotherTriggerRequest.RunID, receivedTrigger.RunID)
		assert.Equal(t, anotherTriggerRequest.Trigger.Type, receivedTrigger.Trigger.Type)
		assert.Equal(t, anotherTriggerRequest.Trigger.Http.Method, receivedTrigger.Trigger.Http.Method)
		assert.Equal(t, anotherTriggerRequest.Trigger.Http.Url, receivedTrigger.Trigger.Http.Url)

		require.Len(t, receivedTrigger.Trigger.Http.Headers, 1)
		assert.Equal(t, anotherTriggerRequest.Trigger.Http.Headers[0].Key, receivedTrigger.Trigger.Http.Headers[0].Key)
		assert.Equal(t, anotherTriggerRequest.Trigger.Http.Headers[0].Value, receivedTrigger.Trigger.Http.Headers[0].Value)
	})
}
