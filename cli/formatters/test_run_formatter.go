package formatters

import (
	"bytes"
	"fmt"

	"github.com/fatih/color"
	"github.com/kubeshop/tracetest/cli/openapi"
)

const (
	PASSED_TEST_ICON = "✔"
	FAILED_TEST_ICON = "✘"
)

func FormatTestRunOutput(test openapi.Test, run openapi.TestRun) string {
	if *run.State == "FAILED" {
		return color.RedString(fmt.Sprintf("Failed to execute test: %s", *run.LastErrorState))
	}

	if run.Result.AllPassed == nil || !*run.Result.AllPassed {
		return formatFailedTest(test, run)
	}

	return formatSuccessfulTest(test)
}

func formatSuccessfulTest(test openapi.Test) string {
	message := fmt.Sprintf("%s %s\n", PASSED_TEST_ICON, *test.Name)
	return getColoredText(true, message)
}

type spanAssertionResult struct {
	allPassed bool
	results   []assertionResult
}

type assertionResult struct {
	assertion     string
	observedValue *string
	passed        bool
}

func formatFailedTest(test openapi.Test, run openapi.TestRun) string {
	var buffer bytes.Buffer

	message := fmt.Sprintf("%s %s\n", FAILED_TEST_ICON, *test.Name)
	message = getColoredText(false, message)
	buffer.WriteString(message)
	for _, specResult := range run.Result.Results {
		results := make(map[string]spanAssertionResult, 0)
		allPassed := true

		for _, result := range specResult.Results {
			assertionQuery := fmt.Sprintf(
				"%s %s %s",
				*result.Assertion.Attribute,
				*result.Assertion.Comparator,
				*result.Assertion.Expected,
			)

			for _, spanResult := range result.SpanResults {
				spanResults, ok := results[*spanResult.SpanId]
				if !ok {
					spanResults = spanAssertionResult{
						allPassed: true,
						results:   make([]assertionResult, 0),
					}
				}

				spanAssertionPassed := spanResult.Passed != nil && *spanResult.Passed

				spanResults.results = append(spanResults.results, assertionResult{
					assertion:     assertionQuery,
					observedValue: spanResult.ObservedValue,
					passed:        spanAssertionPassed,
				})

				if !spanAssertionPassed {
					spanResults.allPassed = false
					allPassed = false
				}

				results[*spanResult.SpanId] = spanResults
			}
		}

		icon := getStateIcon(allPassed)
		message := fmt.Sprintf("\t%s %s\n", icon, *specResult.Selector.Query)
		message = getColoredText(allPassed, message)
		buffer.WriteString(message)

		for spanId, spanResult := range results {
			icon := getStateIcon(spanResult.allPassed)
			message := fmt.Sprintf("\t\t%s #%s\n", icon, spanId)
			message = getColoredText(spanResult.allPassed, message)
			buffer.WriteString(message)

			for _, assertionResult := range spanResult.results {
				icon := getStateIcon(assertionResult.passed)
				var message string
				if assertionResult.observedValue != nil {
					message = fmt.Sprintf("\t\t\t%s %s (%s)\n", icon, assertionResult.assertion, *assertionResult.observedValue)
				} else {
					message = fmt.Sprintf("\t\t\t%s %s\n", icon, assertionResult.assertion)
				}
				message = getColoredText(assertionResult.passed, message)

				buffer.WriteString(message)
			}
		}
	}

	return buffer.String()
}

func getStateIcon(passed bool) string {
	if passed {
		return PASSED_TEST_ICON
	}

	return FAILED_TEST_ICON
}

func getColoredText(passed bool, text string) string {
	if passed {
		return color.GreenString(text)
	}

	return color.RedString(text)
}
