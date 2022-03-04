package tempodb

import (
	"context"
	"fmt"

	tempopb "github.com/kubeshop/tracetest/server/go/internal/proto-gen-go/tempo-idl"
	"github.com/kubeshop/tracetest/server/go/tracedb"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/otel/trace"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type TempoTraceDB struct {
	conn  *grpc.ClientConn
	query tempopb.QuerierClient
}

func New(config *configgrpc.GRPCClientSettings) (tracedb.TraceDB, error) {
	opts, err := config.ToDialOptions(nil, componenttest.NewNopTelemetrySettings())
	if err != nil {
		return nil, fmt.Errorf("tempodb grpc config: %w", err)
	}

	conn, err := grpc.Dial(config.Endpoint, opts...)
	if err != nil {
		return nil, fmt.Errorf("tempodb grpc dial: %w", err)
	}
	return &TempoTraceDB{
		conn:  conn,
		query: tempopb.NewQuerierClient(conn),
	}, nil
}

func (ttd *TempoTraceDB) GetTraceByID(ctx context.Context, traceID string) (*v1.TracesData, error) {
	trID, err := trace.TraceIDFromHex(traceID)
	if err != nil {
		return nil, err
	}
	resp, err := ttd.query.FindTraceByID(ctx, &tempopb.TraceByIDRequest{
		TraceID: []byte(trID[:]),
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return nil, fmt.Errorf("tempo FindTraceByID %w", err)
		}
		if st.Message() == "trace not found" {
			return nil, tracedb.ErrTraceNotFound
		}
		return nil, fmt.Errorf("tempo err: %w", err)
	}

	if len(resp.Trace.Batches) == 0 {
		return nil, tracedb.ErrTraceNotFound
	}
	fmt.Printf("tempo resp: %#v\n", resp.Trace.Batches)
	return &v1.TracesData{
		ResourceSpans: resp.GetTrace().GetBatches(),
	}, nil
}

func (ttd *TempoTraceDB) Close() error {
	err := ttd.conn.Close()
	if err != nil {
		return fmt.Errorf("GRPC close: %w", err)
	}
	return nil
}
