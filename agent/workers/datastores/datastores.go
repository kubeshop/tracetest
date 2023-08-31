package datastores

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

type DataStore interface {
	Connect(ctx context.Context) error
	Ready() bool
	ShouldRetry() bool
	GetTraceID() trace.TraceID
	GetTraceByID(ctx context.Context, traceID string) (traces.Trace, error)
	Close() error
	GetEndpoints() string
}

type TestableDataStore interface {
	DataStore
	TestConnection(ctx context.Context) model.ConnectionResult
}

type noopDataStore struct{}

func (db *noopDataStore) GetTraceByID(ctx context.Context, traceID string) (t traces.Trace, err error) {
	return traces.Trace{}, nil
}

func (db *noopDataStore) GetTraceID() trace.TraceID {
	return IDGen.TraceID()
}
func (db *noopDataStore) Connect(ctx context.Context) error { return nil }
func (db *noopDataStore) Close() error                      { return nil }
func (db *noopDataStore) ShouldRetry() bool                 { return false }
func (db *noopDataStore) Ready() bool                       { return true }
func (db *noopDataStore) GetEndpoints() string              { return "" }
func (db *noopDataStore) TestConnection(ctx context.Context) model.ConnectionResult {
	return model.ConnectionResult{}
}

type dataStoreFactory struct{}

func Factory() func(ds datastore.DataStore) (DataStore, error) {
	f := dataStoreFactory{}

	return f.New
}

func (f *dataStoreFactory) getDatastoreInstance(ds datastore.DataStore) (DataStore, error) {
	var tdb DataStore
	var err error

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
	default:
		return &noopDataStore{}, nil
	}

	if err != nil {
		return nil, err
	}

	if tdb == nil {
		return nil, fmt.Errorf("data store unknown: %s", ds.Type)
	}

	return tdb, err
}

func (f *dataStoreFactory) New(ds datastore.DataStore) (DataStore, error) {
	datastore, err := f.getDatastoreInstance(ds)

	if err != nil {
		return nil, err
	}

	err = datastore.Connect(context.Background())
	if err != nil {
		return nil, fmt.Errorf("cannot connect to datasource: %w", err)
	}

	return datastore, nil
}

type realDataStore struct{}

func (db *realDataStore) ShouldRetry() bool { return true }
func (db *realDataStore) GetTraceID() trace.TraceID {
	return IDGen.TraceID()
}
