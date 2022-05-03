package assertions_test

import (
	"testing"

	"github.com/kubeshop/tracetest/assertions"
	"github.com/kubeshop/tracetest/assertions/comparator"
	"github.com/kubeshop/tracetest/traces"
	"github.com/stretchr/testify/assert"
)

func TestAssertion(t *testing.T) {

	cases := []struct {
		name           string
		testDef        assertions.TestDefinition
		trace          traces.Trace
		expectedResult assertions.TestResult
	}{
		{
			name: "simple assertion works",
			testDef: assertions.TestDefinition{
				"selector": []assertions.Assertion{
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
						"tracetest.span.duration": "2000",
					},
				},
			},
			expectedResult: assertions.TestResult{
				"selector": assertions.AssertionResult{
					Assertion: assertions.Assertion{
						Attribute:  "tracetest.span.duration",
						Comparator: comparator.Eq,
						Value:      "2000",
					},
					AssertionSpanResults: []assertions.AssertionSpanResults{
						{ActualValue: "2000", CompareErr: nil},
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

			for expectedSel, expectedAR := range cl.expectedResult {
				actualAR, ok := actual[expectedSel]
				assert.True(t, ok, "expected selector %s not found", expectedSel)
				assert.Equal(t, expectedAR.Assertion, actualAR.Assertion)
				assert.Len(t, actualAR.AssertionSpanResults, len(expectedAR.AssertionSpanResults))

				for i, expectedSpanRes := range expectedAR.AssertionSpanResults {
					actualSpanRes := actualAR.AssertionSpanResults[i]
					assert.Equal(t, expectedSpanRes.ActualValue, actualSpanRes.ActualValue)
					assert.Equal(t, expectedSpanRes.CompareErr, actualSpanRes.CompareErr)
				}
			}

		})
	}
}
