package testsuite

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/sqlutil"
	"github.com/kubeshop/tracetest/server/test"
)

type testSuiteStepRunRepository interface {
	GetTestSuiteRunSteps(_ context.Context, _ id.ID, runID int) ([]test.Run, error)
}

func NewRunRepository(db *sql.DB, stepsRepository testSuiteStepRunRepository) *RunRepository {
	return &RunRepository{
		db:              db,
		stepsRepository: stepsRepository,
	}
}

type RunRepository struct {
	db              *sql.DB
	stepsRepository testSuiteStepRunRepository
}

const createTestSuiteRunQuery = `
INSERT INTO test_suite_runs (
	"id",
	"test_suite_id",
	"test_suite_version",

	-- timestamps
	"created_at",
	"completed_at",

	-- trigger params
	"state",
	"current_test",

	-- result info
	"last_error",
	"pass",
	"fail",
	"all_steps_required_gates_passed",

	"metadata",

	-- variable_set
	"variable_set",

	"tenant_id"
) VALUES (
	nextval('` + runSequenceName + `'), -- id
	$1,   -- test_suite_id
	$2,   -- test_suite_version

	-- timestamps
	$3,              -- created_at
	to_timestamp(0), -- completed_at

	-- trigger params
	$4, -- state
	$5, -- currentStep

	-- result info
	NULL, -- last_error
	0,    -- pass
	0,    -- fail
	TRUE, -- all_steps_required_gates_passed

	$6, -- metadata
	$7, -- variable_set
	$8 -- tenant_id
)
RETURNING "id"`

const (
	createSequenceQuery = `CREATE SEQUENCE IF NOT EXISTS "` + runSequenceName + `";`
	dropSequenceQuery   = `DROP SEQUENCE IF EXISTS "` + runSequenceName + `";`
	runSequenceName     = "%sequence_name%"
)

func md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func replaceTestSuiteRunSequenceName(sql string, ID id.ID) string {
	// postgres doesn't like uppercase chars in sequence names.
	// transactionID might contain uppercase chars, and we cannot lowercase them
	// because they might lose their uniqueness.
	// md5 creates a unique, lowercase hash.
	seqName := "runs_test_suite_" + md5Hash(ID.String()) + "_seq"
	return strings.ReplaceAll(sql, runSequenceName, seqName)
}

func (td *RunRepository) CreateRun(ctx context.Context, tr TestSuiteRun) (TestSuiteRun, error) {
	jsonMetadata, err := json.Marshal(tr.Metadata)
	if err != nil {
		return TestSuiteRun{}, fmt.Errorf("failed to marshal test_suite run metadata: %w", err)
	}

	jsonVariableSet, err := json.Marshal(tr.VariableSet)
	if err != nil {
		return TestSuiteRun{}, fmt.Errorf("failed to marshal test_suite run variable set: %w", err)
	}

	tx, err := td.db.BeginTx(ctx, nil)
	if err != nil {
		return TestSuiteRun{}, fmt.Errorf("sql beginTx: %w", err)
	}

	_, err = tx.ExecContext(ctx, replaceTestSuiteRunSequenceName(createSequenceQuery, tr.TestSuiteID))
	if err != nil {
		tx.Rollback()
		return TestSuiteRun{}, fmt.Errorf("sql exec: %w", err)
	}

	tenantID := sqlutil.TenantID(ctx)

	var runID int
	err = tx.QueryRowContext(
		ctx,
		replaceTestSuiteRunSequenceName(createTestSuiteRunQuery, tr.TestSuiteID),
		tr.TestSuiteID,
		tr.TestSuiteVersion,
		tr.CreatedAt,
		tr.State,
		tr.CurrentTest,
		jsonMetadata,
		jsonVariableSet,
		tenantID,
	).Scan(&runID)
	if err != nil {
		tx.Rollback()
		return TestSuiteRun{}, fmt.Errorf("sql exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return TestSuiteRun{}, fmt.Errorf("commit: %w", err)
	}

	tr.ID = runID

	return tr, nil
}

const updateTestSuiteRunQuery = `
UPDATE test_suite_runs SET

	-- timestamps
	"completed_at" = $1,

	-- trigger params
	"state" = $2,
	"current_test" = $3,

	-- result info
	"last_error" = $4,
	"pass" = $5,
	"fail" = $6,
	"all_steps_required_gates_passed" = $7,

	"metadata" = $8,

	-- variable_set
	"variable_set" = $9

WHERE id = $10 AND test_suite_id = $11
`

func (td *RunRepository) UpdateRun(ctx context.Context, tr TestSuiteRun) error {
	tx, err := td.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("sql beginTx: %w", err)
	}

	jsonMetadata, err := json.Marshal(tr.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal test_suite run metadata: %w", err)
	}

	jsonVariableSet, err := json.Marshal(tr.VariableSet)
	if err != nil {
		return fmt.Errorf("failed to marshal test_suite run variableSet: %w", err)
	}
	var lastError *string
	if tr.LastError != nil {
		e := tr.LastError.Error()
		lastError = &e
	}

	pass, fail := tr.ResultsCount()
	allStepsRequiredGatesPassed := tr.StepsGatesValidation()

	query, params := sqlutil.Tenant(
		ctx,
		updateTestSuiteRunQuery,
		tr.CompletedAt,
		tr.State,
		tr.CurrentTest,
		lastError,
		pass,
		fail,
		allStepsRequiredGatesPassed,
		jsonMetadata,
		jsonVariableSet,
		strconv.Itoa(tr.ID),
		tr.TestSuiteID,
	)
	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, params...)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("sql exec: %w", err)
	}

	return td.setTestSuiteRunSteps(ctx, tx, tr)
}

func (td *RunRepository) setTestSuiteRunSteps(ctx context.Context, tx *sql.Tx, tr TestSuiteRun) error {
	// delete existing steps
	query, params := sqlutil.Tenant(ctx, "DELETE FROM test_suite_run_steps WHERE test_suite_run_id = $1 AND test_suite_run_test_suite_id = $2", strconv.Itoa(tr.ID), tr.TestSuiteID)
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, params...)
	if err != nil {
		return fmt.Errorf("cannot reset test_suite run steps: %w", err)
	}

	if len(tr.Steps) == 0 {
		return tx.Commit()
	}

	tenantID := sqlutil.TenantID(ctx)

	values := []string{}
	for _, run := range tr.Steps {
		if run.ID == 0 {
			// step not set, skip
			continue
		}

		if tenantID == nil {
			values = append(
				values,
				fmt.Sprintf("('%d', '%s', %d, '%s', NULL)", tr.ID, tr.TestSuiteID, run.ID, run.TestID),
			)
		} else {
			values = append(
				values,
				fmt.Sprintf("('%d', '%s', %d, '%s', '%s')", tr.ID, tr.TestSuiteID, run.ID, run.TestID, *tenantID),
			)
		}
	}

	sql := "INSERT INTO test_suite_run_steps VALUES " + strings.Join(values, ", ")
	_, err = tx.ExecContext(ctx, sql)
	if err != nil {
		return fmt.Errorf("cannot save test_suite run steps: %w", err)
	}
	return tx.Commit()
}

func (td *RunRepository) DeleteTestSuiteRun(ctx context.Context, tr TestSuiteRun) error {
	tx, err := td.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("sql beginTx: %w", err)
	}

	query, params := sqlutil.Tenant(ctx, "DELETE FROM test_suite_run_steps WHERE test_suite_run_id = $1 AND test_suite_run_test_suite_id = $2", tr.ID, tr.TestSuiteID)
	_, err = tx.ExecContext(ctx, query, params...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("delete test_suite run steps: %w", err)
	}

	query, params = sqlutil.Tenant(ctx, "DELETE FROM test_suite_runs WHERE id = $1 AND test_suite_id = $2", tr.ID, tr.TestSuiteID)
	_, err = tx.ExecContext(ctx, query, params...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("delete test_suite runs: %w", err)
	}

	return tx.Commit()
}

const selectTestSuiteRunQuery = `
SELECT
	"id",
	"test_suite_id",
	"test_suite_version",

	"created_at",
	"completed_at",

	"state",
	"current_test",

	"last_error",
	"pass",
	"fail",
	"all_steps_required_gates_passed",

	"metadata",

	"variable_set"
FROM test_suite_runs
`

func (td *RunRepository) GetTestSuiteRun(ctx context.Context, ID id.ID, runID int) (TestSuiteRun, error) {
	query, params := sqlutil.Tenant(ctx, selectTestSuiteRunQuery+" WHERE id = $1 AND test_suite_id = $2", strconv.Itoa(runID), ID)
	stmt, err := td.db.Prepare(query)
	if err != nil {
		return TestSuiteRun{}, fmt.Errorf("prepare: %w", err)
	}

	run, err := td.readRunRow(stmt.QueryRowContext(ctx, params...))
	if err != nil {
		return TestSuiteRun{}, err
	}
	run.Steps, err = td.stepsRepository.GetTestSuiteRunSteps(ctx, run.TestSuiteID, run.ID)
	if err != nil {
		return TestSuiteRun{}, err
	}
	return run, nil
}

func (td *RunRepository) GetLatestRunByTestSuiteVersion(ctx context.Context, ID id.ID, version int) (TestSuiteRun, error) {
	sortQuery := " ORDER BY created_at DESC LIMIT 1"
	query, params := sqlutil.Tenant(ctx, selectTestSuiteRunQuery+" WHERE test_suite_id = $1 AND test_suite_version = $2", ID, version)
	stmt, err := td.db.Prepare(query + sortQuery)
	if err != nil {
		return TestSuiteRun{}, fmt.Errorf("prepare: %w", err)
	}

	run, err := td.readRunRow(stmt.QueryRowContext(ctx, params...))
	if err != nil {
		return TestSuiteRun{}, err
	}
	run.Steps, err = td.stepsRepository.GetTestSuiteRunSteps(ctx, run.TestSuiteID, run.ID)
	if err != nil {
		return TestSuiteRun{}, err
	}
	return run, nil
}

func (td *RunRepository) GetTestSuiteRuns(ctx context.Context, ID id.ID, take, skip int32) ([]TestSuiteRun, error) {
	sortQuery := " ORDER BY created_at DESC LIMIT $2 OFFSET $3"
	query, params := sqlutil.Tenant(ctx, selectTestSuiteRunQuery+" WHERE test_suite_id = $1", ID.String(), take, skip)
	stmt, err := td.db.Prepare(query + sortQuery)
	if err != nil {
		return []TestSuiteRun{}, fmt.Errorf("prepare: %w", err)
	}

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return []TestSuiteRun{}, fmt.Errorf("query: %w", err)
	}

	var runs []TestSuiteRun
	for rows.Next() {
		run, err := td.readRunRow(rows)
		if err != nil {
			return []TestSuiteRun{}, err
		}

		run.Steps, err = td.stepsRepository.GetTestSuiteRunSteps(ctx, run.TestSuiteID, run.ID)
		if err != nil {
			return []TestSuiteRun{}, err
		}

		runs = append(runs, run)
	}

	return runs, nil
}

func (td *RunRepository) readRunRow(row scanner) (TestSuiteRun, error) {
	r := TestSuiteRun{}

	var (
		jsonVariableSet,
		jsonMetadata []byte

		lastError *string
		pass      sql.NullInt32
		fail      sql.NullInt32

		allStepsRequiredGatesPassed *bool
	)

	err := row.Scan(
		&r.ID,
		&r.TestSuiteID,
		&r.TestSuiteVersion,
		&r.CreatedAt,
		&r.CompletedAt,
		&r.State,
		&r.CurrentTest,
		&lastError,
		&pass,
		&fail,
		&allStepsRequiredGatesPassed,
		&jsonMetadata,
		&jsonVariableSet,
	)
	if err != nil {
		return TestSuiteRun{}, fmt.Errorf("cannot read row: %w", err)
	}

	err = json.Unmarshal(jsonMetadata, &r.Metadata)
	if err != nil {
		return TestSuiteRun{}, fmt.Errorf("cannot parse Metadata: %w", err)
	}

	err = json.Unmarshal(jsonVariableSet, &r.VariableSet)
	if err != nil {
		return TestSuiteRun{}, fmt.Errorf("cannot parse VariableSet: %w", err)
	}

	if lastError != nil && *lastError != "" {
		r.LastError = fmt.Errorf(*lastError)
	}

	if pass.Valid {
		r.Pass = int(pass.Int32)
	}

	if fail.Valid {
		r.Fail = int(fail.Int32)
	}

	// checks if the flag exists, if it doesn't we use the fail field to determine if all steps passed
	if allStepsRequiredGatesPassed == nil {
		failed := r.Fail == 0
		allStepsRequiredGatesPassed = &failed
	}

	r.AllStepsRequiredGatesPassed = *allStepsRequiredGatesPassed

	return r, nil
}
