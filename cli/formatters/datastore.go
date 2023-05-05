package formatters

import (
	"github.com/alexeyco/simpletable"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"

	"gopkg.in/yaml.v2"
)

type DatastoreFormatter struct{}

var _ ResourceFormatter = DatastoreFormatter{}

func NewDatastoreFormatter() DatastoreFormatter {
	return DatastoreFormatter{}
}

func (f DatastoreFormatter) ToTable(file *file.File) (*simpletable.Header, *simpletable.Body, error) {
	rawDatastore, err := f.ToStruct(file)
	if err != nil {
		return nil, nil, err
	}

	datastoreResource := rawDatastore.(openapi.DataStoreResource)
	row, err := f.getTableRow(datastoreResource)
	if err != nil {
		return nil, nil, err
	}

	body := simpletable.Body{}
	body.Cells = [][]*simpletable.Cell{row}

	return f.getTableHeader(), &body, nil
}

func (f DatastoreFormatter) ToListTable(file *file.File) (*simpletable.Header, *simpletable.Body, error) {
	return nil, nil, nil
}

func (f DatastoreFormatter) ToStruct(file *file.File) (interface{}, error) {
	var datastoreResource openapi.DataStoreResource

	err := yaml.Unmarshal([]byte(file.Contents()), &datastoreResource)
	if err != nil {
		return nil, err
	}

	return datastoreResource, nil
}

func (f DatastoreFormatter) ToListStruct(file *file.File) ([]interface{}, error) {
	return nil, nil
}

func (f DatastoreFormatter) getTableHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "ID"},
			{Text: "NAME"},
			{Text: "DEFAULT"},
		},
	}
}

func (f DatastoreFormatter) getTableRow(t openapi.DataStoreResource) ([]*simpletable.Cell, error) {
	return []*simpletable.Cell{
		{Text: *t.Spec.Id},
		{Text: t.Spec.Name},
		{Text: f.getDefaultMark(*t.Spec)},
	}, nil
}

func (f DatastoreFormatter) getDefaultMark(dataStore openapi.DataStore) string {
	if dataStore.Default != nil && *dataStore.Default {
		return "*"
	}

	return ""
}
