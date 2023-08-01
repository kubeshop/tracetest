package variableset

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
		INSERT INTO variable_sets (
			"id",
			"name",
			"description",
			"created_at",
			"values",
			"tenant_id"
		) VALUES ($1, $2, $3, $4, $5, $6)`

	updateQuery = `
		UPDATE variable_sets SET
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
	FROM variable_sets e
`

	deleteQuery = "DELETE FROM variable_sets WHERE id = $1"

	idExistsQuery = `
		SELECT COUNT(*) > 0 as exists FROM variable_sets WHERE id = $1
	`

	countQuery = `
		SELECT COUNT(*) FROM variable_sets e
	`
)

func (*Repository) SortingFields() []string {
	return []string{"name", "createdAt"}
}

func (r *Repository) SetID(variableSet VariableSet, id id.ID) VariableSet {
	variableSet.ID = id
	return variableSet
}

func (r *Repository) Create(ctx context.Context, variableSet VariableSet) (VariableSet, error) {
	if !variableSet.HasID() {
		variableSet.ID = variableSet.Slug()
	}

	variableSet.CreatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	return r.insertIntoEnvironments(ctx, variableSet)
}

func (r *Repository) Update(ctx context.Context, variableSet VariableSet) (VariableSet, error) {
	oldEnvironment, err := r.Get(ctx, variableSet.ID)
	if err != nil {
		return VariableSet{}, err
	}

	// keep the same creation date to keep sort order
	variableSet.CreatedAt = oldEnvironment.CreatedAt
	variableSet.ID = oldEnvironment.ID

	return r.updateIntoEnvironments(ctx, variableSet, oldEnvironment.ID)
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

func (r *Repository) List(ctx context.Context, take, skip int, query, sortBy, sortDirection string) ([]VariableSet, error) {
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
		return []VariableSet{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return []VariableSet{}, err
	}

	variableSets := []VariableSet{}

	for rows.Next() {
		variableSet, err := r.readEnvironmentRow(ctx, rows)
		if err != nil {
			return []VariableSet{}, err
		}

		variableSets = append(variableSets, variableSet)
	}

	return variableSets, nil
}

func (r *Repository) Count(ctx context.Context, query string) (int, error) {
	return r.countEnvironments(ctx, query)
}

func (r *Repository) Get(ctx context.Context, id id.ID) (VariableSet, error) {
	query, params := sqlutil.Tenant(ctx, getQuery+" WHERE e.id = $1", id)
	row := r.db.QueryRowContext(ctx, query, params...)

	variableSet, err := r.readEnvironmentRow(ctx, row)
	if err != nil {
		return VariableSet{}, err
	}

	return variableSet, nil
}

func (r *Repository) Exists(ctx context.Context, id id.ID) (bool, error) {
	exists := false

	query, params := sqlutil.Tenant(ctx, idExistsQuery, id)
	row := r.db.QueryRowContext(ctx, query, params...)

	err := row.Scan(&exists)

	return exists, err
}

func (r *Repository) readEnvironmentRow(ctx context.Context, row scanner) (VariableSet, error) {
	variableSet := VariableSet{}

	var (
		jsonValues []byte
	)
	err := row.Scan(
		&variableSet.ID,
		&variableSet.Name,
		&variableSet.Description,
		&variableSet.CreatedAt,
		&jsonValues,
	)

	switch err {
	case sql.ErrNoRows:
		return VariableSet{}, err
	case nil:
		err = json.Unmarshal(jsonValues, &variableSet.Values)
		if err != nil {
			return VariableSet{}, fmt.Errorf("cannot parse variableSet: %w", err)
		}

		return variableSet, nil
	default:
		return VariableSet{}, err
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

func (r *Repository) insertIntoEnvironments(ctx context.Context, variableSet VariableSet) (VariableSet, error) {
	stmt, err := r.db.Prepare(insertQuery)
	if err != nil {
		return VariableSet{}, fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()

	jsonValues, err := json.Marshal(variableSet.Values)
	if err != nil {
		return VariableSet{}, fmt.Errorf("encoding error: %w", err)
	}

	tenantID := sqlutil.TenantID(ctx)

	_, err = stmt.ExecContext(
		ctx,
		variableSet.ID,
		variableSet.Name,
		variableSet.Description,
		variableSet.CreatedAt,
		jsonValues,
		tenantID,
	)

	if err != nil {
		return VariableSet{}, fmt.Errorf("sql exec: %w", err)
	}

	return variableSet, nil
}

func (r *Repository) updateIntoEnvironments(ctx context.Context, variableSet VariableSet, oldId id.ID) (VariableSet, error) {
	jsonValues, err := json.Marshal(variableSet.Values)
	if err != nil {
		return VariableSet{}, fmt.Errorf("encoding error: %w", err)
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return VariableSet{}, err
	}
	defer tx.Rollback()

	query, params := sqlutil.Tenant(ctx, updateQuery, oldId,
		variableSet.Name,
		variableSet.Description,
		variableSet.CreatedAt,
		jsonValues,
	)

	_, err = tx.ExecContext(
		ctx,
		query,
		params...,
	)
	if err != nil {
		return VariableSet{}, fmt.Errorf("sql exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return VariableSet{}, fmt.Errorf("commit: %w", err)
	}

	return variableSet, nil
}

func (r *Repository) Provision(ctx context.Context, variableSet VariableSet) error {
	_, err := r.Create(ctx, variableSet)
	return err
}
