package actions

import (
	"context"
	"fmt"

	"github.com/kubeshop/tracetest/cli/file"
	"github.com/kubeshop/tracetest/cli/openapi"
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
func (actions transactionActions) Apply(ctx context.Context, file file.File) (*file.File, error) {
	if file.HasID() {
		rawTransaction, err := actions.formatter.ToStruct(&file)
		if err != nil {
			return nil, fmt.Errorf("failed to parse file into transaction: %w", err)
		}

		transaction := rawTransaction.(openapi.TransactionResource)
		return actions.resourceClient.Update(ctx, file, *transaction.Spec.Id)
	}

	return actions.resourceClient.Create(ctx, file)
}

// Delete implements ResourceActions
func (actions transactionActions) Delete(ctx context.Context, id string) (string, error) {
	return "Transaction successfully deleted", actions.resourceClient.Delete(ctx, id)
}

// FileType implements ResourceActions
func (actions transactionActions) FileType() yaml.FileType {
	return yaml.FileTypeTransaction
}

// Get implements ResourceActions
func (actions transactionActions) Get(ctx context.Context, id string) (*file.File, error) {
	return actions.resourceClient.Get(ctx, id)
}

// GetID implements ResourceActions
func (actions transactionActions) GetID(file *file.File) (string, error) {
	transaction, err := actions.formatter.ToStruct(file)
	if err != nil {
		return "", fmt.Errorf("could not convert file into struct: %w", err)
	}

	return *transaction.(openapi.TransactionResource).Spec.Id, nil
}

// List implements ResourceActions
func (actions transactionActions) List(ctx context.Context, args utils.ListArgs) (*file.File, error) {
	return actions.resourceClient.List(ctx, args)
}

// Name implements ResourceActions
func (actions transactionActions) Name() string {
	return "transaction"
}
