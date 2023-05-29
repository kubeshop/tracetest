package testdb_test

import (
	"context"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/testmock"
)

func getDB() (model.Repository, func()) {
	db := testmock.GetTestingDatabase()

	clean := func() {
		defer db.Close()
		err := db.Drop()
		if err != nil {
			panic(err)
		}
	}

	return db, clean
}

func createTestWithName(t *testing.T, db model.Repository, name string) model.Test {
	t.Helper()
	test := model.Test{
		Name:        name,
		Description: "description",
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
			HTTP: &model.HTTPRequest{
				URL: "http://localhost:3030/hello-instrumented",
			},
		},
	}

	updated, err := db.CreateTest(context.TODO(), test)
	if err != nil {
		panic(err)
	}
	return updated
}

func createTest(t *testing.T, db model.Repository) model.Test {
	return createTestWithName(t, db, "first test")
}

func createRun(t *testing.T, db model.Repository, test model.Test) model.Run {
	t.Helper()
	run := model.Run{
		TraceID:   testdb.IDGen.TraceID(),
		SpanID:    testdb.IDGen.SpanID(),
		CreatedAt: time.Now(),
	}
	updated, err := db.CreateRun(context.TODO(), test, run)
	if err != nil {
		panic(err)
	}

	return updated
}
