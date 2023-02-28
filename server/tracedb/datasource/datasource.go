package datasource

import (
	"context"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"google.golang.org/grpc"
)

type SupportedDataSource string

var (
	GRPC SupportedDataSource = "grpc"
	HTTP SupportedDataSource = "http"
)

type HttpCallback func(ctx context.Context, traceID string, client *HttpClient) (model.Trace, error)
type GrpcCallback func(ctx context.Context, traceID string, connection *grpc.ClientConn) (model.Trace, error)

type Callbacks struct {
	GRPC GrpcCallback
	HTTP HttpCallback
}

type DataSource interface {
	Endpoint() string
	Connect(ctx context.Context) error
	Ready() bool
	GetTraceByID(ctx context.Context, traceID string) (model.Trace, error)
	TestConnection(ctx context.Context) connection.ConnectionTestStepResult
	Close() error
}

type noopDataSource struct{}

func (dataSource *noopDataSource) GetTraceByID(ctx context.Context, traceID string) (t model.Trace, err error) {
	return model.Trace{}, nil
}
func (db *noopDataSource) Endpoint() string                  { return "" }
func (db *noopDataSource) Connect(ctx context.Context) error { return nil }
func (db *noopDataSource) Close() error                      { return nil }
func (db *noopDataSource) Ready() bool                       { return true }
func (db *noopDataSource) TestConnection(ctx context.Context) connection.ConnectionTestStepResult {
	return connection.ConnectionTestStepResult{}
}

func New(name string, cfg *model.BaseClientConfig, callbacks Callbacks) DataSource {
	sourceType := SupportedDataSource(cfg.Type)

	switch sourceType {
	default:
	case GRPC:
		return NewGrpcClient(name, &cfg.Grpc, callbacks.GRPC)
	case HTTP:
		return NewHttpClient(name, &cfg.Http, callbacks.HTTP)
	}

	return &noopDataSource{}
}
