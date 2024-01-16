package wizard

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/server/pkg/sqlutil"
)

type Repository interface {
	Get(context.Context) (*Wizard, error)
	Update(context.Context, *Wizard) error
}

type wizardRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &wizardRepository{
		db: db,
	}
}

type scanner interface {
	Scan(dest ...interface{}) error
}

var (
	defaultWizard = &Wizard{
		Steps: []Step{{
			ID:    "tracing_backend",
			State: StepStatusPending,
		}, {
			ID:    "create_test",
			State: StepStatusPending,
		}},
	}

	selectQuery = `SELECT steps FROM "wizards"`
	updateQuery = `UPDATE "wizards" SET steps = $1`
)

func (wr *wizardRepository) Get(ctx context.Context) (*Wizard, error) {
	query, params := sqlutil.Tenant(ctx, selectQuery)

	wizard, err := readRow(wr.db.QueryRowContext(ctx, query, params...))
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return defaultWizard, nil
	}

	return &wizard, nil
}

func (wr *wizardRepository) Update(ctx context.Context, wizard *Wizard) error {
	jsonSteps, err := json.Marshal(wizard.Steps)
	if err != nil {
		return fmt.Errorf("cannot marshal steps: %w", err)
	}

	query, params := sqlutil.Tenant(ctx, updateQuery, jsonSteps)

	_, err = wr.db.ExecContext(ctx, query, params...)
	if err != nil {
		return fmt.Errorf("cannot update wizard: %w", err)
	}

	return nil
}

func readRow(row scanner) (Wizard, error) {
	wizard := Wizard{
		Steps: []Step{},
	}

	var (
		jsonSteps []byte
	)

	err := row.Scan(&jsonSteps)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Wizard{}, err
		}

		return Wizard{}, fmt.Errorf("cannot read row: %w", err)
	}

	err = json.Unmarshal(jsonSteps, &wizard.Steps)
	return wizard, err
}
