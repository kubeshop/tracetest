package transaction

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/sqlutil"
	"github.com/kubeshop/tracetest/server/test"
)

func NewRepository(db *sql.DB, stepRepository transactionStepRepository) *Repository {
	repo := &Repository{
		db:             db,
		stepRepository: stepRepository,
	}

	return repo
}

type transactionStepRepository interface {
	GetTransactionSteps(_ context.Context, _ id.ID, version int) ([]test.Test, error)
}

type Repository struct {
	db             *sql.DB
	stepRepository transactionStepRepository
}

// needed for test
func (r *Repository) DB() *sql.DB {
	return r.db
}

func (r *Repository) SetID(t Transaction, id id.ID) Transaction {
	t.ID = id
	return t
}

func (r *Repository) Provision(ctx context.Context, transaction Transaction) error {
	_, err := r.Create(ctx, transaction)

	return err
}

func (r *Repository) Create(ctx context.Context, transaction Transaction) (Transaction, error) {
	setVersion(&transaction, 1)
	if transaction.CreatedAt == nil {
		setCreatedAt(&transaction, time.Now())
	}

	t, err := r.insertIntoTransactions(ctx, transaction)
	if err != nil {
		return Transaction{}, err
	}

	removeNonAugmentedFields(&t)
	return t, nil
}

func (r *Repository) Update(ctx context.Context, transaction Transaction) (Transaction, error) {
	oldTransaction, err := r.GetLatestVersion(ctx, transaction.ID)
	if err != nil {
		return Transaction{}, fmt.Errorf("could not get latest test version while updating test: %w", err)
	}

	if transaction.GetVersion() == 0 {
		setVersion(&transaction, oldTransaction.GetVersion())
	}

	// keep the same creation date to keep sort order
	transaction.CreatedAt = oldTransaction.CreatedAt
	transactionToUpdate := BumpTransactionVersionIfNeeded(oldTransaction, transaction)

	if oldTransaction.GetVersion() == transactionToUpdate.GetVersion() {
		// No change in the version, so nothing changes and it doesn't need to persist it
		return transactionToUpdate, nil
	}

	t, err := r.insertIntoTransactions(ctx, transactionToUpdate)
	if err != nil {
		return Transaction{}, err
	}

	removeNonAugmentedFields(&t)
	return t, nil
}

func (r *Repository) checkIDExists(ctx context.Context, id id.ID) error {
	exists, err := r.IDExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id id.ID) error {
	if err := r.checkIDExists(ctx, id); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	query, params := sqlutil.Tenant(ctx, "DELETE FROM transaction_steps WHERE transaction_id = $1", id)
	_, err = tx.ExecContext(ctx, query, params...)
	if err != nil {
		tx.Rollback()
		return err
	}

	q, params := sqlutil.Tenant(ctx, "DELETE FROM transaction_run_steps WHERE transaction_run_id IN (SELECT id FROM transaction_runs WHERE transaction_id = $1)", id)
	_, err = tx.ExecContext(ctx, q, params...)
	if err != nil {
		tx.Rollback()
		return err
	}

	q, params = sqlutil.Tenant(ctx, "DELETE FROM transaction_runs WHERE transaction_id = $1", id)
	_, err = tx.ExecContext(ctx, q, params...)
	if err != nil {
		tx.Rollback()
		return err
	}

	q, params = sqlutil.Tenant(ctx, "DELETE FROM transactions WHERE id = $1", id)
	_, err = tx.ExecContext(ctx, q, params...)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *Repository) IDExists(ctx context.Context, id id.ID) (bool, error) {
	exists := false
	query, params := sqlutil.Tenant(ctx, "SELECT COUNT(*) > 0 as exists FROM transactions WHERE id = $1", id)
	row := r.db.QueryRowContext(ctx, query, params...)

	err := row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("cannot check id existance: %w", err)
	}

	return exists, nil
}

const getTransactionSQL = `
	SELECT
		%s
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

func (r *Repository) SortingFields() []string {
	return []string{
		"name",
		"created",
		"last_run",
	}
}

func (r *Repository) ListAugmented(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]Transaction, error) {
	return r.list(ctx, take, skip, query, sortBy, sortDirection, true)
}

func (r *Repository) List(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]Transaction, error) {
	return r.list(ctx, take, skip, query, sortBy, sortDirection, false)
}

func (r *Repository) list(ctx context.Context, take, skip int, query, sortBy, sortDirection string, augmented bool) ([]Transaction, error) {
	q, params := listQuery(ctx, querySelect(), query, []any{take, skip})

	sortingFields := map[string]string{
		"created":  "t.created_at",
		"name":     "t.name",
		"last_run": "last_transaction_run_time",
	}

	q = sqlutil.Sort(q, sortBy, sortDirection, "created", sortingFields)
	q += " LIMIT $1 OFFSET $2"

	stmt, err := r.db.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return nil, err
	}

	transactions, err := r.readRows(ctx, rows, augmented)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return transactions, nil
}

func queryCount() string {
	return fmt.Sprintf(getTransactionSQL, "COUNT(*)")
}

func querySelect() string {
	return fmt.Sprintf(getTransactionSQL, `
		t.id,
		t.version,
		t.name,
		t.description,
		t.created_at,
		(
			SELECT step_ids_query.ids FROM (
				SELECT ts.transaction_id, ts.transaction_version, string_agg(ts.test_id, ',') as ids FROM transaction_steps ts
				GROUP BY transaction_id, transaction_version
				HAVING ts.transaction_id = t.id AND ts.transaction_version = t.version
			) as step_ids_query
		) as step_ids,
		(SELECT COUNT(*) FROM transaction_runs tr WHERE tr.transaction_id = t.id) as total_runs,
		last_transaction_run.created_at as last_transaction_run_time,
		last_transaction_run.pass as last_test_run_pass,
		last_transaction_run.fail as last_test_run_fail
	`)
}

func (r *Repository) Count(ctx context.Context, query string) (int, error) {
	sql, params := listQuery(ctx, queryCount(), query, []any{})

	count := 0
	err := r.db.
		QueryRowContext(ctx, sql, params...).
		Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) GetAugmented(ctx context.Context, id id.ID) (Transaction, error) {
	return r.get(ctx, id, true)
}

func (r *Repository) Get(ctx context.Context, id id.ID) (Transaction, error) {
	return r.get(ctx, id, false)
}

func (r *Repository) get(ctx context.Context, id id.ID, augmented bool) (Transaction, error) {
	query, params := sqlutil.Tenant(ctx, querySelect()+" WHERE t.id = $1", id)
	stmt, err := r.db.Prepare(query + "ORDER BY t.version DESC LIMIT 1")
	if err != nil {
		return Transaction{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	transaction, err := r.readRow(ctx, stmt.QueryRowContext(ctx, params...), augmented)
	if err != nil {
		return Transaction{}, err
	}

	return transaction, nil
}

func (r *Repository) GetVersion(ctx context.Context, id id.ID, version int) (Transaction, error) {
	query, params := sqlutil.Tenant(ctx, querySelect()+" WHERE t.id = $1 AND t.version = $2", id, version)
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return Transaction{}, fmt.Errorf("prepare 1: %w", err)
	}
	defer stmt.Close()

	transaction, err := r.readRow(ctx, stmt.QueryRowContext(ctx, params...), true)
	if err != nil {
		return Transaction{}, err
	}

	return transaction, nil
}

func listQuery(ctx context.Context, baseSQL, query string, params []any) (string, []any) {
	paramNumber := len(params) + 1
	condition := fmt.Sprintf(" AND (t.name ilike $%d OR t.description ilike $%d)", paramNumber, paramNumber)

	sql := baseSQL + transactionMaxVersionQuery
	sql, params = sqlutil.Search(sql, condition, query, params)
	sql, params = sqlutil.Tenant(ctx, sql, params...)

	return sql, params
}

func (r *Repository) GetLatestVersion(ctx context.Context, id id.ID) (Transaction, error) {
	query, params := sqlutil.Tenant(ctx, querySelect()+" WHERE t.id = $1", id)
	stmt, err := r.db.Prepare(query + " ORDER BY t.version DESC LIMIT 1")
	if err != nil {
		return Transaction{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	transaction, err := r.readRow(ctx, stmt.QueryRowContext(ctx, params...), true)
	if err != nil {
		return Transaction{}, err
	}

	return transaction, nil
}

const insertIntoTransactionsQuery = `
INSERT INTO transactions (
	"id",
	"version",
	"name",
	"description",
	"created_at",
	"tenant_id"
) VALUES ($1, $2, $3, $4, $5, $6)`

func (r *Repository) insertIntoTransactions(ctx context.Context, transaction Transaction) (Transaction, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	defer tx.Rollback()
	if err != nil {
		return Transaction{}, fmt.Errorf("sql begin: %w", err)
	}

	stmt, err := tx.Prepare(insertIntoTransactionsQuery)
	if err != nil {
		return Transaction{}, fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()

	tenantID := sqlutil.TenantID(ctx)

	_, err = stmt.ExecContext(
		ctx,
		transaction.ID,
		transaction.GetVersion(),
		transaction.Name,
		transaction.Description,
		transaction.GetCreatedAt(),
		tenantID,
	)
	if err != nil {
		return Transaction{}, fmt.Errorf("sql exec: %w", err)
	}

	return r.setTransactionSteps(ctx, tx, transaction)
}

func (r *Repository) setTransactionSteps(ctx context.Context, tx *sql.Tx, transaction Transaction) (Transaction, error) {
	// delete existing steps
	query, params := sqlutil.Tenant(ctx, "DELETE FROM transaction_steps WHERE transaction_id = $1 AND transaction_version = $2", transaction.ID, transaction.GetVersion())
	stmt, err := tx.Prepare(query)
	if err != nil {
		return Transaction{}, err
	}

	_, err = stmt.ExecContext(ctx, params...)
	if err != nil {
		return Transaction{}, err
	}

	if len(transaction.StepIDs) == 0 {
		return transaction, tx.Commit()
	}

	tenantID := sqlutil.TenantID(ctx)

	values := []string{}
	for i, testID := range transaction.StepIDs {
		stepNumber := i + 1

		if tenantID == nil {
			values = append(
				values,
				fmt.Sprintf("('%s', %d, '%s', %d, NULL)", transaction.ID, transaction.GetVersion(), testID, stepNumber),
			)
		} else {
			values = append(
				values,
				fmt.Sprintf("('%s', %d, '%s', %d, '%s')", transaction.ID, transaction.GetVersion(), testID, stepNumber, *tenantID),
			)
		}
	}

	sql := "INSERT INTO transaction_steps VALUES " + strings.Join(values, ", ")
	_, err = tx.ExecContext(ctx, sql)
	if err != nil {
		return Transaction{}, fmt.Errorf("cannot save transaction steps: %w", err)
	}
	return transaction, tx.Commit()
}

func (r *Repository) readRows(ctx context.Context, rows *sql.Rows, augmented bool) ([]Transaction, error) {
	transactions := []Transaction{}

	for rows.Next() {
		transaction, err := r.readRow(ctx, rows, augmented)
		if err != nil {
			return []Transaction{}, fmt.Errorf("cannot read rows: %w", err)
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func (r *Repository) readRow(ctx context.Context, row scanner, augmented bool) (Transaction, error) {
	transaction := Transaction{
		Summary: &test.Summary{},
	}

	var (
		lastRunTime *time.Time
		stepIDs     *string

		pass, fail *int
		version    int
	)

	err := row.Scan(
		&transaction.ID,
		&version,
		&transaction.Name,
		&transaction.Description,
		&transaction.CreatedAt,
		&stepIDs,
		&transaction.Summary.Runs,
		&lastRunTime,
		&pass,
		&fail,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return transaction, err
		}

		return Transaction{}, fmt.Errorf("cannot read row: %w", err)
	}

	if stepIDs != nil && *stepIDs != "" {
		ids := strings.Split(*stepIDs, ",")
		transaction.StepIDs = make([]id.ID, len(ids))
		for i, sid := range ids {
			transaction.StepIDs[i] = id.ID(sid)
		}
	}

	if version != 0 {
		transaction.Version = &version
	}

	if lastRunTime != nil {
		transaction.Summary.LastRun.Time = *lastRunTime
	}
	if pass != nil {
		transaction.Summary.LastRun.Passes = *pass
	}
	if fail != nil {
		transaction.Summary.LastRun.Fails = *fail
	}

	if !augmented {
		removeNonAugmentedFields(&transaction)
	} else {
		steps, err := r.stepRepository.GetTransactionSteps(ctx, transaction.ID, *transaction.Version)
		if err != nil {
			return Transaction{}, fmt.Errorf("cannot read row: %w", err)
		}

		transaction.Steps = steps
	}

	return transaction, nil
}

func removeNonAugmentedFields(t *Transaction) {
	t.CreatedAt = nil
	t.Version = nil
	t.Summary = nil
}

func BumpTransactionVersionIfNeeded(original, updated Transaction) Transaction {
	transactionHasChanged := transactionHasChanged(original, updated)
	if transactionHasChanged {
		setVersion(&updated, original.GetVersion()+1)
	}

	return updated
}

func transactionHasChanged(in, updated Transaction) bool {
	jsons := []struct {
		Name        string
		Description string
		Steps       []string
	}{
		{
			Name:        in.Name,
			Description: in.Description,
			Steps:       stepIDsToStringSlice(in),
		},
		{
			Name:        updated.Name,
			Description: updated.Description,
			Steps:       stepIDsToStringSlice(updated),
		},
	}

	inJson, _ := json.Marshal(jsons[0])
	updatedJson, _ := json.Marshal(jsons[1])

	return string(inJson) != string(updatedJson)
}

func stepIDsToStringSlice(in Transaction) []string {
	steps := make([]string, len(in.StepIDs))
	for i, stepID := range in.StepIDs {
		steps[i] = stepID.String()
	}

	return steps
}
