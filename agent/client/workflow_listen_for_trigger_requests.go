package client

import (
	"context"
	"fmt"
	"log"

	retry "github.com/avast/retry-go"
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
			if isEndOfFileError(err) || isCancelledError(err) {
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
				log.Println("could not get message from trigger stream: %w", err)
				continue
			}

			// TODO: get context from request
			err = c.triggerListener(context.Background(), &resp)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}()
	return nil
}
