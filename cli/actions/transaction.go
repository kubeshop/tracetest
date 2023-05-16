package actions

import (
	"context"

	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/utils"
	"github.com/kubeshop/tracetest/server/model/yaml"
)

type transactionActions struct {
	resourceArgs
}

var _ ResourceActions = &transactionActions{}

func NewTransactionsActions(options ...ResourceArgsOption) transactionActions {
	args := NewResourceArgs(options...)

	return transactionActions{
		resourceArgs: args,
	}
}

// Apply implements ResourceActions
func (actions transactionActions) Apply(context.Context, file.File) (*file.File, error) {
	panic("unimplemented")
}

// Delete implements ResourceActions
func (actions transactionActions) Delete(context.Context, string) (string, error) {
	panic("unimplemented")
}

// FileType implements ResourceActions
func (actions transactionActions) FileType() yaml.FileType {
	return yaml.FileTypeTransaction
}

// Get implements ResourceActions
func (actions transactionActions) Get(context.Context, string) (*file.File, error) {
	panic("unimplemented")
}

// GetID implements ResourceActions
func (actions transactionActions) GetID(file *file.File) (string, error) {
	panic("unimplemented")
}

// List implements ResourceActions
func (actions transactionActions) List(context.Context, utils.ListArgs) (*file.File, error) {
	panic("unimplemented")
}

// Name implements ResourceActions
func (actions transactionActions) Name() string {
	return "transaction"
}
