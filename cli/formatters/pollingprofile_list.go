package formatters

import (
	"encoding/json"
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
)

type pollingProfileList struct {
	config config.Config
}

func PollingProfileList(config config.Config) pollingProfileList {
	return pollingProfileList{
		config: config,
	}
}

func (f pollingProfileList) Format(pollingProfiles []openapi.PollingProfile) string {
	switch CurrentOutput {
	case Pretty:
		return f.pretty(pollingProfiles)
	case JSON:
		return f.json(pollingProfiles)
	}

	return ""
}

func (f pollingProfileList) json(pollingProfiles []openapi.PollingProfile) string {
	bytes, err := json.Marshal(pollingProfiles)
	if err != nil {
		panic(fmt.Errorf("could not marshal output json: %w", err))
	}

	return string(bytes)
}

func (f pollingProfileList) pretty(pollingProfiles []openapi.PollingProfile) string {
	if len(pollingProfiles) == 0 {
		return "No Polling Profiles"
	}

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "ID"},
			{Text: "NAME"},
			{Text: "STRATEGY"},
		},
	}

	for _, t := range pollingProfiles {
		table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
			{Text: t.Spec.Id},
			{Text: t.Spec.Name},
			{Text: t.Spec.Strategy},
		})
	}

	table.SetStyle(simpletable.StyleCompactLite)

	return table.String()
}
