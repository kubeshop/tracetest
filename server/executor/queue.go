package executor

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test"
)

const (
	jobCountHeader string = "X-Tracetest-Job-Count"
)

type headers map[string]string

func (h headers) Get(key string) string {
	return h[key]
}

func (h headers) GetInt(key string) int {
	if h[key] == "" {
		return 0
	}

	v, err := strconv.Atoi(h[key])
	if err != nil {
		panic(fmt.Errorf("cannot convert header %s to int: %w", key, err))
	}
	return v
}

func (h headers) Set(key, value string) {
	if h == nil {
		h = make(map[string]string)
	}
	h[key] = value
}

func (h headers) SetInt(key string, value int) {
	h.Set(key, fmt.Sprintf("%d", value))
}

type Job struct {
	ctxHeaders     headers
	Test           test.Test
	Run            test.Run
	PollingProfile pollingprofile.PollingProfile
	DataStore      datastore.DataStore
}

func (j Job) EnqueueCount() int {
	count := j.ctxHeaders.GetInt(jobCountHeader)
	if count == 0 {
		return 1
	}

	return count
}

func (j Job) increaseEnqueueCount() Job {
	j.ctxHeaders.SetInt(jobCountHeader, j.EnqueueCount()+1)

	return j
}

type Enqueuer interface {
	Enqueue(context.Context, Job)
}

type QueueItemProcessor interface {
	ProcessItem(context.Context, Job)
}

type pollingProfileGetter interface {
	Get(context.Context, id.ID) (pollingprofile.PollingProfile, error)
}

type dataStoreGetter interface {
	Get(context.Context, id.ID) (datastore.DataStore, error)
}

type testGetter interface {
	GetAugmented(context.Context, id.ID) (test.Test, error)
}

type testRunGetter interface {
	GetRun(_ context.Context, testID id.ID, runID int) (test.Run, error)
}

type Queue struct {
	runs            testRunGetter
	tests           testGetter
	pollingProfiles pollingProfileGetter
	dataStores      dataStoreGetter

	itemProcessor QueueItemProcessor
	driver        QueueDriver
}

type QueueBuilder struct {
	runs            testRunGetter
	tests           testGetter
	pollingProfiles pollingProfileGetter
	dataStores      dataStoreGetter
}

func NewQueueBuilder() *QueueBuilder {
	return &QueueBuilder{}
}

func (qb *QueueBuilder) WithRunGetter(runs testRunGetter) *QueueBuilder {
	qb.runs = runs
	return qb
}

func (qb *QueueBuilder) WithTestGetter(tests testGetter) *QueueBuilder {
	qb.tests = tests
	return qb
}

func (qb *QueueBuilder) WithPollingProfileGetter(pollingProfiles pollingProfileGetter) *QueueBuilder {
	qb.pollingProfiles = pollingProfiles
	return qb
}

func (qb *QueueBuilder) WithDataStoreGetter(dataStore dataStoreGetter) *QueueBuilder {
	qb.dataStores = dataStore
	return qb
}

func (qb *QueueBuilder) Build(driver QueueDriver, itemProcessor QueueItemProcessor) *Queue {
	queue := &Queue{
		runs:            qb.runs,
		tests:           qb.tests,
		pollingProfiles: qb.pollingProfiles,
		dataStores:      qb.dataStores,
		driver:          driver,
		itemProcessor:   itemProcessor,
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
	job = job.increaseEnqueueCount()

	r.driver.Enqueue(Job{
		ctxHeaders:     job.ctxHeaders,
		Test:           test.Test{ID: job.Test.ID},
		Run:            test.Run{ID: job.Run.ID},
		PollingProfile: pollingprofile.PollingProfile{ID: job.PollingProfile.ID},
		DataStore:      datastore.DataStore{ID: job.DataStore.ID},
	})
}

func (r Queue) Listen(job Job) {
	// this is called when a new job is put in the queue and we need to process it
	ctx := context.Background()
	// TODO - carry over headers

	test, err := r.tests.GetAugmented(ctx, job.Test.ID)
	if err != nil {
		panic(err)
	}

	run, err := r.runs.GetRun(ctx, test.ID, job.Run.ID)
	if err != nil {
		panic(err)
	}

	pp, err := r.pollingProfiles.Get(ctx, job.PollingProfile.ID)
	if err != nil {
		panic(err)
	}

	ds, err := r.dataStores.Get(ctx, job.DataStore.ID)
	if err != nil {
		panic(err)
	}

	r.itemProcessor.ProcessItem(ctx, Job{
		ctxHeaders:     job.ctxHeaders,
		Test:           test,
		Run:            run,
		PollingProfile: pp,
		DataStore:      ds,
	})
}

func NewInMemoryQueueDriver(name string) *InMemoryQueueDriver {
	return &InMemoryQueueDriver{
		queue: make(chan Job),
		exit:  make(chan bool),
		name:  name,
	}
}

type InMemoryQueueDriver struct {
	queue chan Job
	exit  chan bool
	q     *Queue
	name  string
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
			log.Printf("[InMemoryQueueDriver - %s] start", r.name)
			for {
				select {
				case <-r.exit:
					log.Printf("[InMemoryQueueDriver - %s] exit", r.name)
					return
				case job := <-r.queue:
					r.q.Listen(job)
				}
			}
		}()
	}
}

func (r InMemoryQueueDriver) Stop() {
	log.Printf("[InMemoryQueueDriver - %s] stopping", r.name)
	r.exit <- true
}
