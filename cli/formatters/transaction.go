package formatters

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/alexeyco/simpletable"
	"github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
)

type TransactionFormatter struct{}

var _ ResourceFormatter = TransactionFormatter{}

func NewTransactionsFormatter() TransactionFormatter {
	return TransactionFormatter{}
}

// ToListStruct implements ResourceFormatter
func (f TransactionFormatter) ToListStruct(file *file.File) ([]interface{}, error) {
	var transactionResourceList openapi.TransactionResourceList
	nullableList := openapi.NewNullableTransactionResourceList(&transactionResourceList)

	err := nullableList.UnmarshalJSON([]byte(file.Contents()))
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, len(transactionResourceList.Items))
	for i, item := range transactionResourceList.Items {
		items[i] = item
	}

	return items, nil
}

// ToListTable implements ResourceFormatter
func (f TransactionFormatter) ToListTable(file *file.File) (*simpletable.Header, *simpletable.Body, error) {
	rawTransactionList, err := f.ToListStruct(file)
	if err != nil {
		return nil, nil, err
	}

	body := simpletable.Body{}
	for _, rawDemo := range rawTransactionList {
		transactionResource := rawDemo.(openapi.TransactionResource)
		row, err := f.getTableRow(transactionResource)
		if err != nil {
			return nil, nil, err
		}

		body.Cells = append(body.Cells, row)
	}

	return f.getTableHeader(), &body, nil
}

// ToStruct implements ResourceFormatter
func (TransactionFormatter) ToStruct(file *file.File) (interface{}, error) {
	var resource openapi.TransactionResource

	err := yaml.Unmarshal([]byte(file.Contents()), &resource)
	if err != nil {
		// try JSON
		err = json.Unmarshal([]byte(file.Contents()), &resource)
		if err != nil {
			return nil, err
		}
	}

	return resource, nil
}

// ToTable implements ResourceFormatter
func (f TransactionFormatter) ToTable(file *file.File) (*simpletable.Header, *simpletable.Body, error) {
	rawTransaction, err := f.ToStruct(file)
	if err != nil {
		return nil, nil, err
	}

	transactionResource := rawTransaction.(openapi.TransactionResource)
	row, err := f.getTableRow(transactionResource)
	if err != nil {
		return nil, nil, err
	}

	body := simpletable.Body{}
	body.Cells = [][]*simpletable.Cell{row}

	return f.getTableHeader(), &body, nil
}

func (f TransactionFormatter) getTableHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "ID"},
			{Text: "NAME"},
			{Text: "VERSION"},
			{Text: "STEPS"},
			{Text: "LAST RUN TIME"},
			{Text: "LAST RUN SUCCESSES"},
			{Text: "LAST RUN FAILURES"},
		},
	}
}

func (f TransactionFormatter) getTableRow(t openapi.TransactionResource) ([]*simpletable.Cell, error) {
	version := 1
	runs := 0
	lastRunTime := ""
	passes := 0
	fails := 0
	if t.Spec.Version != nil {
		version = int(*t.Spec.Version)
	}

	if t.Spec.Summary != nil {
		if t.Spec.Summary.Runs != nil {
			runs = int(*t.Spec.Summary.Runs)
		}

		if t.Spec.Summary.LastRun != nil {
			lastRunTime = t.Spec.Summary.LastRun.GetTime().String()
		}

		if t.Spec.Summary.LastRun != nil {
			passes = int(*t.Spec.Summary.LastRun.Passes)
			fails = int(*t.Spec.Summary.LastRun.Fails)
		}
	}

	return []*simpletable.Cell{
		{Text: *t.Spec.Id},
		{Text: *t.Spec.Name},
		{Text: strconv.Itoa(version)},
		{Text: fmt.Sprint(runs)},
		{Text: lastRunTime},
		{Text: fmt.Sprint(passes)},
		{Text: fmt.Sprint(fails)},
	}, nil
}

// FormatRunResult implements RunnableResourceFormatter
func (f TransactionFormatter) FormatRunResult(any) (string, error) {
	return "", nil
}
