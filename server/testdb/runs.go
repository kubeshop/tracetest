package testdb

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kubeshop/tracetest/server/id"
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
	"outputs",
	"last_error",
	"pass",
	"fail",

	"metadata",

	-- environments
	"env_snapshot"
) VALUES (
	nextval('` + runSequenceName + `'), -- id
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
	'[]', -- outputs
	NULL, -- last_error
	0,    -- pass
	0,    -- fail

	$7, -- metadata

	$8 -- envSnapshot
)
RETURNING "id"`

const (
	createSequeceQuery = `CREATE SEQUENCE IF NOT EXISTS "` + runSequenceName + `";`
	dropSequeceQuery   = `DROP SEQUENCE IF EXISTS "` + runSequenceName + `";`
	runSequenceName    = "%sequence_name%"
)

func dropSequece(ctx context.Context, tx *sql.Tx, testID id.ID) error {
	_, err := tx.ExecContext(
		ctx,
		replaceRunSequenceName(createSequeceQuery, testID),
	)

	return err
}

func md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func replaceRunSequenceName(sql string, testID id.ID) string {
	// postgres doesn't like uppercase chars in sequence names.
	// testID might contain uppercase chars, and we cannot lowercase them
	// because they might lose their uniqueness.
	// md5 creates a unique, lowercase hash.
	seqName := "runs_test_" + md5Hash(testID.String()) + "_seq"
	return strings.ReplaceAll(sql, runSequenceName, seqName)
}

func (td *postgresDB) CreateRun(ctx context.Context, test model.Test, run model.Run) (model.Run, error) {
	run.TestID = test.ID
	run.State = model.RunStateCreated
	run.TestVersion = test.Version
	if run.CreatedAt.IsZero() {
		run.CreatedAt = time.Now()
	}

	jsonMetadata, err := json.Marshal(run.Metadata)
	if err != nil {
		return model.Run{}, fmt.Errorf("encoding error: %w", err)
	}

	jsonEnvSnapshot, err := json.Marshal(run.EnvSnapshot)
	if err != nil {
		return model.Run{}, fmt.Errorf("encoding error: %w", err)
	}

	tx, err := td.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Run{}, fmt.Errorf("sql beginTx: %w", err)
	}

	_, err = tx.ExecContext(ctx, replaceRunSequenceName(createSequeceQuery, test.ID))
	if err != nil {
		tx.Rollback()
		return model.Run{}, fmt.Errorf("sql exec: %w", err)
	}

	var runID int
	err = tx.QueryRowContext(
		ctx,
		replaceRunSequenceName(createRunQuery, test.ID),
		test.ID,
		test.Version,
		run.CreatedAt,
		run.State,
		run.TraceID.String(),
		run.SpanID.String(),
		jsonMetadata,
		jsonEnvSnapshot,
	).Scan(&runID)
	if err != nil {
		tx.Rollback()
		return model.Run{}, fmt.Errorf("sql exec: %w", err)
	}

	tx.Commit()

	return td.GetRun(ctx, test.ID, runID)
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
	"env_snapshot" = $18

WHERE id = $16 AND test_id = $17
`

func (td *postgresDB) UpdateRun(ctx context.Context, r model.Run) error {
	stmt, err := td.db.Prepare(updateRunQuery)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	jsonTriggerResults, err := json.Marshal(r.TriggerResult)
	if err != nil {
		return fmt.Errorf("trigger results encoding error: %w", err)
	}

	jsonTestResults, err := json.Marshal(r.Results)
	if err != nil {
		return fmt.Errorf("test results encoding error: %w", err)
	}

	jsonTrace, err := json.Marshal(r.Trace)
	if err != nil {
		return fmt.Errorf("trace encoding error: %w", err)
	}

	jsonOutputs, err := json.Marshal(r.Outputs)
	if err != nil {
		return fmt.Errorf("outputs encoding error: %w", err)
	}

	jsonMetadata, err := json.Marshal(r.Metadata)
	if err != nil {
		return fmt.Errorf("encoding error: %w", err)
	}

	jsonEnvSnapshot, err := json.Marshal(r.EnvSnapshot)
	if err != nil {
		return fmt.Errorf("encoding error: %w", err)
	}

	var lastError *string
	if r.LastError != nil {
		e := r.LastError.Error()
		lastError = &e
	}

	pass, fail := count(r)

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
		jsonOutputs,
		lastError,
		pass,
		fail,
		jsonMetadata,
		r.ID,
		r.TestID,
		jsonEnvSnapshot,
	)
	if err != nil {
		return fmt.Errorf("sql exec: %w", err)
	}

	return nil
}

func count(r model.Run) (pass, fail int) {
	r.Results.Results.ForEach(func(_ model.SpanQuery, ars []model.AssertionResult) error {
		for _, ar := range ars {
			for _, rs := range ar.Results {
				if rs.CompareErr == nil {
					pass++
				} else {
					fail++
				}
			}
		}
		return nil
	})

	return
}

func (td *postgresDB) DeleteRun(ctx context.Context, r model.Run) error {
	stmt, err := td.db.Prepare("DELETE FROM test_runs WHERE id = $1 AND test_id = $2")
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, r.ID, r.TestID)
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
	"outputs",
	"last_error",

	"metadata",
	"env_snapshot"
FROM test_runs
`

func (td *postgresDB) GetRun(ctx context.Context, testID id.ID, runID int) (model.Run, error) {
	stmt, err := td.db.Prepare(selectRunQuery + " WHERE id = $1 AND test_id = $2")
	if err != nil {
		return model.Run{}, err
	}
	defer stmt.Close()

	run, err := readRunRow(stmt.QueryRowContext(ctx, runID, testID.String()))
	if err != nil {
		return model.Run{}, fmt.Errorf("cannot read row: %w", err)
	}
	return run, nil
}

func (td *postgresDB) GetTestRuns(ctx context.Context, test model.Test, take, skip int32) (model.List[model.Run], error) {
	const condition = " WHERE test_id = $1"
	stmt, err := td.db.Prepare(selectRunQuery + condition + " ORDER BY created_at DESC LIMIT $2 OFFSET $3")
	if err != nil {
		return model.List[model.Run]{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, test.ID, take, skip)
	if err != nil {
		return model.List[model.Run]{}, err
	}
	var runs []model.Run

	for rows.Next() {
		run, err := readRunRow(rows)
		if err != nil {
			return model.List[model.Run]{}, fmt.Errorf("cannot read row: %w", err)
		}
		runs = append(runs, run)
	}

	var count int
	err = td.db.
		QueryRowContext(ctx, "SELECT COUNT(*) FROM test_runs"+condition, test.ID).
		Scan(&count)
	if err != nil {
		return model.List[model.Run]{}, err
	}

	return model.List[model.Run]{
		Items:      runs,
		TotalCount: count,
	}, nil
}

func (td *postgresDB) GetRunByTraceID(ctx context.Context, traceID trace.TraceID) (model.Run, error) {
	stmt, err := td.db.Prepare(selectRunQuery + " WHERE trace_id = $1")
	if err != nil {
		return model.Run{}, err
	}
	defer stmt.Close()

	run, err := readRunRow(stmt.QueryRowContext(ctx, traceID.String()))
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
		jsonOutputs,
		jsonEnvSnapshot,
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
		&jsonOutputs,
		&lastError,
		&jsonMetadata,
		&jsonEnvSnapshot,
	)

	switch err {
	case sql.ErrNoRows:
		return model.Run{}, ErrNotFound
	case nil:
		err = json.Unmarshal(jsonTriggerResults, &r.TriggerResult)
		if err != nil {
			return model.Run{}, fmt.Errorf("cannot parse TriggerResult: %w", err)
		}

		err = json.Unmarshal(jsonTestResults, &r.Results)
		if err != nil {
			return model.Run{}, fmt.Errorf("cannot parse Results: %w", err)
		}

		if jsonTrace != nil {
			err = json.Unmarshal(jsonTrace, &r.Trace)
			if err != nil {
				return model.Run{}, fmt.Errorf("cannot parse Trace: %w", err)
			}
		}

		err = json.Unmarshal(jsonOutputs, &r.Outputs)
		if err != nil {
			return model.Run{}, fmt.Errorf("cannot parse Outputs: %w", err)
		}

		err = json.Unmarshal(jsonMetadata, &r.Metadata)
		if err != nil {
			return model.Run{}, fmt.Errorf("cannot parse Metadata: %w", err)
		}

		err = json.Unmarshal(jsonEnvSnapshot, &r.EnvSnapshot)
		if err != nil {
			return model.Run{}, fmt.Errorf("cannot parse EnvSnapshot: %w", err)
		}

		if traceID != "" {
			r.TraceID, err = trace.TraceIDFromHex(traceID)
			if err != nil {
				return model.Run{}, fmt.Errorf("cannot parse TraceID: %w", err)
			}
		}

		if spanID != "" {
			r.SpanID, err = trace.SpanIDFromHex(spanID)
			if err != nil {
				return model.Run{}, fmt.Errorf("cannot parse SpanID: %w", err)
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
