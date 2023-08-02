package client

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/kubeshop/tracetest/agent/proto"
)

func (c *Client) startShutdownListener(ctx context.Context) error {
	client := proto.NewOrchestratorClient(c.conn)

	request, err := c.getConnectionRequest()
	if err != nil {
		return err
	}

	stream, err := client.RegisterShutdownListener(ctx, request)
	if err != nil {
		return fmt.Errorf("could not open agent stream: %w", err)
	}

	go func() {
		for {
			resp := proto.ShutdownRequest{}
			err := stream.RecvMsg(&resp)
			if err == io.EOF {
				return
			}

			if err != nil {
				log.Fatal("could not get message from trigger stream: %w", err)
			}

			// TODO: get context from request
			c.shutdownListener(context.Background(), &resp)
		}
	}()
	return nil
}
