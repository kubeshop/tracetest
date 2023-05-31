package formatters

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/pterm/pterm"
)

const (
	PASSED_TEST_ICON = "✔"
	FAILED_TEST_ICON = "✘"
)

type testRun struct {
	config        config.Config
	colorsEnabled bool
	padding       int
}

type testRunFormatterOption func(*testRun)

func WithPadding(padding int) testRunFormatterOption {
	return func(tr *testRun) {
		tr.padding = padding
	}
}

func TestRun(config config.Config, colorsEnabled bool, options ...testRunFormatterOption) testRun {
	testRun := testRun{
		config:        config,
		colorsEnabled: colorsEnabled,
	}

	for _, option := range options {
		option(&testRun)
	}

	return testRun
}

type TestRunOutput struct {
	HasResults bool            `json:"-"`
	Test       openapi.Test    `json:"test"`
	Run        openapi.TestRun `json:"testRun"`
	RunWebURL  string          `json:"testRunWebUrl"`
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
	result := struct {
		RunWebURL string                   `json:"testRunWebUrl"`
		Results   openapi.AssertionResults `json:"results"`
	}{
		RunWebURL: f.GetRunLink(output.Test.GetId(), output.Run.GetId()),
		Results:   *output.Run.Result,
	}
	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(fmt.Errorf("could not marshal output json: %w", err))
	}

	return string(bytes) + "\n"
}

func (f testRun) pretty(output TestRunOutput) string {
	if utils.RunStateIsFailed(output.Run.GetState()) {
		return f.getColoredText(false, fmt.Sprintf("%s\n%s",
			f.formatMessage("%s %s (%s)",
				FAILED_TEST_ICON,
				*output.Test.Name,
				output.RunWebURL,
			),
			f.formatMessage("\tReason: %s\n",
				*output.Run.LastErrorState,
			),
		))
	}

	if !output.HasResults {
		return f.formatSuccessfulTest(output.Test, output.Run)
	}

	if output.Run.Result.AllPassed == nil || !*output.Run.Result.AllPassed {
		return f.formatFailedTest(output.Test, output.Run)
	}

	return f.formatSuccessfulTest(output.Test, output.Run)
}

func (f testRun) formatSuccessfulTest(test openapi.Test, run openapi.TestRun) string {
	var buffer bytes.Buffer

	link := f.GetRunLink(test.GetId(), run.GetId())
	message := f.formatMessage("%s %s (%s)\n", PASSED_TEST_ICON, *test.Name, link)
	message = f.getColoredText(true, message)
	buffer.WriteString(message)

	for i, specResult := range run.Result.Results {
		if len(test.Specs) <= i {
			break // guard clause: this means that the server sent more results than specs
		}

		title := f.getTestSpecTitle(test.Specs[i].GetName(), specResult)
		message := f.formatMessage("\t%s %s\n", PASSED_TEST_ICON, title)
		message = f.getColoredText(true, message)
		buffer.WriteString(message)
	}

	return buffer.String()
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

	link := f.GetRunLink(test.GetId(), run.GetId())
	message := f.formatMessage("%s %s (%s)\n", FAILED_TEST_ICON, *test.Name, link)
	message = f.getColoredText(false, message)
	buffer.WriteString(message)
	for i, specResult := range run.Result.Results {
		results := make(map[string]spanAssertionResult, 0)
		allPassed := true

		for _, result := range specResult.Results {
			for _, spanResult := range result.SpanResults {
				// meta assertions such as tracetest.selected_spans.count don't have a spanID,
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

		if len(test.Specs) <= i {
			break // guard clause: this means that the server sent more results than specs
		}

		title := f.getTestSpecTitle(test.Specs[i].GetName(), specResult)
		icon := f.getStateIcon(allPassed)
		message := f.formatMessage("\t%s %s\n", icon, title)
		message = f.getColoredText(allPassed, message)
		buffer.WriteString(message)

		baseLink := f.GetRunLink(test.GetId(), run.GetId())

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
	message := f.formatMessage("\t\t%s #%s\n", icon, spanId)
	message = f.getColoredText(spanResult.allPassed, message)
	buffer.WriteString(message)

	for _, assertionResult := range spanResult.results {
		icon := f.getStateIcon(assertionResult.passed)
		var message string
		if assertionResult.observedValue != nil {
			message = f.formatMessage("\t\t\t%s %s (%s)", icon, assertionResult.assertion, *assertionResult.observedValue)
		} else {
			message = f.formatMessage("\t\t\t%s %s", icon, assertionResult.assertion)
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

func (f testRun) formatMessage(format string, args ...interface{}) string {
	message := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s%s", f.getPadding(), message)
}

func (f testRun) getPadding() string {
	padding := ""
	for i := 0; i < f.padding; i++ {
		padding = padding + "\t"
	}

	return padding
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

func (f testRun) GetRunLink(testID, runID string) string {
	return fmt.Sprintf("%s://%s/test/%s/run/%s/test", f.config.Scheme, f.config.Endpoint, testID, runID)
}

func (f testRun) getDeepLink(baseLink string, index int, spanID string) string {
	link := fmt.Sprintf("%s?selectedAssertion=%d", baseLink, index)
	if spanID != "meta" {
		link = fmt.Sprintf("%s&selectedSpan=%s", link, spanID)
	}

	return link
}

func (f testRun) getSelectorQuery(specResult openapi.AssertionResultsResultsInner) string {
	if hasQuery := specResult.Selector.HasQuery(); hasQuery {
		return *specResult.Selector.Query
	}

	return "All Spans"
}

func (f testRun) getTestSpecTitle(specName string, specResult openapi.AssertionResultsResultsInner) string {
	if specName != "" {
		return specName
	}

	return f.getSelectorQuery(specResult)

}
