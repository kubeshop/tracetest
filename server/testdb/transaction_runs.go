package testdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
)

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
	"steps",
	"step_runs",
	"current_test",

	-- result info
	"last_error",

	"metadata",

	-- environment
	"environment"
) VALUES (
	nextval('` + runSequenceName + `'), -- id
	$1,   -- transaction_id
	$2,   -- transaction_version

	-- timestamps
	$3,              -- created_at
	to_timestamp(0), -- completed_at

	-- trigger params
	$4, -- state
	$5, -- steps
	$6, -- stepRuns
	$7, -- currentStep

	-- result info
	NULL, -- last_error

	$8, -- metadata
	$9 -- environment
)
RETURNING "id"`

func replaceTransactionRunSequenceName(sql string, transactionID id.ID) string {
	// postgres doesn't like uppercase chars in sequence names.
	// transactionID might contain uppercase chars, and we cannot lowercase them
	// because they might lose their uniqueness.
	// md5 creates a unique, lowercase hash.
	seqName := "runs_transaction_" + md5Hash(transactionID.String()) + "_seq"
	return strings.ReplaceAll(sql, runSequenceName, seqName)
}

func (td *postgresDB) CreateTransactionRun(ctx context.Context, transactionRun model.TransactionRun) (model.TransactionRun, error) {
	jsonSteps, err := json.Marshal(transactionRun.Steps)
	if err != nil {
		return model.TransactionRun{}, fmt.Errorf("failed to marshal transaction steps: %w", err)
	}

	jsonStepRuns, err := json.Marshal(transactionRun.StepRuns)
	if err != nil {
		return model.TransactionRun{}, fmt.Errorf("failed to marshal transaction step runs: %w", err)
	}

	jsonMetadata, err := json.Marshal(transactionRun.Metadata)
	if err != nil {
		return model.TransactionRun{}, fmt.Errorf("failed to marshal transaction run metadata: %w", err)
	}

	jsonEnvironment, err := json.Marshal(transactionRun.Environment)
	if err != nil {
		return model.TransactionRun{}, fmt.Errorf("failed to marshal transaction run environment: %w", err)
	}

	tx, err := td.db.BeginTx(ctx, nil)
	if err != nil {
		return model.TransactionRun{}, fmt.Errorf("sql beginTx: %w", err)
	}

	_, err = tx.ExecContext(ctx, replaceTransactionRunSequenceName(createSequeceQuery, transactionRun.TransactionID))
	if err != nil {
		tx.Rollback()
		return model.TransactionRun{}, fmt.Errorf("sql exec: %w", err)
	}

	var runID int
	err = tx.QueryRowContext(
		ctx,
		replaceTransactionRunSequenceName(createTransactionRunQuery, transactionRun.TransactionID),
		transactionRun.TransactionID,
		transactionRun.TransactionVersion,
		transactionRun.CreatedAt,
		transactionRun.State,
		jsonSteps,
		jsonStepRuns,
		transactionRun.CurrentTest,
		jsonMetadata,
		jsonEnvironment,
	).Scan(&runID)
	if err != nil {
		tx.Rollback()
		return model.TransactionRun{}, fmt.Errorf("sql exec: %w", err)
	}

	tx.Commit()

	return td.GetTransactionRun(ctx, string(transactionRun.TransactionID), runID)
}

const updateTransactionRunQuery = `
UPDATE transaction_runs SET

	-- timestamps
	"completed_at" = $1,

	-- trigger params
	"state" = $2,
	"steps" = $3,
	"step_runs" = $4,
	"current_test" = $5,

	-- result info
	"last_error" = $6,

	"metadata" = $7,

	-- environment
	"environment" = $8

WHERE id = $9 AND transaction_id = $10
`

func (td *postgresDB) UpdateTransactionRun(ctx context.Context, transactionRun model.TransactionRun) error {
	stmt, err := td.db.Prepare(updateTransactionRunQuery)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	jsonSteps, err := json.Marshal(transactionRun.Steps)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction steps: %w", err)
	}

	jsonStepRuns, err := json.Marshal(transactionRun.StepRuns)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction step runs: %w", err)
	}

	jsonMetadata, err := json.Marshal(transactionRun.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction run metadata: %w", err)
	}

	jsonEnvironment, err := json.Marshal(transactionRun.Environment)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction run environment: %w", err)
	}

	_, err = stmt.ExecContext(
		ctx,
		transactionRun.CompletedAt,
		transactionRun.State,
		jsonSteps,
		jsonStepRuns,
		transactionRun.CurrentTest,
		transactionRun.LastError,
		jsonMetadata,
		jsonEnvironment,
		transactionRun.ID,
		transactionRun.TransactionID,
	)

	if err != nil {
		return fmt.Errorf("sql exec: %w", err)
	}

	return nil
}

func (td *postgresDB) DeleteTransactionRun(ctx context.Context, transactionRun model.TransactionRun) error {
	stmt, err := td.db.Prepare("DELETE FROM transaction_runs WHERE id = $1 AND transaction_id = $2")
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, transactionRun.ID, transactionRun.TransactionID)
	if err != nil {
		return fmt.Errorf("sql exec: %w", err)
	}

	return nil
}

const selectTransactionRunQuery = `
SELECT
	"id",
	"transaction_id",
	"transaction_version",

	"created_at",
	"completed_at",

	"state",
	"steps",
	"step_runs",
	"current_test",

	"last_error",

	"metadata",

	"environment"
FROM transaction_runs
`

func (td *postgresDB) GetTransactionRun(ctx context.Context, transactionId string, runId int) (model.TransactionRun, error) {
	stmt, err := td.db.Prepare(selectTransactionRunQuery + " WHERE id = $1 AND transaction_id = $2")
	if err != nil {
		return model.TransactionRun{}, fmt.Errorf("prepare: %w", err)
	}

	run, err := readTransactionRow(stmt.QueryRowContext(ctx, runId, transactionId))
	if err != nil {
		return model.TransactionRun{}, err
	}
	return run, nil
}

func (td *postgresDB) GetTransactionsRuns(ctx context.Context, transactionId string) ([]model.TransactionRun, error) {
	stmt, err := td.db.Prepare(selectTransactionRunQuery + " WHERE transaction_id = $1")
	if err != nil {
		return []model.TransactionRun{}, fmt.Errorf("prepare: %w", err)
	}

	rows, err := stmt.QueryContext(ctx, transactionId)
	if err != nil {
		return []model.TransactionRun{}, fmt.Errorf("query: %w", err)
	}

	var runs []model.TransactionRun
	for rows.Next() {
		run, err := readTransactionRow(rows)
		if err != nil {
			return []model.TransactionRun{}, err
		}
		runs = append(runs, run)
	}

	return runs, nil
}

func readTransactionRow(row scanner) (model.TransactionRun, error) {
	r := model.TransactionRun{}

	var (
		jsonSteps,
		jsonStepRuns,
		jsonEnvironment,
		jsonMetadata []byte

		lastError *string
	)

	err := row.Scan(
		&r.ID,
		&r.TransactionID,
		&r.TransactionVersion,
		&r.CreatedAt,
		&r.CompletedAt,
		&r.State,
		&jsonSteps,
		&jsonStepRuns,
		&r.CurrentTest,
		&lastError,
		&jsonMetadata,
		&jsonEnvironment,
	)

	switch err {
	case sql.ErrNoRows:
		return model.TransactionRun{}, ErrNotFound
	case nil:
		err = json.Unmarshal(jsonSteps, &r.Steps)
		if err != nil {
			return model.TransactionRun{}, fmt.Errorf("cannot parse transaction steps: %w", err)
		}

		err = json.Unmarshal(jsonStepRuns, &r.StepRuns)
		if err != nil {
			return model.TransactionRun{}, fmt.Errorf("cannot parse transaction step runs: %w", err)
		}

		err = json.Unmarshal(jsonMetadata, &r.Metadata)
		if err != nil {
			return model.TransactionRun{}, fmt.Errorf("cannot parse Metadata: %w", err)
		}

		err = json.Unmarshal(jsonEnvironment, &r.Environment)
		if err != nil {
			return model.TransactionRun{}, fmt.Errorf("cannot parse Environment: %w", err)
		}

		if lastError != nil && *lastError != "" {
			r.LastError = fmt.Errorf(*lastError)
		}

		return r, nil

	default:
		return model.TransactionRun{}, fmt.Errorf("read run row: %w", err)
	}
}
