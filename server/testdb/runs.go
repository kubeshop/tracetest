package testdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/otel/trace"
)

var _ model.RunRepository = &postgresDB{}

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
	"last_error",

	"metadata"
) VALUES (
	nextval('test_runs_id_seq'), -- id
	$1,   -- test_id
	$2,   -- test_version

	-- timestamps
	$3,              -- created_at
	to_timestamp(0), -- service_triggered_at
	to_timestamp(0), -- service_trigger_completed_at
	to_timestamp(0), -- obtained_trace_at
	to_timestamp(0), -- completed_at

	-- trigger params
	$4, -- state
	$5, -- trace_id
	$6, -- span_id

	-- result info
	'{}', -- trigger_results
	'{}', -- test_results
	NULL, -- trace
	NULL, -- last_error

	$7 -- metadata
)
RETURNING "id"`

func (td *postgresDB) CreateRun(ctx context.Context, test model.Test, run model.Run) (model.Run, error) {
	stmt, err := td.db.Prepare(createRunQuery)
	if err != nil {
		return model.Run{}, fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()

	run.TestID = test.ID
	run.State = model.RunStateCreated
	run.TestVersion = test.Version
	run.CreatedAt = time.Now()

	jsonMetadata, err := json.Marshal(run.Metadata)
	if err != nil {
		return model.Run{}, fmt.Errorf("encoding error: %w", err)
	}

	err = stmt.QueryRowContext(
		ctx,
		test.ID,
		test.Version,
		run.CreatedAt,
		run.State,
		run.TraceID.String(),
		run.SpanID.String(),
		jsonMetadata,
	).Scan(&run.ID)
	if err != nil {
		return model.Run{}, fmt.Errorf("sql exec: %w", err)
	}

	return run, nil
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
	"last_error" = $11,

	"metadata" = $12

WHERE id = $13
`

func (td *postgresDB) UpdateRun(ctx context.Context, r model.Run) error {
	stmt, err := td.db.Prepare(updateRunQuery)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	jsonTriggerResults, err := json.Marshal(r.TriggerResult)
	if err != nil {
		return fmt.Errorf("encoding error: %w", err)
	}

	jsonTestResults, err := json.Marshal(r.Results)
	if err != nil {
		return fmt.Errorf("encoding error: %w", err)
	}

	jsonTrace, err := json.Marshal(r.Trace)
	if err != nil {
		return fmt.Errorf("encoding error: %w", err)
	}

	jsonMetadata, err := json.Marshal(r.Metadata)
	if err != nil {
		return fmt.Errorf("encoding error: %w", err)
	}

	var lastError *string
	if r.LastError != nil {
		e := r.LastError.Error()
		lastError = &e
	}

	_, err = stmt.ExecContext(
		ctx,
		r.ServiceTriggeredAt,
		r.ServiceTriggerCompletedAt,
		r.ObtainedTraceAt,
		r.CompletedAt,
		r.State,
		r.TraceID.String(),
		r.SpanID.String(),
		jsonTriggerResults,
		jsonTestResults,
		jsonTrace,
		lastError,
		jsonMetadata,
		r.ID,
	)
	if err != nil {
		return fmt.Errorf("sql exec: %w", err)
	}

	return nil
}

func (td *postgresDB) DeleteRun(ctx context.Context, r model.Run) error {
	stmt, err := td.db.Prepare("DELETE FROM test_runs WHERE id = $1")
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, r.ID)
	if err != nil {
		return fmt.Errorf("sql exec: %w", err)
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
	"last_error",

	"metadata"
FROM test_runs
`

// TODO require test id as well
func (td *postgresDB) GetRun(ctx context.Context, id int) (model.Run, error) {
	stmt, err := td.db.Prepare(selectRunQuery + " WHERE id = $1")
	if err != nil {
		return model.Run{}, err
	}
	defer stmt.Close()

	run, err := readRunRow(stmt.QueryRowContext(ctx, id))
	if err != nil {
		return model.Run{}, fmt.Errorf("cannot read row: %w", err)
	}
	return run, nil
}

func (td *postgresDB) GetRunByShortID(ctx context.Context, shortID string) (model.Run, error) {
	panic("not implemented")
	stmt, err := td.db.Prepare(selectRunQuery + " WHERE id = $1")
	if err != nil {
		return model.Run{}, err
	}
	defer stmt.Close()

	run, err := readRunRow(stmt.QueryRowContext(ctx, shortID))
	if err != nil {
		return model.Run{}, fmt.Errorf("cannot read row: %w", err)
	}
	return run, nil
}

func (td *postgresDB) GetTestRuns(ctx context.Context, test model.Test, take, skip int32) ([]model.Run, error) {
	stmt, err := td.db.Prepare(selectRunQuery + " WHERE test_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, test.ID, take, skip)
	if err != nil {
		return nil, err
	}
	var runs []model.Run

	for rows.Next() {
		run, err := readRunRow(rows)
		if err != nil {
			return nil, fmt.Errorf("cannot read row: %w", err)
		}
		runs = append(runs, run)
	}

	return runs, nil
}

func (td *postgresDB) GetRunByTraceID(ctx context.Context, traceID trace.TraceID) (model.Run, error) {
	stmt, err := td.db.Prepare(selectRunQuery + " WHERE trace_id = $1")
	if err != nil {
		return model.Run{}, err
	}
	defer stmt.Close()

	run, err := readRunRow(stmt.QueryRowContext(ctx, traceID))
	if err != nil {
		return model.Run{}, fmt.Errorf("cannot read row: %w", err)
	}
	return run, nil
}

func readRunRow(row scanner) (model.Run, error) {
	r := model.Run{}

	var (
		jsonTriggerResults,
		jsonTestResults,
		jsonTrace,
		jsonMetadata []byte

		lastError *string
		traceID,
		spanID string
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
		&lastError,
		&jsonMetadata,
	)

	switch err {
	case sql.ErrNoRows:
		return model.Run{}, ErrNotFound
	case nil:
		err = json.Unmarshal(jsonTriggerResults, &r.TriggerResult)
		if err != nil {
			return model.Run{}, fmt.Errorf("cannot parse TriggerResult: %s", err)
		}

		err = json.Unmarshal(jsonTestResults, &r.Results)
		if err != nil {
			return model.Run{}, fmt.Errorf("cannot parse Results: %s", err)
		}

		if jsonTrace != nil {
			err = json.Unmarshal(jsonTrace, &r.Trace)
			if err != nil {
				return model.Run{}, fmt.Errorf("cannot parse Trace: %s", err)
			}
		}

		err = json.Unmarshal(jsonMetadata, &r.Metadata)
		if err != nil {
			return model.Run{}, fmt.Errorf("cannot parse Metadata: %s", err)
		}

		if traceID != "" {
			r.TraceID, err = trace.TraceIDFromHex(traceID)
			if err != nil {
				return model.Run{}, fmt.Errorf("cannot parse TraceID: %s", err)
			}
		}

		if spanID != "" {
			r.SpanID, err = trace.SpanIDFromHex(spanID)
			if err != nil {
				return model.Run{}, fmt.Errorf("cannot parse SpanID: %s", err)
			}
		}

		if lastError != nil && *lastError != "" {
			r.LastError = fmt.Errorf(*lastError)
		}

		return r, nil

	default:
		return model.Run{}, fmt.Errorf("read run row: %w", err)
	}
}
