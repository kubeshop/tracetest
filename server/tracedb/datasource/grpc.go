package datasource

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"github.com/pkg/errors"
	"go.opentelemetry.io/collector/component/componenttest"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	name     string
	config   *model.GRPCClientSettings
	conn     *grpc.ClientConn
	callback GrpcCallback
}

func NewGrpcClient(name string, config *model.GRPCClientSettings, callback GrpcCallback) DataSource {
	return &GrpcClient{
		name:     name,
		config:   config,
		callback: callback,
	}
}

func (client *GrpcClient) Ready() bool {
	return client.conn != nil
}

func (client *GrpcClient) GetTraceByID(ctx context.Context, traceID string) (model.Trace, error) {
	return client.callback(ctx, traceID, client.conn)
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
