package configresource

import (
	"context"
	"database/sql"
	"fmt"
	"os"

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

func (c Config) IsAnalyticsEnabled() bool {
	if os.Getenv("TRACETEST_DEV") != "" {
		return false
	}

	return c.AnalyticsEnabled
}

type option func(*Repository)

func WithPublisher(p publisher) option {
	return func(r *Repository) {
		r.publisher = p
	}
}

func NewRepository(db *sql.DB, opts ...option) *Repository {
	repo := &Repository{
		db: db,
	}

	for _, opt := range opts {
		opt(repo)
	}

	return repo

}

const ResourceID = "/app/config/update"

type publisher interface {
	Publish(resourceID string, message any)
}

type Repository struct {
	db        *sql.DB
	publisher publisher
}

func (r *Repository) publish(config Config) {
	if r.publisher == nil {
		return
	}

	r.publisher.Publish(ResourceID, config)
}

func (r *Repository) Current(ctx context.Context) Config {
	defaultConfig := Config{
		AnalyticsEnabled: true,
	}

	cfg, err := r.Get(ctx, id.ID("current"))
	if err != nil {
		// TODO: log error
		return defaultConfig
	}

	return cfg
}

const selectQuery = `SELECT "analytics_enabled" FROM config`

func (r *Repository) Get(ctx context.Context, _ id.ID) (Config, error) {
	cfg := Config{
		ID:   id.ID("current"),
		Name: "Config",
	}

	err := r.db.
		QueryRowContext(ctx, selectQuery).
		Scan(
			&cfg.AnalyticsEnabled,
		)

	if err != nil {
		return Config{}, fmt.Errorf("sql query: %w", err)
	}

	return cfg, nil
}

const (
	deleteQuery = "DELETE FROM config"
	insertQuery = `INSERT INTO config ("analytics_enabled") VALUES ($1)`
)

func (r *Repository) Update(ctx context.Context, updated Config) (Config, error) {
	// enforce ID and name
	updated = Config{
		ID:               id.ID("current"),
		Name:             "Config",
		AnalyticsEnabled: updated.AnalyticsEnabled,
	}

	tx, err := r.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return Config{}, err
	}

	_, err = tx.ExecContext(ctx, deleteQuery)
	if err != nil {
		return Config{}, fmt.Errorf("sql exec delete: %w", err)
	}

	_, err = tx.ExecContext(ctx, insertQuery, updated.AnalyticsEnabled)
	if err != nil {
		return Config{}, fmt.Errorf("sql exec insert: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return Config{}, fmt.Errorf("commit: %w", err)
	}

	r.publish(updated)

	return updated, nil
}
