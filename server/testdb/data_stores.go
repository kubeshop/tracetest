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

var _ model.DataStoreRepository = &postgresDB{}

const (
	insertIntoDataStoresQuery = `
		INSERT INTO data_stores (
			"id",
			"name",
			"type",
			"is_default",
			"values",
			"created_at"
		) VALUES ($1, $2, $3, $4, $5, $6)`

	updateIntoDataStoresQuery = `
		UPDATE data_stores SET
			"id" = $2,
			"name" = $3,
			"type" = $4,
			"is_default" = $5,
			"values" = $6,
			"created_at" = $7
		WHERE id = $1
	`

	updateAllDefaultDataStoresQuery = `
		UPDATE data_stores SET "is_default" = false
	`

	getFromDataStoresQuery = `
	SELECT
		d.id,
		d.name,
		d.type,
		d.is_default,
		d.values,
		d.created_at
	FROM data_stores d
`

	deleteFromDataStoresQuery = "DELETE FROM data_stores WHERE id = $1"

	idExistsFromDataStoresQuery = `
		SELECT COUNT(*) > 0 as exists FROM data_stores WHERE id = $1
	`

	countFromDataStoresQuery = `
		SELECT COUNT(*) FROM data_stores d
	`
)

func (td *postgresDB) CreateDataStore(ctx context.Context, dataStore model.DataStore) (model.DataStore, error) {
	dataStore.ID = IDGen.ID().String()
	dataStore.CreatedAt = time.Now()

	return td.insertIntoDataStores(ctx, dataStore)
}

func (td *postgresDB) UpdateDataStore(ctx context.Context, dataStore model.DataStore) (model.DataStore, error) {
	oldDataStore, err := td.GetDataStore(ctx, dataStore.ID)
	if err != nil {
		return model.DataStore{}, fmt.Errorf("could not get the data store while updating: %w", err)
	}

	// keep the same creation date to keep sort order
	dataStore.CreatedAt = oldDataStore.CreatedAt
	dataStore.ID = oldDataStore.ID

	return td.updateIntoDataStores(ctx, dataStore, oldDataStore.ID)
}

func (td *postgresDB) DeleteDataStore(ctx context.Context, dataStore model.DataStore) error {
	tx, err := td.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("sql BeginTx: %w", err)
	}

	_, err = tx.ExecContext(ctx, deleteFromDataStoresQuery, dataStore.ID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("sql error: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("sql Commit: %w", err)
	}

	return nil
}

func (td *postgresDB) DefaultDataStore(ctx context.Context) (model.DataStore, error) {
	stmt, err := td.db.Prepare(getFromDataStoresQuery + " WHERE d.is_default = true ORDER BY created_at DESC LIMIT 1")

	if err != nil {
		return model.DataStore{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	dataStore, err := td.readDataStoreRow(ctx, stmt.QueryRowContext(ctx))
	// if no default is found, assume nothing is configured, return empty DS without error
	if err != nil && err != ErrNotFound {
		return model.DataStore{}, err
	}

	return dataStore, nil
}

func (td *postgresDB) GetDataStore(ctx context.Context, id string) (model.DataStore, error) {
	stmt, err := td.db.Prepare(getFromDataStoresQuery + " WHERE d.id = $1")

	if err != nil {
		return model.DataStore{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	dataStore, err := td.readDataStoreRow(ctx, stmt.QueryRowContext(ctx, id))
	if err != nil {
		return model.DataStore{}, err
	}

	return dataStore, nil
}

func (td *postgresDB) GetDataStores(ctx context.Context, take, skip int32, query, sortBy, sortDirection string) (model.List[model.DataStore], error) {
	hasSearchQuery := query != ""
	cleanSearchQuery := "%" + strings.ReplaceAll(query, " ", "%") + "%"
	params := []any{take, skip}

	sql := getFromDataStoresQuery

	const condition = "WHERE (d.name ilike $3) "
	if hasSearchQuery {
		params = append(params, cleanSearchQuery)
		sql += condition
	}

	sortingFields := map[string]string{
		"created": "d.created_at",
		"name":    "d.name",
	}

	sql = sortQuery(sql, sortBy, sortDirection, sortingFields)
	sql += ` LIMIT $1 OFFSET $2 `

	stmt, err := td.db.Prepare(sql)
	if err != nil {
		return model.List[model.DataStore]{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, params...)
	if err != nil {
		return model.List[model.DataStore]{}, err
	}

	dataStores := []model.DataStore{}

	for rows.Next() {
		dataStore, err := td.readDataStoreRow(ctx, rows)
		if err != nil {
			return model.List[model.DataStore]{}, err
		}

		dataStores = append(dataStores, dataStore)
	}

	count, err := td.countDataStores(ctx, condition, cleanSearchQuery)
	if err != nil {
		return model.List[model.DataStore]{}, err
	}

	return model.List[model.DataStore]{
		Items:      dataStores,
		TotalCount: count,
	}, nil
}

func (td *postgresDB) DataStoreIDExists(ctx context.Context, id string) (bool, error) {
	exists := false

	row := td.db.QueryRowContext(ctx, idExistsFromDataStoresQuery, id)

	err := row.Scan(&exists)

	return exists, err
}

func (td *postgresDB) readDataStoreRow(ctx context.Context, row scanner) (model.DataStore, error) {
	dataStore := model.DataStore{}

	var (
		jsonValues []byte
	)
	err := row.Scan(
		&dataStore.ID,
		&dataStore.Name,
		&dataStore.Type,
		&dataStore.IsDefault,
		&jsonValues,
		&dataStore.CreatedAt,
	)

	switch err {
	case sql.ErrNoRows:
		return model.DataStore{}, ErrNotFound
	case nil:
		err = json.Unmarshal(jsonValues, &dataStore.Values)
		if err != nil {
			return model.DataStore{}, fmt.Errorf("cannot parse data store: %w", err)
		}

		return dataStore, nil
	default:
		return model.DataStore{}, err
	}
}

func (td *postgresDB) countDataStores(ctx context.Context, condition, cleanSearchQuery string) (int, error) {
	var (
		count  int
		params []any
	)

	sql := countFromDataStoresQuery
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

func (td *postgresDB) insertIntoDataStores(ctx context.Context, dataStore model.DataStore) (model.DataStore, error) {
	tx, err := td.db.BeginTx(ctx, nil)
	if err != nil {
		return model.DataStore{}, fmt.Errorf("sql BeginTx: %w", err)
	}

	if dataStore.IsDefault {
		_, err = tx.ExecContext(ctx, updateAllDefaultDataStoresQuery)

		if err != nil {
			tx.Rollback()
			return model.DataStore{}, fmt.Errorf("sql exec: %w", err)
		}
	}

	jsonValues, err := json.Marshal(dataStore.Values)
	if err != nil {
		return model.DataStore{}, fmt.Errorf("encoding error: %w", err)
	}

	_, err = tx.ExecContext(ctx, insertIntoDataStoresQuery,
		dataStore.ID,
		dataStore.Name,
		dataStore.Type,
		dataStore.IsDefault,
		jsonValues,
		dataStore.CreatedAt)

	if err != nil {
		tx.Rollback()
		return model.DataStore{}, fmt.Errorf("sql exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return model.DataStore{}, fmt.Errorf("commit: %w", err)
	}

	return dataStore, nil
}

func (td *postgresDB) updateIntoDataStores(ctx context.Context, dataStore model.DataStore, oldId string) (model.DataStore, error) {
	tx, err := td.db.BeginTx(ctx, nil)
	if err != nil {
		return model.DataStore{}, fmt.Errorf("sql BeginTx: %w", err)
	}

	if dataStore.IsDefault {
		_, err = tx.ExecContext(ctx, updateAllDefaultDataStoresQuery)

		if err != nil {
			tx.Rollback()
			return model.DataStore{}, fmt.Errorf("sql exec: %w", err)
		}
	}

	jsonValues, err := json.Marshal(dataStore.Values)
	if err != nil {
		return model.DataStore{}, fmt.Errorf("encoding error: %w", err)
	}

	_, err = tx.ExecContext(ctx, updateIntoDataStoresQuery,
		oldId,
		dataStore.ID,
		dataStore.Name,
		dataStore.Type,
		dataStore.IsDefault,
		jsonValues,
		dataStore.CreatedAt)

	if err != nil {
		tx.Rollback()
		return model.DataStore{}, fmt.Errorf("sql exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return model.DataStore{}, fmt.Errorf("commit: %w", err)
	}

	return dataStore, nil
}
