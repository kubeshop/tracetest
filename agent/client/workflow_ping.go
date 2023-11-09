package client

import (
	"context"
	"log"
	"time"

	retry "github.com/avast/retry-go"
	"github.com/kubeshop/tracetest/agent/proto"
)

func (c *Client) startHearthBeat(ctx context.Context) error {
	client := proto.NewOrchestratorClient(c.conn)
	ticker := time.NewTicker(c.config.PingPeriod)

	go func() {
		for range ticker.C {
			_, err := client.Ping(ctx, c.sessionConfig.AgentIdentification)
			if err != nil && isConnectionError(err) {
				err = retry.Do(func() error {
					return c.reconnect()
				})
				if err == nil {
					// everything was reconnect, so we can exist this goroutine
					// as there's another one running in parallel
					return
				}

				log.Fatal(err)
			}

			if err != nil {
				log.Println("could not get message from ping stream: %w", err)
				continue
			}
		}
	}()

	return nil
}
