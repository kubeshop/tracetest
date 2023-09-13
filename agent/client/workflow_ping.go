package client

import (
	"context"
	"time"

	"github.com/kubeshop/tracetest/agent/proto"
)

func (c *Client) startHearthBeat(ctx context.Context) error {
	client := proto.NewOrchestratorClient(c.conn)
	ticker := time.NewTicker(2 * time.Minute)

	go func() {
		for range ticker.C {
			client.Ping(ctx, c.sessionConfig.AgentIdentification)
		}
	}()

	return nil
}
