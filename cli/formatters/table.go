package formatters

import (
	"github.com/alexeyco/simpletable"
	"github.com/kubeshop/tracetest/cli/file"
)

type Table struct {
	toTableFn ToTable
}

var _ FormatterInterface = Table{}

func NewTable(toTableFn ToTable) Table {
	return Table{toTableFn}
}

func (t Table) Type() string {
	return "pretty"
}

func (t Table) Format(file *file.File) (string, error) {
	table := simpletable.New()

	header, body, err := t.toTableFn(file)
	if err != nil {
		return "", err
	}

	table.Header = header
	table.Body = body

	table.SetStyle(simpletable.StyleCompactLite)
	return table.String(), nil
}
