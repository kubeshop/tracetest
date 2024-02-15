package pipeline

import (
	"context"

	"github.com/alitto/pond"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type Enqueuer[T any] interface {
	Enqueue(context.Context, T)
}

type QueueItemProcessor[T any] interface {
	ProcessItem(context.Context, T)
}
type Listener[T any] interface {
	Listen(context.Context, T)
}

type QueueDriver[T any] interface {
	Enqueue(context.Context, T)
	SetListener(Listener[T])
}

type namedDriver interface {
	Name() string
}

type Queue[T any] struct {
	name             string
	driver           QueueDriver[T]
	itemProcessor    QueueItemProcessor[T]
	enqueueHistogram metric.Int64Histogram
	listenHistogram  metric.Int64Histogram

	EnqueuePreprocessorFn func(context.Context, T) T
	ListenPreprocessorFn  func(context.Context, T) (context.Context, T)

	workerPool *pond.WorkerPool
}

const (
	QueueWorkerCount      = 20
	QueueWorkerBufferSize = QueueWorkerCount * 1_000 // 1k jobs per worker
)

func NewQueue[T any](driver QueueDriver[T], itemProcessor QueueItemProcessor[T]) *Queue[T] {
	queue := &Queue[T]{
		itemProcessor: itemProcessor,
		workerPool:    pond.New(QueueWorkerCount, QueueWorkerBufferSize),
	}

	if namedDriver, ok := driver.(namedDriver); ok {
		queue.name = namedDriver.Name()
	}

	queue.SetDriver(driver)

	return queue
}

func (q *Queue[T]) InitializeMetrics(meter metric.Meter) {
	q.enqueueHistogram, _ = meter.Int64Histogram("messaging.enqueue")
	q.listenHistogram, _ = meter.Int64Histogram("messaging.listen")
}

func (q *Queue[T]) SetDriver(driver QueueDriver[T]) {
	q.driver = driver
	driver.SetListener(q)
}

func (q Queue[T]) Enqueue(ctx context.Context, item T) {
	select {
	default:
	case <-ctx.Done():
		return
	}

	// use a worker to enqueue the job in case the driver takes a bit to actually enqueue
	// this way we release the caller as soon as possible
	q.workerPool.Submit(func() {
		if q.EnqueuePreprocessorFn != nil {
			item = q.EnqueuePreprocessorFn(ctx, item)
		}

		q.enqueueHistogram.Record(ctx, 1, metric.WithAttributes(
			attribute.String("queue.name", q.name),
		))
		q.driver.Enqueue(ctx, item)
	})
}

func (q Queue[T]) Listen(ctx context.Context, item T) {
	if q.ListenPreprocessorFn != nil {
		ctx, item = q.ListenPreprocessorFn(ctx, item)
	}

	// Process the item with cancellation monitoring
	select {
	default:
	case <-ctx.Done():
		return
	}

	q.workerPool.Submit(func() {
		q.listenHistogram.Record(ctx, 1, metric.WithAttributes(
			attribute.String("queue.name", q.name),
		))
		q.itemProcessor.ProcessItem(ctx, item)
	})
}

func (q *Queue[T]) Stop() {
	q.workerPool.StopAndWait()
}
