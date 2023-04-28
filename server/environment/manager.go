package environment

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/resourcemanager"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return Repository{db}
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
			"values"
		) VALUES ($1, $2, $3, $4, $5)`

	updateQuery = `
		UPDATE environments SET
			"id" = $2,
			"name" = $3,
			"description" = $4,
			"created_at" = $5,
			"values" = $6
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

var _ resourcemanager.Create[Environment] = &Repository{}
var _ resourcemanager.Delete[Environment] = &Repository{}
var _ resourcemanager.Get[Environment] = &Repository{}
var _ resourcemanager.List[Environment] = &Repository{}
var _ resourcemanager.Update[Environment] = &Repository{}
var _ resourcemanager.IDSetter[Environment] = &Repository{}
var _ resourcemanager.Provision[Environment] = &Repository{}

func (*Repository) SortingFields() []string {
	return []string{"name", "createdAt"}
}

func (r *Repository) SetID(environment Environment, id id.ID) Environment {
	environment.ID = id
	return environment
}

func (r *Repository) Create(ctx context.Context, environment Environment) (Environment, error) {
	environment.ID = environment.Slug()
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

	_, err = tx.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("sql Commit: %w", err)
	}

	return nil
}

func (r *Repository) List(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]Environment, error) {
	if take == 0 {
		take = 20
	}

	params := []any{take, skip}

	sql := getQuery

	const condition = "WHERE (e.name ilike $3 OR e.description ilike $3) "
	if query != "" {
		sql += condition
		params = append(params, query)
	}

	sortingFields := map[string]string{
		"created": "e.created_at",
		"name":    "e.name",
	}

	sql = sortQuery(sql, sortBy, sortDirection, sortingFields)
	sql += ` LIMIT $1 OFFSET $2 `

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

func sortQuery(sql, sortBy, sortDirection string, sortingFields map[string]string) string {
	sortField, ok := sortingFields[sortBy]

	if !ok {
		sortField = sortingFields["created"]
	}

	dir := "DESC"
	if strings.ToLower(sortDirection) == "asc" {
		dir = "ASC"
	}

	return fmt.Sprintf("%s ORDER BY %s %s", sql, sortField, dir)
}

func (r *Repository) Get(ctx context.Context, id id.ID) (Environment, error) {
	stmt, err := r.db.Prepare(getQuery + " WHERE e.id = $1")

	if err != nil {
		return Environment{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	environment, err := r.readEnvironmentRow(ctx, stmt.QueryRowContext(ctx, id))
	if err != nil {
		return Environment{}, err
	}

	return environment, nil
}

func (r *Repository) EnvironmentIDExists(ctx context.Context, id string) (bool, error) {
	exists := false

	row := r.db.QueryRowContext(ctx, idExistsQuery, id)

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

	sql := countQuery + query

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

	_, err = stmt.ExecContext(
		ctx,
		environment.ID,
		environment.Name,
		environment.Description,
		environment.CreatedAt,
		jsonValues,
	)

	if err != nil {
		return Environment{}, fmt.Errorf("sql exec: %w", err)
	}

	return environment, nil
}

func (r *Repository) updateIntoEnvironments(ctx context.Context, environment Environment, oldId id.ID) (Environment, error) {
	stmt, err := r.db.Prepare(updateQuery)
	if err != nil {
		return Environment{}, fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()

	jsonValues, err := json.Marshal(environment.Values)
	if err != nil {
		return Environment{}, fmt.Errorf("encoding error: %w", err)
	}

	_, err = stmt.ExecContext(
		ctx,
		oldId,
		environment.Slug(),
		environment.Name,
		environment.Description,
		environment.CreatedAt,
		jsonValues,
	)

	if err != nil {
		return Environment{}, fmt.Errorf("sql exec: %w", err)
	}

	return environment, nil
}

func (r *Repository) Provision(ctx context.Context, environment Environment) error {
	_, err := r.Create(ctx, environment)
	return err
}
