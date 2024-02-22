package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/server/telemetry"
	"go.uber.org/zap"
)

func (c *Client) startPollerListener(ctx context.Context) error {
	logger := c.logger.Named("pollerListener")
	logger.Debug("Starting")

	client := proto.NewOrchestratorClient(c.conn)

	stream, err := client.RegisterPollerAgent(ctx, c.sessionConfig.AgentIdentification)
	if err != nil {
		return fmt.Errorf("could not open agent stream: %w", err)
	}

	go func() {
		for {
			resp := proto.PollingRequest{}
			err := stream.RecvMsg(&resp)
			if err != nil {
				logger.Error("could not get message from poller stream", zap.Error(err))
			}
			if isEndOfFileError(err) || isCancelledError(err) {
				logger.Debug("poller stream closed")
				return
			}

			reconnected, err := c.handleDisconnectionError(err)
			if reconnected {
				logger.Debug("reconnected to poller stream")
				return
			}

			if err != nil {
				logger.Error("could not get message from poller stream", zap.Error(err))
				log.Println("could not get message from poller stream: %w", err)
				time.Sleep(1 * time.Second)
				continue
			}

			ctx, err := telemetry.ExtractContextFromStream(stream)
			if err != nil {
				logger.Error("could not extract context from stream", zap.Error(err))
				log.Println("could not extract context from stream %w", err)
			}

			err = c.pollListener(ctx, &resp)
			if err != nil {
				logger.Error("could not handle poll request", zap.Error(err))
				fmt.Println(err.Error())
			}
		}
	}()
	return nil
}
