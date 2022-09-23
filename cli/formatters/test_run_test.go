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

func TestJSON(t *testing.T) {
	in := formatters.TestRunOutput{
		Test: openapi.Test{
			Id:   strp("9876543"),
			Name: strp("Testcase 1"),
		},
		Run: openapi.TestRun{
			Id:      strp("123456"),
			ShortId: strp("ShortID"),
			State:   strp("FINISHED"),
			Result: &openapi.AssertionResults{
				AllPassed: boolp(true),
			},
		},
	}

	formatter := formatters.TestRun(config.Config{
		Scheme:   "http",
		Endpoint: "localhost:11633",
	}, false)

	formatters.SetOutput(formatters.JSON)
	actual := formatter.Format(in)

	expected := `{"test":{"id":"9876543","name":"Testcase 1"},"testRun":{"id":"123456", "shortId": "ShortID","result":{"allPassed":true},"state":"FINISHED"},"testRunWebUrl":"http://localhost:11633/api/r/ShortID"}`

	assert.JSONEq(t, expected, actual)
	formatters.SetOutput(formatters.DefaultOutput)
}

func TestSuccessfulTestRunOutput(t *testing.T) {
	in := formatters.TestRunOutput{
		Test: openapi.Test{
			Id:   strp("9876543"),
			Name: strp("Testcase 1"),
		},
		Run: openapi.TestRun{
			Id:      strp("123456"),
			ShortId: strp("ShortID"),
			State:   strp("FINISHED"),
			Result: &openapi.AssertionResults{
				AllPassed: boolp(true),
			},
		},
	}
	formatter := formatters.TestRun(config.Config{
		Scheme:   "http",
		Endpoint: "localhost:11633",
	}, false)
	output := formatter.Format(in)

	assert.Equal(t, "✔ Testcase 1 (http://localhost:11633/api/r/ShortID)\n", output)
}

func TestFailingTestOutput(t *testing.T) {
	in := formatters.TestRunOutput{
		Test: openapi.Test{
			Id:   strp("9876543"),
			Name: strp("Testcase 2"),
		},
		Run: openapi.TestRun{
			Id:      strp("123456"),
			ShortId: strp("shortID"),
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
		},
	}

	formatter := formatters.TestRun(config.Config{
		Scheme:   "http",
		Endpoint: "localhost:11633",
	}, false)
	output := formatter.Format(in)
	expectedOutput := `✘ Testcase 2 (http://localhost:11633/api/r/shortID)
	✔ span[name = "my span"]
		✔ #123456
			✔ tracetest.span.duration <= 200ms (157ms)
	✘ span[name = "my other span"]
		✘ #456789
			✔ tracetest.span.duration <= 200ms (68ms)
			✘ http.status = 200 (404) (http://localhost:11633/api/r/shortID?selectedAssertion=1&spanId=456789)
`
	assert.Equal(t, expectedOutput, output)
}
