package testdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/kubeshop/tracetest/server/model"
)

var _ model.RunRepository = &postgresDB{}

func (td *postgresDB) GetSpec(ctx context.Context, test model.Test) (model.OrderedMap[model.SpanQuery, []model.Assertion], error) {
	stmt, err := td.db.Prepare("SELECT definition FROM definitions WHERE test_id = $1 and test_version = $2")
	if err != nil {
		return model.OrderedMap[model.SpanQuery, []model.Assertion]{}, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	def, err := readDefinitionRow(stmt.QueryRowContext(ctx, test.ID, test.Version))
	if err != nil {
		return model.OrderedMap[model.SpanQuery, []model.Assertion]{}, err
	}
	return def, nil
}

func (td *postgresDB) SetSpec(ctx context.Context, t model.Test, d model.OrderedMap[model.SpanQuery, []model.Assertion]) error {
	sql := `UPDATE definitions SET definition = $3 WHERE test_id = $1 AND test_version = $2`
	if _, err := td.GetSpec(ctx, t); err == ErrNotFound {
		sql = `INSERT INTO definitions (test_id, test_version, "definition") VALUES ($1, $2, $3)`
	}

	stmt, err := td.db.Prepare(sql)
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	def, err := encodeDefinition(d)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, t.ID, t.Version, def)
	if err != nil {
		return fmt.Errorf("sql exec: %w", err)
	}

	return nil
}

func encodeDefinition(d model.OrderedMap[model.SpanQuery, []model.Assertion]) (string, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return "", fmt.Errorf("json Marshal: %w", err)
	}

	return string(b), nil
}

func decodeDefinition(b []byte) (model.OrderedMap[model.SpanQuery, []model.Assertion], error) {
	var def model.OrderedMap[model.SpanQuery, []model.Assertion]
	err := json.Unmarshal(b, &def)
	if err != nil {
		return model.OrderedMap[model.SpanQuery, []model.Assertion]{}, err
	}
	return def, nil
}

func readDefinitionRow(row scanner) (model.OrderedMap[model.SpanQuery, []model.Assertion], error) {
	var b []byte
	err := row.Scan(&b)
	switch err {
	case sql.ErrNoRows:
		return model.OrderedMap[model.SpanQuery, []model.Assertion]{}, ErrNotFound
	case nil:
		return decodeDefinition(b)
	default:
		return model.OrderedMap[model.SpanQuery, []model.Assertion]{}, err
	}
}
