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
			req := proto.StopRequest{}
			err := stream.RecvMsg(&req)
			if err != nil {
				logger.Error("could not get message from stop stream", zap.Error(err))
			}
			if isEndOfFileError(err) || isCancelledError(err) {
				logger.Debug("stop stream closed")
				return
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

			// we want a new context per request, not to reuse the one from the stream
			ctx := telemetry.InjectMetadataIntoContext(context.Background(), req.Metadata)
			err = c.stopListener(ctx, &req)
			if err != nil {
				logger.Error("could not handle stop request", zap.Error(err))
				fmt.Println(err.Error())
			}
		}
	}()
	return nil
}
