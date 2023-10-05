package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/kubeshop/tracetest/agent/proto"
)

func (c *Client) startTriggerListener(ctx context.Context) error {
	client := proto.NewOrchestratorClient(c.conn)

	stream, err := client.RegisterTriggerAgent(ctx, c.sessionConfig.AgentIdentification)
	if err != nil {
		return fmt.Errorf("could not open agent stream: %w", err)
	}

	go func() {
		for {
			resp := proto.TriggerRequest{}
			err := stream.RecvMsg(&resp)
			if errors.Is(err, io.EOF) || isCancelledError(err) {
				return
			}

			if err != nil {
				log.Fatal("could not get message from trigger stream: %w", err)
			}

			// TODO: get context from request
			c.triggerListener(context.Background(), &resp)
		}
	}()
	return nil
}
