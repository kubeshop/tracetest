package client

import (
	"context"
	"errors"
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
			if errors.Is(err, io.EOF) || isCancelledError(err) {
				return
			}

			if err != nil {
				c.reconnect()
			}

			if c.dataStoreConnectionListener == nil {
				log.Println("warning: datastore connection listener is nil")
				continue
			}

			// TODO: Get ctx from request
			err = c.dataStoreConnectionListener(context.Background(), &req)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}()
	return nil
}
