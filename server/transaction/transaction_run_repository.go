package transaction

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/sqlutil"
	"github.com/kubeshop/tracetest/server/test"
)

type transactionStepRunRepository interface {
	GetTransactionRunSteps(_ context.Context, _ id.ID, runID int) ([]test.Run, error)
}

func NewRunRepository(db *sql.DB, stepsRepository transactionStepRunRepository) *RunRepository {
	return &RunRepository{
		db:              db,
		stepsRepository: stepsRepository,
	}
}

type RunRepository struct {
	db              *sql.DB
	stepsRepository transactionStepRunRepository
}

const createTransactionRunQuery = `
INSERT INTO transaction_runs (
	"id",
	"transaction_id",
	"transaction_version",

	-- timestamps
	"created_at",
	"completed_at",

	-- trigger params
	"state",
	"current_test",

	-- result info
	"last_error",
	"pass",
	"fail",
	"all_steps_required_gates_passed",

	"metadata",

	-- environment
	"environment",

	"tenant_id"
) VALUES (
	nextval('` + runSequenceName + `'), -- id
	$1,   -- transaction_id
	$2,   -- transaction_version

	-- timestamps
	$3,              -- created_at
	to_timestamp(0), -- completed_at

	-- trigger params
	$4, -- state
	$5, -- currentStep

	-- result info
	NULL, -- last_error
	0,    -- pass
	0,    -- fail
	TRUE, -- all_steps_required_gates_passed

	$6, -- metadata
	$7, -- environment
	$8 -- tenant_id
)
RETURNING "id"`

const (
	createSequenceQuery = `CREATE SEQUENCE IF NOT EXISTS "` + runSequenceName + `";`
	dropSequenceQuery   = `DROP SEQUENCE IF EXISTS "` + runSequenceName + `";`
	runSequenceName     = "%sequence_name%"
)

func md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func replaceTransactionRunSequenceName(sql string, transactionID id.ID) string {
	// postgres doesn't like uppercase chars in sequence names.
	// transactionID might contain uppercase chars, and we cannot lowercase them
	// because they might lose their uniqueness.
	// md5 creates a unique, lowercase hash.
	seqName := "runs_transaction_" + md5Hash(transactionID.String()) + "_seq"
	return strings.ReplaceAll(sql, runSequenceName, seqName)
}

func (td *RunRepository) CreateRun(ctx context.Context, tr TransactionRun) (TransactionRun, error) {
	jsonMetadata, err := json.Marshal(tr.Metadata)
	if err != nil {
		return TransactionRun{}, fmt.Errorf("failed to marshal transaction run metadata: %w", err)
	}

	jsonEnvironment, err := json.Marshal(tr.Environment)
	if err != nil {
		return TransactionRun{}, fmt.Errorf("failed to marshal transaction run environment: %w", err)
	}

	tx, err := td.db.BeginTx(ctx, nil)
	if err != nil {
		return TransactionRun{}, fmt.Errorf("sql beginTx: %w", err)
	}

	_, err = tx.ExecContext(ctx, replaceTransactionRunSequenceName(createSequenceQuery, tr.TransactionID))
	if err != nil {
		tx.Rollback()
		return TransactionRun{}, fmt.Errorf("sql exec: %w", err)
	}

	tenantID := sqlutil.TenantID(ctx)

	var runID int
	err = tx.QueryRowContext(
		ctx,
		replaceTransactionRunSequenceName(createTransactionRunQuery, tr.TransactionID),
		tr.TransactionID,
		tr.TransactionVersion,
		tr.CreatedAt,
		tr.State,
		tr.CurrentTest,
		jsonMetadata,
		jsonEnvironment,
		tenantID,
	).Scan(&runID)
	if err != nil {
		tx.Rollback()
		return TransactionRun{}, fmt.Errorf("sql exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return TransactionRun{}, fmt.Errorf("commit: %w", err)
	}

	tr.ID = runID

	return tr, nil
}

const updateTransactionRunQuery = `
UPDATE transaction_runs SET

	-- timestamps
	"completed_at" = $1,

	-- trigger params
	"state" = $2,
	"current_test" = $3,

	-- result info
	"last_error" = $4,
	"pass" = $5,
	"fail" = $6,
	"all_steps_required_gates_passed" = $7,

	"metadata" = $8,

	-- environment
	"environment" = $9

WHERE id = $10 AND transaction_id = $11
`

func (td *RunRepository) UpdateRun(ctx context.Context, tr TransactionRun) error {
	tx, err := td.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("sql beginTx: %w", err)
	}

	jsonMetadata, err := json.Marshal(tr.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction run metadata: %w", err)
	}

	jsonEnvironment, err := json.Marshal(tr.Environment)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction run environment: %w", err)
	}
	var lastError *string
	if tr.LastError != nil {
		e := tr.LastError.Error()
		lastError = &e
	}

	pass, fail := tr.ResultsCount()
	allStepsRequiredGatesPassed := tr.StepsGatesValidation()

	query, params := sqlutil.Tenant(
		ctx,
		updateTransactionRunQuery,
		tr.CompletedAt,
		tr.State,
		tr.CurrentTest,
		lastError,
		pass,
		fail,
		allStepsRequiredGatesPassed,
		jsonMetadata,
		jsonEnvironment,
		tr.ID,
		tr.TransactionID,
	)
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, params...)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("sql exec: %w", err)
	}

	return td.setTransactionRunSteps(ctx, tx, tr)
}

func (td *RunRepository) setTransactionRunSteps(ctx context.Context, tx *sql.Tx, tr TransactionRun) error {
	// delete existing steps
	query, params := sqlutil.Tenant(ctx, "DELETE FROM transaction_run_steps WHERE transaction_run_id = $1 AND transaction_run_transaction_id = $2", tr.ID, tr.TransactionID)
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, params...)
	if err != nil {
		return err
	}

	if len(tr.Steps) == 0 {
		return tx.Commit()
	}

	tenantID := sqlutil.TenantID(ctx)

	values := []string{}
	for _, run := range tr.Steps {
		if run.ID == 0 {
			// step not set, skip
			continue
		}

		if tenantID == nil {
			values = append(
				values,
				fmt.Sprintf("('%d', '%s', %d, '%s', NULL)", tr.ID, tr.TransactionID, run.ID, run.TestID),
			)
		} else {
			values = append(
				values,
				fmt.Sprintf("('%d', '%s', %d, '%s', '%s')", tr.ID, tr.TransactionID, run.ID, run.TestID, *tenantID),
			)
		}
	}

	sql := "INSERT INTO transaction_run_steps VALUES " + strings.Join(values, ", ")
	_, err = tx.ExecContext(ctx, sql)
	if err != nil {
		return fmt.Errorf("cannot save transaction run steps: %w", err)
	}
	return tx.Commit()
}

func (td *RunRepository) DeleteTransactionRun(ctx context.Context, tr TransactionRun) error {
	tx, err := td.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("sql beginTx: %w", err)
	}

	query, params := sqlutil.Tenant(ctx, "DELETE FROM transaction_run_steps WHERE transaction_run_id = $1 AND transaction_run_transaction_id = $2", tr.ID, tr.TransactionID)
	_, err = tx.ExecContext(ctx, query, params...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("delete transaction run steps: %w", err)
	}

	query, params = sqlutil.Tenant(ctx, "DELETE FROM transaction_runs WHERE id = $1 AND transaction_id = $2", tr.ID, tr.TransactionID)
	_, err = tx.ExecContext(ctx, query, params...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("delete transaction runs: %w", err)
	}

	return tx.Commit()
}

const selectTransactionRunQuery = `
SELECT
	"id",
	"transaction_id",
	"transaction_version",

	"created_at",
	"completed_at",

	"state",
	"current_test",

	"last_error",
	"pass",
	"fail",
	"all_steps_required_gates_passed",

	"metadata",

	"environment"
FROM transaction_runs
`

func (td *RunRepository) GetTransactionRun(ctx context.Context, transactionID id.ID, runID int) (TransactionRun, error) {
	query, params := sqlutil.Tenant(ctx, selectTransactionRunQuery+" WHERE id = $1 AND transaction_id = $2", strconv.Itoa(runID), transactionID)
	stmt, err := td.db.Prepare(query)
	if err != nil {
		return TransactionRun{}, fmt.Errorf("prepare: %w", err)
	}

	run, err := td.readRunRow(stmt.QueryRowContext(ctx, params...))
	if err != nil {
		return TransactionRun{}, err
	}
	run.Steps, err = td.stepsRepository.GetTransactionRunSteps(ctx, run.TransactionID, run.ID)
	if err != nil {
		return TransactionRun{}, err
	}
	return run, nil
}

func (td *RunRepository) GetLatestRunByTransactionVersion(ctx context.Context, transactionID id.ID, version int) (TransactionRun, error) {
	sortQuery := "ORDER BY created_at DESC LIMIT 1"
	query, params := sqlutil.Tenant(ctx, selectTransactionRunQuery+" WHERE transaction_id = $1 AND transaction_version = $2", transactionID, version)
	stmt, err := td.db.Prepare(query + sortQuery)
	if err != nil {
		return TransactionRun{}, fmt.Errorf("prepare: %w", err)
	}

	run, err := td.readRunRow(stmt.QueryRowContext(ctx, params...))
	if err != nil {
		return TransactionRun{}, err
	}
	run.Steps, err = td.stepsRepository.GetTransactionRunSteps(ctx, run.TransactionID, run.ID)
	if err != nil {
		return TransactionRun{}, err
	}
	return run, nil
}

func (td *RunRepository) GetTransactionsRuns(ctx context.Context, transactionID id.ID, take, skip int32) ([]TransactionRun, error) {
	sortQuery := "ORDER BY created_at DESC LIMIT $2 OFFSET $3"
	query, params := sqlutil.Tenant(ctx, selectTransactionRunQuery+" WHERE transaction_id = $1", transactionID.String(), take, skip)
	stmt, err := td.db.Prepare(query + sortQuery)
	if err != nil {
		return []TransactionRun{}, fmt.Errorf("prepare: %w", err)
	}

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return []TransactionRun{}, fmt.Errorf("query: %w", err)
	}

	var runs []TransactionRun
	for rows.Next() {
		run, err := td.readRunRow(rows)
		if err != nil {
			return []TransactionRun{}, err
		}

		run.Steps, err = td.stepsRepository.GetTransactionRunSteps(ctx, run.TransactionID, run.ID)
		if err != nil {
			return []TransactionRun{}, err
		}

		runs = append(runs, run)
	}

	return runs, nil
}

func (td *RunRepository) readRunRow(row scanner) (TransactionRun, error) {
	r := TransactionRun{}

	var (
		jsonEnvironment,
		jsonMetadata []byte

		lastError *string
		pass      sql.NullInt32
		fail      sql.NullInt32

		allStepsRequiredGatesPassed *bool
	)

	err := row.Scan(
		&r.ID,
		&r.TransactionID,
		&r.TransactionVersion,
		&r.CreatedAt,
		&r.CompletedAt,
		&r.State,
		&r.CurrentTest,
		&lastError,
		&pass,
		&fail,
		&allStepsRequiredGatesPassed,
		&jsonMetadata,
		&jsonEnvironment,
	)
	if err != nil {
		return TransactionRun{}, fmt.Errorf("cannot read row: %w", err)
	}

	err = json.Unmarshal(jsonMetadata, &r.Metadata)
	if err != nil {
		return TransactionRun{}, fmt.Errorf("cannot parse Metadata: %w", err)
	}

	err = json.Unmarshal(jsonEnvironment, &r.Environment)
	if err != nil {
		return TransactionRun{}, fmt.Errorf("cannot parse Environment: %w", err)
	}

	if lastError != nil && *lastError != "" {
		r.LastError = fmt.Errorf(*lastError)
	}

	if pass.Valid {
		r.Pass = int(pass.Int32)
	}

	if fail.Valid {
		r.Fail = int(fail.Int32)
	}

	// checks if the flag exists, if it doesn't we use the fail field to determine if all steps passed
	if allStepsRequiredGatesPassed == nil {
		failed := r.Fail == 0
		allStepsRequiredGatesPassed = &failed
	}

	r.AllStepsRequiredGatesPassed = *allStepsRequiredGatesPassed

	return r, nil
}
