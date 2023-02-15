package tracedb

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
)

type TraceDB interface {
	Connect(ctx context.Context) error
	Ready() bool
	ShouldRetry() bool
	MinSpanCount() int
	GetTraceByID(ctx context.Context, traceID string) (model.Trace, error)
	TestConnection(ctx context.Context) connection.ConnectionTestResult
	Close() error
}

type noopTraceDB struct{}

func (db *noopTraceDB) GetTraceByID(ctx context.Context, traceID string) (t model.Trace, err error) {
	return model.Trace{}, nil
}

func (db *noopTraceDB) Connect(ctx context.Context) error { return nil }
func (db *noopTraceDB) Close() error                      { return nil }
func (db *noopTraceDB) ShouldRetry() bool                 { return false }
func (db *noopTraceDB) Ready() bool                       { return true }
func (db *noopTraceDB) MinSpanCount() int                 { return 0 }
func (db *noopTraceDB) TestConnection(ctx context.Context) connection.ConnectionTestResult {
	return connection.ConnectionTestResult{}
}

func WithFallback(fn func(ds model.DataStore) (TraceDB, error), fallbackDS model.DataStore) func(ds model.DataStore) (TraceDB, error) {
	return func(ds model.DataStore) (TraceDB, error) {
		if ds.IsZero() {
			ds = fallbackDS
		}
		return fn(ds)
	}
}

type traceDBFactory struct {
	repo model.Repository
}

func Factory(repo model.Repository) func(ds model.DataStore) (TraceDB, error) {
	f := traceDBFactory{
		repo: repo,
	}

	return f.New
}

func (f *traceDBFactory) New(ds model.DataStore) (tdb TraceDB, err error) {
	switch ds.Type {
	case openapi.JAEGER:
		tdb, err = newJaegerDB(ds.Values.Jaeger)
	case openapi.TEMPO:
		tdb, err = newTempoDB(ds.Values.Tempo)
	case openapi.ELASTIC_APM:
		tdb, err = newElasticSearchDB(ds.Values.ElasticApm)
	case openapi.OPEN_SEARCH:
		tdb, err = newOpenSearchDB(ds.Values.OpenSearch)
	case openapi.SIGNAL_FX:
		tdb, err = newSignalFXDB(ds.Values.SignalFx)
	case openapi.NEW_RELIC:
	case openapi.LIGHTSTEP:
	case openapi.DATADOG:
	case openapi.OTLP:
		tdb, err = newCollectorDB(f.repo)
	default:
		return &noopTraceDB{}, nil
	}

	if err != nil {
		return nil, err
	}

	err = tdb.Connect(context.Background())
	if err != nil {
		return nil, fmt.Errorf("cannot connect to datasource: %w", err)
	}

	return tdb, nil

}

type realTraceDB struct{}

func (db *realTraceDB) ShouldRetry() bool { return true }
func (db *realTraceDB) MinSpanCount() int { return 0 }
