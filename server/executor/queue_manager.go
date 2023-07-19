package executor

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test"
)

type Job struct {
	ctxHeaders               map[string]string
	Test                     test.Test
	Run                      test.Run
	InitialJobConfigurations InitialJobConfigurations
}

type InitialJobConfigurations struct {
	PollingProfileID id.ID
	DataStoreID      id.ID
}

type Enqueuer interface {
	Enqueue(context.Context, Job)
}

type QueueItemProcessor interface {
	ProcessItem(context.Context, Job)
}

type Queue struct {
	runs          test.RunRepository
	tests         test.Repository
	itemProcessor QueueItemProcessor
	driver        QueueDriver
}

type QueueBuilder struct {
	runs  test.RunRepository
	tests test.Repository
}

func NewQueueBuilder(runs test.RunRepository, tests test.Repository) *QueueBuilder {
	return &QueueBuilder{
		runs,
		tests,
	}
}

func (qb *QueueBuilder) Build(driver QueueDriver, itemProcessor QueueItemProcessor) *Queue {
	queue := &Queue{
		runs:          qb.runs,
		tests:         qb.tests,
		driver:        driver,
		itemProcessor: itemProcessor,
	}

	driver.SetQueue(queue)

	return queue
}

func (q *Queue) SetDriver(driver QueueDriver) {
	q.driver = driver
	driver.SetQueue(q)
}

type QueueDriver interface {
	Enqueue(Job)
	SetQueue(*Queue)
}

func (r Queue) Enqueue(ctx context.Context, job Job) {
	// TODO: carry context propagation
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
		InitialJobConfigurations: job.InitialJobConfigurations,
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

func (r InMemoryQueueDriver) Start() {
	for i := 0; i < 5; i++ {
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
