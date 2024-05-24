package client

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/telemetry"
)

func (c *Client) SendTrace(ctx context.Context, response *proto.PollingResponse) error {
	client := proto.NewOrchestratorClient(c.conn)

	response.AgentIdentification = c.sessionConfig.AgentIdentification
	response.Metadata = telemetry.ExtractMetadataFromContext(ctx)

	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, err := client.SendPolledSpans(ctx, response)
	if err != nil {
		return fmt.Errorf("could not send polled spans result request: %w", err)
	}

	return nil
}
