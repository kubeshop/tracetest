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

func strp(in string) *string {
	return &in
}

func TestSpecOrder(t *testing.T) {
	input := []openapi.TestSpec{
		{
			Name:     "name 1",
			Selector: "selector 1",
			SelectorParsed: openapi.Selector{
				Query: "selector 1",
			},
			Assertions: []string{
				`attr:attr1 = 1`,
				`attr:attr2 = 2`,
			},
		},
		{
			Name:     "name 2",
			Selector: "selector 2",
			SelectorParsed: openapi.Selector{
				Query: "selector 2",
			},
			Assertions: []string{
				`attr:attr3 = 3`,
				`attr:attr4 = 4`,
			},
		},
	}

	expectedJSON := `[
		{
			"name": "name 1",
			"selector": "selector 1",
			"selectorParsed": {
				"query": "selector 1"
			},
			"assertions": ["attr:attr1 = 1", "attr:attr2 = 2"]
		},
		{
			"name": "name 2",
			"selector": "selector 2",
			"selectorParsed": {
				"query": "selector 2"
			},
			"assertions": ["attr:attr3 = 3", "attr:attr4 = 4"]
		}
	]`

	// try multiple times to hit the map iteration randomization
	attempts := 50
	for i := 0; i < attempts; i++ {
		maps := mappings.New(traces.ConversionConfig{}, comparator.DefaultRegistry())
		definition := maps.In.Definition(input)
		actual := maps.Out.Specs(definition)
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
		Results: []openapi.AssertionResultsResultsInner{
			{
				Selector: openapi.Selector{
					Query: "selector 1",
				},
				Results: []openapi.AssertionResult{
					{
						Assertion: `attr:attr1 = 1`,
					},
					{
						Assertion: `attr:attr2 = 2`,
					},
				},
			},
			{
				Selector: openapi.Selector{
					Query: "selector 2",
				},
				Results: []openapi.AssertionResult{
					{
						Assertion: `attr:attr3 = 3`,
					},
					{
						Assertion: `attr:attr4 = 4`,
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
						"assertion": "attr:attr1 = 1"
					},
					{
						"assertion": "attr:attr2 = 2"
					}
				]
			},
			{
				"selector": {
					"query": "selector 2"
				},
				"results": [
					{
						"assertion": "attr:attr3 = 3"
					},
					{
						"assertion": "attr:attr4 = 4"
					}
				]
			}
		]
	}`

	// try multiple times to hit the map iteration randomization
	attempts := 50
	for i := 0; i < attempts; i++ {
		maps := mappings.New(traces.ConversionConfig{}, comparator.DefaultRegistry())

		result, err := maps.In.Result(input)
		require.NoError(t, err)

		actual := maps.Out.Result(result)
		actualJSON, err := json.Marshal(actual)

		require.NoError(t, err)
		// we just need this to fail once to detect a regression,
		// so we don't even care in which attempt we failed.
		// just treat this as a unique test fail
		require.JSONEq(t, expectedJSON, string(actualJSON))
	}
}
