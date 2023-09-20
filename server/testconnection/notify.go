package testconnection

import (
	"context"

	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"go.opentelemetry.io/otel/trace"
)

type dsTestConnectionNotify struct {
	dsTestListener *Listener
	tracer         trace.Tracer
}

func NewDsTestConnectionNotify(
	dsTestListener *Listener,
	tracer trace.Tracer,
) *dsTestConnectionNotify {
	return &dsTestConnectionNotify{
		dsTestListener: dsTestListener,
		tracer:         tracer,
	}
}

func (w *dsTestConnectionNotify) SetOutputQueue(queue pipeline.Enqueuer[Job]) {
	// noop
}

func (w *dsTestConnectionNotify) ProcessItem(ctx context.Context, job Job) {
	_, pollingSpan := w.tracer.Start(ctx, "dsTestConnectionNotify.ProcessItem")
	defer pollingSpan.End()

	w.dsTestListener.Notify(job)
}
