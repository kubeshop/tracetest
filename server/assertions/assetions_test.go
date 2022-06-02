package assertions_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/assertions"
	"github.com/kubeshop/tracetest/server/assertions/comparator"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssertion(t *testing.T) {

	spanID := id.NewRandGenerator().SpanID()
	cases := []struct {
		name              string
		testDef           model.OrderedMap[model.SpanQuery, []model.Assertion]
		trace             traces.Trace
		expectedResult    model.OrderedMap[model.SpanQuery, []model.AssertionResult]
		expectedAllPassed bool
	}{
		{
			name: "CanAssert",
			testDef: (model.OrderedMap[model.SpanQuery, []model.Assertion]{}).MustAdd(`span[service.name="Pokeshop"]`, []model.Assertion{
				{
					Attribute:  "tracetest.span.duration",
					Comparator: comparator.Eq,
					Value:      "2000",
				},
			}),
			trace: traces.Trace{
				RootSpan: traces.Span{
					ID: spanID,
					Attributes: traces.Attributes{
						"service.name":            "Pokeshop",
						"tracetest.span.duration": "2000",
					},
				},
			},
			expectedAllPassed: true,
			expectedResult: (model.OrderedMap[model.SpanQuery, []model.AssertionResult]{}).MustAdd(`span[service.name="Pokeshop"]`, []model.AssertionResult{
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
			}),
		},
		// https://github.com/kubeshop/tracetest/issues/617
		{
			name: "ContainsWithJSON",
			testDef: (model.OrderedMap[model.SpanQuery, []model.Assertion]{}).MustAdd(`span[service.name="Pokeshop"]`, []model.Assertion{
				{
					Attribute:  "http.response.body",
					Comparator: comparator.Contains,
					Value:      "52",
				},
				{
					Attribute:  "tracetest.span.duration",
					Comparator: comparator.Lt,
					Value:      "2001",
				},
			}),
			trace: traces.Trace{
				RootSpan: traces.Span{
					ID: spanID,
					Attributes: traces.Attributes{
						"service.name":            "Pokeshop",
						"http.response.body":      `{"id":52}`,
						"tracetest.span.duration": "2000",
					},
				},
			},
			expectedAllPassed: true,
			expectedResult: (model.OrderedMap[model.SpanQuery, []model.AssertionResult]{}).MustAdd(`span[service.name="Pokeshop"]`, []model.AssertionResult{
				{
					Assertion: model.Assertion{
						Attribute:  "http.response.body",
						Comparator: comparator.Contains,
						Value:      "52",
					},
					Results: []model.SpanAssertionResult{
						{
							SpanID:        spanID,
							ObservedValue: `{"id":52}`,
							CompareErr:    nil,
						},
					},
				},
				{
					Assertion: model.Assertion{
						Attribute:  "tracetest.span.duration",
						Comparator: comparator.Lt,
						Value:      "2001",
					},
					Results: []model.SpanAssertionResult{
						{
							SpanID:        spanID,
							ObservedValue: "2000",
							CompareErr:    nil,
						},
					},
				},
			}),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cl := c
			t.Parallel()

			actual, allPassed := assertions.Assert(cl.testDef, cl.trace)

			assert.Equal(t, cl.expectedAllPassed, allPassed)

			cl.expectedResult.Map(func(expectedSel model.SpanQuery, expectedAssertionResults []model.AssertionResult) {
				actualAssertionResults := actual.Get(expectedSel)
				assert.NotEmpty(t, actualAssertionResults, `expected selector "%s" not found`, expectedSel)
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
			})

		})
	}
}
