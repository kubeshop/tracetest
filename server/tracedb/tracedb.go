package tracedb

import (
	"context"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
)

var ErrTraceNotFound = errors.New("trace not found")

const (
	JAEGER_BACKEND    string = "jaeger"
	TEMPO_BACKEND     string = "tempo"
	LIGHTSTEP_BACKEND string = "lightstep"
)

type TraceDB interface {
	SendSpan(context.Context, trace.Span) error
	GetTraceByIdentification(context.Context, traces.TraceIdentification) (traces.Trace, error)
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
	case selectedDataStore.Type == LIGHTSTEP_BACKEND:
		db, err = newLightstepDB(selectedDataStore.Lightstep)
	}

	return
}
