package formatters_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/formatters"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/stretchr/testify/assert"
)

func strp(in string) *string {
	return &in
}

func boolp(in bool) *bool {
	return &in
}

func TestSuccessfulTestOutput(t *testing.T) {
	test := openapi.Test{
		Id:   strp("9876543"),
		Name: strp("Testcase 1"),
	}

	run := openapi.TestRun{
		Id:    strp("123456"),
		State: strp("FINISHED"),
		Result: &openapi.AssertionResults{
			AllPassed: boolp(true),
		},
	}

	formatter := formatters.NewTestRunFormatter(config.Config{
		Scheme:   "http",
		Endpoint: "localhost:8080",
	}, false)
	output := formatter.FormatTestRunOutput(test, run)

	assert.Equal(t, "✔ Testcase 1 (http://localhost:8080/test/9876543/run/123456)\n", output)
}

func TestFailingTestOutput(t *testing.T) {
	test := openapi.Test{
		Id:   strp("9876543"),
		Name: strp("Testcase 2"),
	}

	run := openapi.TestRun{
		Id: strp("123456"),
		Result: &openapi.AssertionResults{
			AllPassed: boolp(false),
			Results: []openapi.AssertionResultsResults{
				{
					Selector: &openapi.Selector{
						Query: strp(`span[name = "my span"]`),
					},
					Results: []openapi.AssertionResult{
						{
							Assertion: &openapi.Assertion{
								Attribute:  strp("tracetest.span.duration"),
								Comparator: strp("<="),
								Expected:   strp("200ms"),
							},
							AllPassed: boolp(true),
							SpanResults: []openapi.AssertionSpanResult{
								{
									SpanId:        strp("123456"),
									ObservedValue: strp("157ms"),
									Passed:        boolp(true),
									Error:         nil,
								},
							},
						},
					},
				},
				{
					Selector: &openapi.Selector{
						Query: strp(`span[name = "my other span"]`),
					},
					Results: []openapi.AssertionResult{
						{
							Assertion: &openapi.Assertion{
								Attribute:  strp("tracetest.span.duration"),
								Comparator: strp("<="),
								Expected:   strp("200ms"),
							},
							AllPassed: boolp(true),
							SpanResults: []openapi.AssertionSpanResult{
								{
									SpanId:        strp("456789"),
									ObservedValue: strp("68ms"),
									Passed:        boolp(true),
									Error:         nil,
								},
							},
						},
						{
							Assertion: &openapi.Assertion{
								Attribute:  strp("http.status"),
								Comparator: strp("="),
								Expected:   strp("200"),
							},
							AllPassed: boolp(true),
							SpanResults: []openapi.AssertionSpanResult{
								{
									SpanId:        strp("456789"),
									ObservedValue: strp("404"),
									Passed:        boolp(false),
									Error:         nil,
								},
							},
						},
					},
				},
			},
		},
	}

	formatter := formatters.NewTestRunFormatter(config.Config{
		Scheme:   "http",
		Endpoint: "localhost:8080",
	}, false)
	output := formatter.FormatTestRunOutput(test, run)
	expectedOutput := `✘ Testcase 2 (http://localhost:8080/test/9876543/run/123456)
	✔ span[name = "my span"]
		✔ #123456
			✔ tracetest.span.duration <= 200ms (157ms)
	✘ span[name = "my other span"]
		✘ #456789
			✔ tracetest.span.duration <= 200ms (68ms)
			✘ http.status = 200 (404)
`
	assert.Equal(t, expectedOutput, output)
}
