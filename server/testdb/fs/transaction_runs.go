package fs

import (
	"context"
	"database/sql"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
)

func (td *fsDB) CreateTransactionRun(ctx context.Context, tr model.TransactionRun) (model.TransactionRun, error) {
	panic("CreateTransactionRun not implemented")
}

func (td *fsDB) UpdateTransactionRun(ctx context.Context, tr model.TransactionRun) error {
	panic("UpdateTransactionRun not implemented")
}

func (td *fsDB) setTransactionRunSteps(ctx context.Context, tx *sql.Tx, tr model.TransactionRun) error {
	panic("setTransactionRunSteps not implemented")
}

func (td *fsDB) DeleteTransactionRun(ctx context.Context, tr model.TransactionRun) error {
	panic("DeleteTransactionRun not implemented")
}

func (td *fsDB) GetTransactionRun(ctx context.Context, transactionID id.ID, runID int) (model.TransactionRun, error) {
	panic("GetTransactionRun not implemented")
}

func (td *fsDB) GetTransactionsRuns(ctx context.Context, transactionID id.ID, take, skip int32) ([]model.TransactionRun, error) {
	panic("GetTransactionsRuns not implemented")
}
