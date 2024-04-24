package poller

import (
	"fmt"
	"time"

	gocache "github.com/Code-Hex/go-generics-cache"
)

type SentSpansCache struct {
	cache *gocache.Cache[string, *traceSpans]
}

type traceSpans struct {
	LastSpanTime time.Time
	Spans        map[string]bool
}

func newTraceSpans() *traceSpans {
	return &traceSpans{
		LastSpanTime: time.Time{},
		Spans:        make(map[string]bool),
	}
}

func NewSentSpansCache() *SentSpansCache {
	cache := &SentSpansCache{
		cache: gocache.New[string, *traceSpans](),
	}

	go cache.startCleanupWorker()

	return cache
}

func (c *SentSpansCache) Set(traceID, spanID string) {
	trace, found := c.cache.Get(traceID)
	if !found {
		trace = newTraceSpans()
	}

	trace.Spans[spanID] = true
	trace.LastSpanTime = time.Now()
	c.cache.Set(traceID, trace)
}

func (c *SentSpansCache) Get(traceID, spanID string) bool {
	m, found := c.cache.Get(traceID)
	if !found {
		return false
	}

	_, found = m.Spans[spanID]
	return found
}

func (c *SentSpansCache) startCleanupWorker() {
	ticker := time.NewTicker(1 * time.Minute)
	for {
		<-ticker.C
		now := time.Now()
		numCleanedupEntries := 0
		for _, key := range c.cache.Keys() {
			trace, _ := c.cache.Get(key)
			if trace.LastSpanTime.Before(now) {
				c.cache.Delete(key)
				numCleanedupEntries++
			}
		}

		if numCleanedupEntries > 0 {
			fmt.Printf("%d cache entries cleaned up\n", numCleanedupEntries)
		}
	}
}
