package workers

import (
	"sync"

	gocache "github.com/Code-Hex/go-generics-cache"
)

func key(testID string, runID int32) string {
	return testID + string(runID)
}

type cancelChannelMap struct {
	mutex       sync.Mutex
	internalMap *gocache.Cache[string, chan bool]
}

// Get implements TraceCache.
func (c *cancelChannelMap) Get(key string) (chan bool, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.internalMap.Get(key)
}

// Append implements TraceCache.
func (c *cancelChannelMap) Set(key string, cancelFn chan bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.internalMap.Set(key, cancelFn)
}

func (c *cancelChannelMap) Del(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.internalMap.Delete(key)
}

func NewCancelChannelMap() *cancelChannelMap {
	return &cancelChannelMap{
		internalMap: gocache.New[string, chan bool](),
	}
}
