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

func TestOTLPConnectionResultTrace(t *testing.T) {
	server := mocks.NewGrpcServer()
	defer server.Stop()

	client, err := client.Connect(context.Background(), server.Addr())
	require.NoError(t, err)

	err = client.Start(context.Background())
	require.NoError(t, err)

	now := time.Now()
	result := &proto.OTLPConnectionTestResponse{
		RequestID:           "request-id",
		AgentIdentification: &proto.AgentIdentification{},
		SpanCount:           10,
		LastSpanTimestamp:   now.UnixMilli(),
	}

	err = client.SendOTLPConnectionResult(context.Background(), result)
	require.NoError(t, err)

	receivedResponse := server.GetLastOTLPConnectionResponse()

	assert.Equal(t, result.RequestID, receivedResponse.Data.RequestID)
	assert.Equal(t, result.SpanCount, receivedResponse.Data.SpanCount)
	assert.Equal(t, result.LastSpanTimestamp, receivedResponse.Data.LastSpanTimestamp)
}
