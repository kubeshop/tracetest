package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/server/telemetry"
)

func (c *Client) startOTLPConnectionTestListener(ctx context.Context) error {
	client := proto.NewOrchestratorClient(c.conn)

	stream, err := client.RegisterOTLPConnectionTestListener(ctx, c.sessionConfig.AgentIdentification)
	if err != nil {
		return fmt.Errorf("could not open agent stream: %w", err)
	}

	go func() {
		for {
			req := proto.OTLPConnectionTestRequest{}
			err := stream.RecvMsg(&req)
			if isEndOfFileError(err) || isCancelledError(err) {
				return
			}

			reconnected, err := c.handleDisconnectionError(err)
			if reconnected {
				return
			}

			if err != nil {
				log.Println("could not get message from otlp connection stream: %w", err)
				time.Sleep(1 * time.Second)
				continue
			}

			ctx, err := telemetry.ExtractContextFromStream(stream)
			if err != nil {
				log.Println("could not extract context from stream %w", err)
			}

			err = c.otlpConnectionTestListener(ctx, &req)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}()
	return nil
}
