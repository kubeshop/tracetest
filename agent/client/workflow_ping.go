package client

import (
	"context"
	"log"
	"time"

	"github.com/kubeshop/tracetest/agent/proto"
)

func (c *Client) startHearthBeat(ctx context.Context) error {
	client := proto.NewOrchestratorClient(c.conn)
	ticker := time.NewTicker(c.config.PingPeriod)

	go func() {
		for range ticker.C {
			_, err := client.Ping(ctx, c.sessionConfig.AgentIdentification)
			if err != nil {
				log.Println(err)
				c.reconnect()
			}
		}
	}()

	return nil
}
