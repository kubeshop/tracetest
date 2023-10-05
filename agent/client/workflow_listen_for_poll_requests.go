package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/kubeshop/tracetest/agent/proto"
)

func (c *Client) startPollerListener(ctx context.Context) error {
	client := proto.NewOrchestratorClient(c.conn)

	stream, err := client.RegisterPollerAgent(ctx, c.sessionConfig.AgentIdentification)
	if err != nil {
		return fmt.Errorf("could not open agent stream: %w", err)
	}

	go func() {
		for {
			resp := proto.PollingRequest{}
			err := stream.RecvMsg(&resp)
			if errors.Is(err, io.EOF) || isCancelledError(err) {
				return
			}

			if err != nil {
				log.Fatal("could not get message from polling stream: %w", err)
			}

			// TODO: Get ctx from request
			c.pollListener(context.Background(), &resp)
		}
	}()
	return nil
}
