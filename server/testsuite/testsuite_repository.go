package testsuite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/sqlutil"
	"github.com/kubeshop/tracetest/server/test"
)

func NewRepository(db *sql.DB, stepRepository testSuiteStepRepository) *Repository {
	repo := &Repository{
		db:             db,
		stepRepository: stepRepository,
	}

	return repo
}

type testSuiteStepRepository interface {
	GetTestSuiteSteps(_ context.Context, _ id.ID, version int) ([]test.Test, error)
}

type Repository struct {
	db             *sql.DB
	stepRepository testSuiteStepRepository
}

// needed for test
func (r *Repository) DB() *sql.DB {
	return r.db
}

func (r *Repository) SetID(t TestSuite, id id.ID) TestSuite {
	t.ID = id
	return t
}

func (r *Repository) Provision(ctx context.Context, testSuite TestSuite) error {
	_, err := r.Create(ctx, testSuite)

	return err
}

func (r *Repository) Create(ctx context.Context, testSuite TestSuite) (TestSuite, error) {
	setVersion(&testSuite, 1)
	if testSuite.CreatedAt == nil {
		setCreatedAt(&testSuite, time.Now())
	}

	t, err := r.insertIntoTestSuites(ctx, testSuite)
	if err != nil {
		return TestSuite{}, err
	}

	removeNonAugmentedFields(&t)
	return t, nil
}

func (r *Repository) Update(ctx context.Context, testSuite TestSuite) (TestSuite, error) {
	oldTestSuite, err := r.GetLatestVersion(ctx, testSuite.ID)
	if err != nil {
		return TestSuite{}, fmt.Errorf("could not get latest test version while updating test: %w", err)
	}

	if testSuite.GetVersion() == 0 {
		setVersion(&testSuite, oldTestSuite.GetVersion())
	}

	// keep the same creation date to keep sort order
	testSuite.CreatedAt = oldTestSuite.CreatedAt
	transactionToUpdate := BumpTestSuiteVersionIfNeeded(oldTestSuite, testSuite)

	if oldTestSuite.GetVersion() == transactionToUpdate.GetVersion() {
		// No change in the version, so nothing changes and it doesn't need to persist it
		return transactionToUpdate, nil
	}

	t, err := r.insertIntoTestSuites(ctx, transactionToUpdate)
	if err != nil {
		return TestSuite{}, err
	}

	removeNonAugmentedFields(&t)
	return t, nil
}

func (r *Repository) checkIDExists(ctx context.Context, id id.ID) error {
	exists, err := r.IDExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id id.ID) error {
	if err := r.checkIDExists(ctx, id); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	query, params := sqlutil.Tenant(ctx, "DELETE FROM test_suite_steps WHERE test_suite_id = $1", id)
	_, err = tx.ExecContext(ctx, query, params...)
	if err != nil {
		tx.Rollback()
		return err
	}

	q, params := sqlutil.Tenant(ctx, "DELETE FROM test_suite_run_steps WHERE test_suite_run_id IN (SELECT id FROM test_suite_runs WHERE test_suite_id = $1)", id)
	_, err = tx.ExecContext(ctx, q, params...)
	if err != nil {
		tx.Rollback()
		return err
	}

	q, params = sqlutil.Tenant(ctx, "DELETE FROM test_suite_runs WHERE test_suite_id = $1", id)
	_, err = tx.ExecContext(ctx, q, params...)
	if err != nil {
		tx.Rollback()
		return err
	}

	q, params = sqlutil.Tenant(ctx, "DELETE FROM test_suites WHERE id = $1", id)
	_, err = tx.ExecContext(ctx, q, params...)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *Repository) IDExists(ctx context.Context, id id.ID) (bool, error) {
	exists := false
	query, params := sqlutil.Tenant(ctx, "SELECT COUNT(*) > 0 as exists FROM test_suites WHERE id = $1", id)
	row := r.db.QueryRowContext(ctx, query, params...)

	err := row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("cannot check id existance: %w", err)
	}

	return exists, nil
}

const getTestSuiteSQL = `
	SELECT
		%s
	FROM test_suites t
	LEFT OUTER JOIN (
		SELECT MAX(id) as id, test_suite_id FROM test_suite_runs GROUP BY test_suite_id
	) as ltr ON ltr.test_suite_id = t.id
	LEFT OUTER JOIN
		test_suite_runs last_test_suite_run
	ON last_test_suite_run.test_suite_id = ltr.test_suite_id AND last_test_suite_run.id = ltr.id
`

const testSuiteMaxVersionQuery = `
	INNER JOIN (
		SELECT id as idx, max(version) as latest_version FROM test_suites GROUP BY idx
	) as latest_test_suites ON latest_test_suites.idx = t.id

	WHERE t.version = latest_test_suites.latest_version `

func (r *Repository) SortingFields() []string {
	return []string{
		"name",
		"created",
		"last_run",
	}
}

func (r *Repository) ListAugmented(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]TestSuite, error) {
	return r.list(ctx, take, skip, query, sortBy, sortDirection, true)
}

func (r *Repository) List(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]TestSuite, error) {
	return r.list(ctx, take, skip, query, sortBy, sortDirection, false)
}

func (r *Repository) list(ctx context.Context, take, skip int, query, sortBy, sortDirection string, augmented bool) ([]TestSuite, error) {
	q, params := listQuery(ctx, querySelect(), query, []any{take, skip})

	sortingFields := map[string]string{
		"created":  "t.created_at",
		"name":     "t.name",
		"last_run": "last_test_suite_run_time",
	}

	q = sqlutil.Sort(q, sortBy, sortDirection, "created", sortingFields)
	q += " LIMIT $1 OFFSET $2"

	stmt, err := r.db.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return nil, err
	}

	suites, err := r.readRows(ctx, rows, augmented)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return suites, nil
}

func queryCount() string {
	return fmt.Sprintf(getTestSuiteSQL, "COUNT(*)")
}

func querySelect() string {
	return fmt.Sprintf(getTestSuiteSQL, `
		t.id,
		t.version,
		t.name,
		t.description,
		t.created_at,
		(
			SELECT step_ids_query.ids FROM (
				SELECT ts.test_suite_id, ts.test_suite_version, string_agg(ts.test_id, ',') as ids FROM test_suite_steps ts
				GROUP BY test_suite_id, test_suite_version
				HAVING ts.test_suite_id = t.id AND ts.test_suite_version = t.version
			) as step_ids_query
		) as step_ids,
		(SELECT COUNT(*) FROM test_suite_runs tr WHERE tr.test_suite_id = t.id) as total_runs,
		last_test_suite_run.created_at as last_test_suite_run_time,
		last_test_suite_run.pass as last_test_run_pass,
		last_test_suite_run.fail as last_test_run_fail
	`)
}

func (r *Repository) Count(ctx context.Context, query string) (int, error) {
	sql, params := listQuery(ctx, queryCount(), query, []any{})

	count := 0
	err := r.db.
		QueryRowContext(ctx, sql, params...).
		Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) GetAugmented(ctx context.Context, id id.ID) (TestSuite, error) {
	return r.get(ctx, id, true)
}

func (r *Repository) Get(ctx context.Context, id id.ID) (TestSuite, error) {
	return r.get(ctx, id, false)
}

func (r *Repository) get(ctx context.Context, id id.ID, augmented bool) (TestSuite, error) {
	query, params := sqlutil.TenantWithPrefix(ctx, querySelect()+" WHERE t.id = $1", "t.", id)
	stmt, err := r.db.Prepare(query + "ORDER BY t.version DESC LIMIT 1")
	if err != nil {
		return TestSuite{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	testSuite, err := r.readRow(ctx, stmt.QueryRowContext(ctx, params...), augmented)
	if err != nil {
		return TestSuite{}, err
	}

	return testSuite, nil
}

func (r *Repository) GetVersion(ctx context.Context, id id.ID, version int) (TestSuite, error) {
	query, params := sqlutil.TenantWithPrefix(ctx, querySelect()+" WHERE t.id = $1 AND t.version = $2", "t.", id, version)
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return TestSuite{}, fmt.Errorf("prepare 1: %w", err)
	}
	defer stmt.Close()

	testSuite, err := r.readRow(ctx, stmt.QueryRowContext(ctx, params...), true)
	if err != nil {
		return TestSuite{}, err
	}

	return testSuite, nil
}

func listQuery(ctx context.Context, baseSQL, query string, params []any) (string, []any) {
	paramNumber := len(params) + 1
	condition := fmt.Sprintf(" AND (t.name ilike $%d OR t.description ilike $%d)", paramNumber, paramNumber)

	sql := baseSQL + testSuiteMaxVersionQuery
	sql, params = sqlutil.Search(sql, condition, query, params)
	sql, params = sqlutil.TenantWithPrefix(ctx, sql, "t.", params...)

	return sql, params
}

func (r *Repository) GetLatestVersion(ctx context.Context, id id.ID) (TestSuite, error) {
	query, params := sqlutil.TenantWithPrefix(ctx, querySelect()+" WHERE t.id = $1", "t.", id)
	stmt, err := r.db.Prepare(query + " ORDER BY t.version DESC LIMIT 1")
	if err != nil {
		return TestSuite{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	testSuite, err := r.readRow(ctx, stmt.QueryRowContext(ctx, params...), true)
	if err != nil {
		return TestSuite{}, err
	}

	return testSuite, nil
}

const insertIntoTestSuiteQuery = `
INSERT INTO test_suites (
	"id",
	"version",
	"name",
	"description",
	"created_at",
	"tenant_id"
) VALUES ($1, $2, $3, $4, $5, $6)`

func (r *Repository) insertIntoTestSuites(ctx context.Context, suite TestSuite) (TestSuite, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	defer tx.Rollback()
	if err != nil {
		return TestSuite{}, fmt.Errorf("sql begin: %w", err)
	}

	stmt, err := tx.Prepare(insertIntoTestSuiteQuery)
	if err != nil {
		return TestSuite{}, fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()

	tenantID := sqlutil.TenantID(ctx)

	_, err = stmt.ExecContext(
		ctx,
		suite.ID,
		suite.GetVersion(),
		suite.Name,
		suite.Description,
		suite.GetCreatedAt(),
		tenantID,
	)
	if err != nil {
		return TestSuite{}, fmt.Errorf("sql exec: %w", err)
	}

	return r.setTestSuiteSteps(ctx, tx, suite)
}

func (r *Repository) setTestSuiteSteps(ctx context.Context, tx *sql.Tx, suite TestSuite) (TestSuite, error) {
	// delete existing steps
	query, params := sqlutil.Tenant(ctx, "DELETE FROM test_suite_steps WHERE test_suite_id = $1 AND test_suite_version = $2", suite.ID, suite.GetVersion())
	stmt, err := tx.Prepare(query)
	if err != nil {
		return TestSuite{}, err
	}

	_, err = stmt.ExecContext(ctx, params...)
	if err != nil {
		return TestSuite{}, err
	}

	if len(suite.StepIDs) == 0 {
		return suite, tx.Commit()
	}

	tenantID := sqlutil.TenantID(ctx)

	values := []string{}
	for i, testID := range suite.StepIDs {
		stepNumber := i + 1

		if tenantID == nil {
			values = append(
				values,
				fmt.Sprintf("('%s', %d, '%s', %d, NULL)", suite.ID, suite.GetVersion(), testID, stepNumber),
			)
		} else {
			values = append(
				values,
				fmt.Sprintf("('%s', %d, '%s', %d, '%s')", suite.ID, suite.GetVersion(), testID, stepNumber, *tenantID),
			)
		}
	}

	sql := "INSERT INTO test_suite_steps VALUES " + strings.Join(values, ", ")
	_, err = tx.ExecContext(ctx, sql)
	if err != nil {
		return TestSuite{}, fmt.Errorf("cannot save test suite steps: %w", err)
	}
	return suite, tx.Commit()
}

func (r *Repository) readRows(ctx context.Context, rows *sql.Rows, augmented bool) ([]TestSuite, error) {
	transactions := []TestSuite{}

	for rows.Next() {
		transaction, err := r.readRow(ctx, rows, augmented)
		if err != nil {
			return []TestSuite{}, fmt.Errorf("cannot read rows: %w", err)
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func (r *Repository) readRow(ctx context.Context, row scanner, augmented bool) (TestSuite, error) {
	suite := TestSuite{
		Summary: &test.Summary{},
	}

	var (
		lastRunTime *time.Time
		stepIDs     *string

		pass, fail *int
		version    int
	)

	err := row.Scan(
		&suite.ID,
		&version,
		&suite.Name,
		&suite.Description,
		&suite.CreatedAt,
		&stepIDs,
		&suite.Summary.Runs,
		&lastRunTime,
		&pass,
		&fail,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return suite, err
		}

		return TestSuite{}, fmt.Errorf("cannot read row: %w", err)
	}

	if stepIDs != nil && *stepIDs != "" {
		ids := strings.Split(*stepIDs, ",")
		suite.StepIDs = make([]id.ID, len(ids))
		for i, sid := range ids {
			suite.StepIDs[i] = id.ID(sid)
		}
	}

	if version != 0 {
		suite.Version = &version
	}

	if lastRunTime != nil {
		suite.Summary.LastRun.Time = *lastRunTime
	}
	if pass != nil {
		suite.Summary.LastRun.Passes = *pass
	}
	if fail != nil {
		suite.Summary.LastRun.Fails = *fail
	}

	if !augmented {
		removeNonAugmentedFields(&suite)
	} else {
		steps, err := r.stepRepository.GetTestSuiteSteps(ctx, suite.ID, *suite.Version)
		if err != nil {
			return TestSuite{}, fmt.Errorf("cannot read row: %w", err)
		}

		suite.Steps = steps
	}

	return suite, nil
}

func removeNonAugmentedFields(t *TestSuite) {
	t.CreatedAt = nil
	t.Version = nil
	t.Summary = nil
}

func BumpTestSuiteVersionIfNeeded(original, updated TestSuite) TestSuite {
	transactionHasChanged := testSuiteHasChanged(original, updated)
	if transactionHasChanged {
		setVersion(&updated, original.GetVersion()+1)
	}

	return updated
}

func testSuiteHasChanged(in, updated TestSuite) bool {
	jsons := []struct {
		Name        string
		Description string
		Steps       []string
	}{
		{
			Name:        in.Name,
			Description: in.Description,
			Steps:       stepIDsToStringSlice(in),
		},
		{
			Name:        updated.Name,
			Description: updated.Description,
			Steps:       stepIDsToStringSlice(updated),
		},
	}

	inJson, _ := json.Marshal(jsons[0])
	updatedJson, _ := json.Marshal(jsons[1])

	return string(inJson) != string(updatedJson)
}

func stepIDsToStringSlice(in TestSuite) []string {
	steps := make([]string, len(in.StepIDs))
	for i, stepID := range in.StepIDs {
		steps[i] = stepID.String()
	}

	return steps
}
