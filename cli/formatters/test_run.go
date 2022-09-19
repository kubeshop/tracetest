package formatters

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/pterm/pterm"
)

const (
	PASSED_TEST_ICON = "✔"
	FAILED_TEST_ICON = "✘"
)

type testRun struct {
	config        config.Config
	colorsEnabled bool
}

func TestRun(config config.Config, colorsEnabled bool) testRun {
	return testRun{
		config:        config,
		colorsEnabled: colorsEnabled,
	}
}

type TestRunOutput struct {
	Test      openapi.Test    `json:"test"`
	Run       openapi.TestRun `json:"testRun"`
	RunWebURL string          `json:"testRunWebUrl"`
}

func (f testRun) Format(output TestRunOutput) string {
	switch CurrentOutput {
	case Pretty:
		return f.pretty(output)
	case JSON:
		return f.json(output)
	}

	return ""
}

func (f testRun) json(output TestRunOutput) string {
	output.RunWebURL = f.getRunLink(output.Test, output.Run)
	bytes, err := json.Marshal(output)
	if err != nil {
		panic(fmt.Errorf("could not marshal output json: %w", err))
	}

	return string(bytes)

}

func (f testRun) pretty(output TestRunOutput) string {
	if output.Run.State != nil && *output.Run.State == "FAILED" {
		return f.getColoredText(false, fmt.Sprintf("Failed to execute test: %s", *output.Run.LastErrorState))
	}

	if output.Run.Result.AllPassed == nil || !*output.Run.Result.AllPassed {
		return f.formatFailedTest(output.Test, output.Run)
	}

	return f.formatSuccessfulTest(output.Test, output.Run)
}

func (f testRun) formatSuccessfulTest(test openapi.Test, run openapi.TestRun) string {
	link := f.getRunLink(test, run)
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
	index         int
	spanID        string
}

func (f testRun) formatFailedTest(test openapi.Test, run openapi.TestRun) string {
	var buffer bytes.Buffer

	link := f.getRunLink(test, run)
	message := fmt.Sprintf("%s %s (%s)\n", FAILED_TEST_ICON, *test.Name, link)
	message = f.getColoredText(false, message)
	buffer.WriteString(message)
	for _, specResult := range run.Result.Results {
		results := make(map[string]spanAssertionResult, 0)
		allPassed := true

		for i, result := range specResult.Results {
			assertionQuery := fmt.Sprintf(
				"%s %s %s",
				*result.Assertion.Attribute,
				*result.Assertion.Comparator,
				*result.Assertion.Expected,
			)

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
					assertion:     assertionQuery,
					observedValue: spanResult.ObservedValue,
					passed:        spanAssertionPassed,
					index:         i,
					spanID:        spanID,
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

		baseLink := f.getRunLink(test, run)

		if metaResult, exists := results["meta"]; exists {
			// meta assertions should be placed at the top
			// of the selector section. That's why we treat it as a special case
			// and remove it from the results map afterwards.

			f.generateSpanResult(&buffer, "meta", metaResult, baseLink)
			delete(results, "meta")
		}

		for spanId, spanResult := range results {
			f.generateSpanResult(&buffer, spanId, spanResult, baseLink)
		}
	}

	return buffer.String()
}

func (f testRun) generateSpanResult(buffer *bytes.Buffer, spanId string, spanResult spanAssertionResult, baseLink string) {
	icon := f.getStateIcon(spanResult.allPassed)
	message := fmt.Sprintf("\t\t%s #%s\n", icon, spanId)
	message = f.getColoredText(spanResult.allPassed, message)
	buffer.WriteString(message)

	for _, assertionResult := range spanResult.results {
		icon := f.getStateIcon(assertionResult.passed)
		var message string
		if assertionResult.observedValue != nil {
			message = fmt.Sprintf("\t\t\t%s %s (%s)", icon, assertionResult.assertion, *assertionResult.observedValue)
		} else {
			message = fmt.Sprintf("\t\t\t%s %s", icon, assertionResult.assertion)
		}

		if !assertionResult.passed {
			link := f.getDeepLink(baseLink, assertionResult.index, assertionResult.spanID)
			message = fmt.Sprintf("%s (%s)", message, link)
		}
		message += "\n"

		message = f.getColoredText(assertionResult.passed, message)

		buffer.WriteString(message)
	}
}

func (f testRun) getStateIcon(passed bool) string {
	if passed {
		return PASSED_TEST_ICON
	}

	return FAILED_TEST_ICON
}

func (f testRun) getColoredText(passed bool, text string) string {
	if !f.colorsEnabled {
		return text
	}

	if passed {
		return pterm.FgGreen.Sprintf(text)
	}

	return pterm.FgRed.Sprintf(text)
}

func (f testRun) getRunLink(test openapi.Test, run openapi.TestRun) string {
	return fmt.Sprintf("%s://%s/test/%s/run/%s", f.config.Scheme, f.config.Endpoint, *test.Id, *run.Id)
}

func (f testRun) getDeepLink(baseLink string, index int, spanID string) string {
	link := fmt.Sprintf("%s?selectedAssertion=%d", baseLink, index)
	if spanID != "meta" {
		link = fmt.Sprintf("%s&spanId=%s", link, spanID)
	}

	return link
}
