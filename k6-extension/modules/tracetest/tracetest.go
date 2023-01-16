package tracetest

import (
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
	"github.com/xoscar/xk6-tracetest-tracing/models"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/output"
)

type Tracetest struct {
	Vu              modules.VU
	bufferLock      sync.Mutex
	buffer          []Job
	periodicFlusher *output.PeriodicFlusher
	logger          logrus.FieldLogger
}

func New() *Tracetest {
	logger := *logrus.New()
	tracetest := &Tracetest{
		buffer: []Job{},
		logger: logger.WithField("component", "xk6-tracetest-tracing"),
	}

	duration := 3 * time.Second
	periodicFlusher, _ := output.NewPeriodicFlusher(duration, tracetest.processQueue)

	tracetest.periodicFlusher = periodicFlusher

	return tracetest
}

func (t *Tracetest) Constructor(call goja.ConstructorCall) *goja.Object {
	rt := t.Vu.Runtime()
	isCliInstalled := t.getIsCliInstalled()

	if !isCliInstalled {
		panic("The tracetest CLI is not installed. Please install it before using this module")
	}

	return rt.ToValue(t).ToObject(rt)
}

func (t *Tracetest) RunTest(testID, traceID string) {
	t.queueJob(Job{
		traceID: traceID,
		testID:  testID,
		jobType: RunTestFromId,
	})
}

func (t *Tracetest) RunFromDefinition(testDefinition, traceID string) {
	t.queueJob(Job{
		traceID:    traceID,
		definition: testDefinition,
		jobType:    RunTestFromDefinition,
	})
}

func (t *Tracetest) SyncRunTest(testID, traceID string) (*models.TracetestRun, error) {
	return t.runFromId(testID, traceID)
}

func (t *Tracetest) SyncRunTestFromDefinition(testDefinition, traceID string) (*models.TracetestRun, error) {
	return t.runFromDefinition(testDefinition, traceID)
}
