package testdb_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefinitions(t *testing.T) {
	db, clean := getDB()
	defer clean()

	test := createTest(t, db)

	spec := (model.OrderedMap[model.SpanQuery, []model.Assertion]{}).MustAdd(`span[service.name="Pokeshop"]`, []model.Assertion{
		{
			Attribute:  "tracetest.span.duration",
			Comparator: comparator.Eq,
			Value:      "2000000000",
		},
	})

	err := db.SetSpec(context.TODO(), test, spec)
	require.NoError(t, err)

	actual, err := db.GetSpec(context.TODO(), test)
	require.NoError(t, err)
	assert.Equal(t, spec, actual)
}
