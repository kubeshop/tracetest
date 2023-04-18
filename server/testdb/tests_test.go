package testdb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTest(t *testing.T) {
	db, clean := getDB()
	defer clean()

	test := model.Test{
		Name:        "first test",
		Description: "description",
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
			HTTP: &model.HTTPRequest{
				URL: "http://localhost:3030/hello-instrumented",
			},
		},
		Outputs: (maps.Ordered[string, model.Output]{}).
			MustAdd("output1", model.Output{
				Selector: model.SpanQuery(`span[name="root"]`),
				Value:    "${attr:myapp.some_attribute}",
			}),
	}

	updated, err := db.CreateTest(context.TODO(), test)
	require.NoError(t, err)

	actual, err := db.GetLatestTestVersion(context.TODO(), updated.ID)
	require.NoError(t, err)
	assert.Equal(t, test.Name, actual.Name)
	assert.Equal(t, test.Description, actual.Description)
	assert.Equal(t, test.ServiceUnderTest, actual.ServiceUnderTest)
	assert.Equal(t, test.Specs, actual.Specs)
	assert.Equal(t, test.Outputs, actual.Outputs)
	assert.False(t, actual.CreatedAt.IsZero())
}

func TestDeleteTest(t *testing.T) {
	db, clean := getDB()
	defer clean()

	test := createTest(t, db)

	err := db.DeleteTest(context.TODO(), test)
	require.NoError(t, err)

	actual, err := db.GetLatestTestVersion(context.TODO(), test.ID)
	assert.ErrorIs(t, err, testdb.ErrNotFound)
	assert.Empty(t, actual)

}

func TestGetLatestTestVersion(t *testing.T) {
	db, clean := getDB()
	defer clean()

	test := createTestWithName(t, db, "1")
	test.Name = "1 v2"
	test.Version = 2

	_, err := db.UpdateTest(context.TODO(), test)
	require.NoError(t, err)

	latestTest, err := db.GetLatestTestVersion(context.TODO(), test.ID)
	assert.NoError(t, err)
	assert.Equal(t, "1 v2", latestTest.Name)
	assert.Equal(t, 2, latestTest.Version)
}

func TestGetTests(t *testing.T) {
	db, clean := getDB()
	defer clean()

	createTestWithName(t, db, "one")
	createTestWithName(t, db, "two")
	createTestWithName(t, db, "three")

	t.Run("Order", func(t *testing.T) {
		actual, err := db.GetTests(context.TODO(), 20, 0, "", "", "")
		require.NoError(t, err)

		assert.Len(t, actual.Items, 3)
		assert.Equal(t, actual.TotalCount, 3)

		// test order
		assert.Equal(t, actual.TotalCount, 3)
		assert.Equal(t, "three", actual.Items[0].Name)
		assert.Equal(t, "two", actual.Items[1].Name)
		assert.Equal(t, "one", actual.Items[2].Name)
	})

	t.Run("Pagination", func(t *testing.T) {
		actual, err := db.GetTests(context.TODO(), 20, 10, "", "", "")
		require.NoError(t, err)

		assert.Equal(t, actual.TotalCount, 3)
		assert.Len(t, actual.Items, 0)
	})

	t.Run("SortByCreated", func(t *testing.T) {
		actual, err := db.GetTests(context.TODO(), 20, 0, "", "created", "")
		require.NoError(t, err)

		// test order
		assert.Equal(t, "three", actual.Items[0].Name)
		assert.Equal(t, "two", actual.Items[1].Name)
		assert.Equal(t, "one", actual.Items[2].Name)
	})

	t.Run("SortByNameAsc", func(t *testing.T) {
		actual, err := db.GetTests(context.TODO(), 20, 0, "", "name", "asc")
		require.NoError(t, err)

		// test order
		assert.Equal(t, "one", actual.Items[0].Name)
		assert.Equal(t, "three", actual.Items[1].Name)
		assert.Equal(t, "two", actual.Items[2].Name)
	})

	t.Run("SortByNameDesc", func(t *testing.T) {
		actual, err := db.GetTests(context.TODO(), 20, 0, "", "name", "desc")
		require.NoError(t, err)

		// test order
		assert.Equal(t, "two", actual.Items[0].Name)
		assert.Equal(t, "three", actual.Items[1].Name)
		assert.Equal(t, "one", actual.Items[2].Name)
	})

	t.Run("SearchByName", func(t *testing.T) {
		_, _ = db.CreateTest(context.TODO(), model.Test{Name: "VerySpecificName"})
		actual, err := db.GetTests(context.TODO(), 10, 0, "specif", "", "")
		require.NoError(t, err)
		assert.Len(t, actual.Items, 1)
		assert.Equal(t, actual.TotalCount, 1)

		assert.Equal(t, "VerySpecificName", actual.Items[0].Name)
	})

	t.Run("SearchByDescription", func(t *testing.T) {
		_, _ = db.CreateTest(context.TODO(), model.Test{Description: "VeryUniqueText"})

		actual, err := db.GetTests(context.TODO(), 10, 0, "nique", "", "")
		require.NoError(t, err)
		assert.Len(t, actual.Items, 1)
		assert.Equal(t, actual.TotalCount, 1)

		assert.Equal(t, "VeryUniqueText", actual.Items[0].Description)
	})
}

func TestGetTestsWithMultipleVersions(t *testing.T) {
	db, clean := getDB()
	defer clean()

	test1 := createTestWithName(t, db, "1")
	test1.Name = "1 v2"

	_, err := db.UpdateTest(context.TODO(), test1)
	require.NoError(t, err)

	test2 := createTestWithName(t, db, "2")
	test2.Name = "2 v2"

	_, err = db.UpdateTest(context.TODO(), test2)
	require.NoError(t, err)

	tests, err := db.GetTests(context.TODO(), 20, 0, "", "", "")
	assert.NoError(t, err)
	assert.Len(t, tests.Items, 2)
	assert.Equal(t, 2, tests.TotalCount)

	for _, test := range tests.Items {
		assert.Equal(t, 2, test.Version)
	}
}

func TestSummary(t *testing.T) {
	db, clean := getDB()
	defer clean()

	createRunWithResult := func(t *testing.T, db model.Repository, test model.Test, d time.Time, pass, fail int) model.Run {
		t.Helper()
		run := model.Run{
			TraceID:   testdb.IDGen.TraceID(),
			SpanID:    testdb.IDGen.SpanID(),
			CreatedAt: d,
		}

		run, err := db.CreateRun(context.TODO(), test, run)
		if err != nil {
			panic(err)
		}

		result := []model.AssertionResult{
			{
				Results: []model.SpanAssertionResult{},
			},
		}
		for i := 0; i < pass; i++ {
			// CompareErr: nil means passed
			result[0].Results = append(result[0].Results, model.SpanAssertionResult{CompareErr: nil})
		}
		for i := 0; i < fail; i++ {
			result[0].Results = append(result[0].Results, model.SpanAssertionResult{CompareErr: fmt.Errorf("err")})
		}
		run.Results = &model.RunResults{
			Results: (maps.Ordered[model.SpanQuery, []model.AssertionResult]{}).
				MustAdd("span", result),
		}

		err = db.UpdateRun(context.TODO(), run)
		if err != nil {
			panic(err)
		}

		return run
	}

	test := createTest(t, db)

	// 1 run
	{
		t1 := time.Date(2022, time.July, 01, 12, 23, 00, 0, time.UTC)
		createRunWithResult(t, db, test, t1, 2, 0)

		tests, err := db.GetTests(context.TODO(), 20, 0, "", "", "")
		require.NoError(t, err)

		require.Len(t, tests.Items, 1)
		assert.Equal(t, tests.Items[0].ID, test.ID)

		assert.Equal(t, 1, tests.Items[0].Summary.Runs)
		assert.WithinDuration(t, t1, tests.Items[0].Summary.LastRun.Time, 0) // hack for comparing times
		assert.Equal(t, 2, tests.Items[0].Summary.LastRun.Passes)
		assert.Equal(t, 0, tests.Items[0].Summary.LastRun.Fails)
	}

	{
		// 2 runs
		t2 := time.Date(2022, time.July, 01, 12, 23, 30, 0, time.UTC)
		createRunWithResult(t, db, test, t2, 1, 1)

		tests, err := db.GetTests(context.TODO(), 20, 0, "", "", "")
		require.NoError(t, err)

		require.Len(t, tests.Items, 1)
		assert.Equal(t, tests.Items[0].ID, test.ID)

		assert.Equal(t, 2, tests.Items[0].Summary.Runs)
		assert.WithinDuration(t, t2, tests.Items[0].Summary.LastRun.Time, 0) // hack for comparing times
		assert.Equal(t, 1, tests.Items[0].Summary.LastRun.Passes)
		assert.Equal(t, 1, tests.Items[0].Summary.LastRun.Fails)
	}
}
