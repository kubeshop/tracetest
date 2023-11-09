package client

import (
	"context"
	"fmt"
	"log"
	"time"

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
			if isEndOfFileError(err) || isCancelledError(err) {
				return
			}

			reconnected, err := c.handleDisconnectionError(err)
			if reconnected {
				return
			}

			if err != nil {
				log.Println("could not get message from poller stream: %w", err)
				time.Sleep(1 * time.Second)
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
