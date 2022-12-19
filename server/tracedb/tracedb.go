package tracedb

import (
	"context"
	"errors"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/model"
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
	GetTraceByID(ctx context.Context, traceID string) (model.Trace, error)
	TestConnection(ctx context.Context) ConnectionTestResult
	Close() error
}

type noopTraceDB struct{}

func (db *noopTraceDB) GetTraceByID(ctx context.Context, traceID string) (t model.Trace, err error) {
	return model.Trace{}, nil
}

func (db *noopTraceDB) Connect(ctx context.Context) error {
	return nil
}

func (db *noopTraceDB) Close() error { return nil }

func New(c config.Config, repository model.RunRepository) (db TraceDB, err error) {
	selectedDataStore, err := c.DataStore()
	if err != nil {
		return nil, err
	}

	if selectedDataStore == nil {
		return &noopTraceDB{}, nil
	}

	return NewFromDataStoreConfig(selectedDataStore, repository)
}

func (db *noopTraceDB) TestConnection(ctx context.Context) ConnectionTestResult {
	return ConnectionTestResult{}
}

func NewFromDataStoreConfig(c *config.TracingBackendDataStoreConfig, repository model.RunRepository) (db TraceDB, err error) {
	err = config.ErrInvalidTraceDBProvider

	switch {
	case c.Type == JAEGER_BACKEND:
		db, err = newJaegerDB(&c.Jaeger)
	case c.Type == TEMPO_BACKEND:
		db, err = newTempoDB(&c.Tempo)
	case c.Type == OPENSEARCH_BACKEND:
		db, err = newOpenSearchDB(c.OpenSearch)
	case c.Type == SIGNALFX:
		db, err = newSignalFXDB(c.SignalFX)
	case c.Type == OTLP:
		db, err = newCollectorDB(repository)
	}

	return
}
