package client

import (
	"context"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"google.golang.org/grpc"
)

type SupportedClients string

var (
	GRPC SupportedClients = "grpc"
	HTTP SupportedClients = "http"
)

type HttpCallback func(ctx context.Context, traceID string, client HttpClient) (model.Trace, error)
type GrpcCallback func(ctx context.Context, traceID string, connection *grpc.ClientConn) (model.Trace, error)

type Callbacks struct {
	GRPC GrpcCallback
	HTTP HttpCallback
}

type Client struct {
	clientType SupportedClients
	grpcClient *GrpcClient
	httpClient *HttpClient
	callbacks  Callbacks
}

func NewClient(name string, config *config.BaseClientConfig, callbacks Callbacks) *Client {
	clientType := SupportedClients(config.Type)

	if clientType == GRPC {
		return &Client{
			clientType: clientType,
			grpcClient: NewGrpcClient(name, &config.Grpc),
			callbacks:  callbacks,
		}
	}

	return &Client{
		clientType: SupportedClients(config.Type),
		callbacks:  callbacks,
		httpClient: NewHttpClient(name, &config.Http),
	}
}

func (client *Client) GetTraceByID(ctx context.Context, traceID string) (model.Trace, error) {
	switch client.clientType {
	case GRPC:
		return client.callbacks.GRPC(ctx, traceID, client.grpcClient.conn)
	case HTTP:
		return client.callbacks.HTTP(ctx, traceID, *client.httpClient)
	}

	return model.Trace{}, nil
}

func (client *Client) Connect(ctx context.Context) error {
	switch client.clientType {
	case GRPC:
		return client.grpcClient.Connect(ctx)
	}

	return nil
}

func (client *Client) Ready() bool {
	switch client.clientType {
	case GRPC:
		return client.grpcClient.conn != nil
	}

	return true
}

func (client *Client) Close() error {
	switch client.clientType {
	case GRPC:
		return client.grpcClient.Close()
	}

	return nil
}

func (client *Client) TestConnection(ctx context.Context) (connection.ConnectionTestResult, error) {
	switch client.clientType {
	case GRPC:
		return client.grpcClient.TestConnection(ctx)
	case HTTP:
		return client.httpClient.TestConnection(ctx)
	}

	return connection.ConnectionTestResult{}, nil
}
