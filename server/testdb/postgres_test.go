package testdb_test

import (
	"context"
	"time"

	"github.com/kubeshop/tracetest/model"
	"github.com/kubeshop/tracetest/testdb"
	"github.com/kubeshop/tracetest/testmock"
)

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

func createTestWithName(db model.Repository, name string) model.Test {
	test := model.Test{
		Name:        name,
		Description: "description",
		ServiceUnderTest: model.ServiceUnderTest{
			Request: model.HTTPRequest{
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

func createTest(db model.Repository) model.Test {
	return createTestWithName(db, "first test")
}

func createRun(db model.Repository, test model.Test) model.Run {
	run := model.Run{
		TraceID:   testdb.IDGen.TraceID(),
		SpanID:    testdb.IDGen.SpanID(),
		CreatedAt: time.Now(),
		Request:   test.ServiceUnderTest.Request,
	}
	updated, err := db.CreateRun(context.TODO(), test, run)
	if err != nil {
		panic(err)
	}

	return updated
}
