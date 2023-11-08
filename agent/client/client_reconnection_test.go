package client_test

import (
	"context"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/agent/client"
	"github.com/kubeshop/tracetest/agent/client/mocks"
	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/stretchr/testify/require"
	"gotest.tools/v3/assert"
)

func TestClientReconnection(t *testing.T) {
	server := mocks.NewGrpcServer()

	client, err := client.Connect(context.Background(), server.Addr(), client.WithPingPeriod(time.Second))
	require.NoError(t, err)

	client.OnTriggerRequest(func(ctx context.Context, tr *proto.TriggerRequest) error { return nil })
	client.OnConnectionClosed(func(ctx context.Context, sr *proto.ShutdownRequest) error { return nil })
	client.OnPollingRequest(func(ctx context.Context, pr *proto.PollingRequest) error { return nil })
	client.OnDataStoreTestConnectionRequest(func(ctx context.Context, dsctr *proto.DataStoreConnectionTestRequest) error { return nil })

	err = client.Start(context.Background())
	require.NoError(t, err)

	err = client.Start(context.Background())
	require.NoError(t, err)

	server.Stop()

	err = client.SendTriggerResponse(context.Background(), &proto.TriggerResponse{RequestID: "my-request-id"})
	require.NotNil(t, err)

	time.Sleep(2 * time.Second)

	server.Restart()

	err = client.SendTriggerResponse(context.Background(), &proto.TriggerResponse{RequestID: "my-request-id"})
	require.NoError(t, err)

	triggerResponse := server.GetLastTriggerResponse()
	require.NotNil(t, triggerResponse)
	assert.Equal(t, "my-request-id", triggerResponse.RequestID)
}
