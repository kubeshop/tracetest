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

func TestPollWorkflow(t *testing.T) {
	server := mocks.NewGrpcServer()

	client, err := client.Connect(context.Background(), server.Addr())
	require.NoError(t, err)

	var receivedPollingRequest *proto.PollingRequest
	client.OnPollingRequest(func(_ context.Context, request *proto.PollingRequest) error {
		receivedPollingRequest = request
		return nil
	})

	err = client.Start(context.Background())
	require.NoError(t, err)

	pollingRequest := &proto.PollingRequest{
		TestID:  "my test",
		RunID:   8,
		TraceID: "my-trace-id",
		Datastore: &proto.DataStore{
			Type: "opensearch",
			Opensearch: &proto.ElasticConfig{
				Addresses:          []string{"http://localhost:9200"},
				Username:           "my user",
				Password:           "dolphins",
				Index:              "traces",
				InsecureSkipVerify: true,
			},
		},
	}

	server.SendPollingRequest(pollingRequest)

	// ensures there's enough time for networking between server and client
	time.Sleep(1 * time.Second)

	assert.NotNil(t, receivedPollingRequest)
	assert.Equal(t, pollingRequest.TestID, receivedPollingRequest.TestID)
	assert.Equal(t, pollingRequest.RunID, receivedPollingRequest.RunID)
	assert.Equal(t, pollingRequest.TraceID, receivedPollingRequest.TraceID)
	assert.Equal(t, pollingRequest.Datastore.Type, receivedPollingRequest.Datastore.Type)
	assert.Equal(t, pollingRequest.Datastore.Opensearch.Addresses, receivedPollingRequest.Datastore.Opensearch.Addresses)
	assert.Equal(t, pollingRequest.Datastore.Opensearch.Username, receivedPollingRequest.Datastore.Opensearch.Username)
	assert.Equal(t, pollingRequest.Datastore.Opensearch.Password, receivedPollingRequest.Datastore.Opensearch.Password)
	assert.Equal(t, pollingRequest.Datastore.Opensearch.Index, receivedPollingRequest.Datastore.Opensearch.Index)
	assert.Equal(t, pollingRequest.Datastore.Opensearch.Certificate, receivedPollingRequest.Datastore.Opensearch.Certificate)
	assert.Equal(t, pollingRequest.Datastore.Opensearch.InsecureSkipVerify, receivedPollingRequest.Datastore.Opensearch.InsecureSkipVerify)

	t.Run("Check if next polling request gets accepted as well", func(t *testing.T) {
		anotherPollingRequest := &proto.PollingRequest{
			TestID:  "other test",
			RunID:   19,
			TraceID: "other-trace-id",
			Datastore: &proto.DataStore{
				Type: "opensearch",
				Opensearch: &proto.ElasticConfig{
					Addresses:          []string{"http://localhost:9200"},
					Username:           "my user 2",
					Password:           "not_dolphins",
					Index:              "traces",
					InsecureSkipVerify: true,
				},
			},
		}

		server.SendPollingRequest(anotherPollingRequest)

		// ensures there's enough time for networking between server and client
		time.Sleep(1 * time.Second)

		assert.NotNil(t, receivedPollingRequest)
		assert.Equal(t, anotherPollingRequest.TestID, receivedPollingRequest.TestID)
		assert.Equal(t, anotherPollingRequest.RunID, receivedPollingRequest.RunID)
		assert.Equal(t, anotherPollingRequest.TraceID, receivedPollingRequest.TraceID)
		assert.Equal(t, anotherPollingRequest.Datastore.Type, receivedPollingRequest.Datastore.Type)
		assert.Equal(t, anotherPollingRequest.Datastore.Opensearch.Addresses, receivedPollingRequest.Datastore.Opensearch.Addresses)
		assert.Equal(t, anotherPollingRequest.Datastore.Opensearch.Username, receivedPollingRequest.Datastore.Opensearch.Username)
		assert.Equal(t, anotherPollingRequest.Datastore.Opensearch.Password, receivedPollingRequest.Datastore.Opensearch.Password)
		assert.Equal(t, anotherPollingRequest.Datastore.Opensearch.Index, receivedPollingRequest.Datastore.Opensearch.Index)
		assert.Equal(t, anotherPollingRequest.Datastore.Opensearch.Certificate, receivedPollingRequest.Datastore.Opensearch.Certificate)
		assert.Equal(t, anotherPollingRequest.Datastore.Opensearch.InsecureSkipVerify, receivedPollingRequest.Datastore.Opensearch.InsecureSkipVerify)
	})
}
