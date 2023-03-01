package testdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/kubeshop/tracetest/server/model"
)

var _ model.EnvironmentRepository = &postgresDB{}

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

func (td *postgresDB) CreateEnvironment(ctx context.Context, environment model.Environment) (model.Environment, error) {
	environment.ID = environment.Slug()
	environment.CreatedAt = time.Now()

	return td.insertIntoEnvironments(ctx, environment)
}

func (td *postgresDB) UpdateEnvironment(ctx context.Context, environment model.Environment) (model.Environment, error) {
	oldEnvironment, err := td.GetEnvironment(ctx, environment.ID)
	if err != nil {
		return model.Environment{}, err
	}

	// keep the same creation date to keep sort order
	environment.CreatedAt = oldEnvironment.CreatedAt
	environment.ID = environment.Slug()

	return td.updateIntoEnvironments(ctx, environment, oldEnvironment.ID)
}

func (td *postgresDB) DeleteEnvironment(ctx context.Context, environment model.Environment) error {
	tx, err := td.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("sql BeginTx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, deleteQuery, environment.ID)

	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("sql Commit: %w", err)
	}

	return nil
}

func (td *postgresDB) GetEnvironments(ctx context.Context, take, skip int32, query, sortBy, sortDirection string) (model.List[model.Environment], error) {
	hasSearchQuery := query != ""
	cleanSearchQuery := "%" + strings.ReplaceAll(query, " ", "%") + "%"
	params := []any{take, skip}

	sql := getQuery

	const condition = "WHERE (e.name ilike $3 OR e.description ilike $3) "
	if hasSearchQuery {
		params = append(params, cleanSearchQuery)
		sql += condition
	}

	sortingFields := map[string]string{
		"created": "e.created_at",
		"name":    "e.name",
	}

	sql = sortQuery(sql, sortBy, sortDirection, sortingFields)
	sql += ` LIMIT $1 OFFSET $2 `

	stmt, err := td.db.Prepare(sql)
	if err != nil {
		return model.List[model.Environment]{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return model.List[model.Environment]{}, err
	}

	environments := []model.Environment{}

	for rows.Next() {
		environment, err := td.readEnvironmentRow(ctx, rows)
		if err != nil {
			return model.List[model.Environment]{}, err
		}

		environments = append(environments, environment)
	}

	count, err := td.countEnvironments(ctx, condition, cleanSearchQuery)
	if err != nil {
		return model.List[model.Environment]{}, err
	}

	return model.List[model.Environment]{
		Items:      environments,
		TotalCount: count,
	}, nil
}

func (td *postgresDB) GetEnvironment(ctx context.Context, id string) (model.Environment, error) {
	stmt, err := td.db.Prepare(getQuery + " WHERE e.id = $1")

	if err != nil {
		return model.Environment{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	environment, err := td.readEnvironmentRow(ctx, stmt.QueryRowContext(ctx, id))
	if err != nil {
		return model.Environment{}, err
	}

	return environment, nil
}

func (td *postgresDB) EnvironmentIDExists(ctx context.Context, id string) (bool, error) {
	exists := false

	row := td.db.QueryRowContext(ctx, idExistsQuery, id)

	err := row.Scan(&exists)

	return exists, err
}

func (td *postgresDB) readEnvironmentRow(ctx context.Context, row scanner) (model.Environment, error) {
	environment := model.Environment{}

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
		return model.Environment{}, ErrNotFound
	case nil:
		err = json.Unmarshal(jsonValues, &environment.Values)
		if err != nil {
			return model.Environment{}, fmt.Errorf("cannot parse environment: %w", err)
		}

		return environment, nil
	default:
		return model.Environment{}, err
	}
}

func (td *postgresDB) countEnvironments(ctx context.Context, condition, cleanSearchQuery string) (int, error) {
	var (
		count  int
		params []any
	)

	sql := countQuery
	if cleanSearchQuery != "" {
		params = []any{cleanSearchQuery}
		sql += strings.ReplaceAll(condition, "$3", "$1")
	}

	err := td.db.
		QueryRowContext(ctx, sql, params...).
		Scan(&count)

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (td *postgresDB) insertIntoEnvironments(ctx context.Context, environment model.Environment) (model.Environment, error) {
	stmt, err := td.db.Prepare(insertQuery)
	if err != nil {
		return model.Environment{}, fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()

	jsonValues, err := json.Marshal(environment.Values)
	if err != nil {
		return model.Environment{}, fmt.Errorf("encoding error: %w", err)
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
		return model.Environment{}, fmt.Errorf("sql exec: %w", err)
	}

	return environment, nil
}

func (td *postgresDB) updateIntoEnvironments(ctx context.Context, environment model.Environment, oldId string) (model.Environment, error) {
	stmt, err := td.db.Prepare(updateQuery)
	if err != nil {
		return model.Environment{}, fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()

	jsonValues, err := json.Marshal(environment.Values)
	if err != nil {
		return model.Environment{}, fmt.Errorf("encoding error: %w", err)
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
		return model.Environment{}, fmt.Errorf("sql exec: %w", err)
	}

	return environment, nil
}
