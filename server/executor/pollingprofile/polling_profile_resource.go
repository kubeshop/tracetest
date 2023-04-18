package pollingprofile

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
)

type Strategy string

const (
	Periodic Strategy = "periodic"
)

const (
	ResourceName       = "PollingProfile"
	ResourceNamePlural = "PollingProfiles"
)

var Operations = []resourcemanager.Operation{
	resourcemanager.OperationGet,
	resourcemanager.OperationUpdate,
}

var DefaultPollingProfile = PollingProfile{
	ID:       id.ID("current"),
	Name:     "default",
	Default:  true,
	Strategy: Periodic,
	Periodic: &PeriodicPollingConfig{
		Timeout:    "1m",
		RetryDelay: "5s",
	},
}

type PollingProfile struct {
	ID       id.ID                  `mapstructure:"id"`
	Name     string                 `mapstructure:"name"`
	Default  bool                   `mapstructure:"default"`
	Strategy Strategy               `mapstructure:"strategy"`
	Periodic *PeriodicPollingConfig `mapstructure:"periodic"`
}

type PeriodicPollingConfig struct {
	RetryDelay string `mapstructure:"retryDelay"`
	Timeout    string `mapstructure:"timeout"`
}

func (ppc *PeriodicPollingConfig) TimeoutDuration() time.Duration {
	d, _ := time.ParseDuration(ppc.Timeout)
	return d
}

func (ppc *PeriodicPollingConfig) RetryDelayDuration() time.Duration {
	d, _ := time.ParseDuration(ppc.RetryDelay)
	return d
}

func (ppc *PeriodicPollingConfig) MaxTracePollRetry() int {
	return int(math.Ceil(float64(ppc.TimeoutDuration()) / float64(ppc.RetryDelayDuration())))
}

func (ppc *PeriodicPollingConfig) Validate() error {
	if ppc == nil {
		return fmt.Errorf("missing periodic polling profile configuration")
	}

	if _, err := time.ParseDuration(ppc.RetryDelay); err != nil {
		return fmt.Errorf("retry delay configuration is invalid: %w", err)
	}

	if _, err := time.ParseDuration(ppc.Timeout); err != nil {
		return fmt.Errorf("timeout configuration is invalid: %w", err)
	}

	return nil
}

func (pp PollingProfile) HasID() bool {
	return pp.ID.String() != ""
}

func (pp PollingProfile) Validate() error {
	if pp.Strategy == Periodic {
		if err := pp.Periodic.Validate(); err != nil {
			return err
		}
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

const (
	insertQuery = `
	INSERT INTO polling_profiles(
		"id",
		"name",
		"default",
		"strategy",
		"periodic"
	)
	VALUES ($1, $2, $3, $4, $5)`
	deleteQuery = `DELETE FROM polling_profiles`
)

func (r *Repository) Update(ctx context.Context, updated PollingProfile) (PollingProfile, error) {
	// enforce ID and default
	updated.ID = "current"
	updated.Default = true

	tx, err := r.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		return PollingProfile{}, err
	}

	_, err = tx.ExecContext(ctx, deleteQuery)
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

	_, err = tx.ExecContext(ctx, insertQuery,
		updated.ID,
		updated.Name,
		updated.Default,
		updated.Strategy,
		periodicJSON,
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
	err := r.db.
		QueryRowContext(ctx, getQuery).
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

func (r *Repository) Provision(ctx context.Context, profile PollingProfile) error {
	_, err := r.Update(ctx, profile)
	return err
}
