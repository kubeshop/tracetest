package tracedb

import (
	"context"
	"errors"

	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

var ErrTraceNotFound = errors.New("trace not found")

type TraceDB interface {
	GetTraceByID(ctx context.Context, traceID string) (*v1.TracesData, error)
	Close() error
}
