package pipeline

import (
	"fmt"
)

func NewInMemoryQueueDriver[T any](name string) *InMemoryQueueDriver[T] {
	return &InMemoryQueueDriver[T]{
		log:   newLoggerFn(fmt.Sprintf("InMemoryQueueDriver - %s", name)),
		queue: make(chan T),
		exit:  make(chan bool),
		name:  name,
	}
}

type InMemoryQueueDriver[T any] struct {
	log      loggerFn
	queue    chan T
	exit     chan bool
	listener Listener[T]
	name     string
}

func (qd *InMemoryQueueDriver[T]) SetListener(l Listener[T]) {
	qd.listener = l
}

func (qd InMemoryQueueDriver[T]) Enqueue(item T) {
	qd.queue <- item
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
				qd.listener.Listen(job)
			}
		}
	}()
}

func (qd InMemoryQueueDriver[T]) Stop() {
	qd.log("stopping")
	qd.exit <- true
}
