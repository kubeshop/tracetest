package collector

import (
	"sync"

	gocache "github.com/Code-Hex/go-generics-cache"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

type TraceCache interface {
	Get(string) ([]*v1.Span, bool)
	Set(string, []*v1.Span)
}

type traceCache struct {
	mutex         sync.Mutex
	internalCache *gocache.Cache[string, []*v1.Span]
}

// Get implements TraceCache.
func (c *traceCache) Get(traceID string) ([]*v1.Span, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.internalCache.Get(traceID)
}

// Set implements TraceCache.
func (c *traceCache) Set(traceID string, spans []*v1.Span) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	existingTraces, _ := c.internalCache.Get(traceID)
	spans = append(existingTraces, spans...)

	c.internalCache.Set(traceID, spans)
}

func NewTraceCache() TraceCache {
	return &traceCache{
		internalCache: gocache.New[string, []*v1.Span](),
	}
}
