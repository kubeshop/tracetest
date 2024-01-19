package workers

import (
	"context"
	"sync"

	gocache "github.com/Code-Hex/go-generics-cache"
	"github.com/kubeshop/tracetest/server/executor"
)

type StoppableProcessRunner func(parentCtx context.Context, testID string, runID int32, worker func(context.Context) error) error

func NewProcessStopper() processStopper {
	return processStopper{
		cancelMap: NewCancelFuncMap(),
	}
}

type processStopper struct {
	cancelMap *cancelFuncMap
}

func (p processStopper) CancelMap() *cancelFuncMap {
	return p.cancelMap
}

func (p processStopper) RunStoppableProcess(parentCtx context.Context, testID string, runID int32, worker func(context.Context) error) error {
	done := make(chan bool)

	// create a subcontext for the worker so when canceled it doesn't affect the parent context
	subcontext, cancelSubctx := context.WithCancel(parentCtx)
	defer cancelSubctx()

	cacheKey := key(testID, runID)
	p.cancelMap.Set(cacheKey, cancelSubctx)
	defer p.cancelMap.Del(cacheKey)

	var err error
	go func() {
		err = worker(subcontext)
		done <- true
	}()

	select {
	case <-done:
		// trigger finished successfully
		break
	case <-subcontext.Done():
		// The context was cancelled.
		err = executor.ErrUserCancelled
		break
	}

	return err
}

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
