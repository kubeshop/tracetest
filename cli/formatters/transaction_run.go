package formatters

import (
	"encoding/json"
	"fmt"

	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/pterm/pterm"
)

type transactionRun struct {
	config        config.Config
	colorsEnabled bool
}

func TransactionRun(config config.Config, colorsEnabled bool) transactionRun {
	return transactionRun{
		config:        config,
		colorsEnabled: colorsEnabled,
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

	if output.HasResults {
		link := output.RunWebURL
		message := fmt.Sprintf("%s %s (%s)\n", PASSED_TEST_ICON, output.Transaction.GetName(), link)
		return f.getColoredText(true, message)
	}

	return ""

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
