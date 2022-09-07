package mappings_test

import (
	"encoding/json"
	"testing"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/http/mappings"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/require"
)

func TestSpecOrder(t *testing.T) {
	input := openapi.TestSpecs{
		Specs: []openapi.TestSpecsSpecs{
			{
				Selector: openapi.Selector{
					Query: "selector 1",
				},
				Assertions: []string{
					"attr1 = 1",
					"attr2 = 2",
				},
			},
			{
				Selector: openapi.Selector{
					Query: "selector 2",
				},
				Assertions: []string{
					"attr3 = 3",
					"attr4 = 4",
				},
			},
		},
	}

	expectedJSON := `{
		"specs": [{
				"selector": {
					"query": "selector 1"
				},
				"assertions": ["attr1 = 1", "attr2 = 2"]
			},
			{
				"selector": {
					"query": "selector 2"
				},
				"assertions": ["attr3 = 3", "attr4 = 4"]
			}
		]
	}`

	// try multiple times to hit the map iteration randomization
	attempts := 50
	for i := 0; i < attempts; i++ {
		maps := mappings.New(traces.ConversionConfig{}, comparator.DefaultRegistry())
		actual := maps.Out.Specs(maps.In.Definition(input))
		actualJSON, err := json.Marshal(actual)

		require.NoError(t, err)
		// we just need this to fail once to detect a regression,
		// so we don't even care in which attempt we failed.
		// just treat this as a unique test fail
		require.JSONEq(t, expectedJSON, string(actualJSON))
	}
}

func TestResultsOrder(t *testing.T) {
	input := openapi.AssertionResults{
		Results: []openapi.AssertionResultsResults{
			{
				Selector: openapi.Selector{
					Query: "selector 1",
				},
				Results: []openapi.AssertionResult{
					{
						Assertion: "attr1 = 1",
					},
					{
						Assertion: "attr2 = 2",
					},
				},
			},
			{
				Selector: openapi.Selector{
					Query: "selector 2",
				},
				Results: []openapi.AssertionResult{
					{
						Assertion: "attr3 = 3",
					},
					{
						Assertion: "attr4 = 4",
					},
				},
			},
		},
	}

	expectedJSON := `{
		"results": [
			{
				"selector": {
					"query": "selector 1"
				},
				"results": [
					{
						"assertion": "attr1 = 1"
					},
					{
						"assertion": "attr2 = 2"
					}
				]
			},
			{
				"selector": {
					"query": "selector 2"
				},
				"results": [
					{
						"assertion": "attr3 = 3"
					},
					{
						"assertion": "attr4 = 4"
					}
				]
			}
		]
	}`

	// try multiple times to hit the map iteration randomization
	attempts := 50
	for i := 0; i < attempts; i++ {
		maps := mappings.New(traces.ConversionConfig{}, comparator.DefaultRegistry())

		actual := maps.Out.Result(maps.In.Result(input))
		actualJSON, err := json.Marshal(actual)

		require.NoError(t, err)
		// we just need this to fail once to detect a regression,
		// so we don't even care in which attempt we failed.
		// just treat this as a unique test fail
		require.JSONEq(t, expectedJSON, string(actualJSON))
	}
}
