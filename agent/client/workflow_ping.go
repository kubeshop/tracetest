package client

import (
	"context"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/agent/proto"
)

func (c *Client) startHearthBeat(ctx context.Context) error {
	client := proto.NewOrchestratorClient(c.conn)
	ticker := time.NewTicker(2 * time.Second)

	go func() {
		for range ticker.C {
			fmt.Println("@@PING")
			client.Ping(ctx, c.sessionConfig.AgentIdentification)
		}
	}()

	return nil
}
