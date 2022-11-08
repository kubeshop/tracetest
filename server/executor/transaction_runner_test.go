package executor_test

import (
	"context"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type simpleTestRunner struct {
	db model.Repository
}

func (r simpleTestRunner) Run(ctx context.Context, test model.Test, metadata model.RunMetadata, env model.Environment) model.Run {
	run := model.NewRun()
	run.State = model.RunStateCreated
	newRun, err := r.db.CreateRun(ctx, test, run)
	if err != nil {
		panic(err)
	}

	go func() {
		time.Sleep(2 * time.Second) // simulate some real work

		newRun.State = model.RunStateFinished
		err = r.db.UpdateRun(ctx, newRun)
		if err != nil {
			panic(err)
		}
	}()

	return newRun
}

func TestTransactionRunner(t *testing.T) {
	ctx := context.Background()
	db, clear := getDB()
	defer clear()

	testRunner := simpleTestRunner{
		db,
	}

	test1, err := db.CreateTest(ctx, model.Test{
		Name: "Test 1",
	})
	require.NoError(t, err)

	test2, err := db.CreateTest(ctx, model.Test{
		Name: "Test 1",
	})
	require.NoError(t, err)

	transaction, err := db.CreateTransaction(ctx, model.Transaction{
		Name:    "transaction",
		Version: 1,
		Steps: []model.Test{
			test1,
			test2,
		},
	})
	require.NoError(t, err)

	metadata := model.RunMetadata{
		"environment": "production",
		"service":     "tracetest",
	}

	env, err := db.CreateEnvironment(ctx, model.Environment{
		Name: "production",
		Values: []model.EnvironmentValue{
			{
				Key:   "url",
				Value: "http://my-service.com",
			},
		},
	})
	require.NoError(t, err)

	config := config.Config{
		PoolingConfig: config.PoolingConfig{
			RetryDelay: "2s",
		},
	}

	runner := executor.NewTransactionRunner(testRunner, db, config)
	runner.Start(5)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	transactionRun := runner.Run(ctxWithTimeout, transaction, metadata, env)

	for !transactionRun.State.IsFinal() {
		transactionRun, err = db.GetTransactionRun(ctxWithTimeout, transactionRun.TransactionID.String(), transactionRun.ID)
		require.NoError(t, err)
		time.Sleep(1 * time.Second)
	}

	assert.Equal(t, model.TransactionRunStateFinished, transactionRun.State)
	assert.Len(t, transactionRun.StepRuns, 2)
	assert.Equal(t, transactionRun.StepRuns[0].State, model.RunStateFinished)
	assert.Equal(t, transactionRun.StepRuns[1].State, model.RunStateFinished)
}

func getDB() (model.Repository, func()) {
	db, err := testmock.GetTestingDatabase("file://../migrations")
	if err != nil {
		panic(err)
	}

	clean := func() {
		err = db.Drop()
		if err != nil {
			panic(err)
		}
	}

	return db, clean
}
