package testdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/kubeshop/tracetest/model"
)

var _ model.DefinitionRepository = &postgresDB{}

func (td *postgresDB) GetDefiniton(ctx context.Context, test model.Test) (model.Definition, error) {
	stmt, err := td.db.Prepare("SELECT definition FROM definitions WHERE test_id = $1")
	if err != nil {
		return model.Definition{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	def, err := readDefinitionRow(stmt.QueryRowContext(ctx, test.ID))
	if err != nil {
		return model.Definition{}, err
	}
	return def, nil
}

func (td *postgresDB) SetDefiniton(ctx context.Context, t model.Test, d model.Definition) error {
	stmt, err := td.db.Prepare(`
		INSERT INTO definition (id, test_id, definition) VALUES ($1, $2, $3)
		ON CONFLICT DO
			UPDATE definition SET definition = $3 WHERE id = $1 AND test_id = $2
	`)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	def, err := encodeDefinition(d)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, t.ID, def)
	if err != nil {
		return fmt.Errorf("sql exec: %w", err)
	}

	return nil
}

func encodeDefinition(d model.Definition) (string, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return "", fmt.Errorf("json Marshal: %w", err)
	}

	return string(b), nil
}

func decodeDefinition(b []byte) (model.Definition, error) {
	var def model.Definition
	err := json.Unmarshal(b, &def)
	if err != nil {
		return model.Definition{}, err
	}
	return def, nil
}

func readDefinitionRow(row scanner) (model.Definition, error) {
	var b []byte
	err := row.Scan(&b)
	switch err {
	case sql.ErrNoRows:
		return model.Definition{}, ErrNotFound
	case nil:
		return decodeDefinition(b)
	default:
		return model.Definition{}, err
	}
}
