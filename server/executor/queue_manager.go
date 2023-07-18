package executor

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test"
)

type Job struct {
	ctxHeaders map[string]string
	Test       test.Test
	Run        test.Run
	// TestID                   id.ID
	// RunID                    int
	initialJobConfigurations initialJobConfigurations
}

type initialJobConfigurations struct {
	pollingProfileID id.ID
	dataStoreID      id.ID
}

type Enqueuer interface {
	Enqueue(Job)
}

type QueueItemProcessor interface {
	ProcessItem(context.Context, Job)
}

type Queue struct {
	runs          test.RunRepository
	tests         test.Repository
	itemProcessor QueueItemProcessor
	driver        queueDriver
}

func NewQueue(runs test.RunRepository, tests test.Repository, driver queueDriver, itemProcessor QueueItemProcessor) *Queue {
	q := &Queue{
		runs:          runs,
		tests:         tests,
		itemProcessor: itemProcessor,
		driver:        driver,
	}

	driver.SetQueue(q)

	return q

}

type queueDriver interface {
	Enqueue(Job)
	SetQueue(*Queue)
}

func (r Queue) Enqueue(job Job) {
	r.driver.Enqueue(Job{
		ctxHeaders: job.ctxHeaders,
		Test:       test.Test{ID: job.Test.ID},
		Run:        test.Run{ID: job.Run.ID},
	})
}

func (r Queue) Listen(job Job) {
	// this is called when a new job is put in the queue and we need to process it
	ctx := context.Background()
	// TODO - carry over headers

	test, err := r.tests.Get(ctx, job.Test.ID)
	if err != nil {
		panic(err)
	}

	run, err := r.runs.GetRun(ctx, test.ID, job.Run.ID)
	if err != nil {
		panic(err)
	}

	r.itemProcessor.ProcessItem(ctx, Job{
		ctxHeaders:               job.ctxHeaders,
		initialJobConfigurations: job.initialJobConfigurations,
		Test:                     test,
		Run:                      run,
	})
}

func NewInMemoryQueueDriver() *InMemoryQueueDriver {
	return &InMemoryQueueDriver{
		queue: make(chan Job),
		exit:  make(chan bool),
	}
}

type InMemoryQueueDriver struct {
	queue chan Job
	exit  chan bool
	q     *Queue
}

func (r *InMemoryQueueDriver) SetQueue(q *Queue) {
	r.q = q
}

func (r InMemoryQueueDriver) Enqueue(job Job) {
	r.queue <- job
}

func (r InMemoryQueueDriver) Start(workers int) {
	for i := 0; i < workers; i++ {
		go func() {
			fmt.Println("persistentRunner start goroutine")
			for {
				select {
				case <-r.exit:
					fmt.Println("persistentRunner exit goroutine")
					return
				case job := <-r.queue:
					r.q.Listen(job)
				}
			}
		}()
	}
}

func (r InMemoryQueueDriver) Stop() {
	fmt.Println("persistentRunner stopping")
	r.exit <- true
}
