package executor_test

import (
	"context"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/testmock"
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

	runner := executor.NewTransactionRunner(testRunner)
	runner.Start(1)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	runner.Run(ctxWithTimeout, transaction, metadata, env)
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
