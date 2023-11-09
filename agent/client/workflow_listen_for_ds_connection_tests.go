package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	retry "github.com/avast/retry-go"
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

			if err != nil && isConnectionError(err) {
				err = retry.Do(func() error {
					return c.reconnect()
				})
				if err == nil {
					// everything was reconnect, so we can exist this goroutine
					// as there's another one running in parallel
					return
				}

				log.Fatal(err)
			}

			if err != nil {
				log.Println("could not get message from data store connection stream: %w", err)
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
