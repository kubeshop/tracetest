package pollingprofile

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/server/id"
)

type Strategy string

const (
	Periodic Strategy = "periodic"
)

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

func Repository(db *sql.DB) *repository {
	return &repository{db}
}

type repository struct {
	db *sql.DB
}

func (r *repository) SetID(profile PollingProfile, id id.ID) PollingProfile {
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

func (r *repository) Create(ctx context.Context, profile PollingProfile) (PollingProfile, error) {
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

func (r *repository) Update(ctx context.Context, profile PollingProfile) (PollingProfile, error) {
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

func (r *repository) Get(ctx context.Context, id id.ID) (PollingProfile, error) {
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

func (r *repository) Delete(ctx context.Context, id id.ID) error {
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

func (r *repository) List(ctx context.Context, take, skip int, query, sortedBy, sortDirection string) ([]PollingProfile, error) {
	rows, err := r.db.QueryContext(ctx, baseSelect)
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

const countQuery = `SELECT COUNT(*) FROM polling_profiles`

func (r *repository) Count(ctx context.Context, query string) (int, error) {
	count := 0

	err := r.db.
		QueryRowContext(ctx, countQuery).
		Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("sql query: %w", err)
	}

	return count, nil
}
