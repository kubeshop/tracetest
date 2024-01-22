package executor

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/kubeshop/tracetest/server/datastore"
	"github.com/kubeshop/tracetest/server/executor/pollingprofile"
	"github.com/kubeshop/tracetest/server/http/middleware"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/testsuite"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
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
	TenantID     string

	Test test.Test
	Run  test.Run

	PollingProfile pollingprofile.PollingProfile
	DataStore      datastore.DataStore
}

type JsonJob struct {
	Headers          *headers `json:"headers"`
	TestSuiteID      string   `json:"test_suite_id"`
	TestSuiteRunID   int      `json:"test_suite_run_id"`
	TestID           string   `json:"test_id"`
	TestVersion      int      `json:"test_version"`
	RunID            int      `json:"run_id"`
	TraceID          string   `json:"trace_id"`
	PollingProfileID string   `json:"polling_profile_id"`
	DataStoreID      string   `json:"data_store_id"`
}

func (job Job) MarshalJSON() ([]byte, error) {
	return json.Marshal(JsonJob{
		Headers:          job.Headers,
		TestSuiteID:      job.TestSuite.ID.String(),
		TestSuiteRunID:   job.TestSuiteRun.ID,
		TestID:           job.Test.ID.String(),
		TestVersion:      job.Run.TestVersion,
		RunID:            job.Run.ID,
		TraceID:          job.Run.TraceID.String(),
		PollingProfileID: job.PollingProfile.ID.String(),
		DataStoreID:      job.DataStore.ID.String(),
	})
}

func (job *Job) UnmarshalJSON(data []byte) error {
	var jj JsonJob
	err := json.Unmarshal(data, &jj)
	if err != nil {
		return err
	}

	traceID, err := trace.TraceIDFromHex(jj.TraceID)
	if err != nil {
		return err
	}

	job.Headers = jj.Headers
	job.TestSuite.ID = id.ID(jj.TestSuiteID)
	job.TestSuiteRun.ID = jj.TestSuiteRunID
	job.Test.ID = id.ID(jj.TestID)
	job.Run.ID = jj.RunID
	job.Run.TestVersion = jj.TestVersion
	job.Run.TraceID = traceID
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

type queueConfigurer[T any] struct {
	cancelRunHandlerFn func(ctx context.Context, run test.Run) error
	subscriptor        subscriptor

	runs  testRunGetter
	tests testGetter

	testSuiteRuns testSuiteRunGetter
	testSuites    testSuiteGetter

	pollingProfiles pollingProfileGetter
	dataStores      dataStoreGetter

	meter metric.Meter

	instanceID string
}

func NewQueueConfigurer() *queueConfigurer[Job] {
	return &queueConfigurer[Job]{}
}

func (qb *queueConfigurer[T]) WithCancelRunHandlerFn(fn func(ctx context.Context, run test.Run) error) *queueConfigurer[T] {
	qb.cancelRunHandlerFn = fn
	return qb
}

func (qb *queueConfigurer[T]) WithSubscriptor(subscriptor subscriptor) *queueConfigurer[T] {
	qb.subscriptor = subscriptor
	return qb
}

func (qb *queueConfigurer[T]) WithRunGetter(runs testRunGetter) *queueConfigurer[T] {
	qb.runs = runs
	return qb
}

func (qb *queueConfigurer[T]) WithInstanceID(id string) *queueConfigurer[T] {
	qb.instanceID = id
	return qb
}

func (qb *queueConfigurer[T]) WithTestGetter(tests testGetter) *queueConfigurer[T] {
	qb.tests = tests
	return qb
}

func (qb *queueConfigurer[T]) WithPollingProfileGetter(pollingProfiles pollingProfileGetter) *queueConfigurer[T] {
	qb.pollingProfiles = pollingProfiles
	return qb
}

func (qb *queueConfigurer[T]) WithDataStoreGetter(dataStore dataStoreGetter) *queueConfigurer[T] {
	qb.dataStores = dataStore
	return qb
}

func (qb *queueConfigurer[T]) WithTestSuiteGetter(suites testSuiteGetter) *queueConfigurer[T] {
	qb.testSuites = suites
	return qb
}

func (qb *queueConfigurer[T]) WithTestSuiteRunGetter(suiteRuns testSuiteRunGetter) *queueConfigurer[T] {
	qb.testSuiteRuns = suiteRuns
	return qb
}

func (qb *queueConfigurer[T]) WithMetricMeter(meter metric.Meter) *queueConfigurer[T] {
	qb.meter = meter
	return qb
}

func (qb *queueConfigurer[T]) Configure(queue *pipeline.Queue[Job]) {
	q := &Queue{
		cancelRunHandlerFn: qb.cancelRunHandlerFn,
		subscriptor:        qb.subscriptor,

		runs:  qb.runs,
		tests: qb.tests,

		testSuiteRuns: qb.testSuiteRuns,
		testSuites:    qb.testSuites,

		pollingProfiles: qb.pollingProfiles,
		dataStores:      qb.dataStores,

		instanceID: qb.instanceID,
	}

	queue.EnqueuePreprocessorFn = q.enqueuePreprocess
	queue.ListenPreprocessorFn = q.listenPreprocess

	queue.InitializeMetrics(qb.meter)

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

	instanceID string
}

func propagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}, tenantPropagator{})
}

func (q Queue) enqueuePreprocess(ctx context.Context, input Job) Job {
	if input.Headers == nil {
		input.Headers = &headers{}
	}
	propagator().Inject(ctx, propagation.MapCarrier(*input.Headers))
	input.Headers.Set("InstanceID", q.instanceID)

	version := 1
	if input.Test.Version != nil {
		version = *input.Test.Version
	}

	return Job{
		Headers: input.Headers,

		Test: test.Test{ID: input.Test.ID},
		Run:  test.Run{ID: input.Run.ID, TestVersion: version, TraceID: input.Run.TraceID},

		TestSuite:    testsuite.TestSuite{ID: input.TestSuite.ID},
		TestSuiteRun: testsuite.TestSuiteRun{ID: input.TestSuiteRun.ID},

		PollingProfile: pollingprofile.PollingProfile{ID: input.PollingProfile.ID},
		DataStore:      datastore.DataStore{ID: input.DataStore.ID},
	}
}

func (q Queue) listenPreprocess(ctx context.Context, job Job) (context.Context, Job) {
	log.Printf("queue: received job for run %d", job.Run.ID)
	// this is called when a new job is put in the queue and we need to process it
	ctx = propagator().Extract(ctx, propagation.MapCarrier(*job.Headers))

	ctx = context.WithValue(ctx, "LastInstanceID", job.Headers.Get("InstanceID"))

	ctx, cancelCtx := context.WithCancelCause(ctx)
	q.listenForUserRequests(ctx, cancelCtx, job)

	return ctx, Job{
		Headers:        job.Headers,
		TenantID:       job.TenantID,
		Test:           q.resolveTest(ctx, job),
		Run:            q.resolveTestRun(ctx, job),
		TestSuite:      q.resolveTestSuite(ctx, job),
		TestSuiteRun:   q.resolveTestSuiteRun(ctx, job),
		PollingProfile: q.resolvePollingProfile(ctx, job),
		DataStore:      q.resolveDataStore(ctx, job),
	}
}

type UserRequestType string

var (
	UserRequestTypeStop                UserRequestType = "stop"
	UserRequestTypeSkipTraceCollection UserRequestType = "skip_trace_collection"
)

type UserRequest struct {
	TenantID string
	TestID   id.ID
	RunID    int
	Type     string
}

func (sr UserRequest) ResourceID(requestType UserRequestType) string {
	runID := (test.Run{ID: sr.RunID, TestID: sr.TestID}).ResourceID()
	runID = strings.ReplaceAll(runID, "/", ".")

	return fmt.Sprintf("%s.%s.%s", sr.TenantID, runID, requestType)
}

var (
	ErrSkipTraceCollection = errors.New("skip trace collection")
)

func (q Queue) listenForUserRequests(ctx context.Context, cancelCtx context.CancelCauseFunc, job Job) {
	if q.subscriptor == nil {
		return
	}

	stopTestCallback := subscription.NewSubscriberFunction(func(m subscription.Message) error {
		cancelCtx(nil)
		request := UserRequest{}
		err := m.DecodeContent(&request)
		if err != nil {
			return fmt.Errorf("cannot decode UserRequest message: %w", err)
		}

		run, err := q.runs.GetRun(ctx, request.TestID, request.RunID)
		if err != nil {
			return fmt.Errorf("failed to get run %d for test %s: %w", request.RunID, request.TestID, err)
		}

		if run.State == test.RunStateStopped {
			return nil
		}

		return q.cancelRunHandlerFn(ctx, run)
	})

	skipPollCallback := subscription.NewSubscriberFunction(func(m subscription.Message) error {
		request := UserRequest{}
		err := m.DecodeContent(&request)
		if err != nil {
			return fmt.Errorf("cannot decode UserRequest message: %w", err)
		}

		run, err := q.runs.GetRun(ctx, request.TestID, request.RunID)
		if err != nil {
			return fmt.Errorf("failed to get run %d for test %s: %w", request.RunID, request.TestID, err)
		}

		if run.State == test.RunStateStopped || run.State.IsFinal() {
			return nil
		}

		cancelCtx(ErrSkipTraceCollection)
		return nil
	})

	userReq := UserRequest{
		TenantID: job.TenantID,
		TestID:   job.Test.ID,
		RunID:    job.Run.ID,
	}

	q.subscriptor.Subscribe(userReq.ResourceID(UserRequestTypeStop), stopTestCallback)
	q.subscriptor.Subscribe(userReq.ResourceID(UserRequestTypeSkipTraceCollection), skipPollCallback)
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
		log.Printf("cannot resolve TestSuite: %s", err.Error())
		return testsuite.TestSuite{}
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
		log.Printf("cannot resolve TestSuiteRun: %s", err.Error())
		return testsuite.TestSuiteRun{}
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
		log.Printf("cannot resolve Test: %s", err.Error())
		return test.Test{}
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
		log.Printf("cannot resolve test run: %s", err.Error())
		return test.Run{}
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
		log.Printf("cannot resolve PollingProfile: %s", err.Error())
		return pollingprofile.PollingProfile{}
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
		log.Printf("cannot resolve DataStore: %s", err.Error())
		return datastore.DataStore{}
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

	instanceID := ctx.Value("instanceID")
	if instanceID != nil {
		carrier.Set("instanceID", instanceID.(string))
	}

	eventProperties := middleware.EventPropertiesFromContext(ctx)
	if eventProperties != "" {
		carrier.Set(string(middleware.EventPropertiesKey), eventProperties)
	}
}

// Extract returns a copy of parent with the baggage from the carrier added.
func (b tenantPropagator) Extract(parent context.Context, carrier propagation.TextMapCarrier) context.Context {
	tenantID := carrier.Get(string(middleware.TenantIDKey))
	if tenantID == "" {
		return parent
	}

	resultingCtx := context.WithValue(parent, middleware.TenantIDKey, tenantID)

	instanceID := carrier.Get("instanceID")
	if instanceID != "" {
		resultingCtx = context.WithValue(resultingCtx, "instanceID", instanceID)
	}

	eventProperties := carrier.Get(string(middleware.EventPropertiesKey))
	if eventProperties != "" {
		resultingCtx = context.WithValue(resultingCtx, middleware.EventPropertiesKey, eventProperties)
	}

	return resultingCtx

}

// Fields returns the keys who's values are set with Inject.
func (b tenantPropagator) Fields() []string {
	return []string{string(middleware.TenantIDKey)}
}
