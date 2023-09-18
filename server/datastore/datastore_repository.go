package datastore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/sqlutil"
)

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

type Repository struct {
	db *sql.DB
}

func (r *Repository) SetID(dataStore DataStore, id id.ID) DataStore {
	dataStore.ID = id
	return dataStore
}

const DataStoreSingleID id.ID = "current"

const insertQuery = `
INSERT INTO data_stores (
	"id",
	"name",
	"type",
	"is_default",
	"values",
	"created_at",
	"tenant_id"
) VALUES ($1, $2, $3, $4, $5, $6, $7)`

const deleteQuery = `DELETE FROM data_stores WHERE "id" = $1`

func newCreateAtDateString() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

func (r *Repository) getCreatedAt(ctx context.Context, dataStore DataStore) (string, error) {
	if dataStore.CreatedAt != "" {
		// client passed date, keeping it
		return dataStore.CreatedAt, nil
	}

	// get datastore on the database or the default one
	oldDataStore, err := r.Get(ctx, dataStore.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// record not found, return a new date
			return newCreateAtDateString(), nil
		}

		return "", err
	}

	// record found, return old date
	return oldDataStore.CreatedAt, nil
}

func (r *Repository) Create(ctx context.Context, updated DataStore) (DataStore, error) {
	updated.ID = DataStoreSingleID
	return r.Update(ctx, updated)
}

func (r *Repository) Update(ctx context.Context, dataStore DataStore) (DataStore, error) {
	// enforce ID and default
	dataStore.ID = DataStoreSingleID
	dataStore.Default = true

	// reuse the created_at field for auditing purposes,
	// unless the client explicitly send it
	createdAt, err := r.getCreatedAt(ctx, dataStore)
	if err != nil {
		return DataStore{}, err
	}

	dataStore.CreatedAt = createdAt

	// since we allow only one datastore, delete the table and keep one record
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return DataStore{}, err
	}
	defer tx.Rollback()

	query, params := sqlutil.Tenant(ctx, deleteQuery, DataStoreSingleID)
	_, err = tx.ExecContext(ctx, query, params...)
	if err != nil {
		return DataStore{}, fmt.Errorf("datastore repository sql exec delete: %w", err)
	}

	valuesJSON, err := json.Marshal(dataStore.Values)
	if err != nil {
		return DataStore{}, fmt.Errorf("could not marshal values field configuration: %w", err)
	}

	params = sqlutil.TenantInsert(ctx,
		dataStore.ID,
		dataStore.Name,
		dataStore.Type,
		dataStore.Default,
		valuesJSON,
		dataStore.CreatedAt,
	)
	_, err = tx.ExecContext(ctx, insertQuery, params...)
	if err != nil {
		return DataStore{}, fmt.Errorf("datastore repository sql exec create: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return DataStore{}, fmt.Errorf("datastore repository commit: %w", err)
	}

	return dataStore, nil
}

func (r *Repository) Delete(ctx context.Context, id id.ID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query, params := sqlutil.Tenant(ctx, deleteQuery, id)
	_, err = tx.ExecContext(ctx, query, params...)
	if err != nil {
		return fmt.Errorf("datastore repository sql exec delete: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

const getQuery = `
SELECT
	"id",
	"name",
	"type",
	"is_default",
	"values",
	"created_at"
FROM data_stores
WHERE "id" = $1`

func (r *Repository) Current(ctx context.Context) (DataStore, error) {
	dataStore, err := r.Get(ctx, "current")
	if err != nil {
		return DataStore{}, fmt.Errorf("datastore repository get current: %w", err)
	}

	return dataStore, nil
}

func (r *Repository) Get(ctx context.Context, id id.ID) (DataStore, error) {
	query, params := sqlutil.Tenant(ctx, getQuery, id)
	row := r.db.QueryRowContext(ctx, query, params...)

	dataStore, err := r.readRow(row)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return DataStore{
			CreatedAt: newCreateAtDateString(),
		}, nil // Assumes an empty datastore
	}
	if err != nil {
		return DataStore{}, fmt.Errorf("datastore repository get sql query: %w", err)
	}

	return dataStore, nil
}

func (r *Repository) List(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]DataStore, error) {
	cfg, err := r.Get(ctx, id.ID("current"))
	if err != nil {
		return []DataStore{}, err
	}

	return []DataStore{cfg}, nil
}

func (r *Repository) Count(ctx context.Context, query string) (int, error) {
	return 1, nil
}

func (*Repository) SortingFields() []string {
	return []string{"name"}
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func (r *Repository) readRow(rowScanner scanner) (DataStore, error) {
	var valuesJSON []byte

	dataStore := DataStore{}

	err := rowScanner.Scan(
		&dataStore.ID,
		&dataStore.Name,
		&dataStore.Type,
		&dataStore.Default,
		&valuesJSON,
		&dataStore.CreatedAt,
	)

	if err != nil {
		return DataStore{}, err
	}

	if string(valuesJSON) != "null" {
		err = json.Unmarshal(valuesJSON, &dataStore.Values)
		if err != nil {
			return DataStore{}, fmt.Errorf("unable to parse data store values: %w", err)
		}
	}

	return dataStore, nil
}

func (r *Repository) Provision(ctx context.Context, dataStore DataStore) error {
	_, err := r.Update(ctx, dataStore)
	return err
}
