package pipeline

import (
	"context"
	"fmt"
	"log"
	"reflect"
)

type loggerFn func(string, ...any)

func newLoggerFn(name string) loggerFn {
	return func(format string, params ...any) {
		log.Printf("[%s] %s", name, fmt.Sprintf(format, params...))
	}
}

type Pipeline[T any] struct {
	steps  []Step[T]   // N
	queues []*Queue[T] // N + 1
}

type workerDriver[T any] interface {
	QueueDriver[T]
	Start()
	Stop()
}

type queueConfigurer[T any] interface {
	Configure(*Queue[T])
}

type StepProcessor[T any] interface {
	QueueItemProcessor[T]
	SetOutputQueue(Enqueuer[T])
}

type InputQueueSetter[T any] interface {
	SetInputQueue(Enqueuer[T])
}

type Step[T any] struct {
	Driver           workerDriver[T]
	Processor        StepProcessor[T]
	InputQueueOffset int
}

func New[T any](cfg queueConfigurer[T], steps ...Step[T]) *Pipeline[T] {
	pipeline := &Pipeline[T]{
		steps:  steps,
		queues: make([]*Queue[T], 0, len(steps)),
	}

	// setup an input queue for each pipeline step
	for _, step := range steps {
		q := NewQueue[T](step.Driver, step.Processor)
		cfg.Configure(q)
		pipeline.queues = append(pipeline.queues, q)
	}

	// set the output queue for each step. the ouput queue of a processor (N) is the input queue of the next one (N+1)
	// the last step has no output queue.
	for i, step := range steps {
		if i == len(pipeline.queues)-1 {
			break
		}

		step.Processor.SetOutputQueue(pipeline.queues[i+1])

		// a processor might need to have a reference to its input queue, to requeue items for example.
		// This can be done if it implements the `InputQueueSetter` interace
		if setter, ok := reflect.ValueOf(step.Processor).Interface().(InputQueueSetter[T]); ok {
			setter.SetInputQueue(pipeline.queues[i+step.InputQueueOffset])
		}
	}

	return pipeline
}

func (p *Pipeline[T]) GetQueueForStep(i int) *Queue[T] {
	if i < 0 || i >= len(p.queues) {
		return nil
	}

	return p.queues[i]
}

func (p *Pipeline[T]) Begin(ctx context.Context, item T) {
	if len(p.queues) < 1 {
		// this is a panic instead of an error because this could only happen during development
		panic(fmt.Errorf("pipeline has no input queue"))
	}
	p.queues[0].Enqueue(ctx, item)
}

func (p *Pipeline[T]) Start() {
	for _, step := range p.steps {
		step.Driver.Start()
	}
}

func (p *Pipeline[T]) Stop() {
	for _, step := range p.steps {
		step.Driver.Stop()
	}
}
