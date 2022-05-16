package assertions_test

import (
	"testing"

	"github.com/kubeshop/tracetest/assertions"
	"github.com/kubeshop/tracetest/assertions/comparator"
	"github.com/kubeshop/tracetest/id"
	"github.com/kubeshop/tracetest/model"
	"github.com/kubeshop/tracetest/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssertion(t *testing.T) {

	spanID := id.NewRandGenerator().SpanID()
	cases := []struct {
		name           string
		testDef        model.Definition
		trace          traces.Trace
		expectedResult model.Results
	}{
		{
			name: "CanAssert",
			testDef: model.Definition{
				`span[service.name="Pokeshop"]`: []model.Assertion{
					{
						Attribute:  "tracetest.span.duration",
						Comparator: comparator.Eq,
						Value:      "2000",
					},
				},
			},
			trace: traces.Trace{
				RootSpan: traces.Span{
					ID: spanID,
					Attributes: traces.Attributes{
						"service.name":            "Pokeshop",
						"tracetest.span.duration": "2000",
					},
				},
			},
			expectedResult: model.Results{
				`span[service.name="Pokeshop"]`: []model.AssertionResult{
					{
						Assertion: model.Assertion{
							Attribute:  "tracetest.span.duration",
							Comparator: comparator.Eq,
							Value:      "2000",
						},
						Results: []model.SpanAssertionResult{
							{
								SpanID:        spanID,
								ObservedValue: "2000",
								CompareErr:    nil,
							},
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cl := c
			t.Parallel()

			actual := assertions.Assert(cl.testDef, cl.trace)

			for expectedSel, expectedAssertionResults := range cl.expectedResult {
				actualAssertionResults, ok := actual[expectedSel]
				assert.True(t, ok, `expected selector "%s" not found`, expectedSel)
				for i := 0; i < len(expectedAssertionResults); i++ {
					expectedAR := expectedAssertionResults[i]
					actualAR := actualAssertionResults[i]

					assert.Equal(t, expectedAR.Assertion, actualAR.Assertion)
					require.Len(t, actualAR.Results, len(expectedAR.Results))

					for i, expectedSpanRes := range expectedAR.Results {
						actualSpanRes := actualAR.Results[i]
						assert.Equal(t, expectedSpanRes.ObservedValue, actualSpanRes.ObservedValue)
						assert.Equal(t, expectedSpanRes.CompareErr, actualSpanRes.CompareErr)
					}
				}
			}

		})
	}
}
