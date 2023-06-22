package testdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/transaction"
)

var _ model.TestRepository = &postgresDB{}

func (td *postgresDB) TestIDExists(ctx context.Context, id id.ID) (bool, error) {
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
	"outputs",
	"created_at"
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

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

	jsonOutputs, err := json.Marshal(test.Outputs)
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
		jsonOutputs,
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
		// No change in the version. Nothing changed so no need to persist it
		return testToUpdate, nil
	}

	return td.insertIntoTests(ctx, testToUpdate)
}

func (td *postgresDB) DeleteTest(ctx context.Context, test model.Test) error {
	queries := []string{
		"DELETE FROM transaction_run_steps WHERE test_run_test_id = $1",
		"DELETE FROM transaction_steps WHERE test_id = $1",
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

const (
	getTestSQL = `
	SELECT
		t.id,
		t.version,
		t.name,
		t.description,
		t.service_under_test,
		t.specs,
		t.outputs,
		t.created_at,
		(SELECT COUNT(*) FROM test_runs tr WHERE tr.test_id = t.id) as total_runs,
		last_test_run.created_at as last_test_run_time,
		last_test_run.pass as last_test_run_pass,
		last_test_run.fail as last_test_run_fail,
		last_test_run.linter as last_test_run_linter
	FROM tests t
	LEFT OUTER JOIN (
		SELECT MAX(id) as id, test_id FROM test_runs GROUP BY test_id
	) as ltr ON ltr.test_id = t.id
	LEFT OUTER JOIN
		test_runs last_test_run
	ON last_test_run.test_id = ltr.test_id AND last_test_run.id = ltr.id
`

	testMaxVersionQuery = `
	INNER JOIN (
		SELECT id as idx, max(version) as latest_version FROM tests GROUP BY idx
		) as latest_tests ON latest_tests.idx = t.id AND t.version = latest_tests.latest_version
	`
)

func sortQuery(sql, sortBy, sortDirection string, sortingFields map[string]string) string {
	sortField, ok := sortingFields[sortBy]

	if !ok {
		sortField = sortingFields["created"]
	}

	dir := "DESC"
	if strings.ToLower(sortDirection) == "asc" {
		dir = "ASC"
	}

	return fmt.Sprintf("%s ORDER BY %s %s", sql, sortField, dir)
}

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

func (td *postgresDB) GetTests(ctx context.Context, take, skip int32, query, sortBy, sortDirection string) (model.List[model.Test], error) {
	hasSearchQuery := query != ""
	cleanSearchQuery := "%" + strings.ReplaceAll(query, " ", "%") + "%"
	params := []any{take, skip}

	sql := getTestSQL + testMaxVersionQuery

	const condition = " AND (t.name ilike $3 OR t.description ilike $3)"
	if hasSearchQuery {
		params = append(params, cleanSearchQuery)
		sql += condition
	}

	sortingFields := map[string]string{
		"created":  "t.created_at",
		"name":     "t.name",
		"last_run": "last_test_run_time",
	}

	sql = sortQuery(sql, sortBy, sortDirection, sortingFields)
	sql += ` LIMIT $1 OFFSET $2 `

	stmt, err := td.db.Prepare(sql)
	if err != nil {
		return model.List[model.Test]{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return model.List[model.Test]{}, err
	}

	tests, err := td.readTestRows(ctx, rows)
	if err != nil {
		return model.List[model.Test]{}, err
	}

	count, err := td.count(ctx, condition, cleanSearchQuery)
	if err != nil {
		return model.List[model.Test]{}, err
	}

	return model.List[model.Test]{
		Items:      tests,
		TotalCount: count,
	}, nil
}

func (td *postgresDB) count(ctx context.Context, condition, cleanSearchQuery string) (int, error) {
	var (
		count  int
		params []any
	)
	countQuery := "SELECT COUNT(*) FROM tests t" + testMaxVersionQuery
	if cleanSearchQuery != "" {
		params = []any{cleanSearchQuery}
		countQuery += strings.ReplaceAll(condition, "$3", "$1")
	}

	err := td.db.
		QueryRowContext(ctx, countQuery, params...).
		Scan(&count)

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (td *postgresDB) readTestRows(ctx context.Context, rows *sql.Rows) ([]model.Test, error) {
	tests := []model.Test{}

	for rows.Next() {
		test, err := td.readTestRow(ctx, rows)
		if err != nil {
			return []model.Test{}, err
		}

		tests = append(tests, test)
	}

	return tests, nil
}

func (td *postgresDB) readTestRow(ctx context.Context, row scanner) (model.Test, error) {
	test := model.Test{}

	var (
		jsonServiceUnderTest,
		jsonSpecs,
		jsonOutputs,
		jsonLinter []byte

		lastRunTime *time.Time

		pass, fail *int
	)
	err := row.Scan(
		&test.ID,
		&test.Version,
		&test.Name,
		&test.Description,
		&jsonServiceUnderTest,
		&jsonSpecs,
		&jsonOutputs,
		&test.CreatedAt,
		&test.Summary.Runs,
		&lastRunTime,
		&pass,
		&fail,
		&jsonLinter,
	)

	switch err {
	case nil:
		err = json.Unmarshal(jsonServiceUnderTest, &test.ServiceUnderTest)
		if err != nil {
			return model.Test{}, fmt.Errorf("cannot parse trigger: %w", err)
		}

		err = json.Unmarshal(jsonSpecs, &test.Specs)
		if err != nil {
			return model.Test{}, fmt.Errorf("cannot parse specs: %w", err)
		}

		err = json.Unmarshal(jsonOutputs, &test.Outputs)
		if err != nil {
			return model.Test{}, fmt.Errorf("cannot parse outputs: %w", err)
		}

		if lastRunTime != nil {
			test.Summary.LastRun.Time = *lastRunTime
		}
		if pass != nil {
			test.Summary.LastRun.Passes = *pass
		}
		if fail != nil {
			test.Summary.LastRun.Fails = *fail
		}

		var linter model.LinterResult
		err = json.Unmarshal(jsonLinter, &linter)
		if err == nil {
			test.Summary.LastRun.AnalyzerScore = linter.Score
		}

		return test, nil
	case sql.ErrNoRows:
		return model.Test{}, ErrNotFound
	default:
		return model.Test{}, err
	}
}

func (td *postgresDB) GetTransactionSteps(ctx context.Context, transaction transaction.Transaction) ([]model.Test, error) {
	stmt, err := td.db.Prepare(getTestSQL + testMaxVersionQuery + ` INNER JOIN transaction_steps ts ON t.id = ts.test_id
	 WHERE ts.transaction_id = $1 AND ts.transaction_version = $2 ORDER BY ts.step_number ASC`)
	if err != nil {
		return []model.Test{}, fmt.Errorf("prepare 2: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, transaction.ID, transaction.Version)
	if err != nil {
		return []model.Test{}, fmt.Errorf("query context: %w", err)
	}

	steps, err := td.readTestRows(ctx, rows)
	if err != nil {
		return []model.Test{}, fmt.Errorf("read row: %w", err)
	}

	return steps, nil
}
