package client

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/agent/proto"
)

func (c *Client) SendTrace(ctx context.Context, pollingResponse *proto.PollingResponse) error {
	client := proto.NewOrchestratorClient(c.conn)

	pollingResponse.AgentIdentification = c.sessionConfig.AgentIdentification

	_, err := client.SendPolledSpans(ctx, pollingResponse)
	if err != nil {
		return fmt.Errorf("could not send polled spans result request: %w", err)
	}

	return nil
}
