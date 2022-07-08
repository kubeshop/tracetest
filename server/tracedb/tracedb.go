package tracedb

import (
	"context"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/traces"
)

var ErrTraceNotFound = errors.New("trace not found")

const (
	JAEGER_BACKEND string = "jaeger"
	TEMPO_BACKEND  string = "tempo"
)

type TraceDB interface {
	GetTraceByIdentification(ctx context.Context, traceIdentification traces.TraceIdentification) (traces.Trace, error)
	Close() error
}

var ErrInvalidTraceDBProvider = fmt.Errorf("invalid traceDB provider: available options are (jaeger, tempo)")

func New(c config.Config) (db TraceDB, err error) {
	selectedDataStore, err := c.DataStore()
	if err != nil {
		return nil, ErrInvalidTraceDBProvider
	}

	err = ErrInvalidTraceDBProvider

	switch {
	case selectedDataStore.Type == JAEGER_BACKEND:
		db, err = newJaegerDB(&selectedDataStore.Jaeger)
	case selectedDataStore.Type == TEMPO_BACKEND:
		db, err = newTempoDB(&selectedDataStore.Tempo)
	}

	return
}
