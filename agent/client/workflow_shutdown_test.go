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

func TestShutdownFlow(t *testing.T) {
	server := mocks.NewGrpcServer()

	client, err := client.Connect(context.Background(), server.Addr())
	require.NoError(t, err)

	var called bool = false
	var reason string = ""
	client.OnConnectionClosed(func(ctx context.Context, sr *proto.ShutdownRequest) error {
		called = true
		reason = sr.Reason
		return nil
	})

	err = client.Start(context.Background())
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	server.TerminateConnection("shutdown requested by user")

	time.Sleep(1 * time.Second)
	assert.True(t, called, "client.OnConnectionClosed should have been called")
	assert.Equal(t, "shutdown requested by user", reason)
}
