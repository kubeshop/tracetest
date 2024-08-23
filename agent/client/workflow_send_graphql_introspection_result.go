package client

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/agent/proto"
	"github.com/kubeshop/tracetest/agent/telemetry"
)

func (c *Client) SendGraphqlIntrospectionResult(ctx context.Context, response *proto.GraphqlIntrospectResponse) error {
	client := proto.NewOrchestratorClient(c.conn)

	response.AgentIdentification = c.sessionConfig.AgentIdentification
	response.Metadata = telemetry.ExtractMetadataFromContext(ctx)

	_, err := client.SendGraphqlIntrospectResult(ctx, response)
	if err != nil {
		return fmt.Errorf("could not send graphql introspection result request: %w", err)
	}

	return nil
}
