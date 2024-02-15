package pipeline

import (
	"context"
	"fmt"
)

func NewInMemoryQueueDriver[T any](name string) *InMemoryQueueDriver[T] {
	return &InMemoryQueueDriver[T]{
		log:   newLoggerFn(fmt.Sprintf("InMemoryQueueDriver - %s", name)),
		queue: make(chan queueMessage[T]),
		exit:  make(chan bool),
		name:  name,
	}
}

type queueMessage[T any] struct {
	ctx     context.Context
	message T
}

type InMemoryQueueDriver[T any] struct {
	log      loggerFn
	queue    chan queueMessage[T]
	exit     chan bool
	listener Listener[T]
	name     string
}

func (qd *InMemoryQueueDriver[T]) SetListener(l Listener[T]) {
	qd.listener = l
}

func (qd InMemoryQueueDriver[T]) Enqueue(ctx context.Context, item T) {
	qd.queue <- queueMessage[T]{ctx, item}
}

func (qd InMemoryQueueDriver[T]) Start() {
	go func() {
		qd.log("start")
		for {
			select {
			case <-qd.exit:
				qd.log("exit")
				return
			case job := <-qd.queue:
				qd.listener.Listen(job.ctx, job.message)
			}
		}
	}()
}

func (qd InMemoryQueueDriver[T]) Stop() {
	qd.log("stopping")
	qd.exit <- true
}
