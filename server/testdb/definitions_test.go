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

	def := (model.OrderedMap[model.SpanQuery, []model.Assertion]{}).MustAdd(`span[service.name="Pokeshop"]`, []model.Assertion{
		{
			Attribute:  "tracetest.span.duration",
			Comparator: comparator.Eq,
			Value:      "2000000000",
		},
	})

	err := db.SetDefiniton(context.TODO(), test, def)
	require.NoError(t, err)

	actual, err := db.GetDefiniton(context.TODO(), test)
	require.NoError(t, err)
	assert.Equal(t, def, actual)
}
