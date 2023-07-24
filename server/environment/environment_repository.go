package environment

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/pkg/sqlutil"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

type scanner interface {
	Scan(dest ...interface{}) error
}

const (
	insertQuery = `
		INSERT INTO environments (
			"id",
			"name",
			"description",
			"created_at",
			"values",
			"tenant_id"
		) VALUES ($1, $2, $3, $4, $5, $6)`

	updateQuery = `
		UPDATE environments SET
			"name" = $2,
			"description" = $3,
			"created_at" = $4,
			"values" = $5
		WHERE id = $1
	`

	getQuery = `
	SELECT
		e.id,
		e.name,
		e.description,
		e.created_at,
		e.values
	FROM environments e
`

	deleteQuery = "DELETE FROM environments WHERE id = $1"

	idExistsQuery = `
		SELECT COUNT(*) > 0 as exists FROM environments WHERE id = $1
	`

	countQuery = `
		SELECT COUNT(*) FROM environments e
	`
)

func (*Repository) SortingFields() []string {
	return []string{"name", "createdAt"}
}

func (r *Repository) SetID(environment Environment, id id.ID) Environment {
	environment.ID = id
	return environment
}

func (r *Repository) Create(ctx context.Context, environment Environment) (Environment, error) {
	if !environment.HasID() {
		environment.ID = environment.Slug()
	}

	environment.CreatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	return r.insertIntoEnvironments(ctx, environment)
}

func (r *Repository) Update(ctx context.Context, environment Environment) (Environment, error) {
	oldEnvironment, err := r.Get(ctx, environment.ID)
	if err != nil {
		return Environment{}, err
	}

	// keep the same creation date to keep sort order
	environment.CreatedAt = oldEnvironment.CreatedAt
	environment.ID = oldEnvironment.ID

	return r.updateIntoEnvironments(ctx, environment, oldEnvironment.ID)
}

func (r *Repository) Delete(ctx context.Context, id id.ID) error {
	_, err := r.Get(ctx, id)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("sql BeginTx: %w", err)
	}
	defer tx.Rollback()

	query, params := sqlutil.Tenant(ctx, deleteQuery, id)
	_, err = tx.ExecContext(ctx, query, params...)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("sql Commit: %w", err)
	}

	return nil
}

func listQuery(baseSQL, query string, params []any) (string, []any) {
	paramNumber := len(params) + 1
	condition := fmt.Sprintf(" WHERE (e.name ilike $%d OR e.description ilike $%d)", paramNumber, paramNumber)

	sql, params := sqlutil.Search(baseSQL, condition, query, params)

	return sql, params
}

func (r *Repository) List(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]Environment, error) {
	if take == 0 {
		take = 20
	}

	sql, params := listQuery(getQuery, query, []any{take, skip})

	sortingFields := map[string]string{
		"created": "e.created_at",
		"name":    "e.name",
	}

	sql = sqlutil.Sort(sql, sortBy, sortDirection, "created", sortingFields)
	sql += ` LIMIT $1 OFFSET $2 `

	sql, params = sqlutil.Tenant(ctx, sql, params...)
	stmt, err := r.db.Prepare(sql)
	if err != nil {
		return []Environment{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return []Environment{}, err
	}

	environments := []Environment{}

	for rows.Next() {
		environment, err := r.readEnvironmentRow(ctx, rows)
		if err != nil {
			return []Environment{}, err
		}

		environments = append(environments, environment)
	}

	return environments, nil
}

func (r *Repository) Count(ctx context.Context, query string) (int, error) {
	return r.countEnvironments(ctx, query)
}

func (r *Repository) Get(ctx context.Context, id id.ID) (Environment, error) {
	query, params := sqlutil.Tenant(ctx, getQuery+" WHERE e.id = $1", id)
	row := r.db.QueryRowContext(ctx, query, params...)

	environment, err := r.readEnvironmentRow(ctx, row)
	if err != nil {
		return Environment{}, err
	}

	return environment, nil
}

func (r *Repository) Exists(ctx context.Context, id id.ID) (bool, error) {
	exists := false

	query, params := sqlutil.Tenant(ctx, idExistsQuery, id)
	row := r.db.QueryRowContext(ctx, query, params...)

	err := row.Scan(&exists)

	return exists, err
}

func (r *Repository) readEnvironmentRow(ctx context.Context, row scanner) (Environment, error) {
	environment := Environment{}

	var (
		jsonValues []byte
	)
	err := row.Scan(
		&environment.ID,
		&environment.Name,
		&environment.Description,
		&environment.CreatedAt,
		&jsonValues,
	)

	switch err {
	case sql.ErrNoRows:
		return Environment{}, err
	case nil:
		err = json.Unmarshal(jsonValues, &environment.Values)
		if err != nil {
			return Environment{}, fmt.Errorf("cannot parse environment: %w", err)
		}

		return environment, nil
	default:
		return Environment{}, err
	}
}

func (r *Repository) countEnvironments(ctx context.Context, query string) (int, error) {
	var (
		count  int
		params []any
	)

	condition := " WHERE (e.name ilike $1 OR e.description ilike $1)"
	sql, params := sqlutil.Search(countQuery, condition, query, params)
	sql, params = sqlutil.Tenant(ctx, sql, params...)

	err := r.db.
		QueryRowContext(ctx, sql, params...).
		Scan(&count)

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repository) insertIntoEnvironments(ctx context.Context, environment Environment) (Environment, error) {
	stmt, err := r.db.Prepare(insertQuery)
	if err != nil {
		return Environment{}, fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()

	jsonValues, err := json.Marshal(environment.Values)
	if err != nil {
		return Environment{}, fmt.Errorf("encoding error: %w", err)
	}

	tenantID := sqlutil.TenantID(ctx)

	_, err = stmt.ExecContext(
		ctx,
		environment.ID,
		environment.Name,
		environment.Description,
		environment.CreatedAt,
		jsonValues,
		tenantID,
	)

	if err != nil {
		return Environment{}, fmt.Errorf("sql exec: %w", err)
	}

	return environment, nil
}

func (r *Repository) updateIntoEnvironments(ctx context.Context, environment Environment, oldId id.ID) (Environment, error) {
	jsonValues, err := json.Marshal(environment.Values)
	if err != nil {
		return Environment{}, fmt.Errorf("encoding error: %w", err)
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Environment{}, err
	}
	defer tx.Rollback()

	query, params := sqlutil.Tenant(ctx, updateQuery, oldId,
		environment.Name,
		environment.Description,
		environment.CreatedAt,
		jsonValues,
	)

	_, err = tx.ExecContext(
		ctx,
		query,
		params...,
	)
	if err != nil {
		return Environment{}, fmt.Errorf("sql exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return Environment{}, fmt.Errorf("commit: %w", err)
	}

	return environment, nil
}

func (r *Repository) Provision(ctx context.Context, environment Environment) error {
	_, err := r.Create(ctx, environment)
	return err
}
