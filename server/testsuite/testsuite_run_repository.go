package testsuite

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

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
	encodedRun, err := EncodeRun(tr)
	if err != nil {
		return TestSuiteRun{}, fmt.Errorf("cannot encode TestSuiteRun: %w", err)
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

	params := sqlutil.TenantInsert(ctx,
		encodedRun.TestSuiteID,
		encodedRun.TestSuiteVersion,
		encodedRun.CreatedAt,
		encodedRun.State,
		encodedRun.CurrentTest,
		encodedRun.JsonMetadata,
		encodedRun.JsonVariableSet,
	)

	var runID int
	err = tx.QueryRowContext(
		ctx,
		replaceTestSuiteRunSequenceName(createTestSuiteRunQuery, tr.TestSuiteID),
		params...,
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

	encodedRun, err := EncodeRun(tr)
	if err != nil {
		return fmt.Errorf("cannot encode TestSuiteRun: %w", err)
	}

	pass, fail := tr.ResultsCount()
	allStepsRequiredGatesPassed := tr.StepsGatesValidation()

	query, params := sqlutil.Tenant(
		ctx,
		updateTestSuiteRunQuery,
		encodedRun.CompletedAt,
		encodedRun.State,
		encodedRun.CurrentTest,
		encodedRun.LastError,
		pass,
		fail,
		allStepsRequiredGatesPassed,
		encodedRun.JsonMetadata,
		encodedRun.JsonVariableSet,
		strconv.Itoa(encodedRun.ID),
		encodedRun.TestSuiteID,
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
	encodedRun := EncodedTestSuiteRun{}

	var lastErr *string

	err := row.Scan(
		&encodedRun.ID,
		&encodedRun.TestSuiteID,
		&encodedRun.TestSuiteVersion,
		&encodedRun.CreatedAt,
		&encodedRun.CompletedAt,
		&encodedRun.State,
		&encodedRun.CurrentTest,
		&lastErr,
		&encodedRun.Pass,
		&encodedRun.Fail,
		&encodedRun.AllStepsRequiredGatesPassed,
		&encodedRun.JsonMetadata,
		&encodedRun.JsonVariableSet,
	)

	if err != nil {
		return TestSuiteRun{}, fmt.Errorf("cannot read row: %w", err)
	}

	if lastErr != nil {
		encodedRun.LastError = *lastErr
	}

	return encodedRun.ToTestSuiteRun()
}

type EncodedTestSuiteRun struct {
	ID               int
	TestSuiteID      id.ID
	TestSuiteVersion int

	CreatedAt   time.Time
	CompletedAt time.Time

	State       TestSuiteRunState
	CurrentTest int

	JsonSteps,
	JsonStepIDs,
	JsonVariableSet,
	JsonMetadata []byte

	LastError string

	Pass,
	Fail sql.NullInt32

	AllStepsRequiredGatesPassed bool
}

func (r EncodedTestSuiteRun) ToTestSuiteRun() (TestSuiteRun, error) {
	tr := TestSuiteRun{
		ID:               r.ID,
		TestSuiteID:      r.TestSuiteID,
		TestSuiteVersion: r.TestSuiteVersion,

		CreatedAt:   r.CreatedAt,
		CompletedAt: r.CompletedAt,

		State:                       r.State,
		CurrentTest:                 r.CurrentTest,
		AllStepsRequiredGatesPassed: r.AllStepsRequiredGatesPassed,
	}

	if r.Pass.Valid {
		tr.Pass = int(r.Pass.Int32)
	}

	if r.Fail.Valid {
		tr.Fail = int(r.Fail.Int32)
	}

	err := json.Unmarshal(r.JsonMetadata, &tr.Metadata)
	if err != nil {
		return TestSuiteRun{}, fmt.Errorf("cannot parse Metadata: %w", err)
	}

	err = json.Unmarshal(r.JsonVariableSet, &tr.VariableSet)
	if err != nil {
		return TestSuiteRun{}, fmt.Errorf("cannot parse VariableSet: %w", err)
	}

	if !isValidSliceJSON(r.JsonSteps) {
		err = json.Unmarshal(r.JsonSteps, &tr.Steps)
		if err != nil {
			return TestSuiteRun{}, fmt.Errorf("cannot parse Steps: %w", err)
		}
	}

	if !isValidSliceJSON(r.JsonSteps) {
		err = json.Unmarshal(r.JsonStepIDs, &tr.StepIDs)
		if err != nil {
			return TestSuiteRun{}, fmt.Errorf("cannot parse StepIDs: %w", err)
		}
	}

	if r.LastError != "" {
		tr.LastError = fmt.Errorf(r.LastError)
	}

	return tr, nil
}

func isValidSliceJSON(j []byte) bool {
	return slices.Contains([]string{"[]", "", "null"}, string(j))
}

func EncodeRun(run TestSuiteRun) (EncodedTestSuiteRun, error) {
	jsonMetadata, err := json.Marshal(run.Metadata)
	if err != nil {
		return EncodedTestSuiteRun{}, fmt.Errorf("failed to marshal TestSuiteRun metadata: %w", err)
	}

	jsonVariableSet, err := json.Marshal(run.VariableSet)
	if err != nil {
		return EncodedTestSuiteRun{}, fmt.Errorf("failed to marshal TestSuiteRun variable set: %w", err)
	}

	jsonSteps, err := json.Marshal(run.Steps)
	if err != nil {
		return EncodedTestSuiteRun{}, fmt.Errorf("failed to marshal TestSuiteRun steps: %w", err)
	}

	jsonStepIDs, err := json.Marshal(run.StepIDs)
	if err != nil {
		return EncodedTestSuiteRun{}, fmt.Errorf("failed to marshal TestSuiteRun step IDs: %w", err)
	}

	var lastError string
	if run.LastError != nil {
		lastError = run.LastError.Error()
	}

	pass := sql.NullInt32{Int32: int32(run.Pass), Valid: true}
	fail := sql.NullInt32{Int32: int32(run.Fail), Valid: true}

	return EncodedTestSuiteRun{
		ID:               run.ID,
		TestSuiteID:      run.TestSuiteID,
		TestSuiteVersion: run.TestSuiteVersion,

		CreatedAt:   run.CreatedAt,
		CompletedAt: run.CompletedAt,

		State:       run.State,
		CurrentTest: run.CurrentTest,

		JsonSteps:       jsonSteps,
		JsonStepIDs:     jsonStepIDs,
		JsonMetadata:    jsonMetadata,
		JsonVariableSet: jsonVariableSet,

		LastError: lastError,

		Pass: pass,
		Fail: fail,

		AllStepsRequiredGatesPassed: run.AllStepsRequiredGatesPassed,
	}, nil
}
