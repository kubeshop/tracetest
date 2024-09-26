package poller

import (
	"context"
	"slices"
	"sort"

	"github.com/kubeshop/tracetest/agent/collector"
	"github.com/kubeshop/tracetest/agent/tracedb"
	"github.com/kubeshop/tracetest/agent/tracedb/connection"
	"github.com/kubeshop/tracetest/server/pkg/id"
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

func (d *inmemoryDatastore) List(ctx context.Context, take int, skip int) ([]traces.TraceMetadata, error) {
	keys := d.cache.Keys()
	slices.Reverse(keys)

	list := make([]traces.Trace, 0)
	for _, key := range keys {
		spans, found := d.cache.Get(key)
		if !found {
			return nil, connection.ErrTraceNotFound
		}

		trace := traces.FromSpanList(spans)
		if trace.RootSpan.Name == traces.TemporaryRootSpanName {
			continue
		}

		list = append(list, trace)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].RootSpan.StartTime.After(list[j].RootSpan.StartTime)
	})

	list = paginate(list, skip, take)

	metadata := make([]traces.TraceMetadata, len(list))
	for i, trace := range list {
		metadata[i] = trace.ToTraceMetadata()
	}

	return metadata, nil
}

// GetTraceByID implements tracedb.TraceDB.
func (d *inmemoryDatastore) GetTraceByID(ctx context.Context, traceID string) (traces.Trace, error) {
	spans, found := d.cache.Get(traceID)
	if !found || !d.cache.Exists(traceID) {
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

func paginate(x []traces.Trace, skip int, size int) []traces.Trace {
	if skip > len(x) {
		skip = len(x)
	}

	end := skip + size
	if end > len(x) {
		end = len(x)
	}

	return x[skip:end]
}
