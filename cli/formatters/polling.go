package formatters

import (
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"

	"gopkg.in/yaml.v2"
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
	return nil, nil, nil
}

func (f PollingFormatter) ToStruct(file *file.File) (interface{}, error) {
	var PollingResource openapi.PollingProfile

	err := yaml.Unmarshal([]byte(file.Contents()), &PollingResource)
	if err != nil {
		return nil, err
	}

	return PollingResource, nil
}

func (f PollingFormatter) ToListStruct(file *file.File) (interface{}, error) {
	return nil, nil
}

func (f PollingFormatter) getTableHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "ID"},
			{Text: "NAME"},
			{Text: "STRAGETY"},
			{Text: "RETRY DELAY"},
			{Text: "TIMEOUT"},
			{Text: "DEFAULT"},
		},
	}
}

func (f PollingFormatter) getTableRow(t openapi.PollingProfile) ([]*simpletable.Cell, error) {
	return []*simpletable.Cell{
		{Text: t.Spec.Id},
		{Text: t.Spec.Name},
		{Text: t.Spec.Strategy},
		{Text: "TODO"},
		{Text: *t.Spec.Periodic.Timeout},
		{Text: fmt.Sprintf("%t", *t.Spec.Default)},
	}, nil
}
