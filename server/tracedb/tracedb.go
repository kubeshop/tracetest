package tracedb

import (
	"context"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/config"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

var ErrTraceNotFound = errors.New("trace not found")

type TraceDB interface {
	GetTraceByID(ctx context.Context, traceID string) (*v1.TracesData, error)
	Close() error
}

var ErrInvalidTraceDBProvider = fmt.Errorf("invalid traceDB provider")

func New(c config.Config) (db TraceDB, err error) {
	err = ErrInvalidTraceDBProvider
	switch {
	case c.JaegerConnectionConfig != nil:
		db, err = newJaegerDB(c.JaegerConnectionConfig)
	case c.TempoConnectionConfig != nil:
		db, err = newTempoDB(c.TempoConnectionConfig)
	}

	return
}
