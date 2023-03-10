package formatters

import (
	"encoding/json"
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
)

type configFormatter struct {
	config config.Config
}

func ConfigFormatter(config config.Config) configFormatter {
	return configFormatter{
		config: config,
	}
}

func (f configFormatter) Format(config openapi.ConfigurationResource) string {
	switch CurrentOutput {
	case Pretty:
		return f.pretty(config)
	case JSON:
		return f.json(config)
	}

	return ""
}

func (f configFormatter) json(config openapi.ConfigurationResource) string {
	bytes, err := json.Marshal(config)
	if err != nil {
		panic(fmt.Errorf("could not marshal output json: %w", err))
	}

	return string(bytes)
}

func (f configFormatter) pretty(config openapi.ConfigurationResource) string {
	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "ID"},
			{Text: "NAME"},
			{Text: "ANALYTICS ENABLED"},
		},
	}

	table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
		{Text: *config.Spec.Id},
		{Text: *config.Spec.Name},
		{Text: fmt.Sprintf("%t", config.Spec.AnalyticsEnabled)},
	})

	table.SetStyle(simpletable.StyleCompactLite)

	return table.String()
}
