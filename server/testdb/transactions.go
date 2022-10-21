package testdb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
)

var _ model.TestRepository = &postgresDB{}

func (td *postgresDB) CreateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	if !transaction.HasID() {
		transaction.ID = IDGen.ID()
	}

	transaction.Version = 1
	transaction.CreatedAt = time.Now()

	return td.insertIntoTransactions(ctx, transaction)
}

const insertIntoTransactionsQuery = `
INSERT INTO transactions (
	"id",
	"version",
	"name",
	"description",
	"created_at"
) VALUES ($1, $2, $3, $4, $5)`

func (td *postgresDB) insertIntoTransactions(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	tx, err := td.db.Begin()
	if err != nil {
		return model.Transaction{}, fmt.Errorf("sql begin: %w", err)
	}

	stmt, err := tx.Prepare(insertIntoTransactionsQuery)
	if err != nil {
		tx.Rollback()
		return model.Transaction{}, fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		transaction.ID,
		transaction.Version,
		transaction.Name,
		transaction.Description,
		transaction.CreatedAt,
	)
	if err != nil {
		tx.Rollback()
		return model.Transaction{}, fmt.Errorf("sql exec: %w", err)
	}

	return td.setTransactionSteps(tx, transaction)
}

func (td *postgresDB) setTransactionSteps(tx *sql.Tx, transaction model.Transaction) (model.Transaction, error) {
	if len(transaction.Steps) == 0 {
		return transaction, tx.Commit()
	}

	values := []string{}
	for _, test := range transaction.Steps {
		values = append(
			values,
			fmt.Sprintf("('%s', %d, '%s', %d)", transaction.ID, transaction.Version, test.ID, test.Version),
		)
	}

	sql := "INSERT INTO transaction_steps VALUES " + strings.Join(values, ", ")
	_, err := tx.Exec(sql)
	if err != nil {
		tx.Rollback()
		return model.Transaction{}, fmt.Errorf("cannot save transaction steps: %w", err)
	}
	err = tx.Commit()
	return transaction, err
}

func (td *postgresDB) UpdateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	if transaction.Version == 0 {
		transaction.Version = 1
	}

	oldTest, err := td.GetLatestTransactionVersion(ctx, transaction.ID)
	if err != nil {
		return model.Transaction{}, fmt.Errorf("could not get latest test version while updating test: %w", err)
	}

	// keep the same creation date to keep sort order
	transaction.CreatedAt = oldTest.CreatedAt

	// TODO: check this
	// testToUpdate, err := model.BumpTestVersionIfNeeded(oldTest, transaction)
	testToUpdate := transaction
	if err != nil {
		return model.Transaction{}, fmt.Errorf("could not bump test version: %w", err)
	}

	if oldTest.Version == testToUpdate.Version {
		// No change in the version, so nothing changes and it doesn't need to persist it
		return testToUpdate, nil
	}

	return td.insertIntoTransactions(ctx, testToUpdate)
}

func (td *postgresDB) DeleteTransaction(ctx context.Context, transaction model.Transaction) error {
	_, err := td.db.
		ExecContext(ctx, "DELETE FROM transactions WHERE id = $1", transaction.ID)
	return err
}

const getTransactionSQL = `
	SELECT
		t.id,
		t.version,
		t.name,
		t.description,
		t.created_at
	FROM transactions t
`

func (td *postgresDB) GetLatestTransactionVersion(ctx context.Context, id id.ID) (model.Transaction, error) {
	stmt, err := td.db.Prepare(getTransactionSQL + " WHERE t.id = $1 ORDER BY t.version DESC LIMIT 1")
	if err != nil {
		return model.Transaction{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	transaction, err := td.readTransactionRow(ctx, stmt.QueryRowContext(ctx, id))
	if err != nil {
		return model.Transaction{}, err
	}

	return transaction, nil
}
func (td *postgresDB) GetTransactionVersion(_ context.Context, _ id.ID, version int) (model.Transaction, error) {
	return model.Transaction{}, fmt.Errorf("not implemented")
}
func (td *postgresDB) GetTransactions(_ context.Context, take, skip int32, query, sortBy, sortDirection string) (model.List[model.Transaction], error) {
	return model.List[model.Transaction]{}, fmt.Errorf("not implemented")
}

func (td *postgresDB) readTransactionRow(ctx context.Context, row scanner) (model.Transaction, error) {
	transaction := model.Transaction{}

	err := row.Scan(
		&transaction.ID,
		&transaction.Version,
		&transaction.Name,
		&transaction.Description,
		&transaction.CreatedAt,
	)

	switch err {
	case nil:
		return transaction, nil
	case sql.ErrNoRows:
		return model.Transaction{}, ErrNotFound
	default:
		return model.Transaction{}, err
	}
}
