package formatters

import (
	"fmt"
	"strconv"

	"github.com/alexeyco/simpletable"
	"github.com/goccy/go-yaml"
	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
)

type TestFormatter struct {
	serverBaseURL string
}

var _ ResourceFormatter = TestFormatter{}

func NewTestsFormatter(serverBaseURL string) TestFormatter {
	return TestFormatter{
		serverBaseURL: serverBaseURL,
	}
}

// ToListStruct implements ResourceFormatter
func (f TestFormatter) ToListStruct(file *file.File) ([]interface{}, error) {
	var testResourceList openapi.TestResourceList
	err := yaml.UnmarshalWithOptions([]byte(file.Contents()), &testResourceList, yaml.CustomUnmarshaler(nullableTimeUnmarshaller))
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, len(testResourceList.Items))
	for i, item := range testResourceList.Items {
		items[i] = item
	}

	return items, nil
}

// ToListTable implements ResourceFormatter
func (f TestFormatter) ToListTable(file *file.File) (*simpletable.Header, *simpletable.Body, error) {
	rawTestList, err := f.ToListStruct(file)
	if err != nil {
		return nil, nil, err
	}

	body := simpletable.Body{}
	for _, rawTest := range rawTestList {
		testResource := rawTest.(openapi.TestResource)
		row, err := f.getTableRow(testResource)
		if err != nil {
			return nil, nil, err
		}

		body.Cells = append(body.Cells, row)
	}

	return f.getTableHeader(), &body, nil
}

// ToStruct implements ResourceFormatter
func (TestFormatter) ToStruct(file *file.File) (interface{}, error) {
	var resource openapi.TestResource

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
func (f TestFormatter) ToTable(file *file.File) (*simpletable.Header, *simpletable.Body, error) {
	rawTest, err := f.ToStruct(file)
	if err != nil {
		return nil, nil, err
	}

	testResource := rawTest.(openapi.TestResource)
	row, err := f.getTableRow(testResource)
	if err != nil {
		return nil, nil, err
	}

	body := simpletable.Body{}
	body.Cells = [][]*simpletable.Cell{row}

	return f.getTableHeader(), &body, nil
}

func (f TestFormatter) getTableHeader() *simpletable.Header {
	return &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "ID"},
			{Text: "NAME"},
			{Text: "VERSION"},
			{Text: "TRIGGER TYPE"},
			{Text: "RUNS"},
			{Text: "LAST RUN TIME"},
			{Text: "LAST RUN SUCCESSES"},
			{Text: "LAST RUN FAILURES"},
			{Text: "URL"},
		},
	}
}

func (f TestFormatter) getTableRow(t openapi.TestResource) ([]*simpletable.Cell, error) {

	version := 1
	if v := t.Spec.GetVersion(); v > 1 {
		version = int(v)
	}

	lastRunTime := ""
	runs := 0
	passes := 0
	fails := 0
	if summary, ok := t.Spec.GetSummaryOk(); ok {

		runs = int(summary.GetRuns())

		if lastRun, ok := summary.GetLastRunOk(); ok {
			lastRunTime = lastRun.GetTime().String()
			passes = int(lastRun.GetPasses())
			fails = int(lastRun.GetFails())
		}
	}

	triggerType := ""
	if trigger, ok := t.Spec.GetTriggerOk(); ok {
		triggerType = trigger.GetType()
	}

	url := fmt.Sprintf("%s/test/%s", f.serverBaseURL, t.Spec.GetId())

	return []*simpletable.Cell{
		{Text: t.Spec.GetId()},
		{Text: t.Spec.GetName()},
		{Text: strconv.Itoa(version)},
		{Text: triggerType},
		{Text: strconv.Itoa(runs)},
		{Text: lastRunTime},
		{Text: strconv.Itoa(passes)},
		{Text: strconv.Itoa(fails)},
		{Text: url},
	}, nil
}

// FormatRunResult implements RunnableResourceFormatter
func (f TestFormatter) FormatRunResult(any) (string, error) {
	return "", nil
}
