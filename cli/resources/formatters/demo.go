package resources_formatters

import (
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
)

type DemoFormatter struct{}

var _ ResourceFormatter = DemoFormatter{}

func NewDemoFormatter() DemoFormatter {
	return DemoFormatter{}
}

func (f DemoFormatter) ToTable(file *file.File) (*simpletable.Header, *simpletable.Body, error) {
	rawDemo, err := f.ToStruct(file)
	if err != nil {
		return nil, nil, err
	}

	DemoResource := rawDemo.(openapi.Demo)
	row, err := f.getTableRow(DemoResource)
	if err != nil {
		return nil, nil, err
	}

	body := simpletable.Body{}
	body.Cells = [][]*simpletable.Cell{row}

	return f.getTableHeader(), &body, nil
}

func (f DemoFormatter) ToListTable(file *file.File) (*simpletable.Header, *simpletable.Body, error) {
	rawDemoList, err := f.ToListStruct(file)
	if err != nil {
		return nil, nil, err
	}

	body := simpletable.Body{}
	for _, rawDemo := range rawDemoList {
		demo := rawDemo.(openapi.Demo)
		row, err := f.getTableRow(demo)
		if err != nil {
			return nil, nil, err
		}

		body.Cells = append(body.Cells, row)
	}

	return f.getTableHeader(), &body, nil
}

func (f DemoFormatter) ToStruct(file *file.File) (interface{}, error) {
	var demoResource openapi.Demo
	nullableDemo := openapi.NewNullableDemo(&demoResource)

	err := nullableDemo.UnmarshalJSON([]byte(file.Contents()))
	if err != nil {
		return nil, err
	}

	return demoResource, nil
}

func (f DemoFormatter) ToListStruct(file *file.File) ([]interface{}, error) {
	var demoList openapi.DemoList
	nullableList := openapi.NewNullableDemoList(&demoList)

	err := nullableList.UnmarshalJSON([]byte(file.Contents()))
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, len(demoList.Items))
	for i, item := range demoList.Items {
		items[i] = item
	}

	return items, nil
}

func (f DemoFormatter) getTableHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "ID"},
			{Text: "NAME"},
			{Text: "TYPE"},
			{Text: "ENABLED"},
		},
	}
}

func (f DemoFormatter) getTableRow(t openapi.Demo) ([]*simpletable.Cell, error) {
	return []*simpletable.Cell{
		{Text: *t.Spec.Id},
		{Text: *t.Spec.Name},
		{Text: *t.Spec.Type},
		{Text: fmt.Sprintf("%t", t.Spec.Enabled)},
	}, nil
}
