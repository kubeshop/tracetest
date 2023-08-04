package executor

import (
	"fmt"
	"log"
)

type loggerFn func(string, ...any)

func newLoggerFn(name string) loggerFn {
	return func(format string, params ...any) {
		log.Printf("[%s] %s", name, fmt.Sprintf(format, params...))
	}
}

func NewInMemoryQueueDriver(name string) *InMemoryQueueDriver {
	return &InMemoryQueueDriver{
		log:   newLoggerFn(fmt.Sprintf("InMemoryQueueDriver - %s", name)),
		queue: make(chan Job),
		exit:  make(chan bool),
		name:  name,
	}
}

type InMemoryQueueDriver struct {
	log   loggerFn
	queue chan Job
	exit  chan bool
	q     *Queue
	name  string
}

func (qd *InMemoryQueueDriver) SetQueue(q *Queue) {
	qd.q = q
}

func (qd InMemoryQueueDriver) Enqueue(job Job) {
	qd.queue <- job
}

func (qd InMemoryQueueDriver) Start() {
	for i := 0; i < QueueWorkerCount; i++ {
		go func() {
			qd.log("start")
			for {
				select {
				case <-qd.exit:
					qd.log("exit")
					return
				case job := <-qd.queue:
					qd.q.Listen(job)
				}
			}
		}()
	}
}

func (qd InMemoryQueueDriver) Stop() {
	qd.log("stopping")
	qd.exit <- true
}
