package client

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/kubeshop/tracetest/agent/proto"
)

func (c *Client) startDataStoreConnectionTestListener(ctx context.Context) error {
	client := proto.NewOrchestratorClient(c.conn)

	stream, err := client.RegisterDataStoreConnectionTestAgent(ctx, c.sessionConfig.AgentIdentification)
	if err != nil {
		return fmt.Errorf("could not open agent stream: %w", err)
	}

	go func() {
		for {
			req := proto.DataStoreConnectionTestRequest{}
			err := stream.RecvMsg(&req)
			if err == io.EOF {
				return
			}

			if err != nil {
				log.Fatal("could not get message from trigger stream: %w", err)
			}

			// TODO: Get ctx from request
			c.dataStoreConnectionListener(context.Background(), &req)
		}
	}()
	return nil
}
