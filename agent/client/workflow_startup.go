package client

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kubeshop/tracetest/agent/proto"
)

// The startup workflow consists in exchanging information about the server
// and the agent.
//
// Agent sends information about authentication and identification, server responds with
// a configuration object that must be used when connected to that server.
func (c *Client) startup(ctx context.Context) error {
	name, err := c.getName()
	if err != nil {
		return err
	}

	orchestratorClient := proto.NewOrchestratorClient(c.conn)
	request := proto.ConnectRequest{
		ApiKey: c.config.APIKey,
		Name:   name,
	}

	response, err := orchestratorClient.Connect(ctx, &request)
	if err != nil {
		return fmt.Errorf("could not send request to server: %w", err)
	}

	c.sessionConfig = &SessionConfig{
		BatchTimeout: time.Duration(response.Configuration.BatchTimeout) * time.Millisecond,
	}

	return nil
}

// getName retrieves the name of the agent. By default, it is the host name, however,
// it can be overwritten with an environment variable, or a flag.
func (c *Client) getName() (string, error) {
	if name := c.config.AgentName; name != "" {
		return name, nil
	}

	if name := os.Getenv("TRACETEST_AGENT_NAME"); name != "" {
		return name, nil
	}

	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("could not get hostname: %w", err)
	}

	return hostname, nil
}
