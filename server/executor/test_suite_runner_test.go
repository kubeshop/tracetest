package executor_test

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/pkg/pipeline"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/kubeshop/tracetest/server/testsuite"
	"github.com/kubeshop/tracetest/server/variableset"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric/noop"
)

type fakeTestRunner struct {
	db                  test.RunRepository
	subscriptionManager subscription.Manager
	returnErr           bool
	uid                 int
}

func (r *fakeTestRunner) Run(ctx context.Context, testObj test.Test, metadata test.RunMetadata, variableSet variableset.VariableSet, requiredGates *[]testrunner.RequiredGate) test.Run {
	run := test.NewRun()
	run.VariableSet = variableSet
	run.State = test.RunStateCreated
	newRun, err := r.db.CreateRun(ctx, testObj, run)
	if err != nil {
		panic(err)
	}

	go func() {
		run := newRun                      // make a local copy to avoid race conditions
		time.Sleep(100 * time.Millisecond) // simulate some real work

		if r.returnErr {
			run.State = test.RunStateTriggerFailed
			run.LastError = fmt.Errorf("failed to do something")
		} else {
			run.State = test.RunStateFinished
		}

		r.uid++

		run.Outputs = (maps.Ordered[string, test.RunOutput]{}).MustAdd("USER_ID", test.RunOutput{
			Value: strconv.Itoa(r.uid),
		})

		err = r.db.UpdateRun(ctx, run)
		r.subscriptionManager.PublishUpdate(subscription.Message{
			ResourceID: newRun.ResourceID(),
			Type:       "result_update",
			Content:    run,
		})
	}()

	return newRun
}

func TestTestSuiteRunner(t *testing.T) {

	t.Run("NoErrors", func(t *testing.T) {
		runTestSuiteRunnerTest(t, false, func(t *testing.T, actual testsuite.TestSuiteRun) {
			assert.Equal(t, testsuite.TestSuiteStateFinished, actual.State)
			require.Len(t, actual.Steps, 2)
			assert.Equal(t, actual.Steps[0].State, test.RunStateFinished)
			assert.Equal(t, actual.Steps[1].State, test.RunStateFinished)
			assert.Equal(t, "http://my-service.com", actual.VariableSet.Get("url"))

			assert.Equal(t, test.RunOutput{Name: "", Value: "1", SpanID: ""}, actual.Steps[0].Outputs.Get("USER_ID"))

			// this assertion is supposed to test that the output from the previous step
			// is injected in the env for the next. In practice, this depends
			// on the `fakeTestRunner` used here to actually save the environment
			// to the test run, like the real test runner would.
			// see line 27
			assert.Equal(t, "1", actual.Steps[1].VariableSet.Get("USER_ID"))
			assert.Equal(t, test.RunOutput{Name: "", Value: "2", SpanID: ""}, actual.Steps[1].Outputs.Get("USER_ID"))

			assert.Equal(t, "2", actual.VariableSet.Get("USER_ID"))

		})
	})

	t.Run("WithErrors", func(t *testing.T) {
		runTestSuiteRunnerTest(t, true, func(t *testing.T, actual testsuite.TestSuiteRun) {
			assert.Equal(t, testsuite.TestSuiteStateFailed, actual.State)
			require.Len(t, actual.Steps, 1)
			assert.Equal(t, test.RunStateTriggerFailed, actual.Steps[0].State)
		})
	})

}

func getDB() (model.Repository, *sql.DB) {
	rawDB := testmock.GetRawTestingDatabase()
	db := testmock.GetTestingDatabaseFromRawDB(rawDB)

	return db, rawDB
}

func runTestSuiteRunnerTest(t *testing.T, withErrors bool, assert func(t *testing.T, actual testsuite.TestSuiteRun)) {
	ctx := context.Background()
	_, rawDB := getDB()

	subscriptionManager := subscription.NewManager()
	testRepo := test.NewRepository(rawDB)
	runRepo := test.NewRunRepository(rawDB)

	testRunner := &fakeTestRunner{
		runRepo,
		subscriptionManager,
		withErrors,
		0,
	}

	test1, err := testRepo.Create(ctx, test.Test{Name: "Test 1"})
	require.NoError(t, err)

	test2, err := testRepo.Create(ctx, test.Test{Name: "Test 2"})
	require.NoError(t, err)

	transactionsRepo := testsuite.NewRepository(rawDB, testRepo)
	transactionRunRepo := testsuite.NewRunRepository(rawDB, runRepo)
	tran, err := transactionsRepo.Create(ctx, testsuite.TestSuite{
		ID:      id.ID("tran1"),
		Name:    "test_suite",
		StepIDs: []id.ID{test1.ID, test2.ID},
	})
	require.NoError(t, err)

	tran, err = transactionsRepo.GetAugmented(context.TODO(), tran.ID)
	require.NoError(t, err)

	metadata := test.RunMetadata{
		"environment": "production",
		"service":     "tracetest",
	}

	envRepository := variableset.NewRepository(rawDB)
	env, err := envRepository.Create(ctx, variableset.VariableSet{
		Name: "production",
		Values: []variableset.VariableSetValue{
			{
				Key:   "url",
				Value: "http://my-service.com",
			},
		},
	})
	require.NoError(t, err)

	runner := executor.NewTestSuiteRunner(testRunner, transactionRunRepo, subscriptionManager)

	queueBuilder := executor.NewQueueConfigurer().
		WithTestSuiteGetter(transactionsRepo).
		WithTestSuiteRunGetter(transactionRunRepo).
		WithMetricMeter(noop.NewMeterProvider().Meter("noop"))

	pipeline := pipeline.New(queueBuilder,
		pipeline.Step[executor.Job]{Processor: runner, Driver: pipeline.NewInMemoryQueueDriver[executor.Job]("testSuiteRunner")},
	)

	transactionPipeline := executor.NewTestSuitePipeline(pipeline, transactionRunRepo)
	transactionPipeline.Start()

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	transactionRun := transactionPipeline.Run(ctxWithTimeout, tran, metadata, env, nil)

	done := make(chan testsuite.TestSuiteRun, 1)
	sf := subscription.NewSubscriberFunction(func(m subscription.Message) error {
		tr := testsuite.TestSuiteRun{}
		err := mapstructure.Decode(m.Content, &tr)
		if err != nil {
			return fmt.Errorf("cannot decode TestSuiteRun message: %w", err)
		}
		if tr.State.IsFinal() {
			done <- tr
		}

		return nil
	})
	subscriptionManager.Subscribe(transactionRun.ResourceID(), sf)

	select {
	case finalRun := <-done:
		subscriptionManager.Unsubscribe(transactionRun.ResourceID(), sf.ID()) //cleanup to avoid race conditions
		assert(t, finalRun)
	case <-time.After(10 * time.Second):
		t.Log("timeout after 10 second")
		t.FailNow()
	}
	transactionPipeline.Stop()
}
