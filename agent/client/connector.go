package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func Connect(ctx context.Context, endpoint string, opts ...Option) (*Client, error) {
	config := Config{
		PingPeriod: 30 * time.Second,
	}

	client := &Client{endpoint: endpoint, config: config}
	for _, opt := range opts {
		opt(client)
	}

	err := client.connect(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// See https://pkg.go.dev/google.golang.org/grpc/examples/features/retry#section-readme
// these values were calculated to get a max retry of aprox 120s. This script can be used to play with values:
// curl -SL https://gist.githubusercontent.com/schoren/5dd4dcadf133e4c56fa20c0b6a8d67ef/raw/1992d34d27368578d22e87b061ff00de6f7023be/calculate_max_time.sh | bash -s  -- 0.1 2.0 5.0 120
var retryPolicy = `{
	"methodConfig": [{
			"name": [{}],
			"retryPolicy": {
				"maxAttempts": 29,
				"initialBackoff": "0.1s",
				"maxBackoff": "5s",
				"backoffMultiplier": 1.5,
				"retryableStatusCodes": ["UNAVAILABLE"]
			}
	}]
}`

func (c *Client) connect(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	transportCredentials, err := getTransportCredentialsForEndpoint(c.endpoint)
	if err != nil {
		return fmt.Errorf("could not get transport credentials: %w", err)
	}

	conn, err := grpc.DialContext(
		ctx, c.endpoint,
		grpc.WithTransportCredentials(transportCredentials),
		grpc.WithDefaultServiceConfig(retryPolicy),
		grpc.WithIdleTimeout(0), // disable grpc idle timeout
	)
	if err != nil {
		return fmt.Errorf("could not connect to server: %w", err)
	}

	c.conn = conn
	return nil
}

func getTransportCredentialsForEndpoint(endpoint string) (credentials.TransportCredentials, error) {
	_, port, err := net.SplitHostPort(endpoint)
	if err != nil {
		return nil, fmt.Errorf("cannot parse endpoint: %w", err)
	}

	switch port {
	case "443":
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		transportCredentials := credentials.NewTLS(tlsConfig)
		return transportCredentials, nil

	default:
		return insecure.NewCredentials(), nil
	}

}
