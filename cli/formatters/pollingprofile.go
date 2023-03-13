package formatters

import (
	"encoding/json"
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
)

type pollingProfileFormatter struct {
	config config.Config
}

func PollingProfileFormatter(config config.Config) pollingProfileFormatter {
	return pollingProfileFormatter{
		config: config,
	}
}

func (f pollingProfileFormatter) Format(config openapi.PollingProfile) string {
	switch CurrentOutput {
	case Pretty:
		return f.pretty(config)
	case JSON:
		return f.json(config)
	}

	return ""
}

func (f pollingProfileFormatter) json(config openapi.PollingProfile) string {
	bytes, err := json.Marshal(config)
	if err != nil {
		panic(fmt.Errorf("could not marshal output json: %w", err))
	}

	return string(bytes)
}

func (f pollingProfileFormatter) pretty(pollingProfile openapi.PollingProfile) string {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "ID"},
			{Text: "NAME"},
			{Text: "STRATEGY"},
			{Text: "RETRY_DELAY"},
			{Text: "TIMEOUT"},
		},
	}

	table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
		{Text: pollingProfile.Spec.Id},
		{Text: pollingProfile.Spec.Name},
		{Text: pollingProfile.Spec.Strategy},
		{Text: *pollingProfile.Spec.Periodic.RetryDelay},
		{Text: *pollingProfile.Spec.Periodic.Timeout},
	})

	table.SetStyle(simpletable.StyleCompactLite)

	return table.String()
}
