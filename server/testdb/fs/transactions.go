package fs

import (
	"context"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
)

var _ model.TestRepository = &fsDB{}

func (td *fsDB) TransactionIDExists(ctx context.Context, id id.ID) (bool, error) {
	panic("not implemented")
}

func (td *fsDB) CreateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	panic("not implemented")
}

func (td *fsDB) UpdateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	panic("not implemented")
}

func (td *fsDB) DeleteTransaction(ctx context.Context, transaction model.Transaction) error {
	panic("not implemented")
}

func (td *fsDB) GetLatestTransactionVersion(ctx context.Context, id id.ID) (model.Transaction, error) {
	panic("not implemented")
}

func (td *fsDB) GetTransactionVersion(ctx context.Context, id id.ID, version int) (model.Transaction, error) {
	panic("not implemented")
}

func (td *fsDB) GetTransactions(ctx context.Context, take, skip int32, query, sortBy, sortDirection string) (model.List[model.Transaction], error) {
	return model.List[model.Transaction]{}, nil
}
