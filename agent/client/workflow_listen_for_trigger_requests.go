package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/telemetry"
	"go.uber.org/zap"
)

func (c *Client) startTriggerListener(ctx context.Context) error {
	logger := c.logger.Named("triggerListener")
	logger.Debug("Starting")

	client := proto.NewOrchestratorClient(c.conn)

	stream, err := client.RegisterTriggerAgent(ctx, c.sessionConfig.AgentIdentification)
	if err != nil {
		logger.Error("could not open agent stream", zap.Error(err))
		return fmt.Errorf("could not open agent stream: %w", err)
	}

	go func() {
		for {
			req := proto.TriggerRequest{}
			err := stream.RecvMsg(&req)
			if err != nil {
				logger.Error("could not get message from stop stream", zap.Error(err))
			}
			if isEndOfFileError(err) || isCancelledError(err) {
				logger.Debug("stop stream closed")
				return
			}

			reconnected, err := c.handleDisconnectionError(err, &req)
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

			// we want a new context per request, not to reuse the one from the stream
			ctx := telemetry.InjectMetadataIntoContext(context.Background(), req.Metadata)
			go func() {
				err = c.triggerListener(ctx, &req)
				if err != nil {
					logger.Error("could not handle trigger request", zap.Error(err))
					fmt.Println(err.Error())
				}
			}()
		}
	}()
	return nil
}
