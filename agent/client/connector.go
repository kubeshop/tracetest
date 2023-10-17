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
	conn, err := connect(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	client := &Client{conn: conn}
	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

func connect(ctx context.Context, endpoint string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	transportCredentials, err := getTransportCredentialsForEndpoint(endpoint)
	if err != nil {
		return nil, fmt.Errorf("could not get transport credentials: %w", err)
	}

	conn, err := grpc.DialContext(
		ctx, endpoint,
		grpc.WithTransportCredentials(transportCredentials),
	)
	if err != nil {
		return nil, fmt.Errorf("could not connect to server: %w", err)
	}

	return conn, nil
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
