package testdb_test

import (
	"context"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

func TestCreateRun(t *testing.T) {
	db, clean := getDB()
	defer clean()

	test := createTest(t, db)

	run := model.Run{
		TraceID: testdb.IDGen.TraceID(),
		SpanID:  testdb.IDGen.SpanID(),
		Metadata: model.RunMetadata{
			"key": "Value",
		},
	}

	updated, err := db.CreateRun(context.TODO(), test, run)
	require.NoError(t, err)

	actual, err := db.GetRun(context.TODO(), updated.ID)
	require.NoError(t, err)

	assert.NotEmpty(t, actual.ID)
	assert.Equal(t, test.ID, actual.TestID)
	assert.Equal(t, test.Version, actual.TestVersion)
	assert.Equal(t, model.RunStateCreated, actual.State)
	assert.Equal(t, run.TraceID, actual.TraceID)
	assert.Equal(t, run.SpanID, actual.SpanID)
	assert.Equal(t, run.Metadata, actual.Metadata)
}

func TestUpdateRun(t *testing.T) {
	db, clean := getDB()
	defer clean()

	test := createTest(t, db)
	run := createRun(t, db, test)

	run.State = model.RunStateFinished
	run.Trace = &traces.Trace{
		ID: testdb.IDGen.TraceID(),
		RootSpan: traces.Span{
			ID: testdb.IDGen.SpanID(),
			Attributes: traces.Attributes{
				"service.name":            "Pokeshop",
				"tracetest.span.duration": "2000",
			},
		},
	}
	run.Trace.Flat = map[trace.SpanID]*traces.Span{
		run.Trace.RootSpan.ID: &run.Trace.RootSpan,
	}
	run.Results = &model.RunResults{
		AllPassed: true,
		Results: (model.OrderedMap[model.SpanQuery, []model.AssertionResult]{}).MustAdd(`span[service.name="Pokeshop"]`, []model.AssertionResult{
			{
				Assertion: model.Assertion{
					Attribute:  "tracetest.span.duration",
					Comparator: comparator.Eq,
					Value: &model.AssertionExpression{
						LiteralValue: model.LiteralValue{
							Value: "2000",
							Type:  "number",
						},
					},
				},
				Results: []model.SpanAssertionResult{
					{
						SpanID:        &run.Trace.RootSpan.ID,
						ObservedValue: "2000",
						CompareErr:    nil,
					},
				},
			},
		}),
	}

	err := db.UpdateRun(context.TODO(), run)
	require.NoError(t, err)

	actual, err := db.GetRun(context.TODO(), run.ID)
	require.NoError(t, err)

	updatedList, err := db.GetTestRuns(context.TODO(), test, 20, 0)
	require.NoError(t, err)

	// Ignore time fields in this test
	actual.Trace.RootSpan.StartTime = time.Time{}
	actual.Trace.RootSpan.EndTime = time.Time{}

	runFromList := updatedList[0]
	runFromList.Trace.RootSpan.StartTime = time.Time{}
	runFromList.Trace.RootSpan.EndTime = time.Time{}

	assert.Equal(t, run.SpanID.String(), actual.SpanID.String())
	assert.Equal(t, run.State, actual.State)
	assert.Equal(t, run.Trace, actual.Trace)
	assert.Equal(t, run.Results, actual.Results)

	assert.Equal(t, actual, runFromList)
}
