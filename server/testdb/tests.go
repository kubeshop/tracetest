package testdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kubeshop/tracetest/model"
)

var _ model.TestRepository = &postgresDB{}

func (td *postgresDB) CreateTest(ctx context.Context, test model.Test) (model.Test, error) {
	test.ID = IDGen.UUID()
	test.ReferenceRun = nil
	test.Version = 1

	return td.createTest(ctx, test)
}

func (td *postgresDB) createTest(ctx context.Context, test model.Test) (model.Test, error) {
	stmt, err := td.db.Prepare("INSERT INTO tests(id, test, version) VALUES( $1, $2, $3 )")
	if err != nil {
		return model.Test{}, fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()

	b, err := encodeTest(test)
	if err != nil {
		return model.Test{}, fmt.Errorf("encoding error: %w", err)
	}
	_, err = stmt.ExecContext(ctx, test.ID, b, test.Version)
	if err != nil {
		return model.Test{}, fmt.Errorf("sql exec: %w", err)
	}

	err = td.SetDefiniton(ctx, test, test.Definition)
	if err != nil {
		return model.Test{}, fmt.Errorf("setDefinition error: %w", err)
	}

	return test, nil
}

func (td *postgresDB) UpdateTest(ctx context.Context, test model.Test) error {
	oldTest, err := td.GetLatestTestVersion(ctx, test.ID)
	if err != nil {
		return fmt.Errorf("could not get test latest version: %w", err)
	}

	test.Version = oldTest.Version + 1
	_, err = td.createTest(ctx, test)
	if err != nil {
		return fmt.Errorf("could not create new test version: %w", err)
	}

	return nil
}

func (td *postgresDB) DeleteTest(ctx context.Context, test model.Test) error {
	queries := []string{
		"DELETE FROM tests WHERE id = $1",
		"DELETE FROM definitions WHERE test_id = $1",
		"DELETE FROM runs WHERE test_id = $1",
	}

	for _, sql := range queries {
		_, err := td.db.Exec(sql, test.ID)
		if err != nil {
			return fmt.Errorf("sql error: %w", err)
		}
	}

	return nil
}

func (td *postgresDB) GetLatestTestVersion(ctx context.Context, id uuid.UUID) (model.Test, error) {
	stmt, err := td.db.Prepare("SELECT test FROM tests WHERE id = $1 ORDER BY version DESC LIMIT 1")
	if err != nil {
		return model.Test{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	test, err := td.readTestRow(ctx, stmt.QueryRowContext(ctx, id))
	if err != nil {
		return model.Test{}, err
	}

	return test, nil
}

func (td *postgresDB) GetTests(ctx context.Context, take, skip int32) ([]model.Test, error) {
	stmt, err := td.db.Prepare(`SELECT test FROM tests INNER JOIN (
		SELECT id as idx, max(version) as latest_version FROM tests GROUP BY idx
	) as latestTests ON latestTests.idx = id WHERE version = latestTests.latest_version LIMIT $1 OFFSET $2`)
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
		test, err := td.readTestRow(ctx, rows)
		if err != nil {
			return nil, err
		}

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

func (td *postgresDB) readTestRow(ctx context.Context, row scanner) (model.Test, error) {
	var b []byte
	err := row.Scan(&b)
	switch err {
	case sql.ErrNoRows:
		return model.Test{}, ErrNotFound
	case nil:
		test, err := decodeTest(b)
		if err != nil {
			return model.Test{}, err
		}

		defs, err := td.GetDefiniton(ctx, test)
		if err != nil && !errors.Is(err, ErrNotFound) {
			err = fmt.Errorf("aca 1. %w", err)
			return model.Test{}, err
		}
		test.Definition = defs

		return test, nil
	default:
		return model.Test{}, err
	}
}
