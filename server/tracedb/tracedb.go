package tracedb

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
)

var IDGen = id.NewRandGenerator()

type TraceDB interface {
	Connect(ctx context.Context) error
	Ready() bool
	ShouldRetry() bool
	GetTraceID() trace.TraceID
	GetTraceByID(ctx context.Context, traceID string) (traces.Trace, error)
	Close() error
	GetEndpoints() string
}

type TestableTraceDB interface {
	TraceDB
	TestConnection(ctx context.Context) model.ConnectionResult
}

type TraceAugmenter interface {
	AugmentTrace(ctx context.Context, trace *traces.Trace) (*traces.Trace, error)
}

type noopTraceDB struct{}

func (db *noopTraceDB) GetTraceByID(ctx context.Context, traceID string) (t traces.Trace, err error) {
	return traces.Trace{}, nil
}

func (db *noopTraceDB) GetTraceID() trace.TraceID {
	return IDGen.TraceID()
}
func (db *noopTraceDB) Connect(ctx context.Context) error { return nil }
func (db *noopTraceDB) Close() error                      { return nil }
func (db *noopTraceDB) ShouldRetry() bool                 { return false }
func (db *noopTraceDB) Ready() bool                       { return true }
func (db *noopTraceDB) GetEndpoints() string              { return "" }
func (db *noopTraceDB) TestConnection(ctx context.Context) model.ConnectionResult {
	return model.ConnectionResult{}
}

type traceDBFactory struct {
	traceGetter traceGetter
}

type FactoryFunc func(ds datastore.DataStore) (TraceDB, error)

func Factory(tg traceGetter) FactoryFunc {
	f := traceDBFactory{
		traceGetter: tg,
	}

	return f.New
}

func (f *traceDBFactory) getTraceDBInstance(ds datastore.DataStore) (TraceDB, error) {
	var tdb TraceDB
	var err error

	if ds.IsOTLPBasedProvider() {
		tdb, err = newCollectorDB(f.traceGetter)
		return tdb, err
	}

	switch ds.Type {
	case datastore.DataStoreTypeJaeger:
		tdb, err = newJaegerDB(ds.Values.Jaeger)
	case datastore.DataStoreTypeTempo:
		tdb, err = newTempoDB(ds.Values.Tempo)
	case datastore.DataStoreTypeElasticAPM:
		tdb, err = newElasticSearchDB(ds.Values.ElasticApm)
	case datastore.DataStoreTypeOpenSearch:
		tdb, err = newOpenSearchDB(ds.Values.OpenSearch)
	case datastore.DataStoreTypeSignalFX:
		tdb, err = newSignalFXDB(ds.Values.SignalFx)
	case datastore.DataStoreTypeAwsXRay:
		tdb, err = NewAwsXRayDB(ds.Values.AwsXRay)
	case datastore.DatastoreTypeAzureAppInsights:
		tdb, err = NewAzureAppInsightsDB(ds.Values.AzureAppInsights)
	case datastore.DatastoreTypeSumoLogic:
		tdb, err = NewSumoLogicDB(ds.Values.SumoLogic)
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

func (f *traceDBFactory) New(ds datastore.DataStore) (TraceDB, error) {
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
