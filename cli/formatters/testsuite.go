package formatters

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/pterm/pterm"
)

type testSuiteRun struct {
	baseURLFn        func() string
	colorsEnabled    bool
	testRunFormatter testRun
}

func TestSuiteRun(baseURLFn func() string, colorsEnabled bool) testSuiteRun {
	return testSuiteRun{
		baseURLFn:        baseURLFn,
		colorsEnabled:    colorsEnabled,
		testRunFormatter: TestRun(baseURLFn, colorsEnabled, WithPadding(1)),
	}
}

type TestSuiteRunOutput struct {
	HasResults bool                 `json:"-"`
	TestSuite  openapi.TestSuite    `json:"testSuite"`
	Run        openapi.TestSuiteRun `json:"testSuiteRun"`
	RunWebURL  string               `json:"testSuiteRunWebUrl"`
}

func (f testSuiteRun) Format(output TestSuiteRunOutput, format Output) string {
	output.RunWebURL = f.getRunLink(output.TestSuite.GetId(), output.Run.GetId())

	switch format {
	case JSON:
		return f.json(output)
	case Pretty, "":
		return f.pretty(output)
	}

	return ""
}

func (f testSuiteRun) json(output TestSuiteRunOutput) string {
	type stepResult struct {
		Name    string                   `json:"name"`
		Results openapi.AssertionResults `json:"results"`
	}

	type testSuiteResult struct {
		RunWebURL string       `json:"testRunWebUrl"`
		Steps     []stepResult `json:"steps"`
	}

	stepsResults := make([]stepResult, 0, len(output.Run.Steps))

	for i, step := range output.Run.Steps {
		test := output.TestSuite.FullSteps[i]
		stepsResults = append(stepsResults, stepResult{
			Name:    *test.Name,
			Results: *step.Result,
		})
	}

	result := testSuiteResult{
		RunWebURL: output.RunWebURL,
		Steps:     stepsResults,
	}

	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(fmt.Errorf("could not marshal output json: %w", err))
	}

	return string(bytes) + "\n"
}

func (f testSuiteRun) pretty(output TestSuiteRunOutput) string {
	message := fmt.Sprintf("%s %s (%s)\n", PASSED_TEST_ICON, output.TestSuite.GetName(), output.RunWebURL)
	if !output.HasResults {
		return f.getColoredText(true, message)
	}

	allStepsPassed := f.allTestSuiteStepsPassed(output)
	if !allStepsPassed {
		message = fmt.Sprintf("%s %s (%s)\n", FAILED_TEST_ICON, output.TestSuite.GetName(), output.RunWebURL)
	}

	// the test suite name + all steps
	formattedMessages := make([]string, 0, len(output.Run.Steps)+1)
	formattedMessages = append(formattedMessages, f.getColoredText(allStepsPassed, message))

	for i, testRun := range output.Run.Steps {
		test := output.TestSuite.FullSteps[i]
		message := f.testRunFormatter.pretty(TestRunOutput{
			HasResults: true,
			Test:       test,
			Run:        testRun,
			RunWebURL:  f.testRunFormatter.GetRunLink(*test.Id, *testRun.Id),
		})

		formattedMessages = append(formattedMessages, message)
	}

	return strings.Join(formattedMessages, "")
}

func (f testSuiteRun) allTestSuiteStepsPassed(output TestSuiteRunOutput) bool {
	for _, step := range output.Run.Steps {
		if !step.Result.GetAllPassed() {
			return false
		}
	}
	return true
}

func (f testSuiteRun) getColoredText(passed bool, text string) string {
	if !f.colorsEnabled {
		return text
	}

	if passed {
		return pterm.FgGreen.Sprintf(text)
	}

	return pterm.FgRed.Sprintf(text)
}

func (f testSuiteRun) getRunLink(tsID string, runID int32) string {
	return fmt.Sprintf("%s/testsuite/%s/run/%s", f.baseURLFn(), tsID, runID)
}
