package formatters

import (
	"github.com/alexeyco/simpletable"
	"github.com/kubeshop/tracetest/cli/file"
)

type TransactionFormatter struct{}

// ToListStruct implements ResourceFormatter
func (TransactionFormatter) ToListStruct(*file.File) ([]interface{}, error) {
	panic("unimplemented")
}

// ToListTable implements ResourceFormatter
func (TransactionFormatter) ToListTable(*file.File) (*simpletable.Header, *simpletable.Body, error) {
	panic("unimplemented")
}

// ToStruct implements ResourceFormatter
func (TransactionFormatter) ToStruct(*file.File) (interface{}, error) {
	panic("unimplemented")
}

// ToTable implements ResourceFormatter
func (TransactionFormatter) ToTable(*file.File) (*simpletable.Header, *simpletable.Body, error) {
	panic("unimplemented")
}

var _ ResourceFormatter = TransactionFormatter{}

func NewTransactionsFormatter() TransactionFormatter {
	return TransactionFormatter{}
}
