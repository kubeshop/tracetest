package tracedb

import (
	"context"
	"fmt"
	"io"
	"strings"

	pb "github.com/kubeshop/tracetest/server/internal/proto-gen-go/api_v3"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/pkg/errors"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/otel/trace"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type jaegerTraceDB struct {
	config *configgrpc.GRPCClientSettings
	conn   *grpc.ClientConn
	query  pb.QueryServiceClient
}

func newJaegerDB(config *configgrpc.GRPCClientSettings) (TraceDB, error) {
	return &jaegerTraceDB{
		config: config,
	}, nil
}

func (jtd *jaegerTraceDB) Connect(ctx context.Context) error {
	opts, err := jtd.config.ToDialOptions(nil, componenttest.NewNopTelemetrySettings())
	if err != nil {
		return errors.Wrap(ErrInvalidConfiguration, err.Error())
	}

	conn, err := grpc.DialContext(ctx, jtd.config.Endpoint, opts...)
	if err != nil {
		return errors.Wrap(ErrConnectionFailed, err.Error())
	}

	jtd.conn = conn
	jtd.query = pb.NewQueryServiceClient(conn)

	return nil
}

func (jtd *jaegerTraceDB) TestConnection(ctx context.Context) ConnectionTestResult {
	// TODO: when the test fails, we should still keep all messages.
	// Current implementation makes only errors to show, so if everything worked but the last
	// step, there's no feedback that tracetest was able to ping the service and authenticate to it.

	reachable, err := isReachable(jtd.config.Endpoint)
	if !reachable {
		return ConnectionTestResult{
			ConnectivityTestResult: ConnectionTestStepResult{
				OperationDescription: fmt.Sprintf(`Tracetest tried to connect to "%s" and failed`, jtd.config.Endpoint),
				Error:                err,
			},
		}
	}

	err = jtd.Connect(ctx)
	wrappedErr := errors.Unwrap(err)
	if errors.Is(wrappedErr, ErrConnectionFailed) {
		return ConnectionTestResult{
			ConnectivityTestResult: ConnectionTestStepResult{
				OperationDescription: fmt.Sprintf(`Tracetest tried to open a gRPC connection against "%s" and failed`, jtd.config.Endpoint),
				Error:                err,
			},
		}
	}

	_, err = jtd.GetTraceByID(ctx, trace.TraceID{}.String())
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
				OperationDescription: fmt.Sprintf(`Tracetest tried to open a gRPC connection against "%s" and failed`, jtd.config.Endpoint),
				Error:                err,
			},
		}
	}

	if !errors.Is(err, ErrTraceNotFound) {
		return ConnectionTestResult{
			TraceRetrivalTestResult: ConnectionTestStepResult{
				OperationDescription: fmt.Sprintf(`Tracetest tried to fetch a trace from Jaeger endpoint "%s" and got an error`, jtd.config.Endpoint),
				Error:                err,
			},
		}
	}

	return ConnectionTestResult{
		ConnectivityTestResult: ConnectionTestStepResult{
			OperationDescription: fmt.Sprintf(`Tracetest connected to "%s"`, jtd.config.Endpoint),
		},
		AuthenticationTestResult: ConnectionTestStepResult{
			OperationDescription: `Tracetest managed to authenticate with Jaeger`,
		},
		TraceRetrivalTestResult: ConnectionTestStepResult{
			OperationDescription: `Tracetest was able to search for a trace using Jaeger API`,
		},
	}
}

func (jtd *jaegerTraceDB) GetTraceByID(ctx context.Context, traceID string) (model.Trace, error) {
	stream, err := jtd.query.GetTrace(ctx, &pb.GetTraceRequest{
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
				return model.Trace{}, ErrTraceNotFound
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

func (jtd *jaegerTraceDB) Close() error {
	err := jtd.conn.Close()
	if err != nil {
		return fmt.Errorf("GRPC close: %w", err)
	}
	return nil
}
