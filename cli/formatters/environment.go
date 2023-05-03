package formatters

import (
	"github.com/alexeyco/simpletable"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"

	"gopkg.in/yaml.v2"
)

type EnvironmentsFormatter struct{}

var _ ResourceFormatter = EnvironmentsFormatter{}

func NewEnvironmentsFormatter() EnvironmentsFormatter {
	return EnvironmentsFormatter{}
}

func (f EnvironmentsFormatter) ToTable(file *file.File) (*simpletable.Header, *simpletable.Body, error) {
	rawEnvironment, err := f.ToStruct(file)
	if err != nil {
		return nil, nil, err
	}

	environmentResource := rawEnvironment.(openapi.EnvironmentResource)
	row, err := f.getTableRow(environmentResource)
	if err != nil {
		return nil, nil, err
	}

	body := simpletable.Body{}
	body.Cells = [][]*simpletable.Cell{row}

	return f.getTableHeader(), &body, nil
}

func (f EnvironmentsFormatter) ToListTable(file *file.File) (*simpletable.Header, *simpletable.Body, error) {
	rawEnvironmentList, err := f.ToListStruct(file)
	if err != nil {
		return nil, nil, err
	}

	environmentResourceList := rawEnvironmentList.(openapi.EnvironmentResourceList)

	body := simpletable.Body{}
	for _, environmentResource := range environmentResourceList.Items {
		row, err := f.getTableRow(environmentResource)
		if err != nil {
			return nil, nil, err
		}

		body.Cells = append(body.Cells, row)
	}

	return f.getTableHeader(), &body, nil
}

func (f EnvironmentsFormatter) ToStruct(file *file.File) (interface{}, error) {
	var environmentResource openapi.EnvironmentResource

	err := yaml.Unmarshal([]byte(file.Contents()), &environmentResource)
	if err != nil {
		return nil, err
	}

	return environmentResource, nil
}

func (f EnvironmentsFormatter) ToListStruct(file *file.File) (interface{}, error) {
	var environmentResourceList openapi.EnvironmentResourceList

	err := yaml.Unmarshal([]byte(file.Contents()), &environmentResourceList)
	if err != nil {
		return nil, err
	}

	return environmentResourceList, nil
}

func (f EnvironmentsFormatter) getTableHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "ID"},
			{Text: "NAME"},
			{Text: "DESCRIPTION"},
		},
	}
}

func (f EnvironmentsFormatter) getTableRow(t openapi.EnvironmentResource) ([]*simpletable.Cell, error) {
	return []*simpletable.Cell{
		{Text: *t.Spec.Id},
		{Text: *t.Spec.Name},
		{Text: *t.Spec.Description},
	}, nil
}
