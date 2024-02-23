package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kubeshop/tracetest/agent/proto"
	"go.uber.org/zap"
)

// TODO: fix this and add test
func (c *Client) startStopListener(ctx context.Context) error {
	logger := c.logger.Named("stopListener")
	logger.Debug("Starting")

	client := proto.NewOrchestratorClient(c.conn)

	stream, err := client.RegisterStopRequestAgent(ctx, c.sessionConfig.AgentIdentification)
	if err != nil {
		logger.Error("could not open agent stream", zap.Error(err))
		return fmt.Errorf("could not open agent stream: %w", err)
	}

	go func() {
		for {
			resp := proto.StopRequest{}
			err := stream.RecvMsg(&resp)
			if err != nil {
				logger.Error("could not get message from stop stream", zap.Error(err))
			}
			if isEndOfFileError(err) || isCancelledError(err) {
				logger.Debug("stop stream closed")
				continue
			}

			reconnected, err := c.handleDisconnectionError(err)
			if reconnected {
				logger.Warn("reconnected to stop stream")
				return
			}

			if err != nil {
				logger.Error("could not get message from stop stream", zap.Error(err))
				log.Println("could not get message from trigger stream: %w", err)
				time.Sleep(1 * time.Second)
				continue
			}

			// TODO: get context from request
			err = c.stopListener(context.Background(), &resp)
			if err != nil {
				logger.Error("could not handle stop request", zap.Error(err))
				fmt.Println(err.Error())
			}
		}
	}()
	return nil
}
