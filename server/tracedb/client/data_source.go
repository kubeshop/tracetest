package client

import (
	"context"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"google.golang.org/grpc"
)

type SupportedDataSource string

var (
	GRPC SupportedDataSource = "grpc"
	HTTP SupportedDataSource = "http"
)

type HttpCallback func(ctx context.Context, traceID string, client HttpClient) (model.Trace, error)
type GrpcCallback func(ctx context.Context, traceID string, connection *grpc.ClientConn) (model.Trace, error)

type Callbacks struct {
	GRPC GrpcCallback
	HTTP HttpCallback
}

type DataSource struct {
	sourceType SupportedDataSource
	grpcClient *GrpcClient
	httpClient *HttpClient
	callbacks  Callbacks
}

func NewDataSource(name string, config *config.BaseClientConfig, callbacks Callbacks) *DataSource {
	sourceType := SupportedDataSource(config.Type)

	if sourceType == GRPC {
		return &DataSource{
			sourceType: sourceType,
			grpcClient: NewGrpcClient(name, &config.Grpc),
			callbacks:  callbacks,
		}
	}

	return &DataSource{
		sourceType: SupportedDataSource(config.Type),
		callbacks:  callbacks,
		httpClient: NewHttpClient(name, &config.Http),
	}
}

func (dataSource *DataSource) GetTraceByID(ctx context.Context, traceID string) (model.Trace, error) {
	switch dataSource.sourceType {
	case GRPC:
		return dataSource.callbacks.GRPC(ctx, traceID, dataSource.grpcClient.conn)
	case HTTP:
		return dataSource.callbacks.HTTP(ctx, traceID, *dataSource.httpClient)
	}

	return model.Trace{}, nil
}

func (dataSource *DataSource) Connect(ctx context.Context) error {
	switch dataSource.sourceType {
	case GRPC:
		return dataSource.grpcClient.Connect(ctx)
	}

	return nil
}

func (dataSource *DataSource) Ready() bool {
	switch dataSource.sourceType {
	case GRPC:
		return dataSource.grpcClient.conn != nil
	}

	return true
}

func (dataSource *DataSource) Close() error {
	switch dataSource.sourceType {
	case GRPC:
		return dataSource.grpcClient.Close()
	}

	return nil
}

func (dataSource *DataSource) TestConnection(ctx context.Context) (connection.ConnectionTestResult, error) {
	switch dataSource.sourceType {
	case GRPC:
		return dataSource.grpcClient.TestConnection(ctx)
	case HTTP:
		return dataSource.httpClient.TestConnection(ctx)
	}

	return connection.ConnectionTestResult{}, nil
}
