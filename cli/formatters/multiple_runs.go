package formatters

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/cli/openapi"
	"github.com/pterm/pterm"
)

type RunnerGetter[T any] func(resource any) (Runner[T], error)

type Runner[T any] interface {
	FormatResult(resource T, format string) string
}

type multipleRun[T any] struct {
	baseURLFn     func() string
	colorsEnabled bool
}

func MultipleRun[T any](baseURLFn func() string, colorsEnabled bool) multipleRun[T] {
	return multipleRun[T]{
		baseURLFn:     baseURLFn,
		colorsEnabled: colorsEnabled,
	}
}

type MultipleRunOutput[T any] struct {
	Runs         []T
	Resources    []any
	HasResults   bool
	RunGroup     openapi.RunGroup
	RunnerGetter RunnerGetter[T]
}

func (f multipleRun[T]) Format(output MultipleRunOutput[T], format Output) string {
	switch format {
	case JSON:
		return f.json(output)
	case Pretty, "":
		return f.pretty(output)
	}

	return ""
}

type jsonSummary struct {
	RunGroup openapi.RunGroup `json:"runGroup"`
	Runs     []any            `json:"runs"`
}

func (f multipleRun[T]) json(output MultipleRunOutput[T]) string {
	summary := jsonSummary{
		RunGroup: output.RunGroup,
		Runs:     make([]any, 0, len(output.Runs)),
	}

	for i, run := range output.Runs {
		resource := output.Resources[i]
		runner, _ := output.RunnerGetter(resource)
		result := runner.FormatResult(run, "json")

		var output any

		json.Unmarshal([]byte(result), &output)

		summary.Runs = append(summary.Runs, output)
	}

	bytes, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		panic(fmt.Errorf("could not marshal output json: %w", err))
	}

	return string(bytes) + "\n"
}

func (f multipleRun[T]) pretty(output MultipleRunOutput[T]) string {
	runGroupUrl := fmt.Sprintf("%s/run/%s", f.baseURLFn(), output.RunGroup.Id)
	message := fmt.Sprintf("%s - %s - %s\n", PASSED_TEST_ICON, *output.RunGroup.Status, runGroupUrl)
	if !output.HasResults {
		return f.getColoredText(true, message)
	}

	allStepsPassed := output.RunGroup.Summary.GetFailed() == 0
	if !allStepsPassed {
		message = fmt.Sprintf("#%s - %s - %s\n", FAILED_TEST_ICON, *output.RunGroup.Status, runGroupUrl)
	}

	// the test suite name + all steps
	formattedMessages := make([]string, 0, len(output.Runs)+1)
	formattedMessages = append(formattedMessages, f.getColoredText(allStepsPassed, message))

	for i, run := range output.Runs {
		resource := output.Resources[i]
		runner, _ := output.RunnerGetter(resource)
		result := runner.FormatResult(run, "json")

		formattedMessages = append(formattedMessages, result)
	}

	return strings.Join(formattedMessages, "")
}

func (f multipleRun[T]) getColoredText(passed bool, text string) string {
	if !f.colorsEnabled {
		return text
	}

	if passed {
		return pterm.FgGreen.Sprintf(text)
	}

	return pterm.FgRed.Sprintf(text)
}
