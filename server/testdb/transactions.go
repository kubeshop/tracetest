package testdb

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
)

var _ model.TestRepository = &postgresDB{}

func (td *postgresDB) TransactionIDExists(ctx context.Context, id id.ID) (bool, error) {
	exists := false

	row := td.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) > 0 as exists FROM transactions WHERE id = $1",
		id,
	)

	err := row.Scan(&exists)

	return exists, err
}

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
	tx, err := td.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return model.Transaction{}, fmt.Errorf("sql begin: %w", err)
	}

	stmt, err := tx.Prepare(insertIntoTransactionsQuery)
	if err != nil {
		tx.Rollback()
		return model.Transaction{}, fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
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

	return td.setTransactionSteps(ctx, tx, transaction)
}

func (td *postgresDB) setTransactionSteps(ctx context.Context, tx *sql.Tx, transaction model.Transaction) (model.Transaction, error) {
	// delete existing steps
	stmt, err := tx.Prepare("DELETE FROM transaction_runs WHERE transaction_id = $1 AND transaction_version = $2")
	if err != nil {
		tx.Rollback()
		return model.Transaction{}, err
	}

	_, err = stmt.ExecContext(ctx, transaction.ID, transaction.Version)
	if err != nil {
		tx.Rollback()
		return model.Transaction{}, err
	}

	if len(transaction.Steps) == 0 {
		return transaction, tx.Commit()
	}

	values := []string{}
	for i, test := range transaction.Steps {
		stepNumber := i + 1
		values = append(
			values,
			fmt.Sprintf("('%s', %d, '%s', %d)", transaction.ID, transaction.Version, test.ID, stepNumber),
		)
	}

	sql := "INSERT INTO transaction_steps VALUES " + strings.Join(values, ", ")
	_, err = tx.ExecContext(ctx, sql)
	if err != nil {
		tx.Rollback()
		return model.Transaction{}, fmt.Errorf("cannot save transaction steps: %w", err)
	}
	return transaction, tx.Commit()
}

func (td *postgresDB) UpdateTransaction(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	if transaction.Version == 0 {
		transaction.Version = 1
	}

	oldTransaction, err := td.GetLatestTransactionVersion(ctx, transaction.ID)
	if err != nil {
		return model.Transaction{}, fmt.Errorf("could not get latest test version while updating test: %w", err)
	}

	// keep the same creation date to keep sort order
	transaction.CreatedAt = oldTransaction.CreatedAt
	transactionToUpdate := model.BumpTransactionVersionIfNeeded(oldTransaction, transaction)

	if oldTransaction.Version == transactionToUpdate.Version {
		// No change in the version, so nothing changes and it doesn't need to persist it
		return transactionToUpdate, nil
	}

	return td.insertIntoTransactions(ctx, transactionToUpdate)
}

func (td *postgresDB) DeleteTransaction(ctx context.Context, transaction model.Transaction) error {
	tx, err := td.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM transaction_steps WHERE transaction_id = $1", transaction.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM transaction_runs WHERE transaction_id = $1", transaction.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM transactions WHERE id = $1", transaction.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

const getTransactionSQL = `
	SELECT
		t.id,
		t.version,
		t.name,
		t.description,
		t.created_at,
		(SELECT COUNT(*) FROM transaction_runs tr WHERE tr.transaction_id = t.id) as total_runs,
		last_transaction_run.created_at as last_transaction_run_time,
		last_transaction_run.pass as last_test_run_pass,
		last_transaction_run.fail as last_test_run_fail
	FROM transactions t
	LEFT OUTER JOIN (
		SELECT MAX(id) as id, transaction_id FROM transaction_runs GROUP BY transaction_id
	) as ltr ON ltr.transaction_id = t.id
	LEFT OUTER JOIN
		transaction_runs last_transaction_run
	ON last_transaction_run.transaction_id = ltr.transaction_id AND last_transaction_run.id = ltr.id
`

const transactionMaxVersionQuery = `
	INNER JOIN (
		SELECT id as idx, max(version) as latest_version FROM transactions GROUP BY idx
	) as latest_transactions ON latest_transactions.idx = t.id

	WHERE t.version = latest_transactions.latest_version `

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

	transaction.Steps, err = td.getTransactionSteps(ctx, transaction)
	if err != nil {
		return model.Transaction{}, err
	}

	return transaction, nil
}

func (td *postgresDB) GetTransactionVersion(ctx context.Context, id id.ID, version int) (model.Transaction, error) {
	stmt, err := td.db.Prepare(getTransactionSQL + " WHERE t.id = $1 AND t.version = $2")
	if err != nil {
		return model.Transaction{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	transaction, err := td.readTransactionRow(ctx, stmt.QueryRowContext(ctx, id, version))
	if err != nil {
		return model.Transaction{}, err
	}

	transaction.Steps, err = td.getTransactionSteps(ctx, transaction)
	if err != nil {
		return model.Transaction{}, err
	}

	return transaction, nil
}

func (td *postgresDB) GetTransactions(ctx context.Context, take, skip int32, query, sortBy, sortDirection string) (model.List[model.Transaction], error) {
	hasSearchQuery := query != ""
	cleanSearchQuery := "%" + strings.ReplaceAll(query, " ", "%") + "%"
	params := []any{take, skip}

	sql := getTransactionSQL + transactionMaxVersionQuery

	const condition = " AND (t.name ilike $3 OR t.description ilike $3)"
	if hasSearchQuery {
		params = append(params, cleanSearchQuery)
		sql += condition
	}

	sortingFields := map[string]string{
		"created":  "t.created_at",
		"name":     "t.name",
		"last_run": "t.created_at",
	}

	sql = sortQuery(sql, sortBy, sortDirection, sortingFields)
	sql += ` LIMIT $1 OFFSET $2 `

	stmt, err := td.db.Prepare(sql)
	if err != nil {
		return model.List[model.Transaction]{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return model.List[model.Transaction]{}, err
	}

	transactions, err := td.readTransactionRows(ctx, rows)
	if err != nil {
		return model.List[model.Transaction]{}, err
	}

	count, err := td.countTransactions(ctx, condition, cleanSearchQuery)
	if err != nil {
		return model.List[model.Transaction]{}, err
	}

	return model.List[model.Transaction]{
		Items:      transactions,
		TotalCount: count,
	}, nil
}

func (td *postgresDB) countTransactions(ctx context.Context, condition string, cleanSearchQuery string) (int, error) {
	var (
		count  int
		params []any
	)

	countQuery := "SELECT COUNT(*) FROM transactions t" + transactionMaxVersionQuery
	if cleanSearchQuery != "" {
		params = []any{cleanSearchQuery}
		countQuery += strings.ReplaceAll(condition, "$3", "$1")
	}

	err := td.db.
		QueryRowContext(ctx, countQuery, params...).
		Scan(&count)

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (td *postgresDB) readTransactionRows(ctx context.Context, rows *sql.Rows) ([]model.Transaction, error) {
	transactions := []model.Transaction{}

	for rows.Next() {
		transaction, err := td.readTransactionRow(ctx, rows)
		if err != nil {
			return []model.Transaction{}, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (td *postgresDB) readTransactionRow(ctx context.Context, row scanner) (model.Transaction, error) {
	transaction := model.Transaction{}

	var (
		lastRunTime *time.Time

		pass, fail *int
	)

	err := row.Scan(
		&transaction.ID,
		&transaction.Version,
		&transaction.Name,
		&transaction.Description,
		&transaction.CreatedAt,
		&transaction.Summary.Runs,
		&lastRunTime,
		&pass,
		&fail,
	)

	switch err {
	case nil:
		if lastRunTime != nil {
			transaction.Summary.LastRun.Time = *lastRunTime
		}
		if pass != nil {
			transaction.Summary.LastRun.Passes = *pass
		}
		if fail != nil {
			transaction.Summary.LastRun.Fails = *fail
		}

		return transaction, nil
	case sql.ErrNoRows:
		return model.Transaction{}, ErrNotFound
	default:
		return model.Transaction{}, err
	}
}

func (td *postgresDB) getTransactionSteps(ctx context.Context, transaction model.Transaction) ([]model.Test, error) {
	stmt, err := td.db.Prepare(getTestSQL + testMaxVersionQuery + ` INNER JOIN transaction_steps ts ON t.id = ts.test_id
	 WHERE ts.transaction_id = $1 AND ts.transaction_version = $2 ORDER BY ts.step_number ASC`)
	if err != nil {
		return []model.Test{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, transaction.ID, transaction.Version)
	if err != nil {
		return []model.Test{}, fmt.Errorf("query context: %w", err)
	}

	steps, err := td.readTestRows(ctx, rows)
	if err != nil {
		return []model.Test{}, fmt.Errorf("read row: %w", err)
	}

	return steps, nil
}
