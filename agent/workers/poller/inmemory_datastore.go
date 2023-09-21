package poller

import (
	"context"

	"github.com/kubeshop/tracetest/agent/collector"
	"github.com/kubeshop/tracetest/agent/workers/datastores/connection"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/tracedb"
	"github.com/kubeshop/tracetest/server/traces"
	"go.opentelemetry.io/otel/trace"
)

func NewInMemoryDatastore(cache collector.TraceCache) tracedb.TraceDB {
	return &inmemoryDatastore{cache}
}

type inmemoryDatastore struct {
	cache collector.TraceCache
}

// Close implements tracedb.TraceDB.
func (d *inmemoryDatastore) Close() error {
	return nil
}

// Connect implements tracedb.TraceDB.
func (d *inmemoryDatastore) Connect(ctx context.Context) error {
	return nil
}

// GetEndpoints implements tracedb.TraceDB.
func (d *inmemoryDatastore) GetEndpoints() string {
	return ""
}

// GetTraceByID implements tracedb.TraceDB.
func (d *inmemoryDatastore) GetTraceByID(ctx context.Context, traceID string) (traces.Trace, error) {
	spans, found := d.cache.Get(traceID)
	if !found || len(spans) == 0 {
		return traces.Trace{}, connection.ErrTraceNotFound
	}

	return traces.FromSpanList(spans), nil
}

// GetTraceID implements tracedb.TraceDB.
func (d *inmemoryDatastore) GetTraceID() trace.TraceID {
	return id.NewRandGenerator().TraceID()
}

// Ready implements tracedb.TraceDB.
func (d *inmemoryDatastore) Ready() bool {
	return true
}

// ShouldRetry implements tracedb.TraceDB.
func (d *inmemoryDatastore) ShouldRetry() bool {
	return true
}
