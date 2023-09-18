package testconnection

import (
	"context"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/resourcemanager"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/trace"
)

func NewDsTestConnectionRequest(
	dsTestListener *Listener,
	tracer trace.Tracer,
	newTraceDBFn tracedb.FactoryFunc,
	dsRepo resourcemanager.Current[datastore.DataStore],
) *dsTestConnectionRequest {
	return &dsTestConnectionRequest{
		dsTestListener: dsTestListener,
		tracer:         tracer,
		newTraceDBFn:   newTraceDBFn,
		dsRepo:         dsRepo,
	}
}

type dsTestConnectionRequest struct {
	dsTestListener *Listener
	outputQueue    pipeline.Enqueuer[Job]
	tracer         trace.Tracer
	newTraceDBFn   tracedb.FactoryFunc
	dsRepo         resourcemanager.Current[datastore.DataStore]
}

func (w *dsTestConnectionRequest) SetOutputQueue(queue pipeline.Enqueuer[Job]) {
	w.outputQueue = queue
}

func (w *dsTestConnectionRequest) ProcessItem(ctx context.Context, job Job) {
	ctx, pollingSpan := w.tracer.Start(ctx, "triggerResolverWorker.ProcessItem")
	defer pollingSpan.End()

	traceDB, err := getTraceDB(ctx, job.DataStore, w.newTraceDBFn)

	if err != nil {
		handleError(err, pollingSpan)
		return
	}

	if testableTraceDB, ok := traceDB.(tracedb.TestableTraceDB); ok {
		connectionResult := testableTraceDB.TestConnection(ctx)

		job.TestResult = connectionResult
		w.dsTestListener.Notify(job)
	}
}
