package executor

import (
	"context"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type WorkerMiddlewareBuilder interface {
	New(name string, next pipeline.StepProcessor[Job]) pipeline.StepProcessor[Job]
}

type metricWorkerMiddlewareBuilder struct {
	meter metric.Meter
}

type metricWorkerMiddleware struct {
	latencyHistogram metric.Int64Histogram
	next             pipeline.StepProcessor[Job]
}

func NewWorkerMetricMiddlewareBuilder(meter metric.Meter) WorkerMiddlewareBuilder {
	return &metricWorkerMiddlewareBuilder{
		meter: meter,
	}
}

func (b *metricWorkerMiddlewareBuilder) New(name string, next pipeline.StepProcessor[Job]) pipeline.StepProcessor[Job] {
	metricPrefix := fmt.Sprintf("tracetest.worker.%s", name)

	latencyHistogram, _ := b.meter.Int64Histogram(fmt.Sprintf("%s.latency", metricPrefix))

	return &metricWorkerMiddleware{
		latencyHistogram: latencyHistogram,
		next:             next,
	}
}

// ProcessItem implements pipeline.QueueItemProcessor.
func (m *metricWorkerMiddleware) ProcessItem(ctx context.Context, job Job) {
	attributeSet := attribute.NewSet(
		attribute.String("test_id", job.Test.ID.String()),
		attribute.Int("run_id", job.Run.ID),
		attribute.String("run_state", string(job.Run.State)),
	)

	start := time.Now()
	m.next.ProcessItem(ctx, job)
	latency := time.Since(start)

	m.latencyHistogram.Record(ctx, latency.Milliseconds(), metric.WithAttributeSet(attributeSet))
}

func (m *metricWorkerMiddleware) SetOutputQueue(enqueuer pipeline.Enqueuer[Job]) {
	m.next.SetOutputQueue(enqueuer)
}

func (m *metricWorkerMiddleware) SetInputQueue(queue pipeline.Enqueuer[Job]) {
	if inputQueueSetter, ok := m.next.(pipeline.InputQueueSetter[Job]); ok {
		inputQueueSetter.SetInputQueue(queue)
	}
}
