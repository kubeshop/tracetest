package client

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/kubeshop/tracetest/agent/proto"
)

func (c *Client) startTriggerListener() error {
	ctx, cancelCtx := context.WithCancel(context.Background())

	client := proto.NewOrchestratorClient(c.conn)

	request, err := c.getConnectionRequest()
	if err != nil {
		cancelCtx()
		return err
	}

	stream, err := client.RegisterTriggerAgent(ctx, request)
	if err != nil {
		cancelCtx()
		return fmt.Errorf("could not open agent stream: %w", err)
	}

	go func() {
		<-c.done
		cancelCtx()
	}()

	go func() {
		for {
			resp := proto.TriggerRequest{}
			err := stream.RecvMsg(&resp)
			if err == io.EOF {
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
