package app

import (
	"context"
	"reflect"

	"github.com/kubeshop/tracetest/server/executor"
)

type Pipeline struct {
	steps  []PipelineStep    // N
	queues []*executor.Queue // N + 1
}

type workerDriver interface {
	executor.QueueDriver
	Start()
	Stop()
}

type queueBuilder interface {
	Build(executor.QueueDriver, executor.QueueItemProcessor) *executor.Queue
}

type pipelineStep interface {
	executor.QueueItemProcessor
	SetOutputQueue(executor.Enqueuer)
}

type InputQueueSetter interface {
	SetInputQueue(*executor.Queue)
}

type PipelineStep struct {
	driver    workerDriver
	processor pipelineStep
}

func NewPipeline(builder queueBuilder, steps ...PipelineStep) *Pipeline {
	pipeline := &Pipeline{
		steps:  steps,
		queues: make([]*executor.Queue, 0, len(steps)),
	}

	// setup an input queue for each pipeline step
	for _, step := range steps {
		pipeline.queues = append(pipeline.queues, builder.Build(step.driver, step.processor))
	}

	// set the output queue for each step. the ouput queue of a processor (N) is the input queue of the next one (N+1)
	// the last step has no output queue.
	for i, step := range steps {
		if i == len(pipeline.queues)-1 {
			break
		}

		step.processor.SetOutputQueue(pipeline.queues[i+1])

		// a processor might need to have a reference to its input queue, to requeue items for example.
		// This can be done if it implements the `InputQueueSetter` interace
		if setter, ok := reflect.ValueOf(step.processor).Interface().(InputQueueSetter); ok {
			setter.SetInputQueue(pipeline.queues[i])
		}
	}

	return pipeline
}

func (p *Pipeline) Begin(ctx context.Context, job executor.Job) {
	p.queues[0].Enqueue(ctx, job)
}

func (p *Pipeline) Start() {
	for _, step := range p.steps {
		step.driver.Start()
	}
}

func (p *Pipeline) Stop() {
	for _, step := range p.steps {
		step.driver.Stop()
	}
}
