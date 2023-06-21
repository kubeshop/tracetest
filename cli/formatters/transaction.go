package formatters

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

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

func testSpecUnmarshaller(t *openapi.TestSpecs, b []byte) error {
	t.Specs = make([]openapi.TestSpec, 0)
	return yaml.Unmarshal(b, &t.Specs)
}

func nullableTimeUnmarshaller(t *openapi.NullableTime, b []byte) error {
	var timeString string
	err := yaml.Unmarshal(b, &timeString)
	if err != nil {
		return err
	}

	parsedTime, err := time.Parse(time.RFC3339Nano, timeString)
	if err != nil {
		return err
	}

	*t = *openapi.NewNullableTime(&parsedTime)
	return nil
}

// ToListStruct implements ResourceFormatter
func (f TransactionFormatter) ToListStruct(file *file.File) ([]interface{}, error) {
	var transactionResourceList openapi.TransactionResourceList
	err := yaml.UnmarshalWithOptions([]byte(file.Contents()), &transactionResourceList, yaml.CustomUnmarshaler(testSpecUnmarshaller), yaml.CustomUnmarshaler(nullableTimeUnmarshaller))
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

	// this is a hack to overcome a limitation of encoding/json
	// by converting everything to YAML, we can use the custom unmarshaller without any hacks to
	// inject a Unmarshal() function to the openapi.TestSpecs struct.
	file, err := convertJSONFileIntoYAMLFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not convert JSON file into YAML: %w", err)
	}

	err = yaml.UnmarshalWithOptions([]byte(file.Contents()), &resource, yaml.CustomUnmarshaler(testSpecUnmarshaller), yaml.CustomUnmarshaler(nullableTimeUnmarshaller))
	if err != nil {
		return nil, err
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
			{Text: "RUNS"},
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

	steps := len(t.Spec.Steps)

	if t.Spec.Version != nil {
		version = int(*t.Spec.Version)
	}

	if t.Spec.Summary != nil {
		if t.Spec.Summary.Runs != nil {
			runs = int(*t.Spec.Summary.Runs)
		}

		if t.Spec.Summary.LastRun != nil && !t.Spec.Summary.LastRun.GetTime().IsZero() {
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
		{Text: fmt.Sprint(steps)},
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

func convertJSONFileIntoYAMLFile(f *file.File) (*file.File, error) {
	fileContent := f.Contents()
	if strings.HasPrefix(fileContent, "{") && strings.HasSuffix(fileContent, "}") {
		m := make(map[string]interface{}, 0)
		err := json.Unmarshal([]byte(fileContent), &m)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal JSON file: %w", err)
		}

		yamlContent, err := yaml.Marshal(m)
		if err != nil {
			return nil, fmt.Errorf("could not marshal file content to YAML: %w", err)
		}

		file, err := file.NewFromRaw(f.Path(), yamlContent)
		return &file, err
	}

	// already a YAML file
	return f, nil
}
