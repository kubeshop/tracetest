package assertions_test

import (
	"testing"

	"github.com/kubeshop/tracetest/assertions"
	"github.com/kubeshop/tracetest/assertions/comparator"
	"github.com/kubeshop/tracetest/traces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssertion(t *testing.T) {

	cases := []struct {
		name           string
		testDef        assertions.TestDefinition
		trace          traces.Trace
		expectedResult assertions.TestResult
	}{
		{
			name: "CanAssert",
			testDef: assertions.TestDefinition{
				`span[service.name="Pokeshop"]`: []assertions.Assertion{
					{
						Attribute:  "tracetest.span.duration",
						Comparator: comparator.Eq,
						Value:      "2000",
					},
				},
			},
			trace: traces.Trace{
				RootSpan: traces.Span{
					Attributes: traces.Attributes{
						"service.name":            "Pokeshop",
						"tracetest.span.duration": "2000",
					},
				},
			},
			expectedResult: assertions.TestResult{
				`span[service.name="Pokeshop"]`: []assertions.AssertionResult{
					{
						Assertion: assertions.Assertion{
							Attribute:  "tracetest.span.duration",
							Comparator: comparator.Eq,
							Value:      "2000",
						},
						AssertionSpanResults: []assertions.AssertionSpanResult{
							{ActualValue: "2000", CompareErr: nil},
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

			actual := assertions.Assert(cl.trace, cl.testDef)

			for expectedSel, expectedAssertionResults := range cl.expectedResult {
				actualAssertionResults, ok := actual[expectedSel]
				assert.True(t, ok, `expected selector "%s" not found`, expectedSel)
				for i := 0; i < len(expectedAssertionResults); i++ {
					expectedAR := expectedAssertionResults[i]
					actualAR := actualAssertionResults[i]

					assert.Equal(t, expectedAR.Assertion, actualAR.Assertion)
					require.Len(t, actualAR.AssertionSpanResults, len(expectedAR.AssertionSpanResults))

					for i, expectedSpanRes := range expectedAR.AssertionSpanResults {
						actualSpanRes := actualAR.AssertionSpanResults[i]
						assert.Equal(t, expectedSpanRes.ActualValue, actualSpanRes.ActualValue)
						assert.Equal(t, expectedSpanRes.CompareErr, actualSpanRes.CompareErr)
					}
				}
			}

		})
	}
}
