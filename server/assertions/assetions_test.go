package assertions_test

import (
	"testing"

	"github.com/kubeshop/tracetest/server/assertions"
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
		testDef           model.OrderedMap[model.SpanQuery, model.NamedAssertions]
		trace             traces.Trace
		expectedResult    model.OrderedMap[model.SpanQuery, []model.AssertionResult]
		expectedAllPassed bool
	}{
		{
			name: "CanAssert",
			testDef: (model.OrderedMap[model.SpanQuery, model.NamedAssertions]{}).MustAdd(`span[service.name="Pokeshop"]`, model.NamedAssertions{
				Assertions: []model.Assertion{
					`attr:tracetest.span.duration = 2000ns`,
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
					Assertion: `attr:tracetest.span.duration = 2000ns`,
					Results: []model.SpanAssertionResult{
						{
							SpanID:        &spanID,
							ObservedValue: "2us",
							CompareErr:    nil,
						},
					},
				},
			}),
		},
		{
			name: "CanAssertOnSpanMatchCount",
			testDef: (model.OrderedMap[model.SpanQuery, model.NamedAssertions]{}).MustAdd(`span[service.name="Pokeshop"]`, model.NamedAssertions{
				Assertions: []model.Assertion{
					`attr:tracetest.selected_spans.count = 1`,
				},
			}).MustAdd(`span[service.name="NotExists"]`, model.NamedAssertions{
				Assertions: []model.Assertion{
					`attr:tracetest.selected_spans.count = 0`,
				}}),
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
					Assertion: `attr:tracetest.selected_spans.count = 1`,
					Results: []model.SpanAssertionResult{
						{
							SpanID:        &spanID,
							ObservedValue: "1",
							CompareErr:    nil,
						},
					},
				},
			}).MustAdd(`span[service.name="NotExists"]`, []model.AssertionResult{
				{
					Assertion: `attr:tracetest.selected_spans.count = 0`,
					Results: []model.SpanAssertionResult{
						{
							SpanID:        nil,
							ObservedValue: "0",
							CompareErr:    nil,
						},
					},
				},
			}),
		},
		// https://github.com/kubeshop/tracetest/issues/617
		{
			name: "ContainsWithJSON",
			testDef: (model.OrderedMap[model.SpanQuery, model.NamedAssertions]{}).MustAdd(`span[service.name="Pokeshop"]`, model.NamedAssertions{
				Assertions: []model.Assertion{
					`attr:http.response.body contains 52`,
					`attr:tracetest.span.duration <= 21ms`,
				}}),
			trace: traces.Trace{
				RootSpan: traces.Span{
					ID: spanID,
					Attributes: traces.Attributes{
						"service.name":            "Pokeshop",
						"http.response.body":      `{"id":52}`,
						"tracetest.span.duration": "21000000",
					},
				},
			},
			expectedAllPassed: true,
			expectedResult: (model.OrderedMap[model.SpanQuery, []model.AssertionResult]{}).MustAdd(`span[service.name="Pokeshop"]`, []model.AssertionResult{
				{
					Assertion: `attr:http.response.body contains 52`,
					Results: []model.SpanAssertionResult{
						{
							SpanID:        &spanID,
							ObservedValue: `{"id":52}`,
							CompareErr:    nil,
						},
					},
				},
				{
					Assertion: `attr:tracetest.span.duration <= 21ms`,
					Results: []model.SpanAssertionResult{
						{
							SpanID:        &spanID,
							ObservedValue: "21ms",
							CompareErr:    nil,
						},
					},
				},
			}),
		},
		// https://github.com/kubeshop/tracetest/issues/1203
		{
			name: "DurationComparison",
			testDef: (model.OrderedMap[model.SpanQuery, model.NamedAssertions]{}).MustAdd(`span[service.name="Pokeshop"]`, model.NamedAssertions{
				Assertions: []model.Assertion{
					`attr:tracetest.span.duration <= 25ms`,
				}}),
			trace: traces.Trace{
				RootSpan: traces.Span{
					ID: spanID,
					Attributes: traces.Attributes{
						"service.name":            "Pokeshop",
						"http.response.body":      `{"id":52}`,
						"tracetest.span.duration": "25187564", // 25ms
					},
				},
			},
			expectedAllPassed: true,
			expectedResult: (model.OrderedMap[model.SpanQuery, []model.AssertionResult]{}).MustAdd(`span[service.name="Pokeshop"]`, []model.AssertionResult{
				{
					Assertion: `attr:tracetest.span.duration <= 25ms`,
					Results: []model.SpanAssertionResult{
						{
							SpanID:        &spanID,
							ObservedValue: "25ms",
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
			// t.Parallel()

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
