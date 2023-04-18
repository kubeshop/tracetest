package model_test

import (
	"encoding/json"
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpec(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		spec := (model.Test{}).Specs

		spec, err := spec.Add(model.SpanQuery("1"), model.NamedAssertions{
			Assertions: []model.Assertion{model.Assertion(`1 = 1`)},
		})
		require.NoError(t, err)

		spec, err = spec.Add(model.SpanQuery("2"), model.NamedAssertions{
			Assertions: []model.Assertion{model.Assertion(`2 = 2`)},
		})
		require.NoError(t, err)
		assert.Equal(t, 2, spec.Len())

		spec, err = spec.Add(model.SpanQuery("2"), model.NamedAssertions{
			Assertions: []model.Assertion{model.Assertion(`2 = 2`)},
		})
		assert.ErrorContains(t, err, "selector already exists")
		assert.Equal(t, 0, spec.Len())

	})

	generateSpec := func() maps.Ordered[model.SpanQuery, model.NamedAssertions] {
		spec := (model.Test{}).Specs

		spec, _ = spec.Add(model.SpanQuery("1"), model.NamedAssertions{
			Assertions: []model.Assertion{model.Assertion(`1 = 1`)},
		})

		spec, _ = spec.Add(model.SpanQuery("2"), model.NamedAssertions{
			Assertions: []model.Assertion{model.Assertion(`2 = 2`)},
		})

		return spec
	}

	t.Run("Map", func(t *testing.T) {

		spec := generateSpec()

		expected := map[string]model.NamedAssertions{
			"1": {Assertions: []model.Assertion{model.Assertion(`1 = 1`)}},
			"2": {Assertions: []model.Assertion{model.Assertion(`2 = 2`)}},
		}

		actual := make(map[string]model.NamedAssertions)
		spec.ForEach(func(spanQuery model.SpanQuery, asserts model.NamedAssertions) error {
			actual[string(spanQuery)] = asserts
			return nil
		})

		assert.Equal(t, expected, actual)

	})

	t.Run("Get", func(t *testing.T) {

		spec := generateSpec()

		expected := model.NamedAssertions{
			Assertions: []model.Assertion{model.Assertion(`1 = 1`)},
		}
		actual := spec.Get(model.SpanQuery("1"))

		assert.Equal(t, expected, actual)

		assert.Empty(t, spec.Get(model.SpanQuery("3")))

	})

	t.Run("JSON", func(t *testing.T) {

		spec := generateSpec()

		encoded, err := json.Marshal(spec)
		require.NoError(t, err)

		decoded := maps.Ordered[model.SpanQuery, model.NamedAssertions]{}
		err = json.Unmarshal(encoded, &decoded)
		require.NoError(t, err)

		assert.Equal(t, spec, decoded)
	})

}

func TestResults(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		def := (model.RunResults{}).Results

		def, err := def.Add(model.SpanQuery("1"), []model.AssertionResult{{Assertion: model.Assertion(`1 = 1`)}})
		require.NoError(t, err)

		def, err = def.Add(model.SpanQuery("2"), []model.AssertionResult{{Assertion: model.Assertion(`2 = 2`)}})
		require.NoError(t, err)
		assert.Equal(t, 2, def.Len())

		def, err = def.Add(model.SpanQuery("2"), []model.AssertionResult{{Assertion: model.Assertion(`2 = 2`)}})
		assert.ErrorContains(t, err, "selector already exists")
		assert.Equal(t, 0, def.Len())

	})

	generateDef := func() maps.Ordered[model.SpanQuery, []model.AssertionResult] {
		def := (model.RunResults{}).Results

		def, _ = def.Add(model.SpanQuery("1"), []model.AssertionResult{{Assertion: model.Assertion(`1 = 1`)}})
		def, _ = def.Add(model.SpanQuery("2"), []model.AssertionResult{{Assertion: model.Assertion(`2 = 2`)}})

		return def
	}

	t.Run("Map", func(t *testing.T) {

		def := generateDef()

		expected := map[string][]model.AssertionResult{
			"1": {{Assertion: model.Assertion(`1 = 1`)}},
			"2": {{Assertion: model.Assertion(`2 = 2`)}},
		}

		actual := map[string][]model.AssertionResult{}
		def.ForEach(func(spanQuery model.SpanQuery, asserts []model.AssertionResult) error {
			actual[string(spanQuery)] = asserts
			return nil
		})

		assert.Equal(t, expected, actual)

	})

	t.Run("Get", func(t *testing.T) {

		def := generateDef()

		expected := []model.AssertionResult{{Assertion: model.Assertion(`1 = 1`)}}
		actual := def.Get(model.SpanQuery("1"))

		assert.Equal(t, expected, actual)

		assert.Empty(t, def.Get(model.SpanQuery("3")))

	})

	t.Run("JSON", func(t *testing.T) {

		def := generateDef()

		encoded, err := json.Marshal(def)
		require.NoError(t, err)

		decoded := maps.Ordered[model.SpanQuery, []model.AssertionResult]{}
		err = json.Unmarshal(encoded, &decoded)
		require.NoError(t, err)

		assert.Equal(t, def, decoded)
	})

}
