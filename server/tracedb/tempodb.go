package tracedb

import (
	"context"
	"fmt"

	tempopb "github.com/kubeshop/tracetest/server/internal/proto-gen-go/tempo-idl"
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

func (*tempoTraceDB) SendSpan(context.Context, trace.Span) error {
	return nil
}

func (ttd *tempoTraceDB) GetTraceByIdentification(ctx context.Context, identification traces.TraceIdentification) (traces.Trace, error) {
	trID := identification.TraceID
	resp, err := ttd.query.FindTraceByID(ctx, &tempopb.TraceByIDRequest{
		TraceID: []byte(trID[:]),
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return traces.Trace{}, fmt.Errorf("tempo FindTraceByID %w", err)
		}
		if st.Message() == "trace not found" {
			return traces.Trace{}, ErrTraceNotFound
		}
		return traces.Trace{}, fmt.Errorf("tempo err: %w", err)
	}

	fmt.Printf("tempo resp: %#v\n", resp.Trace.Batches)
	if len(resp.Trace.Batches) == 0 {
		return traces.Trace{}, ErrTraceNotFound
	}

	otelTrace := v1.TracesData{
		ResourceSpans: resp.GetTrace().GetBatches(),
	}

	return traces.FromOtel(&otelTrace), nil
}

func (ttd *tempoTraceDB) Close() error {
	err := ttd.conn.Close()
	if err != nil {
		return fmt.Errorf("GRPC close: %w", err)
	}
	return nil
}
