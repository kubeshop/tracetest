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

	c.closeStreamWhenShuttingDown(stream)

	go func() {
		for {
			resp := proto.PollingRequest{}
			err := stream.RecvMsg(&resp)
			if err == io.EOF {
				return
			}

			if errors.Is(err, context.Canceled) {
				// probably stream was closed, so just skip it
				return
			}

			fmt.Println(err)

			if err != nil {
				log.Fatal("could not get message from trigger stream: %w", err)
			}

			// TODO: Get ctx from request
			c.pollListener(context.Background(), &resp)
		}
	}()
	return nil
}
