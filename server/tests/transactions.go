package tests

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/sqlutil"
)

const (
	TransactionResourceName       = "Transaction"
	TransactionResourceNamePlural = "Transactions"
)

type Transaction struct {
	ID          id.ID          `json:"id"`
	CreatedAt   *time.Time     `json:"createdAt,omitempty"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Version     *int           `json:"version,omitempty"`
	StepIDs     []id.ID        `json:"steps"`
	Steps       []model.Test   `json:"fullSteps,omitempty"`
	Summary     *model.Summary `json:"summary,omitempty"`
}

func setVersion(t *Transaction, v int) {
	t.Version = &v
}

func (t Transaction) GetVersion() int {
	if t.Version == nil {
		return 0
	}
	return *t.Version
}

func setCreatedAt(t *Transaction, d time.Time) {
	t.CreatedAt = &d
}

func (t Transaction) GetCreatedAt() time.Time {
	if t.CreatedAt == nil {
		return time.Time{}
	}
	return *t.CreatedAt
}

func (t Transaction) HasID() bool {
	return t.ID != ""
}

func (t Transaction) Validate() error {
	return nil
}

func (t Transaction) NewRun() TransactionRun {

	return TransactionRun{
		TransactionID:      t.ID,
		TransactionVersion: t.GetVersion(),
		CreatedAt:          time.Now(),
		State:              TransactionRunStateCreated,
		Steps:              make([]model.Run, 0, len(t.StepIDs)),
		CurrentTest:        0,
	}
}

type TransactionRunState string

const (
	TransactionRunStateCreated   TransactionRunState = "CREATED"
	TransactionRunStateExecuting TransactionRunState = "EXECUTING"
	TransactionRunStateFailed    TransactionRunState = "FAILED"
	TransactionRunStateFinished  TransactionRunState = "FINISHED"
)

func (rs TransactionRunState) IsFinal() bool {
	return rs == TransactionRunStateFailed || rs == TransactionRunStateFinished
}

type TransactionRun struct {
	ID                 int
	TransactionID      id.ID
	TransactionVersion int

	// Timestamps
	CreatedAt   time.Time
	CompletedAt time.Time

	// steps
	StepIDs []int
	Steps   []model.Run

	// trigger params
	State       TransactionRunState
	CurrentTest int

	// result info
	LastError error
	Pass      int
	Fail      int

	Metadata model.RunMetadata

	// environment
	Environment environment.Environment
}

func (tr TransactionRun) ResourceID() string {
	return fmt.Sprintf("transaction/%s/run/%d", tr.TransactionID, tr.ID)
}

func (tr TransactionRun) ResultsCount() (pass, fail int) {
	if tr.Steps == nil {
		return
	}

	for _, step := range tr.Steps {
		stepPass, stepFail := step.ResultsCount()

		pass += stepPass
		fail += stepFail
	}

	return
}

type stepsRepository interface {
	GetTransactionSteps(context.Context, Transaction) ([]model.Test, error)
	GetTransactionRunSteps(context.Context, TransactionRun) ([]model.Run, error)
}

func NewTransactionsRepository(db *sql.DB, stepsRepository stepsRepository) *TransactionsRepository {
	repo := &TransactionsRepository{
		db:              db,
		stepsRepository: stepsRepository,
	}

	return repo
}

type TransactionsRepository struct {
	db              *sql.DB
	stepsRepository stepsRepository
}

// needed for test
func (r *TransactionsRepository) DB() *sql.DB {
	return r.db
}

func (r *TransactionsRepository) SetID(t Transaction, id id.ID) Transaction {
	t.ID = id
	return t
}

func (r *TransactionsRepository) Provision(ctx context.Context, transaction Transaction) error {
	_, err := r.Create(ctx, transaction)

	return err
}

func (r *TransactionsRepository) Create(ctx context.Context, transaction Transaction) (Transaction, error) {
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

func (r *TransactionsRepository) Update(ctx context.Context, transaction Transaction) (Transaction, error) {
	if transaction.GetVersion() == 0 {
		setVersion(&transaction, 1)
	}

	oldTransaction, err := r.GetLatestVersion(ctx, transaction.ID)
	if err != nil {
		return Transaction{}, fmt.Errorf("could not get latest test version while updating test: %w", err)
	}

	// keep the same creation date to keep sort order
	transaction.CreatedAt = oldTransaction.CreatedAt
	transactionToUpdate := BumpTransactionVersionIfNeeded(oldTransaction, transaction)

	if oldTransaction.Version == transactionToUpdate.Version {
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

func (r *TransactionsRepository) checkIDExists(ctx context.Context, id id.ID) error {
	exists, err := r.IDExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}

	return nil
}

func (r *TransactionsRepository) Delete(ctx context.Context, id id.ID) error {
	if err := r.checkIDExists(ctx, id); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM transaction_steps WHERE transaction_id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	q := "DELETE FROM transaction_run_steps WHERE transaction_run_id IN (SELECT id FROM transaction_runs WHERE transaction_id = $1)"
	_, err = tx.ExecContext(ctx, q, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM transaction_runs WHERE transaction_id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM transactions WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *TransactionsRepository) IDExists(ctx context.Context, id id.ID) (bool, error) {
	exists := false
	row := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) > 0 as exists FROM transactions WHERE id = $1",
		id,
	)

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

func (r *TransactionsRepository) SortingFields() []string {
	return []string{
		"name",
		"created",
		"last_run",
	}
}

func (r *TransactionsRepository) ListAugmented(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]Transaction, error) {
	return r.list(ctx, take, skip, query, sortBy, sortDirection, true)
}

func (r *TransactionsRepository) List(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]Transaction, error) {
	return r.list(ctx, take, skip, query, sortBy, sortDirection, false)
}

func (r *TransactionsRepository) list(ctx context.Context, take, skip int, query, sortBy, sortDirection string, augmented bool) ([]Transaction, error) {
	q, params := listQuery(querySelect(), query, []any{take, skip})

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

func (r *TransactionsRepository) Count(ctx context.Context, query string) (int, error) {
	sql, params := listQuery(queryCount(), query, []any{})

	count := 0
	err := r.db.
		QueryRowContext(ctx, sql, params...).
		Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *TransactionsRepository) GetAugmented(ctx context.Context, id id.ID) (Transaction, error) {
	return r.get(ctx, id, true)
}

func (r *TransactionsRepository) Get(ctx context.Context, id id.ID) (Transaction, error) {
	return r.get(ctx, id, false)
}

func (r *TransactionsRepository) get(ctx context.Context, id id.ID, augmented bool) (Transaction, error) {
	stmt, err := r.db.Prepare(querySelect() + " WHERE t.id = $1 ORDER BY t.version DESC LIMIT 1")
	if err != nil {
		return Transaction{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	transaction, err := r.readRow(ctx, stmt.QueryRowContext(ctx, id), augmented)
	if err != nil {
		return Transaction{}, err
	}

	return transaction, nil
}

func (r *TransactionsRepository) GetVersion(ctx context.Context, id id.ID, version int) (Transaction, error) {
	stmt, err := r.db.Prepare(querySelect() + " WHERE t.id = $1 AND t.version = $2")
	if err != nil {
		return Transaction{}, fmt.Errorf("prepare 1: %w", err)
	}
	defer stmt.Close()

	transaction, err := r.readRow(ctx, stmt.QueryRowContext(ctx, id, version), true)
	if err != nil {
		return Transaction{}, err
	}

	return transaction, nil
}

func listQuery(baseSQL, query string, params []any) (string, []any) {
	condition := " AND (t.name ilike $3 OR t.description ilike $3)"

	sql := baseSQL + transactionMaxVersionQuery
	sql, params = sqlutil.Search(sql, condition, query, params)

	return sql, params
}

func (r *TransactionsRepository) GetLatestVersion(ctx context.Context, id id.ID) (Transaction, error) {
	stmt, err := r.db.Prepare(querySelect() + " WHERE t.id = $1 ORDER BY t.version DESC LIMIT 1")
	if err != nil {
		return Transaction{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	transaction, err := r.readRow(ctx, stmt.QueryRowContext(ctx, id), true)
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
	"created_at"
) VALUES ($1, $2, $3, $4, $5)`

func (r *TransactionsRepository) insertIntoTransactions(ctx context.Context, transaction Transaction) (Transaction, error) {
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

	_, err = stmt.ExecContext(
		ctx,
		transaction.ID,
		transaction.GetVersion(),
		transaction.Name,
		transaction.Description,
		transaction.GetCreatedAt(),
	)
	if err != nil {
		return Transaction{}, fmt.Errorf("sql exec: %w", err)
	}

	return r.setTransactionSteps(ctx, tx, transaction)
}

func (r *TransactionsRepository) setTransactionSteps(ctx context.Context, tx *sql.Tx, transaction Transaction) (Transaction, error) {
	// delete existing steps
	stmt, err := tx.Prepare("DELETE FROM transaction_steps WHERE transaction_id = $1 AND transaction_version = $2")
	if err != nil {
		return Transaction{}, err
	}

	_, err = stmt.ExecContext(ctx, transaction.ID, transaction.GetVersion())
	if err != nil {
		return Transaction{}, err
	}

	if len(transaction.StepIDs) == 0 {
		return transaction, tx.Commit()
	}

	values := []string{}
	for i, testID := range transaction.StepIDs {
		stepNumber := i + 1
		values = append(
			values,
			fmt.Sprintf("('%s', %d, '%s', %d)", transaction.ID, transaction.GetVersion(), testID, stepNumber),
		)
	}

	sql := "INSERT INTO transaction_steps VALUES " + strings.Join(values, ", ")
	_, err = tx.ExecContext(ctx, sql)
	if err != nil {
		return Transaction{}, fmt.Errorf("cannot save transaction steps: %w", err)
	}
	return transaction, tx.Commit()
}

func (r *TransactionsRepository) readRows(ctx context.Context, rows *sql.Rows, augmented bool) ([]Transaction, error) {
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

func (r *TransactionsRepository) readRow(ctx context.Context, row scanner, augmented bool) (Transaction, error) {
	transaction := Transaction{
		Summary: &model.Summary{},
	}

	var (
		lastRunTime *time.Time
		stepIDs     string

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
		return Transaction{}, fmt.Errorf("cannot read row: %w", err)
	}

	if stepIDs != "" {
		ids := strings.Split(stepIDs, ",")
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
		steps, err := r.stepsRepository.GetTransactionSteps(ctx, transaction)
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

func BumpTransactionVersionIfNeeded(in, updated Transaction) Transaction {
	transactionHasChanged := transactionHasChanged(in, updated)
	if transactionHasChanged {
		setVersion(&updated, in.GetVersion()+1)
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
