package testconnection

import (
	"context"

	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/tracedb"
	"go.opentelemetry.io/otel/trace"
)

type dsTestConnectionRequest struct {
	outputQueue  pipeline.Enqueuer[Job]
	tracer       trace.Tracer
	newTraceDBFn tracedb.FactoryFunc
	enabled      bool
}

func NewDsTestConnectionRequest(
	tracer trace.Tracer,
	newTraceDBFn tracedb.FactoryFunc,
	enabled bool,
) *dsTestConnectionRequest {
	return &dsTestConnectionRequest{
		tracer:       tracer,
		newTraceDBFn: newTraceDBFn,
		enabled:      enabled,
	}
}

func (w *dsTestConnectionRequest) SetOutputQueue(queue pipeline.Enqueuer[Job]) {
	w.outputQueue = queue
}

func (w *dsTestConnectionRequest) ProcessItem(ctx context.Context, job Job) {
	if !w.enabled {
		return
	}

	ctx, pollingSpan := w.tracer.Start(ctx, "dsTestConnectionRequest.ProcessItem")
	defer pollingSpan.End()

	traceDB, err := getTraceDB(job.DataStore, w.newTraceDBFn)

	if err != nil {
		handleError(err, pollingSpan)
		return
	}

	if testableTraceDB, ok := traceDB.(tracedb.TestableTraceDB); ok {
		connectionResult := testableTraceDB.TestConnection(ctx)

		job.TestResult = connectionResult
	}

	w.outputQueue.Enqueue(ctx, job)
}
