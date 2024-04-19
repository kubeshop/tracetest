package client

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/telemetry"
)

func (c *Client) SendTriggerResponse(ctx context.Context, response *proto.TriggerResponse) error {
	client := proto.NewOrchestratorClient(c.conn)

	if response.AgentIdentification == nil {
		response.AgentIdentification = c.sessionConfig.AgentIdentification
	}

	response.Metadata = telemetry.ExtractMetadataFromContext(ctx)

	_, err := client.SendTriggerResult(ctx, response)
	if err != nil {
		return fmt.Errorf("could not send trigger result request: %w", err)
	}

	return nil
}
