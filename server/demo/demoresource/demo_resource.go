package demoresource

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
)

type DemoType string

var (
	DemoTypePokeshop           DemoType = "pokeshop"
	DemoTypeOpentelemetryStore DemoType = "otelstore"
)

type Demo struct {
	ID                 id.ID                   `mapstructure:"id"`
	Name               string                  `mapstructure:"name"`
	Type               DemoType                `mapstructure:"type"`
	Enabled            bool                    `mapstructure:"enabled"`
	Pokeshop           *PokeshopDemo           `mapstructure:"pokeshop,omitempty"`
	OpenTelemetryStore *OpenTelemetryStoreDemo `mapstructure:"opentelemetryStore,omitempty"`
}

func (d Demo) HasID() bool {
	return d.ID != ""
}

func (d Demo) Validate() error {
	return nil
}

type PokeshopDemo struct {
	HTTPEndpoint string `mapstructure:"httpEndpoint,omitempty"`
	GRPCEndpoint string `mapstructure:"grpcEndpoint,omitempty"`
}

type OpenTelemetryStoreDemo struct {
	FrontendEndpoint       string `mapstructure:"frontendEndpoint,omitempty"`
	ProductCatalogEndpoint string `mapstructure:"productCatalogEndpoint,omitempty"`
	CartEndpoint           string `mapstructure:"cartEndpoint,omitempty"`
	CheckoutEndpoint       string `mapstructure:"checkoutEndpoint,omitempty"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

var _ resourcemanager.List[Demo] = &Repository{}

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
		"opentelemetry_store"
	)
	VALUES ($1, $2, $3, $4, $5, $6)`

func (r *Repository) Create(ctx context.Context, demo Demo) (Demo, error) {
	if !demo.HasID() {
		demo.ID = id.SlugFromString(demo.Name)
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

	_, err = tx.ExecContext(ctx, insertQuery,
		demo.ID,
		demo.Name,
		demo.Enabled,
		demo.Type,
		pokeshopJSONData,
		openTelemetryStoreJSONData,
	)

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

	_, err = tx.ExecContext(ctx, updateQuery,
		oldDemo.ID,
		demo.Name,
		demo.Enabled,
		demo.Type,
		pokeshopJSONData,
		openTelemetryStoreJSONData,
	)

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
		FROM demos `

	getQuery        = baseSelect + `WHERE "id" = $1`
	getDefaultQuery = baseSelect + `WHERE "default" = true`
)

func (r *Repository) Get(ctx context.Context, id id.ID) (Demo, error) {
	return r.get(ctx, getQuery, id)
}

func (r *Repository) get(ctx context.Context, query string, args ...any) (Demo, error) {
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

	_, err = tx.ExecContext(ctx, deleteQuery, demo.ID)

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

func (r *Repository) List(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]Demo, error) {
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

	err := r.db.
		QueryRowContext(ctx, countQuery).
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

	if len(pokeshopJSON) > 0 {
		pokeshop := PokeshopDemo{}
		err = json.Unmarshal(pokeshopJSON, &pokeshop)
		if err != nil {
			return Demo{}, fmt.Errorf("could not unmarshal pokeshop JSON: %w", err)
		}

		demo.Pokeshop = &pokeshop
	}

	if len(openTelemetryStoreJSON) > 0 {
		openTelemetryStore := OpenTelemetryStoreDemo{}
		err = json.Unmarshal(openTelemetryStoreJSON, &openTelemetryStore)
		if err != nil {
			return Demo{}, fmt.Errorf("could not unmarshal pokeshop JSON: %w", err)
		}

		demo.OpenTelemetryStore = &openTelemetryStore
	}

	return demo, nil
}
