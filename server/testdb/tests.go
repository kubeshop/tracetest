package testdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/kubeshop/tracetest/model"
)

var _ model.TestRepository = &postgresDB{}

func (td *postgresDB) CreateTest(ctx context.Context, test model.Test) (model.Test, error) {
	stmt, err := td.db.Prepare("INSERT INTO tests(id, test) VALUES( $1, $2 )")
	if err != nil {
		return model.Test{}, fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()

	test.ID = IDGen.UUID()

	b, err := encodeTest(test)
	if err != nil {
		return model.Test{}, fmt.Errorf("encoding error: %w", err)
	}
	_, err = stmt.ExecContext(ctx, test.ID, b)
	if err != nil {
		return model.Test{}, fmt.Errorf("sql exec: %w", err)
	}

	return test, nil
}

func (td *postgresDB) UpdateTest(ctx context.Context, test model.Test) error {
	stmt, err := td.db.Prepare("UPDATE tests SET test = $2 WHERE id = $1")
	if err != nil {
		return fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()

	b, err := encodeTest(test)
	if err != nil {
		return fmt.Errorf("encoding error: %w", err)
	}

	_, err = stmt.Exec(test.ID, b)
	if err != nil {
		return fmt.Errorf("sql exec: %w", err)
	}

	return nil
}

func (td *postgresDB) DeleteTest(ctx context.Context, test model.Test) error {
	queries := []string{
		"DELETE FROM tests WHERE id = $1",
		"DELETE FROM definitions WHERE test_id = $1",
		"DELETE FROM results WHERE test_id = $1",
	}

	for _, sql := range queries {
		_, err := td.db.Exec(sql, test.ID)
		if err != nil {
			return fmt.Errorf("sql error: %w", err)
		}
	}

	return nil
}

func (td *postgresDB) GetTest(ctx context.Context, id uuid.UUID) (model.Test, error) {
	stmt, err := td.db.Prepare("SELECT test FROM tests WHERE id = $1")
	if err != nil {
		return model.Test{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	test, err := readTestRow(stmt.QueryRowContext(ctx, id))
	if err != nil {
		return model.Test{}, err
	}

	defs, err := td.GetDefiniton(ctx, test)
	if err != nil {
		return model.Test{}, err
	}
	test.Definition = defs

	return test, nil
}

func (td *postgresDB) GetTests(ctx context.Context, take, skip int32) ([]model.Test, error) {
	stmt, err := td.db.Prepare("SELECT test FROM tests LIMIT $1 OFFSET $2")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, take, skip)
	if err != nil {
		return nil, err
	}

	tests := []model.Test{}

	for rows.Next() {
		test, err := readTestRow(rows)
		if err != nil {
			return nil, err
		}

		defs, err := td.GetDefiniton(ctx, test)
		if err != nil {
			return nil, err
		}
		test.Definition = defs

		tests = append(tests, test)
	}
	return tests, nil
}

func encodeTest(t model.Test) (string, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return "", fmt.Errorf("json Marshal: %w", err)
	}

	return string(b), nil
}

func decodeTest(b []byte) (model.Test, error) {
	var test model.Test
	err := json.Unmarshal(b, &test)
	if err != nil {
		return model.Test{}, err
	}
	return test, nil
}

func readTestRow(row scanner) (model.Test, error) {
	var b []byte
	err := row.Scan(&b)
	switch err {
	case sql.ErrNoRows:
		return model.Test{}, ErrNotFound
	case nil:
		return decodeTest(b)
	default:
		return model.Test{}, err
	}
}
