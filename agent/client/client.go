package client

import (
	"context"

	"google.golang.org/grpc"
)

type Client struct {
	conn *grpc.ClientConn
	done chan bool
}

func (c *Client) start() {
}

func (c *Client) Send(ctx context.Context, request any) (any, error) {
	return nil, nil
}

func (c *Client) Listen(string, func(ctx context.Context, request any) error) error {
	return nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
