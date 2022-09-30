package testdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
)

var _ model.TestRepository = &postgresDB{}

func (td *postgresDB) IDExists(ctx context.Context, id id.ID) (bool, error) {
	exists := false

	row := td.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) > 0 as exists FROM tests WHERE id = $1",
		id,
	)

	err := row.Scan(&exists)

	return exists, err
}

const insertIntoTestsQuery = `
INSERT INTO tests (
	"id",
	"version",
	"name",
	"description",
	"service_under_test",
	"specs",
	"created_at"
) VALUES ($1, $2, $3, $4, $5, $6, $7)`

func (td *postgresDB) CreateTest(ctx context.Context, test model.Test) (model.Test, error) {
	if !test.HasID() {
		test.ID = IDGen.ID()
	}

	test.Version = 1
	test.CreatedAt = time.Now()

	return td.insertIntoTests(ctx, test)
}

func (td *postgresDB) insertIntoTests(ctx context.Context, test model.Test) (model.Test, error) {
	stmt, err := td.db.Prepare(insertIntoTestsQuery)
	if err != nil {
		return model.Test{}, fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()

	jsonServiceUnderTest, err := json.Marshal(test.ServiceUnderTest)
	if err != nil {
		return model.Test{}, fmt.Errorf("encoding error: %w", err)
	}

	jsonSpecs, err := json.Marshal(test.Specs)
	if err != nil {
		return model.Test{}, fmt.Errorf("encoding error: %w", err)
	}

	_, err = stmt.ExecContext(
		ctx,
		test.ID,
		test.Version,
		test.Name,
		test.Description,
		jsonServiceUnderTest,
		jsonSpecs,
		test.CreatedAt,
	)
	if err != nil {
		return model.Test{}, fmt.Errorf("sql exec: %w", err)
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

	return td.insertIntoTests(ctx, testToUpdate)
}

func (td *postgresDB) DeleteTest(ctx context.Context, test model.Test) error {
	queries := []string{
		"DELETE FROM test_runs WHERE test_id = $1",
		"DELETE FROM tests WHERE id = $1",
	}

	tx, err := td.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("sql BeginTx: %w", err)
	}

	for _, sql := range queries {
		_, err := tx.ExecContext(ctx, sql, test.ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("sql error: %w", err)
		}
	}

	dropSequece(ctx, tx, test.ID)

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("sql Commit: %w", err)
	}

	return nil
}

const getTestSQL = `
	SELECT
		t.id,
		t.version,
		t.name,
		t.description,
		t.service_under_test,
		t.specs,
		t.created_at,
		(SELECT COUNT(*) FROM test_runs tr WHERE tr.test_id = t.id) as total_runs
	FROM tests t
`

func (td *postgresDB) GetTestVersion(ctx context.Context, id id.ID, version int) (model.Test, error) {
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

func (td *postgresDB) GetLatestTestVersion(ctx context.Context, id id.ID) (model.Test, error) {
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
		sql += ` AND (
			t.name ilike $3
			OR t.description ilike $3
		)`
	}

	sql += ` ORDER BY t.created_at DESC LIMIT $1 OFFSET $2`

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

func (td *postgresDB) testLastRunSummary(ctx context.Context, test model.Test) (model.Test, error) {
	runs, err := td.GetTestRuns(ctx, test, 1, 0) //get latest run
	if err != nil {
		return model.Test{}, fmt.Errorf("cannot get test last run: %w", err)
	}

	if len(runs) == 0 {
		return test, nil
	}

	lastRun := runs[0]

	test.Summary.LastRun.Time = lastRun.CreatedAt

	test.Summary.LastRun.Passes,
		test.Summary.LastRun.Fails = lastRun.Count()

	return test, nil
}

func (td *postgresDB) readTestRow(ctx context.Context, row scanner) (model.Test, error) {
	test := model.Test{}

	var jsonServiceUnderTest, jsonSpecs []byte
	err := row.Scan(
		&test.ID,
		&test.Version,
		&test.Name,
		&test.Description,
		&jsonServiceUnderTest,
		&jsonSpecs,
		&test.CreatedAt,
		&test.Summary.Runs,
	)

	switch err {
	case sql.ErrNoRows:
		return model.Test{}, ErrNotFound
	case nil:
		err = json.Unmarshal(jsonServiceUnderTest, &test.ServiceUnderTest)
		if err != nil {
			return model.Test{}, err
		}

		err = json.Unmarshal(jsonSpecs, &test.Specs)
		if err != nil {
			return model.Test{}, err
		}

		return td.testLastRunSummary(ctx, test)
	default:
		return model.Test{}, err
	}
}
