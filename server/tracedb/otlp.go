package tracedb

import (
	"context"
	"strings"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"github.com/kubeshop/tracetest/server/traces"
)

type OTLPTraceDB struct {
	realTraceDB
	db model.RunRepository
}

func newCollectorDB(repository model.RunRepository) (TraceDB, error) {
	return &OTLPTraceDB{
		db: repository,
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

func (jtd *OTLPTraceDB) TestConnection(ctx context.Context) model.ConnectionResult {
	return model.ConnectionResult{}
}

// GetTraceByID implements TraceDB
func (tdb *OTLPTraceDB) GetTraceByID(ctx context.Context, id string) (model.Trace, error) {
	run, err := tdb.db.GetRunByTraceID(ctx, traces.DecodeTraceID(id))
	if err != nil && strings.Contains(err.Error(), "record not found") {
		return model.Trace{}, connection.ErrTraceNotFound
	}

	if run.Trace == nil {
		return model.Trace{}, connection.ErrTraceNotFound
	}

	return *run.Trace, nil
}

var _ TraceDB = &OTLPTraceDB{}
