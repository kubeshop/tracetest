package resourcemanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Jeffail/gabs/v2"
	"github.com/alexeyco/simpletable"
	"github.com/goccy/go-yaml"
)

type Format interface {
	BuildRequest(req *http.Request, verb Verb) error
	ContentType() string
	Format(data string, opts ...any) (string, error)
	Unmarshal(data []byte, v interface{}) error
	ToJSON([]byte) ([]byte, error)
	String() string
}

type formatRegistry []Format

func (f formatRegistry) GetWithFallback(format, fallback string) (Format, error) {
	if format == "" && fallback == "" {
		return nil, fmt.Errorf("format and fallback cannot be empty at the same time")
	}

	if format == "" {
		format = fallback
	}

	actualFormat := f.Get(format)
	if actualFormat == nil {
		return nil, fmt.Errorf("unknown format '%s'", format)
	}

	return actualFormat, nil
}

func (f formatRegistry) Get(format string) Format {
	for _, fr := range f {
		if fr.String() == format {
			return fr
		}
	}

	return nil
}

var Formats = formatRegistry{
	jsonFormat{},
	yamlFormat{},
	prettyFormat{},
}

const FormatJSON = "json"

type jsonFormat struct{}

func (j jsonFormat) String() string {
	return FormatJSON
}

func (j jsonFormat) BuildRequest(req *http.Request, _ Verb) error {
	req.Header.Set("Accept", j.ContentType())
	req.Header.Set("X-Tracetest-Augmented", "true")
	return nil
}

func (j jsonFormat) ContentType() string {
	return "application/json"
}

func (j jsonFormat) Format(data string, _ ...any) (string, error) {
	indented := bytes.NewBuffer([]byte{})
	err := json.Indent(indented, []byte(data), "", "  ")
	if err != nil {
		return "", err
	}

	return indented.String(), nil
}

func (j jsonFormat) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (j jsonFormat) ToJSON(in []byte) ([]byte, error) {
	return in, nil
}

const FormatYAML = "yaml"

type yamlFormat struct{}

func (y yamlFormat) String() string {
	return FormatYAML
}

func (y yamlFormat) BuildRequest(req *http.Request, verb Verb) error {
	req.Header.Set("Accept", y.ContentType())
	if verb == VerbList {
		req.Header.Set("Accept", "text/yaml-stream")
	}
	return nil
}

func (y yamlFormat) ContentType() string {
	return "text/yaml"
}

func (y yamlFormat) Format(data string, _ ...any) (string, error) {
	return data, nil
}

func (y yamlFormat) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

func (j yamlFormat) ToJSON(in []byte) ([]byte, error) {
	return yaml.YAMLToJSON(in)
}

const FormatPretty = "pretty"

type prettyFormat struct {
	jsonFormat
}

func (p prettyFormat) String() string {
	return FormatPretty
}

// Format formats data into table using given mappings.
// mappings is required to be a TableConfig.
// The path is a dot-separated list of keys, e.g. "metadata.name". See github.com/Jeffail/gabs.
func (p prettyFormat) Format(data string, opts ...any) (string, error) {
	// we expect only one option - TableConfig
	if len(opts) != 2 {
		return "", fmt.Errorf("expected 2 options, got %d", len(opts))
	}

	tableConfig, ok := opts[0].(TableConfig)
	if !ok {
		return "", fmt.Errorf("expected option to be a []TableCellConfig, got %T", opts[0])
	}

	listPath := ""
	if len(opts) > 1 {
		listPath, ok = opts[1].(string)
		if !ok {
			return "", fmt.Errorf("expected option to be a string, got %T", opts[1])
		}
	}

	if listPath == "" {
		listPath = "items"
	}

	parsed, err := gabs.ParseJSON([]byte(data))
	if err != nil {
		return "", err
	}

	// iterate over given mappings and build table headers
	headers := make([]*simpletable.Cell, 0, len(tableConfig.Cells))
	for _, mapping := range tableConfig.Cells {
		headers = append(headers, &simpletable.Cell{
			Text: mapping.Header,
		})
	}

	// iterate over parsed data and build table body
	body := buildTableBody(parsed, tableConfig, listPath)

	// configure output table
	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)
	table.Header.Cells = headers
	table.Body.Cells = body

	return table.String(), nil
}

func buildTableBody(parsed *gabs.Container, tableConfig TableConfig, listPath string) [][]*simpletable.Cell {
	items := parsed.Path(listPath)
	// if items is nil, we assume that the parsed data is a single item
	if items == nil {
		body := make([][]*simpletable.Cell, 0, 1)
		body = addRowToTableBody(body, tableConfig, parsed)
		return body
	}

	body := make([][]*simpletable.Cell, 0, len(items.Children()))
	for _, item := range items.Children() {
		body = addRowToTableBody(body, tableConfig, item)
	}
	return body
}

func addRowToTableBody(body [][]*simpletable.Cell, tableConfig TableConfig, item *gabs.Container) [][]*simpletable.Cell {
	row := make([]*simpletable.Cell, 0, len(tableConfig.Cells))

	if tableConfig.ItemModifier != nil {
		tableConfig.ItemModifier(item)
	}

	for _, mapping := range tableConfig.Cells {
		value := ""
		if v := item.Path(mapping.Path).Data(); v != nil {
			value = fmt.Sprintf("%v", v)
		}

		row = append(row, &simpletable.Cell{
			Text: value,
		})
	}

	body = append(body, row)
	return body
}

// TableConfig is a configuration for prettyFormat
// Cells is a list of mappings from JSON keys to table headers. See github.com/Jeffail/gabs.
// ItemModifier is an optional function that can modify each item before it's added to the table.
type TableConfig struct {
	Cells        []TableCellConfig
	ItemModifier func(item *gabs.Container) error
}

type TableCellConfig struct {
	Header string
	Path   string
}
