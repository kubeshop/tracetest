package testdb_test

import (
	"context"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/assertions/comparator"
	"github.com/kubeshop/tracetest/model"
	"github.com/kubeshop/tracetest/testdb"
	"github.com/kubeshop/tracetest/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRun(t *testing.T) {
	db, clean := getDB()
	defer clean()

	test := createTest(db)

	run := model.Run{
		TraceID:   testdb.IDGen.TraceID(),
		SpanID:    testdb.IDGen.SpanID(),
		CreatedAt: time.Now(),
		Request:   test.ServiceUnderTest.Request,
	}

	updated, err := db.CreateRun(context.TODO(), test, run)
	require.NoError(t, err)

	actual, err := db.GetRun(context.TODO(), updated.ID)
	require.NoError(t, err)
	assert.Equal(t, model.RunStateCreated, actual.State)
	assert.Equal(t, run.TraceID, actual.TraceID)
	assert.Equal(t, run.SpanID, actual.SpanID)
	assert.Equal(t, run.CreatedAt, actual.CreatedAt)
	assert.Equal(t, run.Request, actual.Request)
}

func TestUpdateRun(t *testing.T) {
	db, clean := getDB()
	defer clean()

	test := createTest(db)
	run := createRun(db, test)

	run.State = model.RunStateFinished
	run.Trace = traces.Trace{
		RootSpan: traces.Span{
			ID: testdb.IDGen.SpanID(),
			Attributes: traces.Attributes{
				"service.name":            "Pokeshop",
				"tracetest.span.duration": "2000",
			},
		},
	}
	run.Results = model.RunResults{
		AllPassed: true,
		Results: model.Results{
			`span[service.name="Pokeshop"]`: []model.AssertionResult{
				{
					Assertion: model.Assertion{
						Attribute:  "tracetest.span.duration",
						Comparator: comparator.Eq,
						Value:      "2000",
					},
					Results: []model.SpanAssertionResult{
						{
							SpanID:        run.Trace.RootSpan.ID,
							ObservedValue: "2000",
							CompareErr:    nil,
						},
					},
				},
			},
		},
	}

	err := db.UpdateRun(context.TODO(), run)
	require.NoError(t, err)

	actual, err := db.GetRun(context.TODO(), run.ID)
	require.NoError(t, err)

	assert.Equal(t, run, actual)
	assert.Equal(t, run.SpanID, actual.SpanID)
	assert.Equal(t, run.CreatedAt, actual.CreatedAt)
	assert.Equal(t, run.Request, actual.Request)
	assert.Equal(t, run.State, actual.State)
	assert.Equal(t, run.Trace, actual.Trace)
	assert.Equal(t, run.Results, actual.Results)

	updatedList, err := db.GetTestRuns(context.TODO(), test, 20, 0)
	require.NoError(t, err)
	assert.Equal(t, actual, updatedList[0])
}
