package executor

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/transaction"
	"go.opentelemetry.io/otel/propagation"
)

const (
	JobCountHeader string = "X-Tracetest-Job-Count"
)

type headers map[string]string

func (h headers) Get(key string) string {
	return h[key]
}

func (h headers) GetInt(key string) int {
	if h.Get(key) == "" {
		return 0
	}

	v, err := strconv.Atoi(h[key])
	if err != nil {
		panic(fmt.Errorf("cannot convert header %s to int: %w", key, err))
	}
	return v
}

func (h headers) GetBool(key string) bool {
	if h.Get(key) == "" {
		return false
	}

	v, err := strconv.ParseBool(h[key])
	if err != nil {
		panic(fmt.Errorf("cannot convert header %s to bool: %w", key, err))
	}

	return v
}

func (h *headers) Set(key, value string) {
	// if h == nil {
	// 	nh := &headers{}
	// 	*h = *nh
	// }
	(*h)[key] = value
}

func (h headers) SetInt(key string, value int) {
	h.Set(key, fmt.Sprintf("%d", value))
}

func (h headers) SetBool(key string, value bool) {
	h.Set(key, fmt.Sprintf("%t", value))
}

type Job struct {
	Headers *headers

	Transaction    transaction.Transaction
	TransactionRun transaction.TransactionRun

	Test test.Test
	Run  test.Run

	PollingProfile pollingprofile.PollingProfile
	DataStore      datastore.DataStore
}

func NewJob() Job {
	return Job{
		Headers: &headers{},
	}
}

func (j Job) EnqueueCount() int {
	return j.Headers.GetInt(JobCountHeader)
}

func (j *Job) IncreaseEnqueueCount() {
	j.Headers.SetInt(JobCountHeader, j.EnqueueCount()+1)
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

type transactionGetter interface {
	GetAugmented(context.Context, id.ID) (transaction.Transaction, error)
}

type transactionRunGetter interface {
	GetTransactionRun(_ context.Context, transactionID id.ID, runID int) (transaction.TransactionRun, error)
}

type Queue struct {
	runs  testRunGetter
	tests testGetter

	transactionRuns transactionRunGetter
	transactions    transactionGetter

	pollingProfiles pollingProfileGetter
	dataStores      dataStoreGetter

	itemProcessor QueueItemProcessor
	driver        QueueDriver
}

type QueueBuilder struct {
	runs  testRunGetter
	tests testGetter

	transactionRuns transactionRunGetter
	transactions    transactionGetter

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

func (qb *QueueBuilder) WithTransactionGetter(transactions transactionGetter) *QueueBuilder {
	qb.transactions = transactions
	return qb
}

func (qb *QueueBuilder) WithTransactionRunGetter(transactionRuns transactionRunGetter) *QueueBuilder {
	qb.transactionRuns = transactionRuns
	return qb
}

func (qb *QueueBuilder) Build(driver QueueDriver, itemProcessor QueueItemProcessor) *Queue {
	queue := &Queue{
		runs:  qb.runs,
		tests: qb.tests,

		transactionRuns: qb.transactionRuns,
		transactions:    qb.transactions,

		pollingProfiles: qb.pollingProfiles,
		dataStores:      qb.dataStores,

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

func (q Queue) Enqueue(ctx context.Context, job Job) {
	if job.Headers == nil {
		job.Headers = &headers{}
	}
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	propagator.Inject(ctx, propagation.MapCarrier(*job.Headers))

	newJob := Job{
		Headers: job.Headers,

		Test: test.Test{ID: job.Test.ID},
		Run:  test.Run{ID: job.Run.ID},

		Transaction:    transaction.Transaction{ID: job.Transaction.ID},
		TransactionRun: transaction.TransactionRun{ID: job.TransactionRun.ID},

		PollingProfile: pollingprofile.PollingProfile{ID: job.PollingProfile.ID},
		DataStore:      datastore.DataStore{ID: job.DataStore.ID},
	}

	q.driver.Enqueue(newJob)
}

func (q Queue) Listen(job Job) {
	// this is called when a new job is put in the queue and we need to process it
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	ctx := propagator.Extract(context.Background(), propagation.MapCarrier(*job.Headers))

	newJob := Job{
		Headers: job.Headers,
	}
	newJob.Test = q.resolveTest(job)
	newJob.Run = q.resolveTestRun(job)

	newJob.Transaction = q.resolveTransaction(job)
	newJob.TransactionRun = q.resolveTransactionRun(job)

	newJob.PollingProfile = q.resolvePollingProfile(job)
	newJob.DataStore = q.resolveDataStore(job)

	q.itemProcessor.ProcessItem(ctx, newJob)
}

func (q Queue) resolveTransaction(job Job) transaction.Transaction {
	if q.transactions == nil {
		return transaction.Transaction{}
	}

	tran, err := q.transactions.GetAugmented(context.Background(), job.Transaction.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return transaction.Transaction{}
	}
	if err != nil {
		panic(err)
	}

	return tran
}
func (q Queue) resolveTransactionRun(job Job) transaction.TransactionRun {
	if q.transactionRuns == nil {
		return transaction.TransactionRun{}
	}

	tranRun, err := q.transactionRuns.GetTransactionRun(context.Background(), job.Transaction.ID, job.TransactionRun.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return transaction.TransactionRun{}
	}
	if err != nil {
		panic(err)
	}

	return tranRun
}

func (q Queue) resolveTest(job Job) test.Test {
	if q.tests == nil {
		return test.Test{}
	}

	t, err := q.tests.GetAugmented(context.Background(), job.Test.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return test.Test{}
	}
	if err != nil {
		panic(err)
	}

	return t
}

func (q Queue) resolveTestRun(job Job) test.Run {
	if q.runs == nil {
		return test.Run{}
	}

	run, err := q.runs.GetRun(context.Background(), job.Test.ID, job.Run.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return test.Run{}
	}
	if err != nil {
		panic(err)
	}

	return run
}

func (q Queue) resolvePollingProfile(job Job) pollingprofile.PollingProfile {
	if q.pollingProfiles == nil {
		return pollingprofile.PollingProfile{}
	}

	profile, err := q.pollingProfiles.Get(context.Background(), job.PollingProfile.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return pollingprofile.PollingProfile{}
	}
	if err != nil {
		panic(err)
	}

	return profile
}

func (q Queue) resolveDataStore(job Job) datastore.DataStore {
	if q.dataStores == nil {
		return datastore.DataStore{}
	}

	ds, err := q.dataStores.Get(context.Background(), job.DataStore.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return datastore.DataStore{}
	}
	if err != nil {
		panic(err)
	}

	return ds
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
