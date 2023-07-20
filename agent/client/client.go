package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

type Config struct {
	APIKey    string
	AgentName string
}

type SessionConfig struct {
	BatchTimeout time.Duration
}

type Client struct {
	conn          *grpc.ClientConn
	config        Config
	sessionConfig *SessionConfig
	done          chan bool
}

func (c *Client) start(ctx context.Context) error {
	return c.startup(ctx)
}

func (c *Client) SessionConfiguration() *SessionConfig {
	if c.sessionConfig == nil {
		return nil
	}

	deferredPtr := *c.sessionConfig
	return &deferredPtr
}

func (c *Client) Close() error {
	c.done <- true
	return c.conn.Close()
}
