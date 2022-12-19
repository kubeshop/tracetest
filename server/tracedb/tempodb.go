package tracedb

import (
	"context"
	"fmt"

	tempopb "github.com/kubeshop/tracetest/server/internal/proto-gen-go/tempo-idl"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/otel/trace"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type tempoTraceDB struct {
	conn  *grpc.ClientConn
	query tempopb.QuerierClient
}

func newTempoDB(config *configgrpc.GRPCClientSettings) (TraceDB, error) {
	opts, err := config.ToDialOptions(nil, componenttest.NewNopTelemetrySettings())
	if err != nil {
		return nil, fmt.Errorf("tempodb grpc config: %w", err)
	}

	conn, err := grpc.Dial(config.Endpoint, opts...)
	if err != nil {
		return nil, fmt.Errorf("tempodb grpc dial: %w", err)
	}

	return &tempoTraceDB{
		conn:  conn,
		query: tempopb.NewQuerierClient(conn),
	}, nil
}

func (ttd *tempoTraceDB) GetTraceByID(ctx context.Context, traceID string) (model.Trace, error) {
	trID, err := trace.TraceIDFromHex(traceID)
	if err != nil {
		return model.Trace{}, err
	}
	resp, err := ttd.query.FindTraceByID(ctx, &tempopb.TraceByIDRequest{
		TraceID: []byte(trID[:]),
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return model.Trace{}, fmt.Errorf("tempo FindTraceByID %w", err)
		}
		if st.Message() == "trace not found" {
			return model.Trace{}, ErrTraceNotFound
		}
		return model.Trace{}, fmt.Errorf("tempo err: %w", err)
	}

	if resp.Trace == nil {
		return model.Trace{}, ErrTraceNotFound
	}

	if len(resp.Trace.Batches) == 0 {
		return model.Trace{}, ErrTraceNotFound
	}

	trace := &v1.TracesData{
		ResourceSpans: resp.GetTrace().GetBatches(),
	}

	return traces.FromOtel(trace), nil
}

func (ttd *tempoTraceDB) Close() error {
	err := ttd.conn.Close()
	if err != nil {
		return fmt.Errorf("GRPC close: %w", err)
	}
	return nil
}
