package testdb_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID/go"
	"github.com/GIT_USER_ID/GIT_REPO_ID/go/testdb"
	"github.com/google/uuid"
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
		Name:        "first test",
		Description: "description",
		ServiceUnderTest: openapi.TestServiceUnderTest{
			Url: "http://localhost:3030/hello-instrumented",
		},
		Assertions: []openapi.Assertion{{}},
		Repeats:    0,
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
			Repeats:    0,
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
	res := openapi.Result{
		Id:          id,
		CreatedAt:   ti,
		CompletedAt: ti,
		Traceid:     "123",
	}
	ctx := context.Background()
	err = db.CreateResult(ctx, &res)
	assert.NoError(t, err)

	gotRes, err := db.GetResult(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, &res, gotRes)

	res2 := openapi.Result{
		Id:          id,
		CreatedAt:   ti,
		CompletedAt: ti,
		Traceid:     "1234",
	}

	err = db.UpdateResult(ctx, &res2)
	assert.NoError(t, err)

	gotRes, err = db.GetResult(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, &res2, gotRes)
}
