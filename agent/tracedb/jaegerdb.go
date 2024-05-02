package tracedb

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	pb "github.com/kubeshop/tracetest/agent/internal/proto-gen-go/api_v3"
	"github.com/kubeshop/tracetest/agent/tracedb/connection"
	"github.com/kubeshop/tracetest/agent/tracedb/datasource"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/traces"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func jaegerDefaultPorts() []string {
	return []string{"16685"}
}

var (
	connectionPool      map[string]TraceDB = make(map[string]TraceDB)
	connectionPoolMutex sync.Mutex
)

func getConnectionFromCache(grpcConfig *datastore.GRPCClientSettings) (TraceDB, bool) {
	connectionPoolMutex.Lock()
	defer connectionPoolMutex.Unlock()

	connection, ok := connectionPool[grpcConfig.Endpoint]
	return connection, ok
}

func saveConnectionToCache(grpcConfig *datastore.GRPCClientSettings, traceDB TraceDB) {
	connectionPoolMutex.Lock()
	defer connectionPoolMutex.Unlock()

	connectionPool[grpcConfig.Endpoint] = traceDB
}

func invalidateConnectionCache(endpoint string) {
	connectionPoolMutex.Lock()
	defer connectionPoolMutex.Unlock()

	delete(connectionPool, endpoint)
}

type jaegerTraceDB struct {
	realTraceDB
	dataSource datasource.DataSource
}

func newJaegerDB(grpcConfig *datastore.GRPCClientSettings) (TraceDB, error) {
	if connection, ok := getConnectionFromCache(grpcConfig); ok {
		return connection, nil
	}

	baseConfig := &datastore.MultiChannelClientConfig{
		Type: datastore.MultiChannelClientTypeGRPC,
		Grpc: grpcConfig,
	}

	dataSource := datasource.New("Jaeger", baseConfig, datasource.Callbacks{
		GRPC: jaegerGrpcGetTraceByID,
	})

	traceDB := &jaegerTraceDB{
		dataSource: dataSource,
	}

	saveConnectionToCache(grpcConfig, traceDB)

	return traceDB, nil
}

func (jtd *jaegerTraceDB) Connect(ctx context.Context) error {
	return jtd.dataSource.Connect(ctx)
}

func (jtd *jaegerTraceDB) GetEndpoints() string {
	return jtd.dataSource.Endpoint()
}

func (jtd *jaegerTraceDB) TestConnection(ctx context.Context) model.ConnectionResult {
	tester := connection.NewTester(
		connection.WithPortLintingTest(connection.PortLinter("Jaeger", jaegerDefaultPorts(), jtd.dataSource.Endpoint())),
		connection.WithConnectivityTest(jtd.dataSource),
		connection.WithPollingTest(connection.TracePollingTestStep(jtd)),
		connection.WithAuthenticationTest(connection.NewTestStep(func(ctx context.Context) (string, error) {
			_, err := jtd.GetTraceByID(ctx, id.NewRandGenerator().TraceID().String())
			if strings.Contains(err.Error(), "authentication handshake failed") {
				return "Tracetest tried to execute a gRPC request but it failed due to authentication issues", err
			}

			return "Tracetest managed to authenticate with Jaeger", nil
		})),
	)

	return tester.TestConnection(ctx)
}

func (jtd *jaegerTraceDB) GetTraceByID(ctx context.Context, traceID string) (traces.Trace, error) {
	trace, err := jtd.dataSource.GetTraceByID(ctx, traceID)
	return trace, err
}

func (jtd *jaegerTraceDB) Ready() bool {
	return jtd.dataSource.Ready()
}

func (jtd *jaegerTraceDB) Close() error {
	invalidateConnectionCache(jtd.dataSource.Endpoint())
	return jtd.dataSource.Close()
}

func jaegerGrpcGetTraceByID(ctx context.Context, traceID string, conn *grpc.ClientConn) (traces.Trace, error) {
	query := pb.NewQueryServiceClient(conn)

	stream, err := query.GetTrace(ctx, &pb.GetTraceRequest{
		TraceId: traceID,
	})
	if err != nil {
		return traces.Trace{}, fmt.Errorf("jaeger get trace: %w", err)
	}

	// jaeger-query v3 API returns otel spans
	var spans []*v1.ResourceSpans
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				return traces.Trace{}, fmt.Errorf("jaeger stream recv: %w", err)
			}
			if strings.Contains(st.Message(), "trace not found") {
				return traces.Trace{}, connection.ErrTraceNotFound
			}
			return traces.Trace{}, fmt.Errorf("jaeger stream recv err: %w", err)
		}

		spans = append(spans, in.ResourceSpans...)
	}

	trace := &v1.TracesData{
		ResourceSpans: spans,
	}

	return traces.FromOtel(trace), nil
}
