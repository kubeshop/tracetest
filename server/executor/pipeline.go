package executor

import (
	"context"
	"fmt"
	"reflect"
)

type Pipeline struct {
	steps  []PipelineStep // N
	queues []*Queue       // N + 1
}

type workerDriver interface {
	QueueDriver
	Start()
	Stop()
}

type queueBuilder interface {
	Build(QueueDriver, QueueItemProcessor) *Queue
}

type pipelineStep interface {
	QueueItemProcessor
	SetOutputQueue(Enqueuer)
}

type InputQueueSetter interface {
	SetInputQueue(Enqueuer)
}

type PipelineStep struct {
	Driver    workerDriver
	Processor pipelineStep
}

func NewPipeline(builder queueBuilder, steps ...PipelineStep) *Pipeline {
	pipeline := &Pipeline{
		steps:  steps,
		queues: make([]*Queue, 0, len(steps)),
	}

	// setup an input queue for each pipeline step
	for _, step := range steps {
		pipeline.queues = append(pipeline.queues, builder.Build(step.Driver, step.Processor))
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
		if setter, ok := reflect.ValueOf(step.Processor).Interface().(InputQueueSetter); ok {
			setter.SetInputQueue(pipeline.queues[i])
		}
	}

	return pipeline
}

func (p *Pipeline) GetQueueForStep(i int) *Queue {
	if i < 0 || i >= len(p.queues) {
		return nil
	}

	return p.queues[i]
}

func (p *Pipeline) Begin(ctx context.Context, job Job) {
	if len(p.queues) < 1 {
		// this is a panic instead of an error because this could only happen during development
		panic(fmt.Errorf("pipeline has no input queue"))
	}
	p.queues[0].Enqueue(ctx, job)
}

func (p *Pipeline) Start() {
	for _, step := range p.steps {
		step.Driver.Start()
	}
}

func (p *Pipeline) Stop() {
	for _, step := range p.steps {
		step.Driver.Stop()
	}
}
