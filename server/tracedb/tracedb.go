package tracedb

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/tracedb/connection"
	"go.opentelemetry.io/otel/trace"
)

var IDGen = id.NewRandGenerator()

type TraceDB interface {
	Connect(ctx context.Context) error
	Ready() bool
	ShouldRetry() bool
	GetTraceID() trace.TraceID
	GetTraceByID(ctx context.Context, traceID string) (model.Trace, error)
	TestConnection(ctx context.Context) connection.ConnectionTestResult
	Close() error
}

type noopTraceDB struct{}

func (db *noopTraceDB) GetTraceByID(ctx context.Context, traceID string) (t model.Trace, err error) {
	return model.Trace{}, nil
}

func (db *noopTraceDB) GetTraceID() trace.TraceID {
	return IDGen.TraceID()
}
func (db *noopTraceDB) Connect(ctx context.Context) error { return nil }
func (db *noopTraceDB) Close() error                      { return nil }
func (db *noopTraceDB) ShouldRetry() bool                 { return false }
func (db *noopTraceDB) Ready() bool                       { return true }
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
	case model.DataStoreTypeJaeger:
		tdb, err = newJaegerDB(ds.Values.Jaeger)
	case model.DataStoreTypeTempo:
		tdb, err = newTempoDB(ds.Values.Tempo)
	case model.DataStoreTypeElasticAPM:
		tdb, err = newElasticSearchDB(ds.Values.ElasticApm)
	case model.DataStoreTypeOpenSearch:
		tdb, err = newOpenSearchDB(ds.Values.OpenSearch)
	case model.DataStoreTypeSignalFX:
		tdb, err = newSignalFXDB(ds.Values.SignalFx)
	case model.DataStoreTypeAwsXRay:
		tdb, err = NewAwsXRayDB(ds.Values.AwsXRay)
	case model.DataStoreTypeNewRelic:
	case model.DataStoreTypeLighStep:
	case model.DataStoreTypeDataDog:
	case model.DataStoreTypeOTLP:
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
func (db *realTraceDB) GetTraceID() trace.TraceID {
	return IDGen.TraceID()
}
