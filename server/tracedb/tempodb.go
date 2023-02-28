package tracedb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kubeshop/tracetest/server/id"
	tempopb "github.com/kubeshop/tracetest/server/internal/proto-gen-go/tempo-idl"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"github.com/kubeshop/tracetest/server/tracedb/datasource"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func tempoDefaultPorts() []string {
	return []string{"9095"}
}

type tempoTraceDB struct {
	realTraceDB
	dataSource datasource.DataSource
}

func newTempoDB(config *model.BaseClientConfig) (TraceDB, error) {
	dataSource := datasource.New("Tempo", config, datasource.Callbacks{
		HTTP: httpGetTraceByID,
		GRPC: grpcGetTraceByID,
	})

	return &tempoTraceDB{
		dataSource: dataSource,
	}, nil
}

func (tdb *tempoTraceDB) Connect(ctx context.Context) error {
	return tdb.dataSource.Connect(ctx)
}

func (ttd *tempoTraceDB) TestConnection(ctx context.Context) connection.ConnectionTestResult {
	tester := connection.NewTester(
		connection.WithPortLintingTest(connection.PortLinter(tempoDefaultPorts(), ttd.dataSource.Endpoint())),
		connection.WithConnectivityTest(ttd.dataSource),
		connection.WithPollingTest(connection.TracePollingTestStep(ttd)),
		connection.WithAuthenticationTest(connection.NewTestStep(func(ctx context.Context) (string, error) {
			_, err := ttd.GetTraceByID(ctx, id.NewRandGenerator().TraceID().String())
			if strings.Contains(err.Error(), "authentication handshake failed") {
				return "Tracetest tried to execute a request but it failed due to authentication issues", err
			}

			return "Tracetest managed to authenticate with Tempo", nil
		})),
	)

	return tester.TestConnection(ctx)
}

func (ttd *tempoTraceDB) Ready() bool {
	return ttd.dataSource.Ready()
}

func (ttd *tempoTraceDB) GetTraceByID(ctx context.Context, traceID string) (model.Trace, error) {
	trace, err := ttd.dataSource.GetTraceByID(ctx, traceID)
	return trace, err
}

func (ttd *tempoTraceDB) Close() error {
	return ttd.dataSource.Close()
}

func grpcGetTraceByID(ctx context.Context, traceID string, conn *grpc.ClientConn) (model.Trace, error) {
	query := tempopb.NewQuerierClient(conn)

	trID, err := trace.TraceIDFromHex(traceID)
	if err != nil {
		return model.Trace{}, err
	}

	resp, err := query.FindTraceByID(ctx, &tempopb.TraceByIDRequest{
		TraceID: []byte(trID[:]),
	})
	if err != nil {
		return model.Trace{}, handleError(err)
	}

	if resp.Trace == nil {
		return model.Trace{}, connection.ErrTraceNotFound
	}

	if len(resp.Trace.Batches) == 0 {
		return model.Trace{}, connection.ErrTraceNotFound
	}

	trace := &v1.TracesData{
		ResourceSpans: resp.GetTrace().GetBatches(),
	}

	return traces.FromOtel(trace), nil
}

type HttpTempoTraceByIDResponse struct {
	Batches []*traces.HttpResourceSpans `json:"batches"`
}

func httpGetTraceByID(ctx context.Context, traceID string, client *datasource.HttpClient) (model.Trace, error) {
	trID, err := trace.TraceIDFromHex(traceID)
	if err != nil {
		return model.Trace{}, err
	}
	resp, err := client.Request(ctx, fmt.Sprintf("/api/traces/%s", trID), http.MethodGet, "")

	if err != nil {
		return model.Trace{}, handleError(err)
	}

	if resp.StatusCode == 404 {
		return model.Trace{}, connection.ErrTraceNotFound
	}

	var body []byte
	if b, err := io.ReadAll(resp.Body); err == nil {
		body = b
	} else {
		fmt.Println(err)
	}

	if resp.StatusCode == 401 {
		return model.Trace{}, fmt.Errorf("tempo err: %w %s", errors.New("authentication handshake failed"), string(body))
	}

	var trace HttpTempoTraceByIDResponse
	err = json.Unmarshal(body, &trace)
	if err != nil {
		return model.Trace{}, err
	}

	return traces.FromHttpOtelResourceSpans(trace.Batches), nil
}

func handleError(err error) error {
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return fmt.Errorf("tempo FindTraceByID %w", err)
		}
		if st.Message() == "trace not found" {
			return connection.ErrTraceNotFound
		}
		return fmt.Errorf("tempo err: %w", err)
	}

	return nil
}
