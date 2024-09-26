package client

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/telemetry"
)

func (c *Client) SendTraceModeResponse(ctx context.Context, response *proto.TraceModeResponse) error {
	client := proto.NewOrchestratorClient(c.conn)

	response.AgentIdentification = c.sessionConfig.AgentIdentification
	response.Metadata = telemetry.ExtractMetadataFromContext(ctx)

	_, err := client.SendTraceModeResponse(ctx, response)
	if err != nil {
		return fmt.Errorf("could not send list traces result request: %w", err)
	}

	return nil
}
