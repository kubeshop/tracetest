package workers

import (
	"context"
	"sync"

	gocache "github.com/Code-Hex/go-generics-cache"
)

type StoppableProcessRunner func(parentCtx context.Context, testID string, runID int32, worker func(context.Context), stopedCallback func(cause string))

func NewProcessStopper() processStopper {
	return processStopper{
		cancelMap: newCancelCauseFuncMap(),
	}
}

type processStopper struct {
	cancelMap *cancelCauseFuncMap
}

func (p processStopper) CancelMap() *cancelCauseFuncMap {
	return p.cancelMap
}

func (p processStopper) RunStoppableProcess(parentCtx context.Context, testID string, runID int32, worker func(context.Context), stopedCallback func(cause string)) {
	done := make(chan bool)

	// create a subcontext for the worker so when canceled it doesn't affect the parent context
	subcontext, cancelSubctx := context.WithCancelCause(parentCtx)
	defer cancelSubctx(nil)

	cacheKey := key(testID, runID)
	p.cancelMap.Set(cacheKey, cancelSubctx)
	defer p.cancelMap.Del(cacheKey)

	go func() {
		worker(subcontext)
		done <- true
	}()

	select {
	case <-done:
		// trigger finished successfully
		break
	case <-subcontext.Done():
		cause := "cancelled"
		if err := context.Cause(subcontext); err != nil {
			cause = err.Error()
		}
		stopedCallback(cause)
		break
	}
}

func key(testID string, runID int32) string {
	return testID + string(runID)
}

type cancelCauseFuncMap struct {
	mutex       sync.Mutex
	internalMap *gocache.Cache[string, context.CancelCauseFunc]
}

// Get implements TraceCache.
func (c *cancelCauseFuncMap) Get(key string) (context.CancelCauseFunc, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.internalMap.Get(key)
}

// Append implements TraceCache.
func (c *cancelCauseFuncMap) Set(key string, cancelFn context.CancelCauseFunc) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.internalMap.Set(key, cancelFn)
}

func (c *cancelCauseFuncMap) Del(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.internalMap.Delete(key)
}

func newCancelCauseFuncMap() *cancelCauseFuncMap {
	return &cancelCauseFuncMap{
		internalMap: gocache.New[string, context.CancelCauseFunc](),
	}
}
