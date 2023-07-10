package testrunner

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/server/pkg/id"
)

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

type Repository struct {
	db *sql.DB
}

func (r *Repository) SetID(testRunner TestRunner, id id.ID) TestRunner {
	testRunner.ID = id
	return testRunner
}

func (r *Repository) Create(ctx context.Context, updated TestRunner) (TestRunner, error) {
	updated.ID = id.ID("current")
	return r.Update(ctx, updated)
}

const (
	insertQuery = `
	INSERT INTO test_runners(
		"id",
		"name",
		"required_gates"
	)
	VALUES ($1, $2, $3)`
	deleteQuery = `DELETE FROM test_runners`
)

func (r *Repository) Update(ctx context.Context, updated TestRunner) (TestRunner, error) {
	// enforce ID and default
	updated.ID = "current"

	tx, err := r.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return TestRunner{}, err
	}

	_, err = tx.ExecContext(ctx, deleteQuery)
	if err != nil {
		return TestRunner{}, fmt.Errorf("sql exec delete: %w", err)
	}

	var requiredGatesJSON []byte
	if updated.RequiredGates != nil {
		requiredGatesJSON, err = json.Marshal(updated.RequiredGates)
		if err != nil {
			return TestRunner{}, fmt.Errorf("could not marshal periodic strategy configuration: %w", err)
		}
	}

	_, err = tx.ExecContext(ctx, insertQuery,
		updated.ID,
		updated.Name,
		requiredGatesJSON,
	)
	if err != nil {
		return TestRunner{}, fmt.Errorf("sql exec insert: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return TestRunner{}, fmt.Errorf("commit: %w", err)
	}

	return updated, nil

}

const (
	getQuery = `
		SELECT
			"name",
			"required_gates"
		FROM test_runners`
)

func (r *Repository) GetDefault(ctx context.Context) TestRunner {
	tr, _ := r.Get(ctx, id.ID("current"))
	return tr
}

func (r *Repository) Get(ctx context.Context, id id.ID) (TestRunner, error) {
	testRunner := TestRunner{
		ID: "current",
	}

	var requiredGatesJSON []byte
	err := r.db.
		QueryRowContext(ctx, getQuery).
		Scan(
			&testRunner.Name,
			&requiredGatesJSON,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DefaultTestRunner, nil
		}
		return TestRunner{}, fmt.Errorf("sql query: %w", err)
	}

	if string(requiredGatesJSON) != "null" {
		var requiredGates RequiredGates
		err = json.Unmarshal(requiredGatesJSON, &requiredGates)
		if err != nil {
			return TestRunner{}, fmt.Errorf("could not unmarshal periodic strategy config: %w", err)
		}

		testRunner.RequiredGates = requiredGates
	}

	return testRunner, nil
}

func (r *Repository) List(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]TestRunner, error) {
	testRunner, err := r.Get(ctx, id.ID("current"))
	if err != nil {
		return []TestRunner{}, err
	}

	return []TestRunner{testRunner}, nil
}

func (r *Repository) Count(ctx context.Context, query string) (int, error) {
	return 1, nil
}

func (*Repository) SortingFields() []string {
	return []string{"name"}
}

func (r *Repository) Provision(ctx context.Context, testRunner TestRunner) error {
	_, err := r.Update(ctx, testRunner)
	return err
}
