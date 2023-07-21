package client

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/kubeshop/tracetest/agent/proto"
)

func (c *Client) startPollerListener() error {
	ctx := context.Background()

	client := proto.NewOrchestratorClient(c.conn)

	request, err := c.getConnectionRequest()
	if err != nil {
		return err
	}

	stream, err := client.RegisterPollerAgent(ctx, request)
	if err != nil {
		return fmt.Errorf("could not open agent stream: %w", err)
	}

	go func() {
		for {
			resp := proto.PollingRequest{}
			err := stream.RecvMsg(&resp)
			if err == io.EOF {
				return
			}

			if err != nil {
				log.Fatal("could not get message from trigger stream: %w", err)
			}

			c.pollListener(&resp)
		}
	}()
	return nil
}
