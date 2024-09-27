package collector

import (
	"slices"
	"sync"

	gocache "github.com/Code-Hex/go-generics-cache"
	"go.opentelemetry.io/otel/trace"
	v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

type TraceCache interface {
	Get(string) ([]*v1.Span, bool)
	Append(string, []*v1.Span)
	RemoveSpans(string, []string)
	Exists(string) bool
}

type traceCache struct {
	mutex         sync.Mutex
	internalCache *gocache.Cache[string, []*v1.Span]
	receivedSpans *gocache.Cache[string, int]
}

// Get implements TraceCache.
func (c *traceCache) Get(traceID string) ([]*v1.Span, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.internalCache.Get(traceID)
}

// Append implements TraceCache.
func (c *traceCache) Append(traceID string, spans []*v1.Span) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	currentNumberSpans, _ := c.receivedSpans.Get(traceID)
	currentNumberSpans += len(spans)

	existingTraces, _ := c.internalCache.Get(traceID)
	spans = append(existingTraces, spans...)

	c.internalCache.Set(traceID, spans)
	c.receivedSpans.Set(traceID, currentNumberSpans)
}

func (c *traceCache) RemoveSpans(traceID string, spanID []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	spans, found := c.internalCache.Get(traceID)
	if !found {
		return
	}

	newSpans := make([]*v1.Span, 0, len(spans))
	for _, span := range spans {
		currentSpanID := trace.SpanID(span.SpanId).String()
		if !slices.Contains(spanID, currentSpanID) {
			newSpans = append(newSpans, span)
		}
	}

	c.internalCache.Set(traceID, newSpans)
}

func (c *traceCache) Exists(traceID string) bool {
	numberSpans, _ := c.receivedSpans.Get(traceID)
	return numberSpans > 0
}

func NewTraceCache() TraceCache {
	return &traceCache{
		internalCache: gocache.New[string, []*v1.Span](),
		receivedSpans: gocache.New[string, int](),
	}
}
