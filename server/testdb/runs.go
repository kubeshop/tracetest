package testdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/kubeshop/tracetest/server/model"
	"go.opentelemetry.io/otel/trace"
)

var _ model.RunRepository = &postgresDB{}

func (td *postgresDB) CreateRun(ctx context.Context, test model.Test, run model.Run) (model.Run, error) {
	stmt, err := td.db.Prepare("INSERT INTO runs(id, test_id, test_version, run) VALUES( $1, $2, $3, $4 )")
	if err != nil {
		return model.Run{}, fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()

	run.ID = IDGen.UUID()
	run.State = model.RunStateCreated
	run.TestVersion = test.Version

	encoded, err := encodeRun(run)
	if err != nil {
		return model.Run{}, fmt.Errorf("encoding error: %w", err)
	}

	_, err = stmt.ExecContext(ctx, run.ID, test.ID, run.TestVersion, encoded)
	if err != nil {
		return model.Run{}, fmt.Errorf("sql exec: %w", err)
	}

	return run, nil
}

func (td *postgresDB) UpdateRun(ctx context.Context, r model.Run) error {
	stmt, err := td.db.Prepare("UPDATE runs SET run = $2 WHERE id = $1")
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	def, err := encodeRun(r)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, r.ID, def)
	if err != nil {
		return fmt.Errorf("sql exec: %w", err)
	}

	return nil
}

func (td *postgresDB) GetRun(ctx context.Context, id uuid.UUID) (model.Run, error) {
	stmt, err := td.db.Prepare("SELECT run FROM runs WHERE id = $1")
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

func (td *postgresDB) GetTestRuns(ctx context.Context, test model.Test, take, skip int32) ([]model.Run, error) {
	stmt, err := td.db.Prepare("SELECT run FROM runs WHERE test_id = $1 ORDER BY run ->> 'createdAt' DESC LIMIT $2 OFFSET $3")
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

func (td *postgresDB) GetRunByTraceID(ctx context.Context, test model.Test, traceID trace.TraceID) (model.Run, error) {
	stmt, err := td.db.Prepare("SELECT run FROM runs WHERE test_id = $1 AND run ->> 'traceId' = $2")
	if err != nil {
		return model.Run{}, err
	}
	defer stmt.Close()

	run, err := readRunRow(stmt.QueryRowContext(ctx, test.ID, traceID))
	if err != nil {
		return model.Run{}, fmt.Errorf("cannot read row: %w", err)
	}
	return run, nil
}

func encodeRun(r model.Run) (string, error) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", fmt.Errorf("json Marshal: %w", err)
	}

	return string(b), nil
}

func decodeRun(b []byte) (model.Run, error) {
	var run model.Run
	err := json.Unmarshal(b, &run)
	if err != nil {
		return model.Run{}, fmt.Errorf("unmarshal run: %w", err)
	}
	return run, nil
}

func readRunRow(row scanner) (model.Run, error) {
	var b []byte
	err := row.Scan(&b)
	switch err {
	case sql.ErrNoRows:
		return model.Run{}, ErrNotFound
	case nil:
		return decodeRun(b)
	default:
		return model.Run{}, fmt.Errorf("read run row: %w", err)
	}
}
