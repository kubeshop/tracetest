package configresource

import (
	"context"
	"database/sql"
	"errors"
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

	list, err := r.List(ctx, 1, 0, "", "", "")
	if err != nil || len(list) != 1 {
		// TODO: log error
		return defaultConfig
	}

	return list[0]
}

func (r *Repository) SetID(cfg Config, id id.ID) Config {
	cfg.ID = id
	return cfg
}

const insertQuery = `
		INSERT INTO configs (
			"id",
			"name",
			"analytics_enabled"
		) VALUES ($1, $2, $3)`

func (r *Repository) Create(ctx context.Context, cfg Config) (Config, error) {
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

func (r *Repository) Update(ctx context.Context, updated Config) (Config, error) {
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

	r.publish(updated)

	return updated, nil
}

const getQuery = baseSelect + `WHERE "id" = $1`

func (r *Repository) Get(ctx context.Context, id id.ID) (Config, error) {
	cfg := Config{}

	err := r.db.
		QueryRowContext(ctx, getQuery, id).
		Scan(
			&cfg.ID,
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

func (r *Repository) Delete(ctx context.Context, id id.ID) error {
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

const (
	baseSelect = `
		SELECT
			"id",
			"name",
			"analytics_enabled"
		FROM configs `
)

func (r *Repository) SortingFields() []string {
	return []string{"id", "name", "analytics_enabled"}
}

func (r *Repository) List(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]Config, error) {
	listQuery := baseSelect

	if sortDirection == "" {
		sortDirection = "ASC"
	}

	if query != "" {
		listQuery = fmt.Sprintf("%s WHERE %s", listQuery, query)
	}

	if sortBy != "" {
		listQuery = fmt.Sprintf("%s ORDER BY %s %s", listQuery, sortBy, sortDirection)
	}

	if take > 0 {
		listQuery = fmt.Sprintf("%s LIMIT %d", listQuery, take)
	}

	if skip > 0 {
		listQuery = fmt.Sprintf("%s OFFSET %d", listQuery, skip)
	}

	rows, err := r.db.QueryContext(ctx, listQuery)
	if err != nil {
		return nil, fmt.Errorf("sql query: %w", err)
	}

	configs := []Config{}
	for rows.Next() {
		cfg := Config{}
		err := rows.Scan(
			&cfg.ID,
			&cfg.Name,
			&cfg.AnalyticsEnabled,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}
			return nil, fmt.Errorf("sql query: %w", err)
		}

		configs = append(configs, cfg)
	}

	return configs, nil
}

const baseCountQuery = `SELECT COUNT(*) FROM configs`

func (r *Repository) Count(ctx context.Context, query string) (int, error) {
	countQuery := baseCountQuery

	if query != "" {
		countQuery = fmt.Sprintf("%s WHERE %s", countQuery, query)
	}

	count := 0

	err := r.db.
		QueryRowContext(ctx, countQuery).
		Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("sql query: %w", err)
	}

	return count, nil
}
