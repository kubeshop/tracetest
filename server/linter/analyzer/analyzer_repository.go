package analyzer

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/server/pkg/id"
)

type Repository struct {
	db *sql.DB
}

var defaultlinter = Linter{
	ID:           id.ID("current"),
	Name:         "analyzer",
	Enabled:      true,
	MinimumScore: 0,
	Plugins: []LinterPlugin{
		{Name: "standards", Enabled: true, Required: true},
		{Name: "security", Enabled: true, Required: true},
		{Name: "common", Enabled: true, Required: true},
	},
}

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
			"plugins"
		) VALUES ($1, $2, $3, $4, $5)`

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

	_, err = tx.ExecContext(ctx, deleteQuery, updated.ID)
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

	_, err = tx.ExecContext(
		ctx,
		insertQuery,
		updated.ID,
		updated.Name,
		updated.Enabled,
		updated.MinimumScore,
		pluginsJSON,
	)
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

	_, err = tx.ExecContext(ctx, deleteQuery, id)
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
	linter := defaultlinter

	var rawPlugins []byte
	err := r.db.
		QueryRowContext(ctx, getQuery).
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
