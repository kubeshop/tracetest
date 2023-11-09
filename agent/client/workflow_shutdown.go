package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/kubeshop/tracetest/agent/proto"
)

func (c *Client) startShutdownListener(ctx context.Context) error {
	client := proto.NewOrchestratorClient(c.conn)

	stream, err := client.RegisterShutdownListener(ctx, c.sessionConfig.AgentIdentification)
	if err != nil {
		return fmt.Errorf("could not open agent stream: %w", err)
	}

	go func() {
		for {
			resp := proto.ShutdownRequest{}
			err := stream.RecvMsg(&resp)

			if errors.Is(err, io.EOF) || isCancelledError(err) {
				return
			}

			if err != nil {
				c.reconnect()
				continue
			}

			if c.shutdownListener == nil {
				log.Fatal("warning: shutdown listener is nil")
			}

			// TODO: get context from request
			err = c.shutdownListener(context.Background(), &resp)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}()
	return nil
}
