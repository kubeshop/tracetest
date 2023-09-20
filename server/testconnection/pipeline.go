package testconnection

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/http/middleware"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Job struct {
	Headers    map[string]string      `json:"headers"`
	ID         string                 `json:"id"`
	DataStore  datastore.DataStore    `json:"datastore"`
	TestResult model.ConnectionResult `json:"testResult"`
}

type DsTestListener interface {
	Notify(Job)
	Subscribe(jobID string, notifier NotifierFn)
	Unsubscribe(jobID string)
}

type DataStoreTestPipeline struct {
	*pipeline.Pipeline[Job]
	dsTestListener DsTestListener
}

type Configurer[T any] struct{}

func (c *Configurer[Job]) Configure(_ *pipeline.Queue[Job]) {}

func NewDataStoreTestPipeline(
	pipeline *pipeline.Pipeline[Job],
	listener DsTestListener,
) *DataStoreTestPipeline {
	return &DataStoreTestPipeline{
		Pipeline:       pipeline,
		dsTestListener: listener,
	}
}

func (p *DataStoreTestPipeline) NewJob(ctx context.Context, datastore datastore.DataStore) Job {
	job := Job{
		Headers: make(map[string]string),
	}
	job.ID = uuid.New().String()
	job.DataStore = datastore
	job.Headers[string(middleware.TenantIDKey)] = middleware.TenantIDFromContext(ctx)

	return job
}

func (p *DataStoreTestPipeline) Run(ctx context.Context, job Job) {
	p.Pipeline.Begin(ctx, job)
}

func (p *DataStoreTestPipeline) Subscribe(jobID string, notifier NotifierFn) {
	p.dsTestListener.Subscribe(jobID, notifier)
}

func (p *DataStoreTestPipeline) Unsubscribe(jobID string) {
	p.dsTestListener.Unsubscribe(jobID)
}

func getTraceDB(ds datastore.DataStore, newTraceDBFn tracedb.FactoryFunc) (tracedb.TraceDB, error) {
	tdb, err := newTraceDBFn(ds)
	if err != nil {
		return nil, fmt.Errorf(`cannot get tracedb from DataStore config with ID "%s": %w`, ds.ID, err)
	}

	return tdb, nil
}

func handleError(err error, span trace.Span) {
	span.RecordError(err)
	span.SetAttributes(attribute.String("tracetest.run.trace_poller.error", err.Error()))
}
