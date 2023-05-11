package misc_formatters_test

import (
	"testing"

	"github.com/kubeshop/tracetest/cli/config"
	formatters "github.com/kubeshop/tracetest/cli/misc/formatters"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	in := formatters.TestRunOutput{
		Test: openapi.Test{
			Id:   openapi.PtrString("9876543"),
			Name: openapi.PtrString("Testcase 1"),
		},
		Run: openapi.TestRun{
			Id:    openapi.PtrString("1"),
			State: openapi.PtrString("FINISHED"),
			Result: &openapi.AssertionResults{
				AllPassed: openapi.PtrBool(true),
			},
		},
	}

	formatter := formatters.TestRun(config.Config{
		Scheme:   "http",
		Endpoint: "localhost:11633",
	}, false)

	formatters.SetOutput(formatters.JSON)
	actual := formatter.Format(in)

	expected := `{"results":{"allPassed":true},"testRunWebUrl":"http://localhost:11633/test/9876543/run/1/test"}`

	assert.JSONEq(t, expected, actual)
	formatters.SetOutput(formatters.DefaultOutput)
}

func TestSuccessfulTestRunOutput(t *testing.T) {
	in := formatters.TestRunOutput{
		Test: openapi.Test{
			Id:   openapi.PtrString("9876543"),
			Name: openapi.PtrString("Testcase 1"),
		},
		Run: openapi.TestRun{
			Id:    openapi.PtrString("1"),
			State: openapi.PtrString("FINISHED"),
			Result: &openapi.AssertionResults{
				AllPassed: openapi.PtrBool(true),
			},
		},
	}
	formatter := formatters.TestRun(config.Config{
		Scheme:   "http",
		Endpoint: "localhost:11633",
	}, false)
	output := formatter.Format(in)

	assert.Equal(t, "✔ Testcase 1 (http://localhost:11633/test/9876543/run/1/test)\n", output)
}

func TestSuccessfulTestRunOutputWithResult(t *testing.T) {
	testSpecName := "Validate span duration"
	in := formatters.TestRunOutput{
		HasResults: false,
		Test: openapi.Test{
			Id:   openapi.PtrString("9876543"),
			Name: openapi.PtrString("Testcase 1"),
			Specs: &openapi.TestSpecs{
				Specs: []openapi.TestSpecsSpecsInner{
					{
						Selector: &openapi.Selector{
							Query: openapi.PtrString(`span[name = "my span"]`),
						},
						Name: *openapi.NewNullableString(&testSpecName),
					},
				},
			},
		},
		Run: openapi.TestRun{
			Id:    openapi.PtrString("1"),
			State: openapi.PtrString("FINISHED"),
			Result: &openapi.AssertionResults{
				AllPassed: openapi.PtrBool(true),
				Results: []openapi.AssertionResultsResultsInner{
					{
						Selector: &openapi.Selector{
							Query: openapi.PtrString(`span[name = "my span"]`),
						},
						Results: []openapi.AssertionResult{
							{
								Assertion: openapi.PtrString(`attr:tracetest.span.duration <= 200ms`),
								AllPassed: openapi.PtrBool(true),
								SpanResults: []openapi.AssertionSpanResult{
									{
										SpanId:        openapi.PtrString("123456"),
										ObservedValue: openapi.PtrString("157ms"),
										Passed:        openapi.PtrBool(true),
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
	expectedOutput := `✔ Testcase 1 (http://localhost:11633/test/9876543/run/1/test)
	✔ Validate span duration
`

	assert.Equal(t, expectedOutput, output)
}

func TestFailingTestOutput(t *testing.T) {
	testSpecName := "Validate span duration"
	in := formatters.TestRunOutput{
		HasResults: true,
		Test: openapi.Test{
			Id:   openapi.PtrString("9876543"),
			Name: openapi.PtrString("Testcase 2"),
			Specs: &openapi.TestSpecs{
				Specs: []openapi.TestSpecsSpecsInner{
					{
						Selector: &openapi.Selector{
							Query: openapi.PtrString(`span[name = "my span"]`),
						},
						Name: *openapi.NewNullableString(&testSpecName),
					},
					{
						Selector: &openapi.Selector{
							Query: openapi.PtrString(`span[name = "my other span"]`),
						},
					},
				},
			},
		},
		Run: openapi.TestRun{
			Id: openapi.PtrString("1"),
			Result: &openapi.AssertionResults{
				AllPassed: openapi.PtrBool(false),
				Results: []openapi.AssertionResultsResultsInner{
					{
						Selector: &openapi.Selector{
							Query: openapi.PtrString(`span[name = "my span"]`),
						},
						Results: []openapi.AssertionResult{
							{
								Assertion: openapi.PtrString(`attr:tracetest.span.duration <= 200ms`),
								AllPassed: openapi.PtrBool(true),
								SpanResults: []openapi.AssertionSpanResult{
									{
										SpanId:        openapi.PtrString("123456"),
										ObservedValue: openapi.PtrString("157ms"),
										Passed:        openapi.PtrBool(true),
										Error:         nil,
									},
								},
							},
						},
					},
					{
						Selector: &openapi.Selector{
							Query: openapi.PtrString(`span[name = "my other span"]`),
						},
						Results: []openapi.AssertionResult{
							{
								Assertion: openapi.PtrString(`attr:http.status = 200`),
								AllPassed: openapi.PtrBool(true),
								SpanResults: []openapi.AssertionSpanResult{
									{
										SpanId:        openapi.PtrString("456789"),
										ObservedValue: openapi.PtrString("404"),
										Passed:        openapi.PtrBool(false),
										Error:         nil,
									},
								},
							},
							{
								Assertion: openapi.PtrString(`attr:tracetest.span.duration <= 200ms`),
								AllPassed: openapi.PtrBool(true),
								SpanResults: []openapi.AssertionSpanResult{
									{
										SpanId:        openapi.PtrString("456789"),
										ObservedValue: openapi.PtrString("68ms"),
										Passed:        openapi.PtrBool(true),
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
	expectedOutput := `✘ Testcase 2 (http://localhost:11633/test/9876543/run/1/test)
	✔ Validate span duration
		✔ #123456
			✔ attr:tracetest.span.duration <= 200ms (157ms)
	✘ span[name = "my other span"]
		✘ #456789
			✘ attr:http.status = 200 (404) (http://localhost:11633/test/9876543/run/1/test?selectedAssertion=1&selectedSpan=456789)
			✔ attr:tracetest.span.duration <= 200ms (68ms)
`
	assert.Equal(t, expectedOutput, output)
}

func TestFailingTestOutputWithPadding(t *testing.T) {
	testSpecName := "Validate span duration"
	in := formatters.TestRunOutput{
		HasResults: true,
		Test: openapi.Test{
			Id:   openapi.PtrString("9876543"),
			Name: openapi.PtrString("Testcase 2"),
			Specs: &openapi.TestSpecs{
				Specs: []openapi.TestSpecsSpecsInner{
					{
						Selector: &openapi.Selector{
							Query: openapi.PtrString(`span[name = "my span"]`),
						},
						Name: *openapi.NewNullableString(&testSpecName),
					},
					{
						Selector: &openapi.Selector{
							Query: openapi.PtrString(`span[name = "my other span"]`),
						},
					},
				},
			},
		},
		Run: openapi.TestRun{
			Id: openapi.PtrString("1"),
			Result: &openapi.AssertionResults{
				AllPassed: openapi.PtrBool(false),
				Results: []openapi.AssertionResultsResultsInner{
					{
						Selector: &openapi.Selector{
							Query: openapi.PtrString(`span[name = "my span"]`),
						},
						Results: []openapi.AssertionResult{
							{
								Assertion: openapi.PtrString(`attr:tracetest.span.duration <= 200ms`),
								AllPassed: openapi.PtrBool(true),
								SpanResults: []openapi.AssertionSpanResult{
									{
										SpanId:        openapi.PtrString("123456"),
										ObservedValue: openapi.PtrString("157ms"),
										Passed:        openapi.PtrBool(true),
										Error:         nil,
									},
								},
							},
						},
					},
					{
						Selector: &openapi.Selector{
							Query: openapi.PtrString(`span[name = "my other span"]`),
						},
						Results: []openapi.AssertionResult{
							{
								Assertion: openapi.PtrString(`attr:http.status = 200`),
								AllPassed: openapi.PtrBool(true),
								SpanResults: []openapi.AssertionSpanResult{
									{
										SpanId:        openapi.PtrString("456789"),
										ObservedValue: openapi.PtrString("404"),
										Passed:        openapi.PtrBool(false),
										Error:         nil,
									},
								},
							},
							{
								Assertion: openapi.PtrString(`attr:tracetest.span.duration <= 200ms`),
								AllPassed: openapi.PtrBool(true),
								SpanResults: []openapi.AssertionSpanResult{
									{
										SpanId:        openapi.PtrString("456789"),
										ObservedValue: openapi.PtrString("68ms"),
										Passed:        openapi.PtrBool(true),
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
	}, false, formatters.WithPadding(1))
	output := formatter.Format(in)
	expectedOutput := `	✘ Testcase 2 (http://localhost:11633/test/9876543/run/1/test)
		✔ Validate span duration
			✔ #123456
				✔ attr:tracetest.span.duration <= 200ms (157ms)
		✘ span[name = "my other span"]
			✘ #456789
				✘ attr:http.status = 200 (404) (http://localhost:11633/test/9876543/run/1/test?selectedAssertion=1&selectedSpan=456789)
				✔ attr:tracetest.span.duration <= 200ms (68ms)
`
	assert.Equal(t, expectedOutput, output)
}
