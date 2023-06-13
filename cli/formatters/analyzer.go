package formatters

import (
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
)

type AnalyzerFormatter struct{}

var _ ResourceFormatter = AnalyzerFormatter{}

func NewAnalyzerFormatter() AnalyzerFormatter {
	return AnalyzerFormatter{}
}

func (f AnalyzerFormatter) ToTable(file *file.File) (*simpletable.Header, *simpletable.Body, error) {
	rawConfig, err := f.ToStruct(file)
	if err != nil {
		return nil, nil, err
	}

	linterResource := rawConfig.(openapi.LinterResource)
	row, err := f.getTableRow(linterResource)
	if err != nil {
		return nil, nil, err
	}

	body := simpletable.Body{}
	body.Cells = [][]*simpletable.Cell{row}

	return f.getTableHeader(), &body, nil
}

func (f AnalyzerFormatter) ToListTable(file *file.File) (*simpletable.Header, *simpletable.Body, error) {
	rawConfigList, err := f.ToListStruct(file)
	if err != nil {
		return nil, nil, err
	}

	body := simpletable.Body{}
	for _, rawConfig := range rawConfigList {
		linterResource := rawConfig.(openapi.LinterResource)
		row, err := f.getTableRow(linterResource)
		if err != nil {
			return nil, nil, err
		}

		body.Cells = append(body.Cells, row)
	}

	return f.getTableHeader(), &body, nil
}

func (f AnalyzerFormatter) ToStruct(file *file.File) (interface{}, error) {
	var linterResource openapi.LinterResource
	err := yaml.Unmarshal([]byte(file.Contents()), &linterResource)
	if err != nil {
		return nil, err
	}

	return linterResource, nil
}

func (f AnalyzerFormatter) ToListStruct(file *file.File) ([]interface{}, error) {
	var analyzerList openapi.LinterResourceList

	err := yaml.Unmarshal([]byte(file.Contents()), &analyzerList)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, len(analyzerList.Items))
	for i, item := range analyzerList.Items {
		items[i] = item
	}

	return items, nil
}

func (f AnalyzerFormatter) getTableHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "ID"},
			{Text: "NAME"},
			{Text: "ENABLED"},
			{Text: "MINIMUM SCORE"},
		},
	}
}

func (f AnalyzerFormatter) getTableRow(t openapi.LinterResource) ([]*simpletable.Cell, error) {
	return []*simpletable.Cell{
		{Text: *t.Spec.Id},
		{Text: *t.Spec.Name},
		{Text: fmt.Sprintf("%t", *t.Spec.Enabled)},
		{Text: fmt.Sprintf("%d", *t.Spec.MinimumScore)},
	}, nil
}
