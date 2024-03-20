package datasource

import (
	"context"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/traces"
	"google.golang.org/grpc"
)

type SupportedDataSource string

var (
	GRPC SupportedDataSource = "grpc"
	HTTP SupportedDataSource = "http"
)

type HttpCallback func(ctx context.Context, traceID string, client *HttpClient) (traces.Trace, error)
type GrpcCallback func(ctx context.Context, traceID string, connection *grpc.ClientConn) (traces.Trace, error)

type Callbacks struct {
	GRPC GrpcCallback
	HTTP HttpCallback
}

type DataSource interface {
	Endpoint() string
	Connect(ctx context.Context) error
	Ready() bool
	GetTraceByID(ctx context.Context, traceID string) (traces.Trace, error)
	TestConnection(ctx context.Context) model.ConnectionTestStep
	Close() error
}

type noopDataSource struct{}

func (dataSource *noopDataSource) GetTraceByID(ctx context.Context, traceID string) (t traces.Trace, err error) {
	return traces.Trace{}, nil
}
func (db *noopDataSource) Endpoint() string                  { return "" }
func (db *noopDataSource) Connect(ctx context.Context) error { return nil }
func (db *noopDataSource) Close() error                      { return nil }
func (db *noopDataSource) Ready() bool                       { return true }
func (db *noopDataSource) TestConnection(ctx context.Context) model.ConnectionTestStep {
	return model.ConnectionTestStep{}
}

func New(name string, cfg *datastore.MultiChannelClientConfig, callbacks Callbacks) DataSource {
	sourceType := SupportedDataSource(cfg.Type)

	switch sourceType {
	default:
	case GRPC:
		return NewGrpcClient(name, cfg.Grpc, callbacks.GRPC)
	case HTTP:
		return NewHttpClient(name, cfg.Http, callbacks.HTTP)
	}

	return &noopDataSource{}
}
