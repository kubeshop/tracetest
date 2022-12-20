package tracedb

import (
	"context"
	"errors"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/openapi"
)

var (
	ErrTraceNotFound        = errors.New("trace not found")
	ErrInvalidConfiguration = errors.New("invalid data store configuration")
	ErrConnectionFailed     = errors.New("could not connect to data store")
)

const (
	JAEGER_BACKEND     string = "jaeger"
	TEMPO_BACKEND      string = "tempo"
	OPENSEARCH_BACKEND string = "opensearch"
	SIGNALFX           string = "signalfx"
	OTLP               string = "otlp"
)

type TraceDB interface {
	Connect(ctx context.Context) error
	ShouldRetry() bool
	MinSpanCount() int
	GetTraceByID(ctx context.Context, traceID string) (model.Trace, error)
	TestConnection(ctx context.Context) ConnectionTestResult
	Close() error
}

type noopTraceDB struct{}

func (db *noopTraceDB) GetTraceByID(ctx context.Context, traceID string) (t model.Trace, err error) {
	return model.Trace{}, nil
}

func (db *noopTraceDB) Connect(ctx context.Context) error { return nil }
func (db *noopTraceDB) Close() error                      { return nil }
func (db *noopTraceDB) ShouldRetry() bool                 { return false }
func (db *noopTraceDB) MinSpanCount() int                 { return 0 }
func (db *noopTraceDB) TestConnection(ctx context.Context) ConnectionTestResult {
	return ConnectionTestResult{}
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

func (f *traceDBFactory) New(ds model.DataStore) (TraceDB, error) {
	switch ds.Type {
	case openapi.JAEGER:
		return newJaegerDB(ds.Values.Jaeger)
	case openapi.TEMPO:
		return newTempoDB(ds.Values.Tempo)
	case openapi.OPEN_SEARCH:
		return newOpenSearchDB(ds.Values.OpenSearch)
	case openapi.SIGNAL_FX:
		return newSignalFXDB(ds.Values.SignalFx)
	case openapi.OTLP:
		return newCollectorDB(f.repo)
	}

	return &noopTraceDB{}, nil
}

type realTraceDB struct{}

func (db *realTraceDB) ShouldRetry() bool { return true }
func (db *realTraceDB) MinSpanCount() int { return 1 }
