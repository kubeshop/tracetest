package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/environment"
	"github.com/kubeshop/tracetest/server/executor/testrunner"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/sqlutil"
	"go.opentelemetry.io/otel/trace"
)

type RunRepository interface {
	CreateRun(context.Context, Test, Run) (Run, error)
	UpdateRun(context.Context, Run) error
	DeleteRun(context.Context, Run) error
	GetRun(_ context.Context, testID id.ID, runID int) (Run, error)
	GetTestRuns(_ context.Context, _ Test, take, skip int32) ([]Run, error)
	GetRunByTraceID(context.Context, trace.TraceID) (Run, error)
	GetLatestRunByTestVersion(context.Context, id.ID, int) (Run, error)

	GetTransactionRunSteps(_ context.Context, _ id.ID, runID int) ([]Run, error)
}

type runRepository struct {
	db *sql.DB
}

func NewRunRepository(db *sql.DB) RunRepository {
	return &runRepository{db}
}

const (
	createSequeceQuery = `CREATE SEQUENCE IF NOT EXISTS "` + runSequenceName + `";`
	dropSequeceQuery   = `DROP SEQUENCE IF EXISTS "` + runSequenceName + `";`
)

const createRunQuery = `
INSERT INTO test_runs (
	"id",
	"test_id",
	"test_version",

	-- timestamps
	"created_at",
	"service_triggered_at",
	"service_trigger_completed_at",
	"obtained_trace_at",
	"completed_at",

	-- trigger params
	"state",
	"trace_id",
	"span_id",

	-- result info
	"trigger_results",
	"test_results",
	"trace",
	"outputs",
	"last_error",
	"pass",
	"fail",

	"metadata",

	-- environment
	"environment",

	-- linter
	"linter",

	-- required gates
	"required_gates_result",

	"tenant_id"
) VALUES (
	nextval('` + runSequenceName + `'), -- id
	$1,   -- test_id
	$2,   -- test_version

	-- timestamps
	$3,              -- created_at
	$4,              -- service_triggered_at
	$5,              -- service_trigger_completed_at
	$6,              -- obtained_trace_at
	to_timestamp(0), -- completed_at

	-- trigger params
	$7, -- state
	$8, -- trace_id
	$9, -- span_id

	-- result info
	$10,  -- trigger_results
	'{}', -- test_results
	$11,  -- trace
	'[]', -- outputs
	NULL, -- last_error
	0,    -- pass
	0,    -- fail

	$12, -- metadata
	$13, -- environment
	$14, -- linter
	$15,  -- required_gates_result
	$16  -- tenant_id
)
RETURNING "id"`

func (r *runRepository) CreateRun(ctx context.Context, test Test, run Run) (Run, error) {
	run.TestID = test.ID
	run.State = RunStateCreated
	run.TestVersion = test.SafeVersion()
	if run.CreatedAt.IsZero() {
		run.CreatedAt = time.Now()
	}

	jsonTriggerResults, err := json.Marshal(run.TriggerResult)
	if err != nil {
		return Run{}, fmt.Errorf("trigger results encoding error: %w", err)
	}

	jsonTrace, err := json.Marshal(run.Trace)
	if err != nil {
		return Run{}, fmt.Errorf("trace encoding error: %w", err)
	}

	jsonMetadata, err := json.Marshal(run.Metadata)
	if err != nil {
		return Run{}, fmt.Errorf("metadata encoding error: %w", err)
	}

	jsonEnvironment, err := json.Marshal(run.Environment)
	if err != nil {
		return Run{}, fmt.Errorf("environment encoding error: %w", err)
	}

	jsonlinter, err := json.Marshal(run.Linter)
	if err != nil {
		return Run{}, fmt.Errorf("linter encoding error: %w", err)
	}

	jsonGatesResult, err := json.Marshal(run.RequiredGatesResult)
	if err != nil {
		return Run{}, fmt.Errorf("required gates result encoding error: %w", err)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Run{}, fmt.Errorf("sql beginTx: %w", err)
	}

	_, err = tx.ExecContext(ctx, replaceRunSequenceName(createSequeceQuery, test.ID))
	if err != nil {
		tx.Rollback()
		return Run{}, fmt.Errorf("sql exec: %w", err)
	}

	tenantID := sqlutil.TenantID(ctx)

	var runID int
	err = tx.QueryRowContext(
		ctx,
		replaceRunSequenceName(createRunQuery, test.ID),
		test.ID,
		test.SafeVersion(),
		run.CreatedAt,
		run.ServiceTriggeredAt,
		run.ServiceTriggerCompletedAt,
		run.ObtainedTraceAt,
		run.State,
		run.TraceID.String(),
		run.SpanID.String(),
		jsonTriggerResults,
		jsonTrace,
		jsonMetadata,
		jsonEnvironment,
		jsonlinter,
		jsonGatesResult,
		tenantID,
	).Scan(&runID)
	if err != nil {
		tx.Rollback()
		return Run{}, fmt.Errorf("sql exec: %w", err)
	}

	tx.Commit()

	return r.GetRun(ctx, test.ID, runID)
}

const updateRunQuery = `
UPDATE test_runs SET

	-- timestamps
	"service_triggered_at" = $1,
	"service_trigger_completed_at" = $2,
	"obtained_trace_at" = $3,
	"completed_at" = $4,

	-- trigger params
	"state" = $5,
	"trace_id" = $6,
	"span_id" = $7,

	-- result info
	"trigger_results" = $8,
	"test_results" = $9,
	"trace" = $10,
	"outputs" = $11,
	"last_error" = $12,
	"pass" = $13,
	"fail" = $14,

	"metadata" = $15,
	"environment" = $18,

	--- linter
	"linter" = $19,

	--- required gates
	"required_gates_result" = $20

WHERE id = $16 AND test_id = $17
`

func (r *runRepository) UpdateRun(ctx context.Context, run Run) error {
	jsonTriggerResults, err := json.Marshal(run.TriggerResult)
	if err != nil {
		return fmt.Errorf("trigger results encoding error: %w", err)
	}

	jsonTestResults, err := json.Marshal(run.Results)
	if err != nil {
		return fmt.Errorf("test results encoding error: %w", err)
	}

	jsonTrace, err := json.Marshal(run.Trace)
	if err != nil {
		return fmt.Errorf("trace encoding error: %w", err)
	}

	jsonOutputs, err := json.Marshal(run.Outputs)
	if err != nil {
		return fmt.Errorf("outputs encoding error: %w", err)
	}

	jsonMetadata, err := json.Marshal(run.Metadata)
	if err != nil {
		return fmt.Errorf("encoding error: %w", err)
	}

	jsonEnvironment, err := json.Marshal(run.Environment)
	if err != nil {
		return fmt.Errorf("encoding error: %w", err)
	}

	jsonLinter, err := json.Marshal(run.Linter)
	if err != nil {
		return fmt.Errorf("encoding error: %w", err)
	}

	jsonGatesResult, err := json.Marshal(run.RequiredGatesResult)
	if err != nil {
		return fmt.Errorf("encoding error: %w", err)
	}

	var lastError *string
	if run.LastError != nil {
		e := run.LastError.Error()
		lastError = &e
	}

	pass, fail := run.ResultsCount()

	query, params := sqlutil.Tenant(
		ctx,
		updateRunQuery,
		run.ServiceTriggeredAt,
		run.ServiceTriggerCompletedAt,
		run.ObtainedTraceAt,
		run.CompletedAt,
		run.State,
		run.TraceID.String(),
		run.SpanID.String(),
		jsonTriggerResults,
		jsonTestResults,
		jsonTrace,
		jsonOutputs,
		lastError,
		pass,
		fail,
		jsonMetadata,
		run.ID,
		run.TestID,
		jsonEnvironment,
		jsonLinter,
		jsonGatesResult,
	)

	_, err = r.db.ExecContext(
		ctx,
		query,
		params...,
	)
	if err != nil {
		return fmt.Errorf("sql exec: %w", err)
	}

	return nil
}

func (r *runRepository) DeleteRun(ctx context.Context, run Run) error {
	queries := []string{
		"DELETE FROM transaction_run_steps WHERE test_run_id = $1 AND test_run_test_id = $2",
		"DELETE FROM test_runs WHERE id = $1 AND test_id = $2",
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("sql BeginTx: %w", err)
	}

	for _, sql := range queries {
		query, params := sqlutil.Tenant(ctx, sql, run.ID, run.TestID)
		_, err := tx.ExecContext(ctx, query, params...)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("sql error: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("sql Commit: %w", err)
	}

	return nil
}

const selectRunQuery = `
SELECT
	"id",
	"test_id",
	"test_version",

	-- timestamps
	"created_at",
	"service_triggered_at",
	"service_trigger_completed_at",
	"obtained_trace_at",
	"completed_at",

	-- trigger params
	"state",
	"trace_id",
	"span_id",

	-- result info
	"trigger_results",
	"test_results",
	"trace",
	"outputs",
	"last_error",
	"metadata",
	"environment",

	-- transaction run
	transaction_run_steps.transaction_run_id,
	transaction_run_steps.transaction_run_transaction_id,
	"linter",
	"required_gates_result"

FROM
	test_runs
		LEFT OUTER JOIN transaction_run_steps
		ON transaction_run_steps.test_run_id = test_runs.id AND transaction_run_steps.test_run_test_id = test_runs.test_id
`

func (r *runRepository) GetRun(ctx context.Context, testID id.ID, runID int) (Run, error) {
	query, params := sqlutil.Tenant(ctx, selectRunQuery+" WHERE id = $1 AND test_id = $2", runID, testID)

	run, err := readRunRow(r.db.QueryRowContext(ctx, query, params...))
	if err != nil {
		return Run{}, fmt.Errorf("cannot read row: %w", err)
	}
	return run, nil
}

func (r *runRepository) GetTestRuns(ctx context.Context, test Test, take, skip int32) ([]Run, error) {
	query, params := sqlutil.Tenant(ctx, selectRunQuery+" WHERE test_id = $1", test.ID, take, skip)
	stmt, err := r.db.Prepare(query + " ORDER BY created_at DESC LIMIT $2 OFFSET $3")
	if err != nil {
		return []Run{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return []Run{}, err
	}

	runs, err := r.readRunRows(ctx, rows)
	if err != nil {
		return []Run{}, err
	}

	return runs, nil
}

func (r *runRepository) GetRunByTraceID(ctx context.Context, traceID trace.TraceID) (Run, error) {
	query, params := sqlutil.Tenant(ctx, selectRunQuery+" WHERE trace_id = $1", traceID.String())
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return Run{}, err
	}
	defer stmt.Close()

	run, err := readRunRow(stmt.QueryRowContext(ctx, params...))
	if err != nil {
		return Run{}, fmt.Errorf("cannot read row: %w", err)
	}
	return run, nil
}

func (r *runRepository) GetLatestRunByTestVersion(ctx context.Context, testID id.ID, version int) (Run, error) {
	query, params := sqlutil.Tenant(ctx, selectRunQuery+" WHERE test_id = $1 AND test_version = $2 ORDER BY created_at DESC LIMIT 1", testID.String(), version)
	stmt, err := r.db.Prepare(query)

	if err != nil {
		return Run{}, err
	}
	defer stmt.Close()

	run, err := readRunRow(stmt.QueryRowContext(ctx, params...))
	if err != nil {
		return Run{}, err
	}
	return run, nil
}

func (r *runRepository) readRunRows(ctx context.Context, rows *sql.Rows) ([]Run, error) {
	var runs []Run

	for rows.Next() {
		run, err := readRunRow(rows)
		if err != nil {
			return []Run{}, fmt.Errorf("cannot read row: %w", err)
		}
		runs = append(runs, run)
	}

	return runs, nil
}

func readRunRow(row scanner) (Run, error) {
	r := Run{}

	var (
		jsonTriggerResults,
		jsonTestResults,
		jsonTrace,
		jsonOutputs,
		jsonEnvironment,
		jsonLinter,
		jsonGatesResult,
		jsonMetadata []byte

		lastError *string
		traceID,
		spanID string

		transactionID,
		transactionRunID sql.NullString
	)

	err := row.Scan(
		&r.ID,
		&r.TestID,
		&r.TestVersion,
		&r.CreatedAt,
		&r.ServiceTriggeredAt,
		&r.ServiceTriggerCompletedAt,
		&r.ObtainedTraceAt,
		&r.CompletedAt,
		&r.State,
		&traceID,
		&spanID,
		&jsonTriggerResults,
		&jsonTestResults,
		&jsonTrace,
		&jsonOutputs,
		&lastError,
		&jsonMetadata,
		&jsonEnvironment,
		&transactionRunID,
		&transactionID,
		&jsonLinter,
		&jsonGatesResult,
	)

	if err != nil {
		return Run{}, err
	}

	err = json.Unmarshal(jsonTriggerResults, &r.TriggerResult)
	if err != nil {
		return Run{}, fmt.Errorf("cannot parse TriggerResult: %w", err)
	}

	err = json.Unmarshal(jsonTestResults, &r.Results)
	if err != nil {
		return Run{}, fmt.Errorf("cannot parse Results: %w", err)
	}

	if jsonTrace != nil {
		err = json.Unmarshal(jsonTrace, &r.Trace)
		if err != nil {
			return Run{}, fmt.Errorf("cannot parse Trace: %w", err)
		}
	}

	if jsonLinter != nil {
		err = json.Unmarshal(jsonLinter, &r.Linter)
		if err != nil {
			return Run{}, fmt.Errorf("cannot parse linter: %w", err)
		}
	}

	err = json.Unmarshal(jsonOutputs, &r.Outputs)
	if err != nil {
		// try with raw outputs
		var rawOutputs []environment.EnvironmentValue
		err = json.Unmarshal(jsonOutputs, &rawOutputs)

		for _, value := range rawOutputs {
			r.Outputs.Add(value.Key, RunOutput{
				Name:   value.Key,
				Value:  value.Value,
				SpanID: "",
			})
		}

		if err != nil {
			return Run{}, fmt.Errorf("cannot parse Outputs: %w", err)
		}
	}

	err = json.Unmarshal(jsonMetadata, &r.Metadata)
	if err != nil {
		return Run{}, fmt.Errorf("cannot parse Metadata: %w", err)
	}

	err = json.Unmarshal(jsonEnvironment, &r.Environment)
	if err != nil {
		return Run{}, fmt.Errorf("cannot parse Environment: %w", err)
	}

	if jsonGatesResult != nil {
		err = json.Unmarshal(jsonGatesResult, &r.RequiredGatesResult)
		if err != nil {
			return Run{}, fmt.Errorf("cannot parse required gates result: %w", err)
		}
	} else {
		// fallback for retro-compatibility
		r.RequiredGatesResult = r.GenerateRequiredGateResult(testrunner.DefaultTestRunner.RequiredGates)
	}

	if traceID != "" {
		r.TraceID, err = trace.TraceIDFromHex(traceID)
		if err != nil {
			return Run{}, fmt.Errorf("cannot parse TraceID: %w", err)
		}
	}

	if spanID != "" {
		r.SpanID, err = trace.SpanIDFromHex(spanID)
		if err != nil {
			return Run{}, fmt.Errorf("cannot parse SpanID: %w", err)
		}
	}

	if lastError != nil && *lastError != "" {
		r.LastError = fmt.Errorf(*lastError)
	}

	if transactionID.Valid && transactionRunID.Valid {
		r.TransactionID = transactionID.String
		r.TransactionRunID = transactionRunID.String
	}

	return r, nil
}

func (r *runRepository) GetTransactionRunSteps(ctx context.Context, id id.ID, runID int) ([]Run, error) {
	query := selectRunQuery + `
WHERE transaction_run_steps.transaction_run_id = $1 AND transaction_run_steps.transaction_run_transaction_id = $2
ORDER BY test_runs.completed_at ASC
`
	query, params := sqlutil.Tenant(ctx, query, runID, id)

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return []Run{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return []Run{}, fmt.Errorf("query context: %w", err)
	}

	steps, err := r.readRunRows(ctx, rows)
	if err != nil {
		return []Run{}, fmt.Errorf("read row: %w", err)
	}

	return steps, nil
}
