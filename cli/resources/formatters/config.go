package resources_formatters

import (
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
)

type ConfigFormatter struct{}

var _ ResourceFormatter = ConfigFormatter{}

func NewConfigFormatter() ConfigFormatter {
	return ConfigFormatter{}
}

func (f ConfigFormatter) ToTable(file *file.File) (*simpletable.Header, *simpletable.Body, error) {
	rawConfig, err := f.ToStruct(file)
	if err != nil {
		return nil, nil, err
	}

	ConfigResource := rawConfig.(openapi.ConfigurationResource)
	row, err := f.getTableRow(ConfigResource)
	if err != nil {
		return nil, nil, err
	}

	body := simpletable.Body{}
	body.Cells = [][]*simpletable.Cell{row}

	return f.getTableHeader(), &body, nil
}

func (f ConfigFormatter) ToListTable(file *file.File) (*simpletable.Header, *simpletable.Body, error) {
	rawConfigList, err := f.ToListStruct(file)
	if err != nil {
		return nil, nil, err
	}

	body := simpletable.Body{}
	for _, rawConfig := range rawConfigList {
		configResource := rawConfig.(openapi.ConfigurationResource)
		row, err := f.getTableRow(configResource)
		if err != nil {
			return nil, nil, err
		}

		body.Cells = append(body.Cells, row)
	}

	return f.getTableHeader(), &body, nil
}

func (f ConfigFormatter) ToStruct(file *file.File) (interface{}, error) {
	var ConfigResource openapi.ConfigurationResource
	nullableConfig := openapi.NewNullableConfigurationResource(&ConfigResource)

	err := nullableConfig.UnmarshalJSON([]byte(file.Contents()))
	if err != nil {
		return nil, err
	}

	return ConfigResource, nil
}

func (f ConfigFormatter) ToListStruct(file *file.File) ([]interface{}, error) {
	var ConfigResourceList openapi.ConfigurationResourceList
	nullableList := openapi.NewNullableConfigurationResourceList(&ConfigResourceList)

	err := nullableList.UnmarshalJSON([]byte(file.Contents()))
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, len(ConfigResourceList.Items))
	for i, item := range ConfigResourceList.Items {
		items[i] = item
	}

	return items, nil
}

func (f ConfigFormatter) getTableHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "ID"},
			{Text: "NAME"},
			{Text: "ANALYTICS ENABLED"},
		},
	}
}

func (f ConfigFormatter) getTableRow(t openapi.ConfigurationResource) ([]*simpletable.Cell, error) {
	return []*simpletable.Cell{
		{Text: *t.Spec.Id},
		{Text: *t.Spec.Name},
		{Text: fmt.Sprintf("%t", t.Spec.AnalyticsEnabled)},
	}, nil
}
