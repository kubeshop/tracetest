package pollingprofile

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/server/http/middleware"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/sqlutil"
)

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

type Repository struct {
	db *sql.DB
}

func (r *Repository) SetID(profile PollingProfile, id id.ID) PollingProfile {
	profile.ID = id
	return profile
}

func (r *Repository) Create(ctx context.Context, updated PollingProfile) (PollingProfile, error) {
	updated.ID = id.ID("current")
	return r.Update(ctx, updated)
}

const (
	insertQuery = `
	INSERT INTO polling_profiles(
		"id",
		"name",
		"default",
		"strategy",
		"periodic",
		"tenant_id"
	)
	VALUES ($1, $2, $3, $4, $5, $6)`
	deleteQuery = `DELETE FROM polling_profiles`
)

func (r *Repository) Update(ctx context.Context, updated PollingProfile) (PollingProfile, error) {
	// enforce ID and default
	updated.ID = "current"
	updated.Default = true

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return PollingProfile{}, err
	}
	defer tx.Rollback()

	query, params := sqlutil.Tenant(ctx, deleteQuery)
	_, err = tx.ExecContext(ctx, query, params...)
	if err != nil {
		return PollingProfile{}, fmt.Errorf("sql exec delete: %w", err)
	}

	var periodicJSON []byte
	if updated.Periodic != nil {
		periodicJSON, err = json.Marshal(updated.Periodic)
		if err != nil {
			return PollingProfile{}, fmt.Errorf("could not marshal periodic strategy configuration: %w", err)
		}
	}

	tenantID := middleware.TenantIDFromContext(ctx)
	_, err = tx.ExecContext(ctx, insertQuery,
		updated.ID,
		updated.Name,
		updated.Default,
		updated.Strategy,
		periodicJSON,
		tenantID,
	)
	if err != nil {
		return PollingProfile{}, fmt.Errorf("sql exec insert: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return PollingProfile{}, fmt.Errorf("commit: %w", err)
	}

	return updated, nil

}

const (
	getQuery = `
		SELECT
			"name",
			"strategy",
			"periodic"
		FROM polling_profiles `
)

func (r *Repository) GetDefault(ctx context.Context) PollingProfile {
	pp, _ := r.Get(ctx, id.ID("current"))
	return pp
}

func (r *Repository) Get(ctx context.Context, id id.ID) (PollingProfile, error) {
	profile := PollingProfile{
		ID:      "current",
		Default: true,
	}

	var periodicJSON []byte
	query, params := sqlutil.Tenant(ctx, getQuery)
	err := r.db.
		QueryRowContext(ctx, query, params...).
		Scan(
			&profile.Name,
			&profile.Strategy,
			&periodicJSON,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DefaultPollingProfile, nil
		}
		return PollingProfile{}, fmt.Errorf("sql query: %w", err)
	}

	if string(periodicJSON) != "null" {
		var periodicConfig PeriodicPollingConfig
		err = json.Unmarshal(periodicJSON, &periodicConfig)
		if err != nil {
			return PollingProfile{}, fmt.Errorf("could not unmarshal periodic strategy config: %w", err)
		}

		profile.Periodic = &periodicConfig
	}

	return profile, nil
}

func (r *Repository) List(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]PollingProfile, error) {
	cfg, err := r.Get(ctx, id.ID("current"))
	if err != nil {
		return []PollingProfile{}, err
	}

	return []PollingProfile{cfg}, nil
}

func (r *Repository) Count(ctx context.Context, query string) (int, error) {
	return 1, nil
}

func (*Repository) SortingFields() []string {
	return []string{"name"}
}

func (r *Repository) Provision(ctx context.Context, profile PollingProfile) error {
	_, err := r.Update(ctx, profile)
	return err
}
