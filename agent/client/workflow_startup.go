package client

import (
	"context"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/agent/proto"
)

// The startup workflow consists in exchanging information about the server
// and the agent.
//
// Agent sends information about authentication and identification, server responds with
// a configuration object that must be used when connected to that server.
func (c *Client) startup(ctx context.Context) error {
	orchestratorClient := proto.NewOrchestratorClient(c.conn)

	request, err := c.getConnectionRequest()
	if err != nil {
		return err
	}

	response, err := orchestratorClient.Connect(ctx, request)
	if err != nil {
		return fmt.Errorf("could not send request to server: %w", err)
	}

	c.sessionConfig = &SessionConfig{
		BatchTimeout:        time.Duration(response.Configuration.BatchTimeout) * time.Millisecond,
		AgentIdentification: response.Identification,
	}

	return nil
}
