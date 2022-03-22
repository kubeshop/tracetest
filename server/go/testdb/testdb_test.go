package testdb_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	openapi "github.com/kubeshop/tracetest/server/go"
	"github.com/kubeshop/tracetest/server/go/testdb"
	"github.com/stretchr/testify/assert"
)

func TestCreateTest(t *testing.T) {
	dsn := "host=localhost user=postgres password=postgres port=5432 sslmode=disable"
	db, err := testdb.New(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	}()
	test := openapi.Test{
		TestId:      "",
		Name:        "first test",
		Description: "description",
		ServiceUnderTest: openapi.TestServiceUnderTest{
			Url: "http://localhost:3030/hello-instrumented",
		},
	}
	ctx := context.Background()
	id, err := db.CreateTest(ctx, &test)
	if err != nil {
		t.Fatal(err)
	}

	gotTest, err := db.GetTest(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, &test, gotTest)
}

func TestUpdateTest(t *testing.T) {
	dsn := "host=localhost user=postgres password=postgres port=5432 sslmode=disable"
	db, err := testdb.New(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	}()
	test := openapi.Test{
		TestId:      "",
		Name:        "first test",
		Description: "description",
		ServiceUnderTest: openapi.TestServiceUnderTest{
			Url: "http://localhost:3030/hello-instrumented",
		},
	}
	ctx := context.Background()
	id, err := db.CreateTest(ctx, &test)
	if err != nil {
		t.Fatal(err)
	}

	test.Name = "updated test"
	err = db.UpdateTest(ctx, &test)
	if err != nil {
		t.Fatal(err)
	}

	gotTest, err := db.GetTest(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, &test, gotTest)
}

func TestGetTest(t *testing.T) {
	dsn := "host=localhost user=postgres password=postgres port=5432 sslmode=disable"
	db, err := testdb.New(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	}()
	ctx := context.Background()
	_, err = db.GetTest(ctx, uuid.New().String())
	assert.Equal(t, openapi.ErrNotFound, err)
}

func TestGetTests(t *testing.T) {
	dsn := "host=localhost user=postgres password=postgres port=5432 sslmode=disable"
	db, err := testdb.New(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	}()

	ctx := context.Background()
	for i := 0; i < 2; i++ {
		test := openapi.Test{
			Name:        strconv.Itoa(i),
			Description: "description",
			ServiceUnderTest: openapi.TestServiceUnderTest{
				Url: "http://localhost:3030/hello-instrumented",
			},
			Assertions: []openapi.Assertion{{}},
		}
		_, err = db.CreateTest(ctx, &test)
		if err != nil {
			t.Fatal(err)
		}
	}
	gotTests, err := db.GetTests(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, gotTests, 2)
}

func TestCreateResults(t *testing.T) {
	dsn := "host=localhost user=postgres password=postgres port=5432 sslmode=disable"
	db, err := testdb.New(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	}()
	ti := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	id := uuid.New().String()
	testID := uuid.New().String()
	res := openapi.TestRunResult{
		ResultId:    id,
		CreatedAt:   ti,
		CompletedAt: ti,
		TraceId:     "123",
	}
	ctx := context.Background()
	err = db.CreateResult(ctx, testID, &res)
	assert.NoError(t, err)

	gotRes, err := db.GetResult(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, &res, gotRes)

	res2 := openapi.TestRunResult{
		ResultId:    id,
		CreatedAt:   ti,
		CompletedAt: ti,
		TraceId:     "1234",
	}

	err = db.UpdateResult(ctx, &res2)
	assert.NoError(t, err)

	gotRes, err = db.GetResult(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, &res2, gotRes)

	gotResults, err := db.GetResultsByTestID(ctx, testID)
	assert.NoError(t, err)
	assert.Equal(t, res2, gotResults[0])
}

func TestCreateAssertions(t *testing.T) {
	dsn := "host=localhost user=postgres password=postgres port=5432 sslmode=disable"
	db, err := testdb.New(dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = db.Drop()
		if err != nil {
			t.Fatal(err)
		}
	}()
	res := openapi.Assertion{
		Selectors: []openapi.SelectorItem{
			{
				LocationName: "SPAN",
				PropertyName: "operation",
				Value:        "POST /users/verify",
				ValueType:    "stringValue",
			},
		},
		SpanAssertions: []openapi.SpanAssertion{
			{
				LocationName:    "SPAN_ATTRIBUTES",
				PropertyName:    "http.status.code",
				ValueType:       "intValue",
				Operator:        "EQUALS",
				ComparisonValue: "200",
			},
		},
	}

	testid := uuid.New().String()
	ctx := context.Background()
	id, err := db.CreateAssertion(ctx, testid, &res)
	assert.NoError(t, err)
	res.AssertionId = id

	gotRes, err := db.GetAssertion(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, &res, gotRes)

	gotAssertions, err := db.GetAssertionsByTestID(ctx, testid)
	assert.NoError(t, err)
	assert.Equal(t, res, gotAssertions[0])
}
