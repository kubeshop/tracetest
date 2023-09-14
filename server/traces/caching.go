package traces

import "sync"

type spanCache map[string]*Span

var cache spanCache = nil
var cacheMutex sync.Mutex

func getCache() spanCache {
	if cache == nil {
		cache = spanCache{}
	}

	return cache
}

func resetCache() {
	cache = nil
}

func (c spanCache) Get(key string) (*Span, bool) {
	value, ok := c[key]
	return value, ok
}

func (c spanCache) Set(key string, value *Span) {
	c[key] = value
}
