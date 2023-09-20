package testconnection

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DsTestListener interface {
	Notify(executor.Job)
	Subscribe(jobID string, notifier NotifierFn)
	Unsubscribe(jobID string)
}

type DataStoreTestPipeline struct {
	*pipeline.Pipeline[executor.Job]
	dsTestListener DsTestListener
}

type Configurer[T any] struct{}

func (c *Configurer[Job]) Configure(_ *pipeline.Queue[Job]) {}

func NewDataStoreTestPipeline(
	pipeline *pipeline.Pipeline[executor.Job],
	listener DsTestListener,
) *DataStoreTestPipeline {
	return &DataStoreTestPipeline{
		Pipeline:       pipeline,
		dsTestListener: listener,
	}
}

func (p *DataStoreTestPipeline) NewJob(datastore datastore.DataStore) executor.Job {
	job := executor.NewJob()
	job.ID = uuid.New().String()
	job.MemoryDataStore = datastore
	job.Run.TraceID = id.NewRandGenerator().TraceID()

	return job
}

func (p *DataStoreTestPipeline) Run(ctx context.Context, job executor.Job) {
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
