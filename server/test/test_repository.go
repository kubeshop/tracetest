package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/sqlutil"
)

type Repository interface {
	List(_ context.Context, take, skip int, query, sortBy, sortDirection string) ([]Test, error)
	Count(context.Context, string) (int, error)
	SortingFields() []string
	Provision(context.Context, Test) error
	SetID(test Test, id id.ID) Test

	// TODO: uncomment as we add new functionality
	// Get(context.Context, id.ID) (Test, error)
	// Create(context.Context, Test) (Test, error)
	// Update(context.Context, Test) (Test, error)
	// Delete(context.Context, Test) error
	// Exists(context.Context, id.ID) (bool, error)
	// GetTestVersion(_ context.Context, _ id.ID, version int) (Test, error)
	// List(_ context.Context, take, skip int32, query, sortBy, sortDirection string) ([]Test, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

func (r *repository) SetID(test Test, id id.ID) Test {
	test.ID = id
	return test
}

func (r *repository) Provision(ctx context.Context, test Test) error {
	return nil
}

func (r *repository) SortingFields() []string {
	return []string{"created", "name", "last_run"}
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
		last_test_run.fail as last_test_run_fail
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

func (r *repository) List(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]Test, error) {
	sql := getTestSQL + testMaxVersionQuery
	params := []any{take, skip}

	const condition = " AND (t.name ilike $3 OR t.description ilike $3)"
	q, params := sqlutil.Search(sql, condition, query, params)

	sortingFields := map[string]string{
		"created":  "t.created_at",
		"name":     "t.name",
		"last_run": "last_test_run_time",
	}

	q = sqlutil.Sort(q, sortBy, sortDirection, "created", sortingFields)
	q += ` LIMIT $1 OFFSET $2 `

	stmt, err := r.db.Prepare(q)
	if err != nil {
		return []Test{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return []Test{}, err
	}

	tests, err := r.readRows(ctx, rows)
	if err != nil {
		return []Test{}, err
	}

	return tests, nil
}

func (r *repository) Count(ctx context.Context, query string) (int, error) {
	countQuery := "SELECT COUNT(*) FROM tests t" + testMaxVersionQuery

	if query != "" {
		countQuery = fmt.Sprintf("%s WHERE %s", countQuery, query)
	}

	count := 0

	err := r.db.
		QueryRowContext(ctx, countQuery).
		Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("sql query: %w", err)
	}

	return count, nil
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func (r *repository) readRows(ctx context.Context, rows *sql.Rows) ([]Test, error) {
	tests := []Test{}

	for rows.Next() {
		test, err := r.readRow(ctx, rows)
		if err != nil {
			return []Test{}, err
		}

		tests = append(tests, test)
	}

	return tests, nil
}

func (r *repository) readRow(ctx context.Context, row scanner) (Test, error) {
	test := Test{}

	var (
		jsonServiceUnderTest,
		jsonSpecs,
		jsonOutputs []byte

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
	)

	if err != nil {
		return Test{}, fmt.Errorf("cannot read row: %w", err)
	}

	err = json.Unmarshal(jsonServiceUnderTest, &test.Trigger)
	if err != nil {
		return Test{}, fmt.Errorf("cannot parse trigger: %w", err)
	}

	err = json.Unmarshal(jsonSpecs, &test.Specs)
	if err != nil {
		return Test{}, fmt.Errorf("cannot parse specs: %w", err)
	}

	err = json.Unmarshal(jsonOutputs, &test.Outputs)
	if err != nil {
		return Test{}, fmt.Errorf("cannot parse outputs: %w", err)
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

	return test, nil
}
