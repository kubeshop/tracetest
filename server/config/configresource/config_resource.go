package configresource

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kubeshop/tracetest/server/id"
)

type Config struct {
	ID   string `mapstructure:"id"`
	Name string `mapstructure:"name"`

	AnalyticsEnabled bool `mapstructure:"analyticsEnabled"`
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

// TODO: add context.Context
func (r *repository) Create(cfg Config) (Config, error) {
	ctx := context.TODO()

	cfg.ID = r.idgen().String()
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
