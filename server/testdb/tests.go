package testdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kubeshop/tracetest/server/model"
)

var _ model.TestRepository = &postgresDB{}

func (td *postgresDB) IDExists(ctx context.Context, id uuid.UUID) (bool, error) {
	exists := false

	row := td.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) > 0 as exists FROM tests WHERE id = $1",
		id.String(),
	)

	err := row.Scan(&exists)

	return exists, err
}

func (td *postgresDB) CreateTest(ctx context.Context, test model.Test) (model.Test, error) {
	if !test.HasID() {
		test.ID = IDGen.UUID()
	}

	test.CreatedAt = time.Now()
	test.ReferenceRun = nil
	test.Version = 1

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

func (td *postgresDB) UpdateTest(ctx context.Context, test model.Test) (model.Test, error) {
	if test.Version == 0 {
		test.Version = 1
	}

	oldTest, err := td.GetLatestTestVersion(ctx, test.ID)
	if err != nil {
		return model.Test{}, fmt.Errorf("could not get latest test version while updating test: %w", err)
	}

	// keep the same creation date to keep sort order
	test.CreatedAt = oldTest.CreatedAt

	testToUpdate, err := model.BumpTestVersionIfNeeded(oldTest, test)
	if err != nil {
		return model.Test{}, fmt.Errorf("could not bump test version: %w", err)
	}

	if oldTest.Version == testToUpdate.Version {
		// No change in the version, so nothing changes and it doesn't need to persist it
		return testToUpdate, nil
	}

	stmt, err := td.db.Prepare("INSERT INTO tests(id, test, version) VALUES( $1, $2, $3 )")
	if err != nil {
		return model.Test{}, fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()

	b, err := encodeTest(testToUpdate)
	if err != nil {
		return model.Test{}, fmt.Errorf("encoding error: %w", err)
	}
	_, err = stmt.ExecContext(ctx, testToUpdate.ID, b, testToUpdate.Version)
	if err != nil {
		return model.Test{}, fmt.Errorf("sql exec: %w", err)
	}

	err = td.SetDefiniton(ctx, testToUpdate, testToUpdate.Definition)
	if err != nil {
		return model.Test{}, fmt.Errorf("setDefinition error: %w", err)
	}

	return testToUpdate, nil
}

func (td *postgresDB) UpdateTestVersion(ctx context.Context, test model.Test) error {
	stmt, err := td.db.Prepare("UPDATE tests SET test = $2 WHERE id = $1 AND version = $3")
	if err != nil {
		return fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()

	b, err := encodeTest(test)
	if err != nil {
		return fmt.Errorf("encoding error: %w", err)
	}

	_, err = stmt.ExecContext(ctx, test.ID, b, test.Version)
	if err != nil {
		return fmt.Errorf("sql exec: %w", err)
	}

	return nil
}

func (td *postgresDB) DeleteTest(ctx context.Context, test model.Test) error {
	queries := []string{
		"DELETE FROM runs WHERE test_id = $1",
		"DELETE FROM definitions WHERE test_id = $1",
		"DELETE FROM tests WHERE id = $1",
	}

	for _, sql := range queries {
		_, err := td.db.ExecContext(ctx, sql, test.ID)
		if err != nil {
			return fmt.Errorf("sql error: %w", err)
		}
	}

	return nil
}

const getTestSQL = `
SELECT t.test, d.definition
FROM tests t
JOIN definitions d ON d.test_id = t.id AND d.test_version = t.version
`

func (td *postgresDB) GetTestVersion(ctx context.Context, id uuid.UUID, version int) (model.Test, error) {
	stmt, err := td.db.Prepare(getTestSQL + " WHERE t.id = $1 AND t.version = $2")
	if err != nil {
		return model.Test{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	test, err := td.readTestRow(ctx, stmt.QueryRowContext(ctx, id, version))
	if err != nil {
		return model.Test{}, err
	}

	return test, nil
}

func (td *postgresDB) GetLatestTestVersion(ctx context.Context, id uuid.UUID) (model.Test, error) {
	stmt, err := td.db.Prepare(getTestSQL + " WHERE t.id = $1 ORDER BY t.version DESC LIMIT 1")
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

func (td *postgresDB) GetTests(ctx context.Context, take, skip int32, query string) ([]model.Test, error) {
	hasSearchQuery := query != ""
	params := []any{take, skip}

	sql := getTestSQL + `
	INNER JOIN (
		SELECT id as idx, max(version) as latest_version FROM tests GROUP BY idx
	) as latestTests ON latestTests.idx = t.id
	WHERE t.version = latestTests.latest_version `
	if hasSearchQuery {
		params = append(params, "%"+strings.ReplaceAll(query, " ", "%")+"%")
		sql += ` AND (t.test ->> 'Name') ilike $3`
	}

	sql += ` ORDER BY (t.test ->> 'CreatedAt')::timestamp DESC LIMIT $1 OFFSET $2`

	stmt, err := td.db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, params...)
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
	var testB, defB []byte
	err := row.Scan(&testB, &defB)
	switch err {
	case sql.ErrNoRows:
		return model.Test{}, ErrNotFound
	case nil:
		test, err := decodeTest(testB)
		if err != nil {
			return model.Test{}, err
		}

		defs, err := decodeDefinition(defB)
		if err != nil {
			return model.Test{}, err
		}

		test.Definition = defs

		return test, nil
	default:
		return model.Test{}, err
	}
}
