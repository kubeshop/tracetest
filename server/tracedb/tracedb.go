package tracedb

import (
	"context"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/traces"
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
	GetTraceByID(ctx context.Context, traceID string) (traces.Trace, error)
	TestConnection(ctx context.Context) ConnectionTestResult
	Close() error
}

var ErrInvalidTraceDBProvider = fmt.Errorf("invalid traceDB provider: available options are (jaeger, tempo)")

func New(c config.Config, repository model.RunRepository) (db TraceDB, err error) {
	selectedDataStore, err := c.DataStore()
	if err != nil {
		return nil, ErrInvalidTraceDBProvider
	}

	return NewFromDataStoreConfig(selectedDataStore, repository)
}

func NewFromDataStoreConfig(c *config.TracingBackendDataStoreConfig, repository model.RunRepository) (db TraceDB, err error) {
	err = ErrInvalidTraceDBProvider

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
