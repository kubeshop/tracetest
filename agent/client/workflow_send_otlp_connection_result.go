package client

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/telemetry"
)

func (c *Client) SendOTLPConnectionResult(ctx context.Context, response *proto.OTLPConnectionTestResponse) error {
	client := proto.NewOrchestratorClient(c.conn)

	response.AgentIdentification = c.sessionConfig.AgentIdentification
	response.Metadata = telemetry.ExtractMetadataFromContext(ctx)

	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, err := client.SendOTLPConnectionTestResult(ctx, response)
	if err != nil {
		return fmt.Errorf("could not send otlp connection result request: %w", err)
	}

	return nil
}
