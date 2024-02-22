package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kubeshop/tracetest/agent/proto"
	"go.uber.org/zap"
)

func (c *Client) startShutdownListener(ctx context.Context) error {
	logger := c.logger.Named("shutdownListener")
	logger.Debug("Starting")

	client := proto.NewOrchestratorClient(c.conn)

	stream, err := client.RegisterShutdownListener(ctx, c.sessionConfig.AgentIdentification)
	if err != nil {
		logger.Error("could not open agent stream", zap.Error(err))
		return fmt.Errorf("could not open agent stream: %w", err)
	}

	go func() {
		for {
			resp := proto.ShutdownRequest{}
			err := stream.RecvMsg(&resp)

			if err != nil {
				logger.Error("could not get message from shutdown stream", zap.Error(err))
			}

			if isEndOfFileError(err) || isCancelledError(err) {
				logger.Debug("shutdown stream closed")
				return
			}

			reconnected, err := c.handleDisconnectionError(err)
			if reconnected {
				logger.Warn("reconnected to shutdown stream")
				return
			}

			if err != nil {
				logger.Error("could not get message from shutdown stream", zap.Error(err))
				log.Println("could not get message from shutdown stream: %w", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// TODO: get context from request
			err = c.shutdownListener(context.Background(), &resp)
			if err != nil {
				logger.Error("could not handle shutdown request", zap.Error(err))
				fmt.Println(err.Error())
			}
		}
	}()
	return nil
}
