package lintern_resource

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

var defaultLintern = Lintern{
	ID:           id.ID("current"),
	Name:         "Config",
	Enabled:      true,
	MinimumScore: 80,
	Plugins: []LinternPlugin{
		{Name: "standards", Enabled: true, Required: true},
		{Name: "security", Enabled: true, Required: true},
	},
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

const (
	insertQuery = `
		INSERT INTO linterns (
			"id",
			"name",
			"enabled",
			"minimum_score"
			"plugins"
		) VALUES ($1, $2, $3, $4, $5)`

	getQuery = `
	SELECT
		l.id,
		l.name,
		l.enabled,
		l.minimum_score,
		l.plugins
	FROM linterns l
`

	deleteQuery = "DELETE FROM linterns WHERE id = $1"
)

func (*Repository) SortingFields() []string {
	return []string{"name"}
}

func (r *Repository) SetID(linterns Lintern, id id.ID) Lintern {
	linterns.ID = id
	return linterns
}

func (r *Repository) Update(ctx context.Context, lintern Lintern) (Lintern, error) {
	// enforce ID and name
	updated := Lintern{
		ID:           id.ID("current"),
		Name:         "Lintern",
		Enabled:      lintern.Enabled,
		MinimumScore: lintern.MinimumScore,
		Plugins:      lintern.Plugins,
	}

	tx, err := r.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return Lintern{}, err
	}

	_, err = tx.ExecContext(ctx, deleteQuery)
	if err != nil {
		return Lintern{}, fmt.Errorf("sql exec delete: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		insertQuery,
		updated.ID,
		updated.Name,
		updated.Enabled,
		updated.MinimumScore,
		updated.Plugins,
	)
	if err != nil {
		return Lintern{}, fmt.Errorf("sql exec insert: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return Lintern{}, fmt.Errorf("commit: %w", err)
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

func (r *Repository) List(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]Lintern, error) {
	lintern, err := r.Get(ctx, id.ID("current"))
	if err != nil {
		return []Lintern{}, err
	}

	return []Lintern{lintern}, nil
}

func (r *Repository) Count(ctx context.Context, query string) (int, error) {
	return 1, nil
}

func (r *Repository) Get(ctx context.Context, id id.ID) (Lintern, error) {
	lintern := defaultLintern

	var rawPlugins []byte
	err := r.db.
		QueryRowContext(ctx, getQuery).
		Scan(
			&lintern.ID,
			&lintern.Name,
			&lintern.Enabled,
			&lintern.MinimumScore,
			&rawPlugins,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lintern, nil
		}
		return Lintern{}, fmt.Errorf("sql query: %w", err)
	}

	if string(rawPlugins) != "null" {
		var plugins []LinternPlugin
		err = json.Unmarshal(rawPlugins, &plugins)
		if err != nil {
			return Lintern{}, fmt.Errorf("could not unmarshal plugins: %w", err)
		}

		lintern.Plugins = plugins
	}

	return lintern, nil
}

func (r *Repository) Provision(ctx context.Context, lintern Lintern) error {
	_, err := r.Update(ctx, lintern)
	return err
}
