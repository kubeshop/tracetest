package client

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"github.com/pkg/errors"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configgrpc"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	name   string
	config *configgrpc.GRPCClientSettings
	conn   *grpc.ClientConn
}

func NewGrpcClient(name string, config *configgrpc.GRPCClientSettings) *GrpcClient {
	return &GrpcClient{
		name:   name,
		config: config,
	}
}

func (client *GrpcClient) Connect(ctx context.Context) error {
	opts, err := client.config.ToDialOptions(nil, componenttest.NewNopTelemetrySettings())
	if err != nil {
		return errors.Wrap(connection.ErrInvalidConfiguration, err.Error())
	}

	conn, err := grpc.DialContext(ctx, client.config.Endpoint, opts...)
	if err != nil {
		return errors.Wrap(connection.ErrConnectionFailed, err.Error())
	}

	client.conn = conn
	return nil
}

func (client *GrpcClient) TestConnection(ctx context.Context) (connection.ConnectionTestResult, error) {
	connectionTestResult := connection.ConnectionTestResult{
		ConnectivityTestResult: connection.ConnectionTestStepResult{
			OperationDescription: fmt.Sprintf(`Tracetest connected to "%s"`, client.config.Endpoint),
		},
	}

	reachable, err := connection.IsReachable(client.config.Endpoint)
	if !reachable {
		return connection.ConnectionTestResult{
			ConnectivityTestResult: connection.ConnectionTestStepResult{
				OperationDescription: fmt.Sprintf(`Tracetest tried to connect to "%s" and failed`, client.config.Endpoint),
				Error:                err,
			},
		}, err
	}

	err = client.Connect(ctx)
	wrappedErr := errors.Unwrap(err)
	if errors.Is(wrappedErr, connection.ErrConnectionFailed) {
		return connection.ConnectionTestResult{
			ConnectivityTestResult: connection.ConnectionTestStepResult{
				OperationDescription: fmt.Sprintf(`Tracetest tried to open a gRPC connection against "%s" and failed`, client.config.Endpoint),
				Error:                err,
			},
		}, err
	}

	return connectionTestResult, nil
}

func (client *GrpcClient) Close() error {
	err := client.conn.Close()
	if err != nil {
		return fmt.Errorf("GRPC close: %w", err)
	}

	return nil
}
