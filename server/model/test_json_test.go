package model_test

import (
	"encoding/json"
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/maps"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONEncoding(t *testing.T) {
	t.Parallel()

	original := model.Test{
		ID:          "ezMn7bE4g",
		Name:        "first test",
		Description: "description",
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
			HTTP: &model.HTTPRequest{
				URL: "http://localhost:3030/hello-instrumented",
			},
		},
		Specs: (maps.Ordered[model.SpanQuery, model.NamedAssertions]{}).
			MustAdd("query", model.NamedAssertions{
				Name: "some assertion",
				Assertions: []model.Assertion{
					"attr:some_attr = 1",
				},
			}).
			MustAdd("unnamed", model.NamedAssertions{ // unnamed
				Assertions: []model.Assertion{
					"attr:some_attr = 1",
				},
			}),
		Outputs: (maps.Ordered[string, model.Output]{}).
			MustAdd("output", model.Output{
				Selector: "selector",
				Value:    "value",
			}),
	}

	encoded, err := json.Marshal(original)
	require.NoError(t, err)

	decoded := model.Test{}
	err = json.Unmarshal(encoded, &decoded)
	require.NoError(t, err)

	assert.Equal(t, original, decoded)
}
