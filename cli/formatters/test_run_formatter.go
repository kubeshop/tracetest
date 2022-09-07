package formatters

import (
	"bytes"
	"fmt"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/pterm/pterm"
)

const (
	PASSED_TEST_ICON = "✔"
	FAILED_TEST_ICON = "✘"
)

type TestRunFormatter struct {
	config        config.Config
	colorsEnabled bool
}

func NewTestRunFormatter(config config.Config, colorsEnabled bool) TestRunFormatter {
	return TestRunFormatter{
		config:        config,
		colorsEnabled: colorsEnabled,
	}
}

func (f TestRunFormatter) FormatTestRunOutput(test openapi.Test, run openapi.TestRun) string {
	if run.State != nil && *run.State == "FAILED" {
		return f.getColoredText(false, fmt.Sprintf("Failed to execute test: %s", *run.LastErrorState))
	}

	if run.Result.AllPassed == nil || !*run.Result.AllPassed {
		return f.formatFailedTest(test, run)
	}

	return f.formatSuccessfulTest(test, run)
}

func (f TestRunFormatter) formatSuccessfulTest(test openapi.Test, run openapi.TestRun) string {
	link := f.getLink(test, run)
	message := fmt.Sprintf("%s %s (%s)\n", PASSED_TEST_ICON, *test.Name, link)
	return f.getColoredText(true, message)
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

func (f TestRunFormatter) formatFailedTest(test openapi.Test, run openapi.TestRun) string {
	var buffer bytes.Buffer

	link := f.getLink(test, run)
	message := fmt.Sprintf("%s %s (%s)\n", FAILED_TEST_ICON, *test.Name, link)
	message = f.getColoredText(false, message)
	buffer.WriteString(message)
	for _, specResult := range run.Result.Results {
		results := make(map[string]spanAssertionResult, 0)
		allPassed := true

		for _, result := range specResult.Results {
			for _, spanResult := range result.SpanResults {
				// meta assertions such as tracetest.selected_spasn.count don't have a spanID
				// so they will be treated differently. To overcome them, we will place all
				// meta assertions under the "spanID = "meta"
				spanID := "meta"
				if spanResult.SpanId != nil {
					spanID = *spanResult.SpanId
				}

				spanResults, ok := results[spanID]
				if !ok {
					spanResults = spanAssertionResult{
						allPassed: true,
						results:   make([]assertionResult, 0),
					}
				}

				spanAssertionPassed := spanResult.Passed != nil && *spanResult.Passed

				spanResults.results = append(spanResults.results, assertionResult{
					assertion:     *result.Assertion,
					observedValue: spanResult.ObservedValue,
					passed:        spanAssertionPassed,
				})

				if !spanAssertionPassed {
					spanResults.allPassed = false
					allPassed = false
				}

				results[spanID] = spanResults
			}
		}

		icon := f.getStateIcon(allPassed)
		message := fmt.Sprintf("\t%s %s\n", icon, *specResult.Selector.Query)
		message = f.getColoredText(allPassed, message)
		buffer.WriteString(message)

		if metaResult, exists := results["meta"]; exists {
			// meta assertions should be placed at the top
			// of the selector section. That's why we treat it as a special case
			// and remove it from the results map afterwards.

			f.generateSpanResult(&buffer, "meta", metaResult)
			delete(results, "meta")
		}

		for spanId, spanResult := range results {
			f.generateSpanResult(&buffer, spanId, spanResult)
		}
	}

	return buffer.String()
}

func (f TestRunFormatter) generateSpanResult(buffer *bytes.Buffer, spanId string, spanResult spanAssertionResult) {
	icon := f.getStateIcon(spanResult.allPassed)
	message := fmt.Sprintf("\t\t%s #%s\n", icon, spanId)
	message = f.getColoredText(spanResult.allPassed, message)
	buffer.WriteString(message)

	for _, assertionResult := range spanResult.results {
		icon := f.getStateIcon(assertionResult.passed)
		var message string
		if assertionResult.observedValue != nil {
			message = fmt.Sprintf("\t\t\t%s %s (%s)\n", icon, assertionResult.assertion, *assertionResult.observedValue)
		} else {
			message = fmt.Sprintf("\t\t\t%s %s\n", icon, assertionResult.assertion)
		}
		message = f.getColoredText(assertionResult.passed, message)

		buffer.WriteString(message)
	}
}

func (f TestRunFormatter) getStateIcon(passed bool) string {
	if passed {
		return PASSED_TEST_ICON
	}

	return FAILED_TEST_ICON
}

func (f TestRunFormatter) getColoredText(passed bool, text string) string {
	if !f.colorsEnabled {
		return text
	}

	if passed {
		return pterm.FgGreen.Sprintf(text)
	}

	return pterm.FgRed.Sprintf(text)
}

func (f TestRunFormatter) getLink(test openapi.Test, run openapi.TestRun) string {
	return fmt.Sprintf("%s://%s/test/%s/run/%s", f.config.Scheme, f.config.Endpoint, *test.Id, *run.Id)
}
