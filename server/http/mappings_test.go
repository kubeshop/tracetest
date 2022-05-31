package http_test

import (
	"encoding/json"
	"testing"

	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/http"
	"github.com/kubeshop/tracetest/server/openapi"
	"github.com/stretchr/testify/require"
)

func TestDefinitionsOrder(t *testing.T) {
	input := openapi.TestDefinition{
		Definitions: []openapi.TestDefinitionDefinitions{
			{
				Selector: "selector 1",
				Assertions: []openapi.Assertion{
					{
						Attribute:  "attr 1",
						Comparator: "=",
						Expected:   "1",
					},
					{
						Attribute:  "attr 2",
						Comparator: "=",
						Expected:   "2",
					},
				},
			},
			{
				Selector: "selector 2",
				Assertions: []openapi.Assertion{
					{
						Attribute:  "attr 3",
						Comparator: "=",
						Expected:   "3",
					},
					{
						Attribute:  "attr 4",
						Comparator: "=",
						Expected:   "4",
					},
				},
			},
		},
	}

	expectedJSON := `{
		"definitions": [{
				"selector": "selector 1",
				"assertions": [{
						"attribute": "attr 1",
						"comparator": "=",
						"expected": "1"
					},
					{
						"attribute": "attr 2",
						"comparator": "=",
						"expected": "2"
					}
				]
			},
			{
				"selector": "selector 2",
				"assertions": [{
						"attribute": "attr 3",
						"comparator": "=",
						"expected": "3"
					},
					{
						"attribute": "attr 4",
						"comparator": "=",
						"expected": "4"
					}
				]
			}
		]
	}`

	// try multiple times to hit the map iteration randomization
	attempts := 50
	for i := 0; i < attempts; i++ {
		m := http.ModelMapper{
			Comparators: comparator.DefaultRegistry(),
		}
		oapi := http.OpenAPIMapper{}

		actual := oapi.Definition(m.Definition(input))
		actualJSON, err := json.Marshal(actual)

		require.NoError(t, err)
		// we just need this to fail once to detect a regression,
		// so we don't even care in which attempt we failed.
		// just treat this as a unique test fail
		require.JSONEq(t, expectedJSON, string(actualJSON))
	}
}
