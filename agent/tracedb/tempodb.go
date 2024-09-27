package tracedb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/golang/protobuf/proto"
	tempopb "github.com/kubeshop/tracetest/agent/internal/proto-gen-go/tempo-idl"
	"github.com/kubeshop/tracetest/agent/tracedb/connection"
	"github.com/kubeshop/tracetest/agent/tracedb/datasource"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func tempoDefaultPorts() []string {
	return []string{"9095", "443", "80", ""}
}

type tempoTraceDB struct {
	realTraceDB
	dataSource datasource.DataSource
}

func newTempoDB(config *datastore.MultiChannelClientConfig) (TraceDB, error) {
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

func (ttd *tempoTraceDB) GetEndpoints() string {
	return ttd.dataSource.Endpoint()
}

func (ttd *tempoTraceDB) TestConnection(ctx context.Context) model.ConnectionResult {
	tester := connection.NewTester(
		connection.WithPortLintingTest(connection.PortLinter("Tempo", tempoDefaultPorts(), ttd.dataSource.Endpoint())),
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

func (ttd *tempoTraceDB) GetTraceByID(ctx context.Context, traceID string) (traces.Trace, error) {
	trace, err := ttd.dataSource.GetTraceByID(ctx, traceID)
	return trace, err
}

func (ttd *tempoTraceDB) Close() error {
	return ttd.dataSource.Close()
}

func grpcGetTraceByID(ctx context.Context, traceID string, conn *grpc.ClientConn) (traces.Trace, error) {
	query := tempopb.NewQuerierClient(conn)

	trID, err := trace.TraceIDFromHex(traceID)
	if err != nil {
		return traces.Trace{}, err
	}

	// This is a temporary solution while we figure out how to paginate the traces from Tempo
	maxSizeOption := grpc.MaxCallRecvMsgSize(1024 * 1024 * 10) // 10MB

	resp, err := query.FindTraceByID(context.Background(), &tempopb.TraceByIDRequest{
		TraceID: []byte(trID[:]),
	}, maxSizeOption)
	if err != nil {
		return traces.Trace{}, handleError(err)
	}

	if resp.Trace == nil {
		return traces.Trace{}, connection.ErrTraceNotFound
	}

	if len(resp.Trace.Batches) == 0 {
		return traces.Trace{}, connection.ErrTraceNotFound
	}

	trace := &v1.TracesData{
		ResourceSpans: resp.GetTrace().GetBatches(),
	}

	return traces.FromOtel(trace), nil
}

type HttpTempoTraceByIDResponse struct {
	Batches []*traces.HttpResourceSpans `json:"batches"`
}

func httpGetTraceByID(ctx context.Context, traceID string, client *datasource.HttpClient) (traces.Trace, error) {
	trID, err := trace.TraceIDFromHex(traceID)
	if err != nil {
		return traces.Trace{}, err
	}
	resp, err := client.Request(ctx, fmt.Sprintf("/api/traces/%s", trID), http.MethodGet, "")
	if err != nil {
		return traces.Trace{}, handleError(err)
	}

	if resp == nil {
		return traces.Trace{}, fmt.Errorf("could not get response")
	}

	if resp.StatusCode == 404 {
		return traces.Trace{}, connection.ErrTraceNotFound
	}

	var body []byte
	if b, err := io.ReadAll(resp.Body); err == nil {
		body = b
	} else {
		fmt.Println(err)
	}

	if resp.StatusCode == 401 {
		return traces.Trace{}, fmt.Errorf("tempo err: %w %s", errors.New("authentication handshake failed"), string(body))
	}

	if contentTypeAcceptsType(resp, "application/protobuf") {
		var protoMessage tempopb.Trace
		err = proto.Unmarshal(body, &protoMessage)
		if err != nil {
			return traces.Trace{}, err
		}

		trace := &v1.TracesData{
			ResourceSpans: protoMessage.GetBatches(),
		}

		return traces.FromOtel(trace), nil
	}
	var trace HttpTempoTraceByIDResponse
	err = json.Unmarshal(body, &trace)
	if err != nil {
		return traces.Trace{}, err
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

func contentTypeAcceptsType(resp *http.Response, accept string) bool {
	for _, headerValue := range resp.Header["Content-Type"] {
		if headerValue == accept {
			return true
		}
	}

	return false
}
