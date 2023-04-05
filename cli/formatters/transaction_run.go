package formatters

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/pterm/pterm"
)

type transactionRun struct {
	config           config.Config
	colorsEnabled    bool
	testRunFormatter testRun
}

func TransactionRun(config config.Config, colorsEnabled bool) transactionRun {
	return transactionRun{
		config:           config,
		colorsEnabled:    colorsEnabled,
		testRunFormatter: TestRun(config, colorsEnabled, WithPadding(1)),
	}
}

type TransactionRunOutput struct {
	HasResults  bool                   `json:"-"`
	Transaction openapi.Transaction    `json:"transaction"`
	Run         openapi.TransactionRun `json:"transactionRun"`
	RunWebURL   string                 `json:"transactionRunWebUrl"`
}

func (f transactionRun) Format(output TransactionRunOutput) string {
	output.RunWebURL = f.getRunLink(output.Transaction.GetId(), output.Run.GetId())
	switch CurrentOutput {
	case Pretty:
		return f.pretty(output)
	case JSON:
		return f.json(output)
	}

	return ""
}

func (f transactionRun) json(output TransactionRunOutput) string {
	type stepResult struct {
		Name    string                   `json:"name"`
		Results openapi.AssertionResults `json:"results"`
	}

	type transactionResult struct {
		RunWebURL string       `json:"testRunWebUrl"`
		Steps     []stepResult `json:"steps"`
	}

	stepsResults := make([]stepResult, 0, len(output.Run.Steps))

	for i, step := range output.Run.Steps {
		test := output.Transaction.Steps[i]
		stepsResults = append(stepsResults, stepResult{
			Name:    *test.Name,
			Results: *step.Result,
		})
	}

	result := transactionResult{
		RunWebURL: output.RunWebURL,
		Steps:     stepsResults,
	}

	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		panic(fmt.Errorf("could not marshal output json: %w", err))
	}

	return string(bytes) + "\n"
}

func (f transactionRun) pretty(output TransactionRunOutput) string {
	if utils.RunStateIsFailed(output.Run.GetState()) {
		errorMessage := ""
		if len(output.Run.Steps) > 0 {
			lastStep := output.Run.Steps[len(output.Run.Steps)-1]
			lastError := lastStep.LastErrorState
			if lastError != nil {
				errorMessage = *lastError
			}
		}

		return f.getColoredText(false, fmt.Sprintf("Failed to execute transaction: %s\n", errorMessage))
	}

	if !output.HasResults {
		return ""
	}

	link := output.RunWebURL
	allStepsPassed := f.allTransactionStepsPassed(output)
	message := fmt.Sprintf("%s %s (%s)\n", PASSED_TEST_ICON, output.Transaction.GetName(), link)

	if !allStepsPassed {
		message = fmt.Sprintf("%s %s (%s)\n", FAILED_TEST_ICON, output.Transaction.GetName(), link)
	}

	// the transaction name + all steps
	formattedMessages := make([]string, 0, len(output.Run.Steps)+1)
	formattedMessages = append(formattedMessages, f.getColoredText(allStepsPassed, message))

	for i, testRun := range output.Run.Steps {
		test := output.Transaction.Steps[i]
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

func (f transactionRun) allTransactionStepsPassed(output TransactionRunOutput) bool {
	for _, step := range output.Run.Steps {
		if !step.Result.GetAllPassed() {
			return false
		}
	}
	return true
}

func (f transactionRun) getColoredText(passed bool, text string) string {
	if !f.colorsEnabled {
		return text
	}

	if passed {
		return pterm.FgGreen.Sprintf(text)
	}

	return pterm.FgRed.Sprintf(text)
}

func (f transactionRun) getRunLink(transactionID, runID string) string {
	return fmt.Sprintf("%s://%s/transaction/%s/run/%s", f.config.Scheme, f.config.Endpoint, transactionID, runID)
}
