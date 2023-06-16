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
	ListAugmented(_ context.Context, take, skip int, query, sortBy, sortDirection string) ([]Test, error)
	Count(context.Context, string) (int, error)
	SortingFields() []string
	Provision(context.Context, Test) error
	SetID(test Test, id id.ID) Test
	Create(context.Context, Test) (Test, error)
	Get(context.Context, id.ID) (Test, error)
	Update(context.Context, Test) (Test, error)

	// TODO: uncomment as we add new functionality
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
	tests, err := r.list(ctx, take, skip, query, sortBy, sortDirection)
	if err != nil {
		return []Test{}, err
	}

	for i, test := range tests {
		test.CreatedAt = nil
		test.Summary = nil
		test.Version = nil
		tests[i] = test
	}

	return tests, err

}

func (r *repository) ListAugmented(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]Test, error) {
	return r.list(ctx, take, skip, query, sortBy, sortDirection)
}

func (r *repository) list(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]Test, error) {
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

const insertQuery = `
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

func (r *repository) Create(ctx context.Context, test Test) (Test, error) {
	if !test.HasID() {
		test.ID = IDGen.ID()
	}

	_, err := r.Get(ctx, test.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			test.Version = 1
			test.CreatedAt = time.Now()

			return r.insertTest(ctx, test)
		}

		return Test{}, err
	}

	return r.Update(ctx, test)
}

func (r *repository) insertTest(ctx context.Context, test Test) (Test, error) {
	stmt, err := r.db.Prepare(insertQuery)
	if err != nil {
		return Test{}, fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()

	triggerJson, err := json.Marshal(test.Trigger)
	if err != nil {
		return Test{}, fmt.Errorf("encoding error: %w", err)
	}

	specsJson, err := json.Marshal(test.Specs)
	if err != nil {
		return Test{}, fmt.Errorf("encoding error: %w", err)
	}

	outputsJson, err := json.Marshal(test.Outputs)
	if err != nil {
		return Test{}, fmt.Errorf("encoding error: %w", err)
	}

	_, err = stmt.ExecContext(
		ctx,
		test.ID,
		test.Version,
		test.Name,
		test.Description,
		triggerJson,
		specsJson,
		outputsJson,
		test.CreatedAt,
	)
	if err != nil {
		return Test{}, fmt.Errorf("sql exec: %w", err)
	}

	return test, nil
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
	version := 0
	createdAt := time.Now()

	test := Test{
		CreatedAt: &createdAt,
		Version:   &version,
		Summary:   &Summary{},
	}

	var (
		jsonTrigger,
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
		&jsonTrigger,
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

	err = json.Unmarshal(jsonTrigger, &test.Trigger)
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

func (r *repository) Get(ctx context.Context, id id.ID) (Test, error) {
	stmt, err := r.db.Prepare(getTestSQL + " WHERE t.id = $1 ORDER BY t.version DESC LIMIT 1")
	if err != nil {
		return Test{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	test, err := r.readRow(ctx, stmt.QueryRowContext(ctx, id))
	if err != nil {
		return Test{}, err
	}

	return test, nil
}

func (r *repository) Update(ctx context.Context, test Test) (Test, error) {
	if test.Version == 0 {
		test.Version = 1
	}

	oldTest, err := r.Get(ctx, test.ID)
	if err != nil {
		return Test{}, fmt.Errorf("could not get latest test version while updating test: %w", err)
	}

	// keep the same creation date to keep sort order
	test.CreatedAt = oldTest.CreatedAt

	testToUpdate, err := bumpTestVersionIfNeeded(oldTest, test)
	if err != nil {
		return Test{}, fmt.Errorf("could not bump test version: %w", err)
	}

	if oldTest.Version == testToUpdate.Version {
		// No change in the version. Nothing changed so no need to persist it
		return testToUpdate, nil
	}

	return r.insertTest(ctx, testToUpdate)
}

func bumpTestVersionIfNeeded(in, updated Test) (Test, error) {
	testHasChanged, err := testHasChanged(in, updated)
	if err != nil {
		return Test{}, err
	}

	if testHasChanged {
		updated.Version = in.Version + 1
	}

	return updated, nil
}

func testHasChanged(oldTest Test, newTest Test) (bool, error) {
	outputsHaveChanged, err := testFieldHasChanged(oldTest.Outputs, newTest.Outputs)
	if err != nil {
		return false, err
	}

	definitionHasChanged, err := testFieldHasChanged(oldTest.Specs, newTest.Specs)
	if err != nil {
		return false, err
	}

	serviceUnderTestHasChanged, err := testFieldHasChanged(oldTest.Trigger, newTest.Trigger)
	if err != nil {
		return false, err
	}

	nameHasChanged := oldTest.Name != newTest.Name
	descriptionHasChanged := oldTest.Description != newTest.Description

	return outputsHaveChanged || definitionHasChanged || serviceUnderTestHasChanged || nameHasChanged || descriptionHasChanged, nil
}

func testFieldHasChanged(oldField interface{}, newField interface{}) (bool, error) {
	oldFieldJSON, err := json.Marshal(oldField)
	if err != nil {
		return false, err
	}

	newFieldJSON, err := json.Marshal(newField)
	if err != nil {
		return false, err
	}

	return string(oldFieldJSON) != string(newFieldJSON), nil
}
