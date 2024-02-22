package client

import (
	"context"
	"log"
	"time"

	"github.com/kubeshop/tracetest/agent/proto"
	"go.uber.org/zap"
)

func (c *Client) startHeartBeat(ctx context.Context) error {
	logger := c.logger.Named("pingListener")
	logger.Debug("Starting")

	client := proto.NewOrchestratorClient(c.conn)
	ticker := time.NewTicker(c.config.PingPeriod)

	go func() {
		for range ticker.C {
			_, err := client.Ping(ctx, c.sessionConfig.AgentIdentification)
			if err != nil {
				log.Println("could not send ping: %w", err)
			}
			if isEndOfFileError(err) || isCancelledError(err) {
				log.Println("ping stream closed")
				return
			}

			reconnected, err := c.handleDisconnectionError(err)
			if reconnected {
				log.Println("reconnected to ping stream")
				return
			}

			if err != nil {
				logger.Error("could not get message from ping stream", zap.Error(err))
				log.Println("could not get message from ping stream: %w", err)
				time.Sleep(1 * time.Second)
				continue
			}
		}
	}()

	return nil
}
