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

func TestDataStoreConnectionTestWorkflow(t *testing.T) {
	ctx := context.Background()
	server := mocks.NewGrpcServer()
	defer server.Stop()

	client, err := client.Connect(ctx, server.Addr(), client.WithInsecure())
	require.NoError(t, err)

	var receivedConnectionTestRequest *proto.DataStoreConnectionTestRequest
	client.OnDataStoreTestConnectionRequest(func(ctx context.Context, otr *proto.DataStoreConnectionTestRequest) error {
		receivedConnectionTestRequest = otr
		return nil
	})

	err = client.Start(ctx)
	require.NoError(t, err)

	connectionTestRequest := &proto.DataStoreConnectionTestRequest{
		RequestID: "request-id",
	}

	server.SendDataStoreConnectionTestRequest(ctx, connectionTestRequest)

	// ensures there's enough time for networking between server and client
	time.Sleep(1 * time.Second)

	assert.NotNil(t, receivedConnectionTestRequest)
	assert.Equal(t, connectionTestRequest.RequestID, "request-id")
}
