package configresource

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kubeshop/tracetest/server/id"
)

type Config struct {
	ID   id.ID  `mapstructure:"id"`
	Name string `mapstructure:"name"`

	AnalyticsEnabled bool `mapstructure:"analyticsEnabled"`
}

func (c Config) HasID() bool {
	return c.ID.String() != ""
}

func (c Config) Validate() error {
	return nil
}

func Repository(db *sql.DB, idgen id.GeneratorFunc) *repository {
	return &repository{db, idgen}
}

type repository struct {
	db    *sql.DB
	idgen id.GeneratorFunc
}

const insertQuery = `
		INSERT INTO configs (
			"id",
			"name",
			"analytics_enabled"
		) VALUES ($1, $2, $3)`

func (r *repository) SetID(cfg Config, id id.ID) Config {
	cfg.ID = id
	return cfg
}

func (r *repository) Create(ctx context.Context, cfg Config) (Config, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Config{}, err
	}

	_, err = tx.ExecContext(ctx, insertQuery,
		cfg.ID,
		cfg.Name,
		cfg.AnalyticsEnabled,
	)

	if err != nil {
		tx.Rollback()
		return Config{}, fmt.Errorf("sql exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return Config{}, fmt.Errorf("commit: %w", err)
	}

	return cfg, nil
}

const updateQuery = `
		UPDATE configs SET
			"name" = $2,
			"analytics_enabled" = $3
		WHERE "id" = $1`

func (r *repository) Update(ctx context.Context, updated Config) (Config, error) {
	cfg, err := r.Get(ctx, updated.ID)
	if err != nil {
		return Config{}, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Config{}, err
	}

	_, err = tx.ExecContext(ctx, updateQuery,
		cfg.ID,
		updated.Name,
		updated.AnalyticsEnabled,
	)

	if err != nil {
		tx.Rollback()
		return Config{}, fmt.Errorf("sql exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return Config{}, fmt.Errorf("commit: %w", err)
	}

	return updated, nil
}

const getQuery = `
		SELECT
			"name",
			"analytics_enabled"
		FROM configs
		WHERE "id" = $1`

func (r *repository) Get(ctx context.Context, id id.ID) (Config, error) {
	cfg := Config{
		ID: id,
	}

	err := r.db.
		QueryRowContext(ctx, getQuery, id).
		Scan(
			&cfg.Name,
			&cfg.AnalyticsEnabled,
		)

	if err != nil {
		return Config{}, fmt.Errorf("sql query: %w", err)
	}

	return cfg, nil

}

const deleteQuery = `
		DELETE FROM configs
		WHERE "id" = $1`

func (r *repository) Delete(ctx context.Context, id id.ID) error {
	cfg, err := r.Get(ctx, id)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, deleteQuery, cfg.ID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("sql exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}
