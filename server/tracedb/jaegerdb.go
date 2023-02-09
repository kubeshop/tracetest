package tracedb

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/id"
	pb "github.com/kubeshop/tracetest/server/internal/proto-gen-go/api_v3"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb/client"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/pkg/errors"
	"go.opentelemetry.io/collector/config/configgrpc"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type jaegerTraceDB struct {
	realTraceDB
	client *client.Client
}

func newJaegerDB(grpcConfig *configgrpc.GRPCClientSettings) (TraceDB, error) {
	baseConfig := &config.BaseClientConfig{
		Type: string(client.GRPC),
		Grpc: *grpcConfig,
	}

	client := client.NewClient("Jaeger", baseConfig, client.Callbacks{
		GRPC: jaegerGrpcGetTraceByID,
	})

	return &jaegerTraceDB{
		client: client,
	}, nil
}

func (jtd *jaegerTraceDB) Connect(ctx context.Context) error {
	return jtd.client.Connect(ctx)
}

func (jtd *jaegerTraceDB) TestConnection(ctx context.Context) connection.ConnectionTestResult {
	connectionTestResult, err := jtd.client.TestConnection(ctx)
	if err != nil {
		return connectionTestResult
	}

	_, err = jtd.GetTraceByID(ctx, id.NewRandGenerator().TraceID().String())
	if strings.Contains(err.Error(), "authentication handshake failed") {
		return connection.ConnectionTestResult{
			ConnectivityTestResult: connectionTestResult.ConnectivityTestResult,
			AuthenticationTestResult: connection.ConnectionTestStepResult{
				OperationDescription: "Tracetest tried to execute a gRPC request but it failed due to authentication issues",
				Error:                err,
			},
		}
	}

	if !errors.Is(err, connection.ErrTraceNotFound) {
		return connection.ConnectionTestResult{
			ConnectivityTestResult:   connectionTestResult.ConnectivityTestResult,
			AuthenticationTestResult: connectionTestResult.AuthenticationTestResult,
			TraceRetrievalTestResult: connection.ConnectionTestStepResult{
				OperationDescription: "Tracetest tried to fetch a trace from Jaeger",
				Error:                err,
			},
		}
	}

	connectionTestResult.AuthenticationTestResult = connection.ConnectionTestStepResult{
		OperationDescription: `Tracetest managed to authenticate with Jaeger`,
	}
	connectionTestResult.TraceRetrievalTestResult = connection.ConnectionTestStepResult{
		OperationDescription: `Tracetest was able to search for a trace using Jaeger API`,
	}
	return connectionTestResult
}

func (jtd *jaegerTraceDB) GetTraceByID(ctx context.Context, traceID string) (model.Trace, error) {
	trace, err := jtd.client.GetTraceByID(ctx, traceID)
	return trace, err
}

func (jtd *jaegerTraceDB) Ready() bool {
	return jtd.client.Ready()
}

func (jtd *jaegerTraceDB) Close() error {
	return jtd.client.Close()
}

func jaegerGrpcGetTraceByID(ctx context.Context, traceID string, conn *grpc.ClientConn) (model.Trace, error) {
	query := pb.NewQueryServiceClient(conn)

	stream, err := query.GetTrace(ctx, &pb.GetTraceRequest{
		TraceId: traceID,
	})
	if err != nil {
		return model.Trace{}, fmt.Errorf("jaeger get trace: %w", err)
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
				return model.Trace{}, fmt.Errorf("jaeger stream recv: %w", err)
			}
			if st.Message() == "trace not found" {
				return model.Trace{}, connection.ErrTraceNotFound
			}
			return model.Trace{}, fmt.Errorf("jaeger stream recv err: %w", err)
		}

		spans = append(spans, in.ResourceSpans...)
	}

	trace := &v1.TracesData{
		ResourceSpans: spans,
	}

	return traces.FromOtel(trace), nil
}
