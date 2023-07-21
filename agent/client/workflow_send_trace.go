package client

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/agent/proto"
)

func (c *Client) SendTrace(ctx context.Context, request *proto.PollingRequest, spans ...*proto.Span) error {
	client := proto.NewOrchestratorClient(c.conn)

	pollingResponse := &proto.PollingResponse{
		TestID:  request.TestID,
		RunID:   request.RunID,
		TraceID: request.TraceID,
		Spans:   spans,
	}

	_, err := client.SendPolledSpans(ctx, pollingResponse)
	if err != nil {
		return fmt.Errorf("could not send trigger result request: %w", err)
	}

	return nil
}
