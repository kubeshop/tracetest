package testdb_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/model/modeltest"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/testdb"
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
		Environment: environment.Environment{
			Name:        "env1",
			Description: "env1",
			Values: []environment.EnvironmentValue{{
				Key:   "key",
				Value: "value",
			}}},
	}

	updated, err := db.CreateRun(context.TODO(), test, run)
	require.NoError(t, err)

	actual, err := db.GetRun(context.TODO(), test.ID, updated.ID)
	require.NoError(t, err)

	assert.NotEmpty(t, actual.ID)
	assert.Equal(t, test.ID, actual.TestID)
	assert.Equal(t, test.Version, actual.TestVersion)
	assert.Equal(t, model.RunStateCreated, actual.State)
	assert.Equal(t, run.TraceID, actual.TraceID)
	assert.Equal(t, run.SpanID, actual.SpanID)
	assert.Equal(t, run.Metadata, actual.Metadata)
	assert.Equal(t, run.Environment, actual.Environment)
}

func TestCreateRunIDsIncrementForTest(t *testing.T) {
	db, clean := getDB()
	defer clean()

	run := model.Run{
		TraceID: testdb.IDGen.TraceID(),
		SpanID:  testdb.IDGen.SpanID(),
	}

	test1 := createTest(t, db)
	test2 := createTest(t, db)

	t1r1, err := db.CreateRun(context.TODO(), test1, run)
	require.NoError(t, err)

	t2r1, err := db.CreateRun(context.TODO(), test2, run)
	require.NoError(t, err)

	t1r2, err := db.CreateRun(context.TODO(), test1, run)
	require.NoError(t, err)

	t2r2, err := db.CreateRun(context.TODO(), test2, run)
	require.NoError(t, err)

	assert.Equal(t, 1, t1r1.ID)
	assert.Equal(t, 2, t1r2.ID)
	assert.Equal(t, 1, t2r1.ID)
	assert.Equal(t, 2, t2r2.ID)
}

func TestUpdateRun(t *testing.T) {
	db, clean := getDB()
	defer clean()

	test := createTest(t, db)
	run := createRun(t, db, test)

	run.State = model.RunStateFinished
	run.Trace = &model.Trace{
		ID: testdb.IDGen.TraceID(),
		RootSpan: model.Span{
			ID: testdb.IDGen.SpanID(),
			Attributes: model.Attributes{
				"service.name":            "Pokeshop",
				"tracetest.span.duration": "2000",
			},
		},
	}
	run.Trace.Flat = map[trace.SpanID]*model.Span{
		run.Trace.RootSpan.ID: &run.Trace.RootSpan,
	}
	run.Results = &model.RunResults{
		AllPassed: true,
		Results: (maps.Ordered[model.SpanQuery, []model.AssertionResult]{}).MustAdd(`span[service.name="Pokeshop"]`, []model.AssertionResult{
			{
				Assertion: model.Assertion(`attr:tracetest.span.duration = 2000`),
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

	run.Outputs = (maps.Ordered[string, model.RunOutput]{}).
		MustAdd("key", model.RunOutput{
			Value: "value",
		})

	err := db.UpdateRun(context.TODO(), run)
	require.NoError(t, err)

	actual, err := db.GetRun(context.TODO(), test.ID, run.ID)
	require.NoError(t, err)

	modeltest.AssertRunEqual(t, run, actual)

	updatedList, err := db.GetTestRuns(context.TODO(), test, 20, 0)
	require.NoError(t, err)

	assert.Len(t, updatedList.Items, 1)
	assert.Equal(t, 1, updatedList.TotalCount)
	modeltest.AssertRunEqual(t, updatedList.Items[0], actual)
}

func TestUpdateRunWithNewIDs(t *testing.T) {
	db, clean := getDB()
	defer clean()

	t1r1 := createRun(t, db, createTest(t, db))
	t2r1 := createRun(t, db, createTest(t, db))

	t1r1.Metadata = model.RunMetadata{"key": "val"}
	db.UpdateRun(context.TODO(), t1r1)

	t1r1Updated, err := db.GetRun(context.TODO(), t1r1.TestID, t1r1.ID)
	require.NoError(t, err)

	t2r1Updated, err := db.GetRun(context.TODO(), t2r1.TestID, t2r1.ID)
	require.NoError(t, err)

	assert.Equal(t, t1r1.Metadata, t1r1Updated.Metadata)
	assert.Equal(t, t2r1.Metadata, t2r1Updated.Metadata)
}

func TestDeleteRun(t *testing.T) {
	db, clean := getDB()
	defer clean()

	t1r1 := createRun(t, db, createTest(t, db))
	t2r1 := createRun(t, db, createTest(t, db))

	db.DeleteRun(context.TODO(), t2r1)

	_, err := db.GetRun(context.TODO(), t1r1.TestID, t1r1.ID)
	require.NoError(t, err)

	_, err = db.GetRun(context.TODO(), t2r1.TestID, t2r1.ID)
	require.ErrorIs(t, err, testdb.ErrNotFound)
}

func TestGetRunByTraceID(t *testing.T) {
	db, clean := getDB()
	defer clean()

	test := createTest(t, db)
	expected := createRun(t, db, test)

	actual, err := db.GetRunByTraceID(context.TODO(), expected.TraceID)
	require.NoError(t, err)

	modeltest.AssertRunEqual(t, expected, actual)
}
