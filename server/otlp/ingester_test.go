package otlp_test

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"testing"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/otlp"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	pb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

type stubConnector struct{}

func (stubConnector) Connect(context.Context) (driver.Conn, error) { return stubConn{}, nil }
func (stubConnector) Driver() driver.Driver                        { return stubDriver{} }

type stubDriver struct{}

func (stubDriver) Open(string) (driver.Conn, error) { return stubConn{}, nil }

type stubConn struct{}

func (stubConn) Prepare(string) (driver.Stmt, error) { return stubStmt{}, nil }
func (stubConn) Close() error                        { return nil }
func (stubConn) Begin() (driver.Tx, error)           { return stubTx{}, nil }

type stubTx struct{}

func (stubTx) Commit() error   { return nil }
func (stubTx) Rollback() error { return nil }

type stubStmt struct{}

func (stubStmt) Close() error  { return nil }
func (stubStmt) NumInput() int { return -1 }
func (stubStmt) Exec([]driver.Value) (driver.Result, error) {
	return driver.RowsAffected(0), nil
}
func (stubStmt) Query([]driver.Value) (driver.Rows, error) { return stubRows{}, nil }

type stubRows struct{}

func (stubRows) Columns() []string {
	return []string{"id", "name", "type", "is_default", "values", "created_at"}
}
func (stubRows) Close() error              { return nil }
func (stubRows) Next([]driver.Value) error { return io.EOF }

type mapRunGetter struct {
	runs map[string]test.Run
}

func (m mapRunGetter) GetRunByTraceID(_ context.Context, traceID trace.TraceID) (test.Run, error) {
	run, found := m.runs[traceID.String()]
	if !found {
		return test.Run{}, sql.ErrNoRows
	}

	return run, nil
}

type recordingPersister struct {
	savedTraceIDs []trace.TraceID
}

func (r *recordingPersister) UpdateTraceSpans(_ context.Context, t *traces.Trace) error {
	r.savedTraceIDs = append(r.savedTraceIDs, t.ID)
	return nil
}

type noopEventEmitter struct{}

func (noopEventEmitter) Emit(context.Context, model.TestRunEvent) error { return nil }

func TestIngestProcessesAwaitingTraceWhenBatchHasFinishedTrace(t *testing.T) {
	finishedTraceID := bytes.Repeat([]byte{0xAA}, 16)
	awaitingTraceID := bytes.Repeat([]byte{0xBB}, 16)

	request := &pb.ExportTraceServiceRequest{
		ResourceSpans: []*v1.ResourceSpans{
			{
				ScopeSpans: []*v1.ScopeSpans{
					{
						Spans: []*v1.Span{
							{TraceId: finishedTraceID, SpanId: []byte{1, 1, 1, 1, 1, 1, 1, 1}, Name: "finished-span"},
							{TraceId: awaitingTraceID, SpanId: []byte{2, 2, 2, 2, 2, 2, 2, 2}, Name: "awaiting-span"},
						},
					},
				},
			},
		},
	}

	runGetter := mapRunGetter{runs: map[string]test.Run{
		trace.TraceID(finishedTraceID).String(): {State: test.RunStateFinished},
		trace.TraceID(awaitingTraceID).String(): {State: test.RunStateAwaitingTrace},
	}}

	persister := &recordingPersister{}
	ingester := otlp.NewIngester(
		persister,
		runGetter,
		noopEventEmitter{},
		datastore.NewRepository(sql.OpenDB(stubConnector{})),
		subscription.NewManager(),
		trace.NewNoopTracerProvider().Tracer(""),
	)

	_, err := ingester.Ingest(context.Background(), request, otlp.RequestTypeHTTP)

	require.NoError(t, err)
	require.Len(t, persister.savedTraceIDs, 1)
	assert.Equal(t, trace.TraceID(awaitingTraceID).String(), persister.savedTraceIDs[0].String())
}
