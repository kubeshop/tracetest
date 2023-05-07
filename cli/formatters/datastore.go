package formatters

import (
	"github.com/alexeyco/simpletable"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
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
	rawList, err := f.ToListStruct(file)
	if err != nil {
		return nil, nil, err
	}

	body := simpletable.Body{}
	for _, raw := range rawList {
		resource := raw.(openapi.DataStoreResource)
		row, err := f.getTableRow(resource)
		if err != nil {
			return nil, nil, err
		}

		body.Cells = append(body.Cells, row)
	}

	return f.getTableHeader(), &body, nil
}

func (f DatastoreFormatter) ToStruct(file *file.File) (interface{}, error) {
	var datastoreResource openapi.DataStoreResource
	nullableDataStore := openapi.NewNullableDataStoreResource(&datastoreResource)

	err := nullableDataStore.UnmarshalJSON([]byte(file.Contents()))
	if err != nil {
		return nil, err
	}

	return datastoreResource, nil
}

func (f DatastoreFormatter) ToListStruct(file *file.File) ([]interface{}, error) {
	var dataStoreList openapi.DataStoreList
	nullableList := openapi.NewNullableDataStoreList(&dataStoreList)

	err := nullableList.UnmarshalJSON([]byte(file.Contents()))
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, len(dataStoreList.Items))
	for i, item := range dataStoreList.Items {
		items[i] = item
	}

	return items, nil
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
