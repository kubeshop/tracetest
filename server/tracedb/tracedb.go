package tracedb

import (
	"context"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/server/config"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

var ErrTraceNotFound = errors.New("trace not found")

const (
	JAEGER_BACKEND string = "jaeger"
	TEMPO_BACKEND  string = "tempo"
)

type TraceDB interface {
	GetTraceByID(ctx context.Context, traceID string) (*v1.TracesData, error)
	Close() error
}

var ErrInvalidTraceDBProvider = fmt.Errorf("invalid traceDB provider")

func New(c config.Config) (db TraceDB, err error) {
	err = ErrInvalidTraceDBProvider
	switch {
	case c.TracingBackend.DataStore.Type == JAEGER_BACKEND:
		db, err = newJaegerDB(&c.TracingBackend.DataStore.Jaeger)
	case c.TracingBackend.DataStore.Type == TEMPO_BACKEND:
		db, err = newTempoDB(&c.TracingBackend.DataStore.Tempo)
	}

	return
}
