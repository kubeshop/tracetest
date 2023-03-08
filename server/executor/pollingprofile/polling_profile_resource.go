package pollingprofile

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
)

type Strategy string

const (
	Periodic Strategy = "periodic"
)

const ResourceName = "PollingProfile"

var Operations = []resourcemanager.Operation{
	resourcemanager.OperationCreate,
	resourcemanager.OperationDelete,
	resourcemanager.OperationGet,
	resourcemanager.OperationList,
	resourcemanager.OperationUpdate,
}

type PollingProfile struct {
	ID       id.ID                  `mapstructure:"id"`
	Name     string                 `mapstructure:"name"`
	Strategy Strategy               `mapstructure:"strategy"`
	Periodic *PeriodicPollingConfig `mapstructure:"periodic"`
}

type PeriodicPollingConfig struct {
	RetryDelay string `mapstructure:"retryDelay"`
	Timeout    string `mapstructure:"timeout"`
}

func (pp PollingProfile) HasID() bool {
	return pp.ID.String() != ""
}

func (pp PollingProfile) Validate() error {
	if pp.Strategy == Periodic && pp.Periodic == nil {
		return fmt.Errorf("missing periodic polling profile configuration")
	}

	return nil
}

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

const insertQuery = `INSERT INTO polling_profiles (
		"id",
		"name",
		"strategy",
		"periodic"
	)
	VALUES ($1, $2, $3, $4)`

func (r *Repository) Create(ctx context.Context, profile PollingProfile) (PollingProfile, error) {
	var (
		periodicJSON []byte
		err          error
	)

	if profile.Periodic != nil {
		periodicJSON, err = json.Marshal(profile.Periodic)
		if err != nil {
			return PollingProfile{}, fmt.Errorf("could not marshal periodic strategy configuration: %w", err)
		}
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return PollingProfile{}, err
	}

	_, err = tx.ExecContext(ctx, insertQuery,
		profile.ID,
		profile.Name,
		profile.Strategy,
		periodicJSON,
	)

	if err != nil {
		tx.Rollback()
		return PollingProfile{}, fmt.Errorf("sql exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return PollingProfile{}, fmt.Errorf("commit: %w", err)
	}

	return profile, nil
}

const updateQuery = `
	UPDATE polling_profiles SET
		"name" = $2,
		"strategy" = $3,
		"periodic" = $4
	WHERE "id" = $1`

func (r *Repository) Update(ctx context.Context, profile PollingProfile) (PollingProfile, error) {
	var (
		periodicJSON []byte
		err          error
	)

	if profile.Periodic != nil {
		periodicJSON, err = json.Marshal(profile.Periodic)
		if err != nil {
			return PollingProfile{}, fmt.Errorf("could not marshal periodic strategy configuration: %w", err)
		}
	}

	oldProfile, err := r.Get(ctx, profile.ID)
	if err != nil {
		return PollingProfile{}, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return PollingProfile{}, err
	}

	_, err = tx.ExecContext(ctx, updateQuery,
		oldProfile.ID,
		profile.Name,
		profile.Strategy,
		periodicJSON,
	)

	if err != nil {
		tx.Rollback()
		return PollingProfile{}, fmt.Errorf("sql exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return PollingProfile{}, fmt.Errorf("commit: %w", err)
	}

	return profile, nil
}

const (
	baseSelect = `
		SELECT
			"id",
			"name",
			"strategy",
			"periodic"
		FROM polling_profiles `

	getQuery = baseSelect + `WHERE "id" = $1`
)

func (r *Repository) Get(ctx context.Context, id id.ID) (PollingProfile, error) {
	profile := PollingProfile{}

	var periodicJSON []byte
	err := r.db.
		QueryRowContext(ctx, getQuery, id).
		Scan(
			&profile.ID,
			&profile.Name,
			&profile.Strategy,
			&periodicJSON,
		)

	if err != nil {
		return PollingProfile{}, fmt.Errorf("sql query: %w", err)
	}

	if len(periodicJSON) > 0 {
		var periodicConfig PeriodicPollingConfig
		err = json.Unmarshal(periodicJSON, &periodicConfig)
		if err != nil {
			return PollingProfile{}, fmt.Errorf("could not unmarshal periodic strategy config: %w", err)
		}

		profile.Periodic = &periodicConfig
	}

	return profile, nil
}

const deleteQuery = `DELETE FROM polling_profiles WHERE "id" = $1`

func (r *Repository) Delete(ctx context.Context, id id.ID) error {
	profile, err := r.Get(ctx, id)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, deleteQuery, profile.ID)

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

func (r Repository) SortingFields() []string {
	return []string{"id", "name", "strategy"}
}

func (r *Repository) List(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]PollingProfile, error) {
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

	profiles := []PollingProfile{}
	for rows.Next() {
		profile := PollingProfile{}
		var periodicJSON []byte
		err := rows.Scan(
			&profile.ID,
			&profile.Name,
			&profile.Strategy,
			&periodicJSON,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}
			return nil, fmt.Errorf("sql query: %w", err)
		}

		if len(periodicJSON) > 0 {
			var periodicConfig PeriodicPollingConfig
			err = json.Unmarshal(periodicJSON, &periodicConfig)
			if err != nil {
				return profiles, fmt.Errorf("could not unmarshal periodic strategy config: %w", err)
			}

			profile.Periodic = &periodicConfig
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}

const baseCountQuery = `SELECT COUNT(*) FROM polling_profiles`

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
