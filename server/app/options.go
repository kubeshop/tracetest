package app

import (
	"github.com/kubeshop/tracetest/testdb"
	"github.com/kubeshop/tracetest/tracedb"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Option func(a *App) error

func WithDB(db testdb.Repository) Option {
	return func(a *App) error {
		a.db = db
		return nil
	}
}

func WithTraceDB(traceDB tracedb.TraceDB) Option {
	return func(a *App) error {
		a.traceDB = traceDB
		return nil
	}
}

func WithTracerProvider(tracerProvider *trace.TracerProvider) Option {
	return func(a *App) error {
		a.tracerProvider = tracerProvider
		return nil
	}
}
