package formatters

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
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
	bytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		panic(fmt.Errorf("could not marshal output json: %w", err))
	}

	return string(bytes) + "\n"
}

func (f transactionRun) pretty(output TransactionRunOutput) string {
	if output.Run.GetState() == "FAILED" {
		return f.getColoredText(false, "Failed to execute transaction")
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
		if step.Result.AllPassed == nil || !*step.Result.AllPassed {
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
