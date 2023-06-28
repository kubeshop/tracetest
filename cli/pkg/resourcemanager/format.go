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
	Format(data string, opts ...any) (string, error)
	Unmarshal(data []byte, v interface{}) error
	String() string
}

type formatRegistry []Format

func (f formatRegistry) Get(format, fallback string) (Format, error) {
	if format == "" && fallback == "" {
		return nil, fmt.Errorf("format and fallback cannot be empty at the same time")
	}

	if format == "" {
		return f.Get(fallback, "")
	}

	for _, fr := range f {
		if fr.String() == format {
			return fr, nil
		}
	}

	return nil, fmt.Errorf("format '%s' not supported", format)
}

var Formats = formatRegistry{
	jsonFormat{},
	yamlFormat{},
	prettyFormat{},
}

type jsonFormat struct{}

func (j jsonFormat) String() string {
	return "json"
}

func (j jsonFormat) BuildRequest(req *http.Request, _ Verb) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Tracetest-Augmented", "true")
	return nil
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

type yamlFormat struct{}

func (y yamlFormat) String() string {
	return "yaml"
}

func (y yamlFormat) BuildRequest(req *http.Request, verb Verb) error {
	req.Header.Set("Content-Type", "text/yaml")
	if verb == VerbList {
		req.Header.Set("Accept", "text/yaml-stream")
	}
	return nil
}

func (y yamlFormat) Format(data string, _ ...any) (string, error) {
	return data, nil
}

func (y yamlFormat) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

type prettyFormat struct {
	jsonFormat
}

func (p prettyFormat) String() string {
	return "pretty"
}

// Format formats data into table using given mappings.
// mappings is required to be a TableConfig.
// The path is a dot-separated list of keys, e.g. "metadata.name". See github.com/Jeffail/gabs.
func (p prettyFormat) Format(data string, opts ...any) (string, error) {
	// we expect only one option - TableConfig
	if len(opts) != 1 {
		return "", fmt.Errorf("expected 1 option, got %d", len(opts))
	}

	tableConfig, ok := opts[0].(TableConfig)
	if !ok {
		return "", fmt.Errorf("expected option to be a []TableCellConfig, got %T", opts[0])
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
	body := buildTableBody(parsed, tableConfig)

	// configure output table
	table := simpletable.New()
	table.SetStyle(simpletable.StyleCompactLite)
	table.Header.Cells = headers
	table.Body.Cells = body

	return table.String(), nil
}

func buildTableBody(parsed *gabs.Container, tableConfig TableConfig) [][]*simpletable.Cell {
	items := parsed.Path("items")
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
