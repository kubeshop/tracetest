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

func (c *Client) startOTLPConnectionTestListener(ctx context.Context) error {
	logger := c.logger.Named("otlpConnectionTestListener")
	logger.Debug("Starting")

	client := proto.NewOrchestratorClient(c.conn)

	stream, err := client.RegisterOTLPConnectionTestListener(ctx, c.sessionConfig.AgentIdentification)
	if err != nil {
		logger.Error("could not open agent stream", zap.Error(err))
		return fmt.Errorf("could not open agent stream: %w", err)
	}

	go func() {
		for {
			req := proto.OTLPConnectionTestRequest{}
			err := stream.RecvMsg(&req)
			if err != nil {
				logger.Error("could not get message from otlp connection stream", zap.Error(err))
			}
			if isEndOfFileError(err) || isCancelledError(err) {
				logger.Debug("otlp connection stream closed")
				return
			}

			reconnected, err := c.handleDisconnectionError(err)
			if reconnected {
				logger.Warn("reconnected to otlp connection stream")
				return
			}

			if err != nil {
				logger.Error("could not get message from otlp connection stream", zap.Error(err))
				log.Println("could not get message from otlp connection stream: %w", err)
				time.Sleep(1 * time.Second)
				continue
			}

			ctx, err := telemetry.ExtractContextFromStream(stream)
			if err != nil {
				logger.Error("could not extract context from stream", zap.Error(err))
				log.Println("could not extract context from stream %w", err)
			}

			err = c.otlpConnectionTestListener(ctx, &req)
			if err != nil {
				logger.Error("could not handle otlp connection test request", zap.Error(err))
				fmt.Println(err.Error())
			}
		}
	}()
	return nil
}
