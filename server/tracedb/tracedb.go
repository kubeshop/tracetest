package tracedb

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/tracedb/datastoreresource"
	"go.opentelemetry.io/otel/trace"
)

var IDGen = id.NewRandGenerator()

type TraceDB interface {
	Connect(ctx context.Context) error
	Ready() bool
	ShouldRetry() bool
	GetTraceID() trace.TraceID
	GetTraceByID(ctx context.Context, traceID string) (model.Trace, error)
	Close() error
}

type TestableTraceDB interface {
	TraceDB
	TestConnection(ctx context.Context) model.ConnectionResult
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
func (db *noopTraceDB) TestConnection(ctx context.Context) model.ConnectionResult {
	return model.ConnectionResult{}
}

type traceDBFactory struct {
	runRepository model.RunRepository
}

func Factory(runRepository model.RunRepository) func(ds datastoreresource.DataStore) (TraceDB, error) {
	f := traceDBFactory{
		runRepository: runRepository,
	}

	return f.New
}

func (f *traceDBFactory) getTraceDBInstance(ds datastoreresource.DataStore) (TraceDB, error) {
	var tdb TraceDB
	var err error

	if ds.IsOTLPBasedProvider() {
		tdb, err = newCollectorDB(f.runRepository)
		return tdb, err
	}

	switch ds.Type {
	case datastoreresource.DataStoreTypeJaeger:
		tdb, err = newJaegerDB(ds.Values.Jaeger)
	case datastoreresource.DataStoreTypeTempo:
		tdb, err = newTempoDB(ds.Values.Tempo)
	case datastoreresource.DataStoreTypeElasticAPM:
		tdb, err = newElasticSearchDB(ds.Values.ElasticApm)
	case datastoreresource.DataStoreTypeOpenSearch:
		tdb, err = newOpenSearchDB(ds.Values.OpenSearch)
	case datastoreresource.DataStoreTypeSignalFX:
		tdb, err = newSignalFXDB(ds.Values.SignalFx)
	case datastoreresource.DataStoreTypeAwsXRay:
		tdb, err = NewAwsXRayDB(ds.Values.AwsXRay)
	default:
		return &noopTraceDB{}, nil
	}

	if err != nil {
		return nil, err
	}

	if tdb == nil {
		return nil, fmt.Errorf("data store unknown: %s", ds.Type)
	}

	return tdb, err
}

func (f *traceDBFactory) New(ds datastoreresource.DataStore) (TraceDB, error) {
	tdb, err := f.getTraceDBInstance(ds)

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
