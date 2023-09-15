package tracedb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
)

type traceGetter interface {
	Get(context.Context, trace.TraceID) (traces.Trace, error)
}

type OTLPTraceDB struct {
	realTraceDB
	traceGetter traceGetter
}

func newCollectorDB(repository traceGetter) (TraceDB, error) {
	return &OTLPTraceDB{
		traceGetter: repository,
	}, nil
}

func (tdb *OTLPTraceDB) Ready() bool {
	return true
}

func (tdb *OTLPTraceDB) Connect(ctx context.Context) error {
	return nil
}

func (tdb *OTLPTraceDB) Close() error {
	// No need to implement this
	return nil
}

func (tdb *OTLPTraceDB) GetEndpoints() string {
	return ""
}

// GetTraceByID implements TraceDB
func (tdb *OTLPTraceDB) GetTraceByID(ctx context.Context, id string) (traces.Trace, error) {
	t, err := tdb.traceGetter.Get(ctx, traces.DecodeTraceID(id))
	if errors.Is(err, sql.ErrNoRows) {
		err = connection.ErrTraceNotFound
	}

	return t, err
}
