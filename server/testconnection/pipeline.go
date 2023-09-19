package testconnection

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

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

func (p *DataStoreTestPipeline) NewJob(datastore datastore.DataStore) Job {
	return Job{
		ID:        uuid.New().String(),
		DataStore: datastore,
	}
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

func getTraceDB(ctx context.Context, ds datastore.DataStore, newTraceDBFn tracedb.FactoryFunc) (tracedb.TraceDB, error) {
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
