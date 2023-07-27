package executor_test

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/kubeshop/tracetest/server/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeTestRunner struct {
	db                  test.RunRepository
	subscriptionManager *subscription.Manager
	returnErr           bool
	uid                 int
}

func (r *fakeTestRunner) Run(ctx context.Context, testObj test.Test, metadata test.RunMetadata, env environment.Environment, requiredGates *[]testrunner.RequiredGate) test.Run {
	run := test.NewRun()
	run.Environment = env
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

func TestTransactionRunner(t *testing.T) {

	t.Run("NoErrors", func(t *testing.T) {
		runTransactionRunnerTest(t, false, func(t *testing.T, actual transaction.TransactionRun) {
			assert.Equal(t, transaction.TransactionRunStateFinished, actual.State)
			require.Len(t, actual.Steps, 2)
			assert.Equal(t, actual.Steps[0].State, test.RunStateFinished)
			assert.Equal(t, actual.Steps[1].State, test.RunStateFinished)
			assert.Equal(t, "http://my-service.com", actual.Environment.Get("url"))

			assert.Equal(t, test.RunOutput{Name: "", Value: "1", SpanID: ""}, actual.Steps[0].Outputs.Get("USER_ID"))

			// this assertion is supposed to test that the output from the previous step
			// is injected in the env for the next. In practice, this depends
			// on the `fakeTestRunner` used here to actually save the environment
			// to the test run, like the real test runner would.
			// see line 27
			assert.Equal(t, "1", actual.Steps[1].Environment.Get("USER_ID"))
			assert.Equal(t, test.RunOutput{Name: "", Value: "2", SpanID: ""}, actual.Steps[1].Outputs.Get("USER_ID"))

			assert.Equal(t, "2", actual.Environment.Get("USER_ID"))

		})
	})

	t.Run("WithErrors", func(t *testing.T) {
		runTransactionRunnerTest(t, true, func(t *testing.T, actual transaction.TransactionRun) {
			assert.Equal(t, transaction.TransactionRunStateFailed, actual.State)
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

func runTransactionRunnerTest(t *testing.T, withErrors bool, assert func(t *testing.T, actual transaction.TransactionRun)) {
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

	transactionsRepo := transaction.NewRepository(rawDB, testRepo)
	transactionRunRepo := transaction.NewRunRepository(rawDB, runRepo)
	tran, err := transactionsRepo.Create(ctx, transaction.Transaction{
		ID:      id.ID("tran1"),
		Name:    "transaction",
		StepIDs: []id.ID{test1.ID, test2.ID},
	})
	require.NoError(t, err)

	tran, err = transactionsRepo.GetAugmented(context.TODO(), tran.ID)
	require.NoError(t, err)

	metadata := test.RunMetadata{
		"environment": "production",
		"service":     "tracetest",
	}

	envRepository := environment.NewRepository(rawDB)
	env, err := envRepository.Create(ctx, environment.Environment{
		Name: "production",
		Values: []environment.EnvironmentValue{
			{
				Key:   "url",
				Value: "http://my-service.com",
			},
		},
	})
	require.NoError(t, err)

	runner := executor.NewTransactionRunner(testRunner, transactionRunRepo, subscriptionManager)

	queueBuilder := executor.NewQueueBuilder().
		WithTransactionGetter(transactionsRepo).
		WithTransactionRunGetter(transactionRunRepo)

	pipeline := executor.NewPipeline(queueBuilder,
		executor.PipelineStep{Processor: runner, Driver: executor.NewInMemoryQueueDriver("runner")},
	)

	transactionPipeline := executor.NewTransactionPipeline(pipeline, transactionRunRepo)
	transactionPipeline.Start()

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	transactionRun := transactionPipeline.Run(ctxWithTimeout, tran, metadata, env, nil)

	done := make(chan transaction.TransactionRun, 1)
	sf := subscription.NewSubscriberFunction(func(m subscription.Message) error {
		tr := m.Content.(transaction.TransactionRun)
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
