package tracedb

import (
	"context"
	"fmt"
	"io"

	pb "github.com/kubeshop/tracetest/internal/proto-gen-go/api_v3"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configgrpc"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type jaegerTraceDB struct {
	conn  *grpc.ClientConn
	query pb.QueryServiceClient
}

func newJaegerDB(config *configgrpc.GRPCClientSettings) (TraceDB, error) {
	opts, err := config.ToDialOptions(nil, componenttest.NewNopTelemetrySettings())
	if err != nil {
		return nil, fmt.Errorf("jaegerdb grpc config: %w", err)
	}

	conn, err := grpc.Dial(config.Endpoint, opts...)
	if err != nil {
		return nil, fmt.Errorf("jaegerdb grpc dial: %w", err)
	}
	return &jaegerTraceDB{
		conn:  conn,
		query: pb.NewQueryServiceClient(conn),
	}, nil
}

func (jtd *jaegerTraceDB) GetTraceByID(ctx context.Context, traceID string) (*v1.TracesData, error) {
	stream, err := jtd.query.GetTrace(ctx, &pb.GetTraceRequest{
		TraceId: traceID,
	})
	if err != nil {
		return nil, fmt.Errorf("jaeger get trace: %w", err)
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
				return nil, fmt.Errorf("jaeger stream recv: %w", err)
			}
			if st.Message() == "trace not found" {
				return nil, ErrTraceNotFound
			}
			return nil, fmt.Errorf("jaeger stream recv err: %w", err)
		}

		spans = append(spans, in.ResourceSpans...)
	}

	return &v1.TracesData{
		ResourceSpans: spans,
	}, nil
}

func (jtd *jaegerTraceDB) Close() error {
	err := jtd.conn.Close()
	if err != nil {
		return fmt.Errorf("GRPC close: %w", err)
	}
	return nil
}
