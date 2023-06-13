package formatters

import (
	"github.com/alexeyco/simpletable"
	"github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
)

type PollingFormatter struct{}

var _ ResourceFormatter = PollingFormatter{}

func NewPollingFormatter() PollingFormatter {
	return PollingFormatter{}
}

func (f PollingFormatter) ToTable(file *file.File) (*simpletable.Header, *simpletable.Body, error) {
	rawPolling, err := f.ToStruct(file)
	if err != nil {
		return nil, nil, err
	}

	PollingResource := rawPolling.(openapi.PollingProfile)
	row, err := f.getTableRow(PollingResource)
	if err != nil {
		return nil, nil, err
	}

	body := simpletable.Body{}
	body.Cells = [][]*simpletable.Cell{row}

	return f.getTableHeader(), &body, nil
}

func (f PollingFormatter) ToListTable(file *file.File) (*simpletable.Header, *simpletable.Body, error) {
	rawPollingList, err := f.ToListStruct(file)
	if err != nil {
		return nil, nil, err
	}

	body := simpletable.Body{}
	for _, rawPolling := range rawPollingList {
		pollingResource := rawPolling.(openapi.PollingProfile)
		row, err := f.getTableRow(pollingResource)
		if err != nil {
			return nil, nil, err
		}

		body.Cells = append(body.Cells, row)
	}

	return f.getTableHeader(), &body, nil
}

func (f PollingFormatter) ToStruct(file *file.File) (interface{}, error) {
	var pollingResource openapi.PollingProfile

	err := yaml.Unmarshal([]byte(file.Contents()), &pollingResource)
	if err != nil {
		return nil, err
	}

	return pollingResource, nil
}

func (f PollingFormatter) ToListStruct(file *file.File) ([]interface{}, error) {
	var pollingProfileList openapi.PollingProfileList

	err := yaml.Unmarshal([]byte(file.Contents()), &pollingProfileList)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, len(pollingProfileList.Items))
	for i, item := range pollingProfileList.Items {
		items[i] = item
	}

	return items, nil
}

func (f PollingFormatter) getTableHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "ID"},
			{Text: "NAME"},
			{Text: "STRAGETY"},
		},
	}
}

func (f PollingFormatter) getTableRow(t openapi.PollingProfile) ([]*simpletable.Cell, error) {
	return []*simpletable.Cell{
		{Text: t.Spec.Id},
		{Text: t.Spec.Name},
		{Text: t.Spec.Strategy},
	}, nil
}
