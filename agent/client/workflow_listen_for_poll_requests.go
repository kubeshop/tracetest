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
				log.Println("could not get message from poller stream: %w", err)
				continue
			}

			// TODO: Get ctx from request
			err = c.pollListener(context.Background(), &resp)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}()
	return nil
}
