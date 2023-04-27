package formatters

import (
	"encoding/json"
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/kubeshop/tracetest/cli/config"
	"github.com/kubeshop/tracetest/cli/openapi"
)

type dataStoreList struct {
	config config.Config
}

func DataStoreList(config config.Config) dataStoreList {
	return dataStoreList{
		config: config,
	}
}

func (f dataStoreList) Format(dataStores []openapi.DataStore) string {
	switch CurrentOutput {
	case Pretty:
		return f.pretty(dataStores)
	case JSON:
		return f.json(dataStores)
	}

	return ""
}

func (f dataStoreList) json(dataStores []openapi.DataStore) string {
	bytes, err := json.Marshal(dataStores)
	if err != nil {
		panic(fmt.Errorf("could not marshal output json: %w", err))
	}

	return string(bytes)
}

func (f dataStoreList) pretty(dataStores []openapi.DataStore) string {
	if len(dataStores) == 0 {
		return "No tests"
	}

	table := simpletable.New()

	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "ID"},
			{Text: "NAME"},
			{Text: "DEFAULT"},
		},
	}

	for _, t := range dataStores {
		table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
			{Text: *t.Id},
			{Text: t.Name},
			{Text: f.getDefaultMark(t)},
		})
	}

	table.SetStyle(simpletable.StyleCompactLite)

	return table.String()
}

func (f dataStoreList) getDefaultMark(dataStore openapi.DataStore) string {
	if dataStore.Default != nil && *dataStore.Default {
		return "*"
	}

	return ""
}
