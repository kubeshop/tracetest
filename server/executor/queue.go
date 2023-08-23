package executor

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/alitto/pond"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/http/middleware"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/testsuite"
	"go.opentelemetry.io/otel/propagation"
)

const (
	QueueWorkerCount      = 10
	QueueWorkerBufferSize = QueueWorkerCount * 100 // 100 jobs per worker

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

	TestSuite    testsuite.TestSuite
	TestSuiteRun testsuite.TestSuiteRun

	Test test.Test
	Run  test.Run

	PollingProfile pollingprofile.PollingProfile
	DataStore      datastore.DataStore
}

type jsonJob struct {
	Headers          *headers `json:"headers"`
	TestSuiteID      string   `json:"test_suite_id"`
	TestSuiteRunID   int      `json:"test_suite_run_id"`
	TestID           string   `json:"test_id"`
	RunID            int      `json:"run_id"`
	PollingProfileID string   `json:"polling_profile_id"`
	DataStoreID      string   `json:"data_store_id"`
}

func (job Job) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonJob{
		Headers:          job.Headers,
		TestSuiteID:      job.TestSuite.ID.String(),
		TestSuiteRunID:   job.TestSuiteRun.ID,
		TestID:           job.Test.ID.String(),
		RunID:            job.Run.ID,
		PollingProfileID: job.PollingProfile.ID.String(),
		DataStoreID:      job.DataStore.ID.String(),
	})
}

func (job *Job) UnmarshalJSON(data []byte) error {
	var jj jsonJob
	err := json.Unmarshal(data, &jj)
	if err != nil {
		return err
	}

	job.Headers = jj.Headers
	job.TestSuite.ID = id.ID(jj.TestSuiteID)
	job.TestSuiteRun.ID = jj.TestSuiteRunID
	job.Test.ID = id.ID(jj.TestID)
	job.Run.ID = jj.RunID
	job.PollingProfile.ID = id.ID(jj.PollingProfileID)
	job.DataStore.ID = id.ID(jj.DataStoreID)

	return nil
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

type testSuiteGetter interface {
	GetAugmented(context.Context, id.ID) (testsuite.TestSuite, error)
}

type testSuiteRunGetter interface {
	GetTestSuiteRun(_ context.Context, transactionID id.ID, runID int) (testsuite.TestSuiteRun, error)
}

type subscriptor interface {
	Subscribe(string, subscription.Subscriber)
}

type Listener interface {
	Listen(Job)
}

type QueueDriver interface {
	Enqueue(Job)
	SetListener(Listener)
}

type QueueBuilder struct {
	cancelRunHandlerFn func(ctx context.Context, run test.Run) error
	subscriptor        subscriptor

	runs  testRunGetter
	tests testGetter

	testSuiteRuns testSuiteRunGetter
	testSuites    testSuiteGetter

	pollingProfiles pollingProfileGetter
	dataStores      dataStoreGetter
}

func NewQueueBuilder() *QueueBuilder {
	return &QueueBuilder{}
}

func (qb *QueueBuilder) WithCancelRunHandlerFn(fn func(ctx context.Context, run test.Run) error) *QueueBuilder {
	qb.cancelRunHandlerFn = fn
	return qb
}

func (qb *QueueBuilder) WithSubscriptor(subscriptor subscriptor) *QueueBuilder {
	qb.subscriptor = subscriptor
	return qb
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

func (qb *QueueBuilder) WithTestSuiteGetter(suites testSuiteGetter) *QueueBuilder {
	qb.testSuites = suites
	return qb
}

func (qb *QueueBuilder) WithTestSuiteRunGetter(suiteRuns testSuiteRunGetter) *QueueBuilder {
	qb.testSuiteRuns = suiteRuns
	return qb
}

func (qb *QueueBuilder) Build(driver QueueDriver, itemProcessor QueueItemProcessor) *Queue {
	queue := &Queue{
		cancelRunHandlerFn: qb.cancelRunHandlerFn,
		subscriptor:        qb.subscriptor,

		runs:  qb.runs,
		tests: qb.tests,

		testSuiteRuns: qb.testSuiteRuns,
		testSuites:    qb.testSuites,

		pollingProfiles: qb.pollingProfiles,
		dataStores:      qb.dataStores,

		driver:        driver,
		itemProcessor: itemProcessor,
		workerPool:    pond.New(QueueWorkerCount, QueueWorkerBufferSize),
	}

	driver.SetListener(queue)

	return queue
}

type Queue struct {
	cancelRunHandlerFn func(ctx context.Context, run test.Run) error
	subscriptor        subscriptor

	runs  testRunGetter
	tests testGetter

	testSuiteRuns testSuiteRunGetter
	testSuites    testSuiteGetter

	pollingProfiles pollingProfileGetter
	dataStores      dataStoreGetter

	itemProcessor QueueItemProcessor
	driver        QueueDriver
	workerPool    *pond.WorkerPool
}

func (q *Queue) SetDriver(driver QueueDriver) {
	q.driver = driver
	driver.SetListener(q)
}

func propagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}, tenantPropagator{})
}

func (q Queue) Enqueue(ctx context.Context, job Job) {
	select {
	default:
	case <-ctx.Done():
		return
	}

	if job.Headers == nil {
		job.Headers = &headers{}
	}
	propagator().Inject(ctx, propagation.MapCarrier(*job.Headers))

	newJob := Job{
		Headers: job.Headers,

		Test: test.Test{ID: job.Test.ID},
		Run:  test.Run{ID: job.Run.ID},

		TestSuite:    testsuite.TestSuite{ID: job.TestSuite.ID},
		TestSuiteRun: testsuite.TestSuiteRun{ID: job.TestSuiteRun.ID},

		PollingProfile: pollingprofile.PollingProfile{ID: job.PollingProfile.ID},
		DataStore:      datastore.DataStore{ID: job.DataStore.ID},
	}

	q.driver.Enqueue(newJob)
}

func (q Queue) Listen(job Job) {
	// this is called when a new job is put in the queue and we need to process it
	ctx := propagator().Extract(context.Background(), propagation.MapCarrier(*job.Headers))

	ctx, cancelCtx := context.WithCancel(ctx)
	q.listenForStopRequests(context.Background(), cancelCtx, job)

	newJob := Job{
		Headers: job.Headers,
	}
	newJob.Test = q.resolveTest(ctx, job)
	// todo: revert when using actual queues
	newJob.Run = q.resolveTestRun(ctx, job)
	// todo: change the otlp server to have its own table
	// newJob.Run = job.Run

	newJob.TestSuite = q.resolveTestSuite(ctx, job)
	newJob.TestSuiteRun = q.resolveTestSuiteRun(ctx, job)

	newJob.PollingProfile = q.resolvePollingProfile(ctx, job)
	newJob.DataStore = q.resolveDataStore(ctx, job)

	// Process the item with cancellation monitoring
	select {
	default:
	case <-ctx.Done():
		return
	}

	q.workerPool.Submit(func() {
		q.itemProcessor.ProcessItem(ctx, newJob)
	})
}

func (q *Queue) Stop() {
	q.workerPool.StopAndWait()
}

type StopRequest struct {
	TestID id.ID
	RunID  int
}

func (sr StopRequest) ResourceID() string {
	runID := (test.Run{ID: sr.RunID, TestID: sr.TestID}).ResourceID()
	return runID + "/stop"
}

func (q Queue) listenForStopRequests(ctx context.Context, cancelCtx context.CancelFunc, job Job) {
	if q.subscriptor == nil {
		return
	}

	sfn := subscription.NewSubscriberFunction(func(m subscription.Message) error {
		cancelCtx()
		stopRequest, ok := m.Content.(StopRequest)
		if !ok {
			return nil
		}

		run, err := q.runs.GetRun(ctx, stopRequest.TestID, stopRequest.RunID)
		if err != nil {
			return fmt.Errorf("failed to get run %d for test %s: %w", stopRequest.RunID, stopRequest.TestID, err)
		}

		if run.State == test.RunStateStopped {
			return nil
		}

		return q.cancelRunHandlerFn(ctx, run)

	})

	q.subscriptor.Subscribe((StopRequest{job.Test.ID, job.Run.ID}).ResourceID(), sfn)
}

func (q Queue) resolveTestSuite(ctx context.Context, job Job) testsuite.TestSuite {
	if q.testSuites == nil {
		return testsuite.TestSuite{}
	}

	tran, err := q.testSuites.GetAugmented(ctx, job.TestSuite.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return testsuite.TestSuite{}
	}
	if err != nil {
		panic(err)
	}

	return tran
}
func (q Queue) resolveTestSuiteRun(ctx context.Context, job Job) testsuite.TestSuiteRun {
	if q.testSuiteRuns == nil {
		return testsuite.TestSuiteRun{}
	}

	tranRun, err := q.testSuiteRuns.GetTestSuiteRun(ctx, job.TestSuite.ID, job.TestSuiteRun.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return testsuite.TestSuiteRun{}
	}
	if err != nil {
		panic(err)
	}

	return tranRun
}

func (q Queue) resolveTest(ctx context.Context, job Job) test.Test {
	if q.tests == nil {
		return test.Test{}
	}

	t, err := q.tests.GetAugmented(ctx, job.Test.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return test.Test{}
	}
	if err != nil {
		panic(err)
	}

	return t
}

func (q Queue) resolveTestRun(ctx context.Context, job Job) test.Run {
	if q.runs == nil {
		return test.Run{}
	}

	run, err := q.runs.GetRun(ctx, job.Test.ID, job.Run.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return test.Run{}
	}
	if err != nil {
		panic(err)
	}

	return run
}

func (q Queue) resolvePollingProfile(ctx context.Context, job Job) pollingprofile.PollingProfile {
	if q.pollingProfiles == nil {
		return pollingprofile.PollingProfile{}
	}

	profile, err := q.pollingProfiles.Get(ctx, job.PollingProfile.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return pollingprofile.PollingProfile{}
	}
	if err != nil {
		panic(err)
	}

	return profile
}

func (q Queue) resolveDataStore(ctx context.Context, job Job) datastore.DataStore {
	if q.dataStores == nil {
		return datastore.DataStore{}
	}

	ds, err := q.dataStores.Get(ctx, job.DataStore.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return datastore.DataStore{}
	}
	if err != nil {
		panic(err)
	}

	return ds
}

type tenantPropagator struct{}

var _ propagation.TextMapPropagator = tenantPropagator{}

func (b tenantPropagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	tenantID := middleware.TenantIDFromContext(ctx)
	if tenantID != "" {
		carrier.Set(string(middleware.TenantIDKey), tenantID)
	}
}

// Extract returns a copy of parent with the baggage from the carrier added.
func (b tenantPropagator) Extract(parent context.Context, carrier propagation.TextMapCarrier) context.Context {
	tenantID := carrier.Get(string(middleware.TenantIDKey))
	if tenantID == "" {
		return parent
	}

	return context.WithValue(parent, middleware.TenantIDKey, tenantID)
}

// Fields returns the keys who's values are set with Inject.
func (b tenantPropagator) Fields() []string {
	return []string{string(middleware.TenantIDKey)}
}
