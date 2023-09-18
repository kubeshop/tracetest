package demo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/sqlutil"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) SetID(demo Demo, id id.ID) Demo {
	demo.ID = id
	return demo
}

const insertQuery = `INSERT INTO demos (
		"id",
		"name",
		"enabled",
		"type",
		"pokeshop",
		"opentelemetry_store",
		"tenant_id"
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

func (r *Repository) Create(ctx context.Context, demo Demo) (Demo, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Demo{}, err
	}

	pokeshopJSONData, err := json.Marshal(demo.Pokeshop)
	if err != nil {
		return Demo{}, fmt.Errorf("could not get JSON data from pokeshop example: %w", err)
	}

	openTelemetryStoreJSONData, err := json.Marshal(demo.OpenTelemetryStore)
	if err != nil {
		return Demo{}, fmt.Errorf("could not get JSON data from opentelemetry store example: %w", err)
	}

	params := sqlutil.TenantInsert(ctx,
		demo.ID,
		demo.Name,
		demo.Enabled,
		demo.Type,
		pokeshopJSONData,
		openTelemetryStoreJSONData,
	)
	_, err = tx.ExecContext(ctx, insertQuery, params...)

	if err != nil {
		tx.Rollback()
		return Demo{}, fmt.Errorf("sql exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return Demo{}, fmt.Errorf("commit: %w", err)
	}

	return demo, nil
}

const updateQuery = `
	UPDATE demos SET
		"name" = $2,
		"enabled" = $3,
		"type" = $4,
		"pokeshop" = $5,
		"opentelemetry_store" = $6
	WHERE "id" = $1`

func (r *Repository) Update(ctx context.Context, demo Demo) (Demo, error) {
	oldDemo, err := r.Get(ctx, demo.ID)
	if err != nil {
		return Demo{}, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Demo{}, err
	}

	pokeshopJSONData, err := json.Marshal(demo.Pokeshop)
	if err != nil {
		return Demo{}, fmt.Errorf("could not get JSON data from pokeshop example: %w", err)
	}

	openTelemetryStoreJSONData, err := json.Marshal(demo.OpenTelemetryStore)
	if err != nil {
		return Demo{}, fmt.Errorf("could not get JSON data from opentelemetry store example: %w", err)
	}

	query, params := sqlutil.Tenant(ctx, updateQuery,
		oldDemo.ID,
		demo.Name,
		demo.Enabled,
		demo.Type,
		pokeshopJSONData,
		openTelemetryStoreJSONData,
	)

	_, err = tx.ExecContext(ctx, query, params...)

	if err != nil {
		tx.Rollback()
		return Demo{}, fmt.Errorf("sql exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return Demo{}, fmt.Errorf("commit: %w", err)
	}

	return demo, nil
}

const (
	baseSelect = `
		SELECT
			"id",
			"name",
			"enabled",
			"type",
			"pokeshop",
			"opentelemetry_store"
		FROM demos`

	getQuery        = baseSelect + ` WHERE "id" = $1`
	getDefaultQuery = baseSelect + ` WHERE "default" = true`
)

func (r *Repository) Get(ctx context.Context, id id.ID) (Demo, error) {
	return r.get(ctx, getQuery, id)
}

func (r *Repository) get(ctx context.Context, query string, args ...any) (Demo, error) {
	query, args = sqlutil.Tenant(ctx, query, args...)
	row := r.db.QueryRowContext(ctx, query, args...)
	return readRow(row)
}

const deleteQuery = `DELETE FROM demos WHERE "id" = $1`

func (r *Repository) Delete(ctx context.Context, id id.ID) error {
	demo, err := r.Get(ctx, id)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query, params := sqlutil.Tenant(ctx, deleteQuery, demo.ID)
	_, err = tx.ExecContext(ctx, query, params...)

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
	return []string{"id", "name", "type"}
}

func listQuery(baseSQL, query string, params []any) (string, []any) {
	paramNumber := len(params) + 1
	condition := fmt.Sprintf(" AND (t.name ilike $%d)", paramNumber)

	sql, params := sqlutil.Search(baseSQL, condition, query, params)

	return sql, params
}

func (r *Repository) List(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]Demo, error) {
	q, params := listQuery(baseSelect, query, []any{take, skip})
	q, params = sqlutil.Tenant(ctx, q, params...)

	sortingFields := map[string]string{
		"id":   "id",
		"name": "name",
		"type": "type",
	}

	q = sqlutil.Sort(q, sortBy, sortDirection, "name", sortingFields)
	q += " LIMIT $1 OFFSET $2"

	stmt, err := r.db.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return nil, fmt.Errorf("sql query: %w", err)
	}

	demos := []Demo{}
	for rows.Next() {
		demo, err := readRow(rows)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}

			return nil, fmt.Errorf("sql query: %w", err)
		}

		demos = append(demos, demo)
	}

	return demos, nil
}

const baseCountQuery = `SELECT COUNT(*) FROM demos`

func (r *Repository) Count(ctx context.Context, query string) (int, error) {
	countQuery := baseCountQuery

	if query != "" {
		countQuery = fmt.Sprintf("%s WHERE %s", countQuery, query)
	}

	count := 0

	countQuery, params := sqlutil.Tenant(ctx, countQuery)
	err := r.db.
		QueryRowContext(ctx, countQuery, params...).
		Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("sql query: %w", err)
	}

	return count, nil
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func readRow(row scanner) (Demo, error) {
	demo := Demo{}
	var pokeshopJSON []byte
	var openTelemetryStoreJSON []byte
	err := row.Scan(
		&demo.ID,
		&demo.Name,
		&demo.Enabled,
		&demo.Type,
		&pokeshopJSON,
		&openTelemetryStoreJSON,
	)

	if err != nil {
		return Demo{}, fmt.Errorf("could not read row: %w", err)
	}

	if string(pokeshopJSON) != "null" {
		pokeshop := PokeshopDemo{}
		err = json.Unmarshal(pokeshopJSON, &pokeshop)
		if err != nil {
			return Demo{}, fmt.Errorf("could not unmarshal pokeshop JSON: %w", err)
		}

		demo.Pokeshop = &pokeshop
	}

	if string(openTelemetryStoreJSON) != "null" {
		openTelemetryStore := OpenTelemetryStoreDemo{}
		err = json.Unmarshal(openTelemetryStoreJSON, &openTelemetryStore)
		if err != nil {
			return Demo{}, fmt.Errorf("could not unmarshal pokeshop JSON: %w", err)
		}

		demo.OpenTelemetryStore = &openTelemetryStore
	}

	return demo, nil
}

func (r *Repository) Provision(ctx context.Context, demo Demo) error {
	_, err := r.Create(ctx, demo)
	return err
}
