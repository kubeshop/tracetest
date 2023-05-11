package misc_formatters

import (
	"encoding/json"
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
)

type testsList struct {
	config config.Config
}

func TestsList(config config.Config) testsList {
	return testsList{
		config: config,
	}
}

func (f testsList) Format(tests []openapi.Test) string {
	switch CurrentOutput {
	case Pretty:
		return f.pretty(tests)
	case JSON:
		return f.json(tests)
	}

	return ""
}

func (f testsList) json(tests []openapi.Test) string {
	bytes, err := json.Marshal(tests)
	if err != nil {
		panic(fmt.Errorf("could not marshal output json: %w", err))
	}

	return string(bytes)
}

func (f testsList) pretty(tests []openapi.Test) string {
	if len(tests) == 0 {
		return "No tests"
	}

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "ID"},
			{Text: "NAME"},
			{Text: "URL"},
		},
	}

	for _, t := range tests {
		table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
			{Text: *t.Id},
			{Text: *t.Name},
			{Text: f.getTestLink(t)},
		})
	}

	table.SetStyle(simpletable.StyleCompactLite)

	return table.String()
}

func (f testsList) getTestLink(test openapi.Test) string {
	return fmt.Sprintf("%s://%s/test/%s", f.config.Scheme, f.config.Endpoint, *test.Id)
}
