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
			if isEndOfFileError(err) || isCancelledError(err) {
				return
			}

			reconnected, err := c.handleDisconnectionError(err)
			if reconnected {
				return
			}

			if err != nil {
				log.Println("could not get message from ping stream: %w", err)
				time.Sleep(1 * time.Second)
				continue
			}
		}
	}()

	return nil
}
