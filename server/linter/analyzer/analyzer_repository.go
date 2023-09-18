package analyzer

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/sqlutil"
)

type Repository struct {
	db *sql.DB
}

var defaultLinter = GetDefaultLinter()

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

const (
	insertQuery = `
		INSERT INTO linters (
			"id",
			"name",
			"enabled",
			"minimum_score",
			"plugins",
			"tenant_id"
		) VALUES ($1, $2, $3, $4, $5, $6)`

	getQuery = `
	SELECT
		l.id,
		l.name,
		l.enabled,
		l.minimum_score,
		l.plugins
	FROM linters l
`

	deleteQuery = "DELETE FROM linters WHERE id = $1"
)

func (r *Repository) GetDefault(ctx context.Context) Linter {
	pp, _ := r.Get(ctx, id.ID("current"))
	return pp
}

func (r *Repository) SetID(linters Linter, id id.ID) Linter {
	linters.ID = id
	return linters
}

func (r *Repository) Create(ctx context.Context, linter Linter) (Linter, error) {
	linter.ID = id.ID("current")
	return r.Update(ctx, linter)
}

func (r *Repository) Update(ctx context.Context, linter Linter) (Linter, error) {
	// enforce ID and name
	updated := Linter{
		ID:           id.ID("current"),
		Name:         linter.Name,
		Enabled:      linter.Enabled,
		MinimumScore: linter.MinimumScore,
		Plugins:      linter.Plugins,
	}

	tx, err := r.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return Linter{}, err
	}

	query, params := sqlutil.Tenant(ctx, deleteQuery, updated.ID)
	_, err = tx.ExecContext(ctx, query, params...)
	if err != nil {
		return Linter{}, fmt.Errorf("sql exec delete: %w", err)
	}

	pluginsJSON := []byte("[]")
	if updated.Plugins != nil {
		pluginsJSON, err = json.Marshal(updated.Plugins)
		if err != nil {
			return Linter{}, fmt.Errorf("could not marshal plugins configuration: %w", err)
		}
	}

	params = sqlutil.TenantInsert(ctx,
		updated.ID,
		updated.Name,
		updated.Enabled,
		updated.MinimumScore,
		pluginsJSON,
	)
	_, err = tx.ExecContext(ctx, insertQuery, params...)
	if err != nil {
		return Linter{}, fmt.Errorf("sql exec insert: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return Linter{}, fmt.Errorf("commit: %w", err)
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id id.ID) error {
	_, err := r.Get(ctx, id)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("sql BeginTx: %w", err)
	}
	defer tx.Rollback()

	query, params := sqlutil.Tenant(ctx, deleteQuery, id)
	_, err = tx.ExecContext(ctx, query, params...)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("sql Commit: %w", err)
	}

	return nil
}

func (r *Repository) List(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]Linter, error) {
	linter, err := r.Get(ctx, id.ID("current"))
	if err != nil {
		return []Linter{}, err
	}

	return []Linter{linter}, nil
}

func (r *Repository) Count(ctx context.Context, query string) (int, error) {
	return 1, nil
}

func (*Repository) SortingFields() []string {
	return []string{"name"}
}

func (r *Repository) Get(ctx context.Context, id id.ID) (Linter, error) {
	linter := defaultLinter

	var rawPlugins []byte
	query, params := sqlutil.Tenant(ctx, getQuery)
	err := r.db.
		QueryRowContext(ctx, query, params...).
		Scan(
			&linter.ID,
			&linter.Name,
			&linter.Enabled,
			&linter.MinimumScore,
			&rawPlugins,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return linter, nil
		}
		return Linter{}, fmt.Errorf("sql query: %w", err)
	}

	if string(rawPlugins) != "null" {
		var plugins []LinterPlugin
		err = json.Unmarshal(rawPlugins, &plugins)
		if err != nil {
			return Linter{}, fmt.Errorf("could not unmarshal plugins: %w", err)
		}

		linter.Plugins = plugins
	}

	return linter, nil
}

func (r *Repository) Provision(ctx context.Context, linter Linter) error {
	_, err := r.Update(ctx, linter)
	return err
}
