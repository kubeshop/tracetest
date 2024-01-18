package workers

import (
	"context"
	"sync"

	gocache "github.com/Code-Hex/go-generics-cache"
)

func key(testID string, runID int32) string {
	return testID + string(runID)
}

type cancelFuncMap struct {
	mutex       sync.Mutex
	internalMap *gocache.Cache[string, context.CancelFunc]
}

// Get implements TraceCache.
func (c *cancelFuncMap) Get(key string) (context.CancelFunc, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.internalMap.Get(key)
}

// Append implements TraceCache.
func (c *cancelFuncMap) Set(key string, cancelFn context.CancelFunc) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.internalMap.Set(key, cancelFn)
}

func (c *cancelFuncMap) Del(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.internalMap.Delete(key)
}

func NewCancelFuncMap() *cancelFuncMap {
	return &cancelFuncMap{
		internalMap: gocache.New[string, context.CancelFunc](),
	}
}
