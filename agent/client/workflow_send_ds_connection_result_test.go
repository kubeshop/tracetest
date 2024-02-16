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

func TestDataStoreConnectionResult(t *testing.T) {
	server := mocks.NewGrpcServer()
	defer server.Stop()

	client, err := client.Connect(context.Background(), server.Addr())
	require.NoError(t, err)

	err = client.Start(context.Background())
	require.NoError(t, err)

	result := &proto.DataStoreConnectionTestResponse{
		RequestID:           "request-id",
		AgentIdentification: &proto.AgentIdentification{},
		Successful:          true,
		Steps: &proto.DataStoreConnectionTestSteps{
			PortCheck: &proto.DataStoreConnectionTestStep{
				Passed: true,
			},
		},
	}

	err = client.SendDataStoreConnectionResult(context.Background(), result)
	require.NoError(t, err)

	receivedResponse := server.GetLastDataStoreConnectionResponse()

	assert.Equal(t, result.RequestID, receivedResponse.Data.RequestID)
	assert.True(t, result.Successful)
	assert.True(t, result.Steps.PortCheck.Passed)
}
