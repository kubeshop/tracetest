package tracedb

import (
	"context"
	"fmt"
	"strings"

	tempopb "github.com/kubeshop/tracetest/server/internal/proto-gen-go/tempo-idl"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/pkg/errors"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/otel/trace"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type tempoTraceDB struct {
	config *configgrpc.GRPCClientSettings
	conn   *grpc.ClientConn
	query  tempopb.QuerierClient
}

func newTempoDB(config *configgrpc.GRPCClientSettings) (TraceDB, error) {
	return &tempoTraceDB{
		config: config,
	}, nil
}

func (tdb *tempoTraceDB) Connect(ctx context.Context) error {
	opts, err := tdb.config.ToDialOptions(nil, componenttest.NewNopTelemetrySettings())
	if err != nil {
		return errors.Wrap(ErrInvalidConfiguration, err.Error())
	}

	conn, err := grpc.DialContext(ctx, tdb.config.Endpoint, opts...)
	if err != nil {
		return errors.Wrap(ErrConnectionFailed, err.Error())
	}

	tdb.conn = conn
	tdb.query = tempopb.NewQuerierClient(conn)

	return nil
}

func (ttd *tempoTraceDB) TestConnection(ctx context.Context) ConnectionTestResult {
	reachable, err := isReachable(ttd.config.Endpoint)
	if !reachable {
		return ConnectionTestResult{
			ConnectivityTestResult: ConnectionTestStepResult{
				OperationDescription: fmt.Sprintf(`Tracetest tried to connect to "%s" and failed`, ttd.config.Endpoint),
				Error:                err,
			},
		}
	}

	err = ttd.Connect(ctx)
	wrappedErr := errors.Unwrap(err)
	if errors.Is(wrappedErr, ErrConnectionFailed) {
		return ConnectionTestResult{
			ConnectivityTestResult: ConnectionTestStepResult{
				OperationDescription: fmt.Sprintf(`Tracetest tried to open a gRPC connection against "%s" and failed`, ttd.config.Endpoint),
				Error:                err,
			},
		}
	}

	_, err = ttd.GetTraceByID(ctx, trace.TraceID{}.String())
	if strings.Contains(err.Error(), "authentication handshake failed") {
		return ConnectionTestResult{
			AuthenticationTestResult: ConnectionTestStepResult{
				OperationDescription: `Tracetest tried to execue a gRPC request but it failed due to authentication issues`,
				Error:                err,
			},
		}
	}

	if strings.Contains(err.Error(), "connection error") {
		return ConnectionTestResult{
			ConnectivityTestResult: ConnectionTestStepResult{
				OperationDescription: fmt.Sprintf(`Tracetest tried to open a gRPC connection against "%s" and failed`, ttd.config.Endpoint),
				Error:                err,
			},
		}
	}

	if !errors.Is(err, ErrTraceNotFound) {
		return ConnectionTestResult{
			TraceRetrivalTestResult: ConnectionTestStepResult{
				OperationDescription: fmt.Sprintf(`Tracetest tried to fetch a trace from Tempo endpoint "%s" and got an error`, ttd.config.Endpoint),
				Error:                err,
			},
		}
	}

	return ConnectionTestResult{
		ConnectivityTestResult: ConnectionTestStepResult{
			OperationDescription: fmt.Sprintf(`Tracetest connected to "%s"`, ttd.config.Endpoint),
		},
		AuthenticationTestResult: ConnectionTestStepResult{
			OperationDescription: `Tracetest managed to authenticate with Tempo`,
		},
		TraceRetrivalTestResult: ConnectionTestStepResult{
			OperationDescription: `Tracetest was able to search for a trace using Tempo API`,
		},
	}
}

func (ttd *tempoTraceDB) GetTraceByID(ctx context.Context, traceID string) (traces.Trace, error) {
	trID, err := trace.TraceIDFromHex(traceID)
	if err != nil {
		return traces.Trace{}, err
	}
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

	if resp.Trace == nil {
		return traces.Trace{}, ErrTraceNotFound
	}

	if len(resp.Trace.Batches) == 0 {
		return traces.Trace{}, ErrTraceNotFound
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
