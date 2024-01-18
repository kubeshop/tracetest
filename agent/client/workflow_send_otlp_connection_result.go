package client

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/agent/proto"
)

func (c *Client) SendDataStoreConnectionResult(ctx context.Context, response *proto.DataStoreConnectionTestResponse) error {
	client := proto.NewOrchestratorClient(c.conn)

	response.AgentIdentification = c.sessionConfig.AgentIdentification

	_, err := client.SendDataStoreConnectionTestResult(ctx, response)
	if err != nil {
		return fmt.Errorf("could not send otlp connection result request: %w", err)
	}

	return nil
}
