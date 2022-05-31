package model_test

import (
	"encoding/json"
	"testing"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefinition(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		def := model.Definition{}

		def, err := def.Add(model.SpanQuery("1"), []model.Assertion{{"1", comparator.Eq, "1"}})
		require.NoError(t, err)

		def, err = def.Add(model.SpanQuery("2"), []model.Assertion{{"2", comparator.Eq, "2"}})
		require.NoError(t, err)
		assert.Equal(t, 2, def.Len())

		def, err = def.Add(model.SpanQuery("2"), []model.Assertion{{"2", comparator.Eq, "2"}})
		assert.ErrorContains(t, err, "selector already exists")
		assert.Equal(t, 0, def.Len())

	})

	generateDef := func() model.Definition {
		def := model.Definition{}

		def, _ = def.Add(model.SpanQuery("1"), []model.Assertion{{"1", comparator.Eq, "1"}})
		def, _ = def.Add(model.SpanQuery("2"), []model.Assertion{{"2", comparator.Eq, "2"}})

		return def
	}

	t.Run("Map", func(t *testing.T) {

		def := generateDef()

		expected := map[string][]model.Assertion{
			"1": {{"1", comparator.Eq, "1"}},
			"2": {{"2", comparator.Eq, "2"}},
		}

		actual := map[string][]model.Assertion{}
		def.Map(func(spanQuery model.SpanQuery, asserts []model.Assertion) {
			actual[string(spanQuery)] = asserts
		})

		assert.Equal(t, expected, actual)

	})

	t.Run("Get", func(t *testing.T) {

		def := generateDef()

		expected := []model.Assertion{{"1", comparator.Eq, "1"}}
		actual := def.Get(model.SpanQuery("1"))

		assert.Equal(t, expected, actual)

		assert.Empty(t, def.Get(model.SpanQuery("3")))

	})

	t.Run("JSON", func(t *testing.T) {

		def := generateDef()

		encoded, err := json.Marshal(def)
		require.NoError(t, err)

		decoded := model.Definition{}
		err = json.Unmarshal(encoded, &decoded)
		require.NoError(t, err)

		assert.Equal(t, def, decoded)
	})

}
