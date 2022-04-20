package testdb

import (
	"context"
	"encoding/json"
	"fmt"

	openapi "github.com/kubeshop/tracetest/server/go"
)

func (td *TestDB) CreateResult(ctx context.Context, testid string, run *openapi.TestRunResult) error {
	stmt, err := td.db.Prepare("INSERT INTO results(id, test_id, result) VALUES( $1, $2, $3 )")
	if err != nil {
		return fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()
	b, err := json.Marshal(run)
	if err != nil {
		return fmt.Errorf("json Marshal: %w", err)
	}
	_, err = stmt.ExecContext(ctx, run.ResultId, testid, b)
	if err != nil {
		return fmt.Errorf("sql exec: %w", err)
	}

	return nil
}

func (td *TestDB) GetResult(ctx context.Context, id string) (*openapi.TestRunResult, error) {
	stmt, err := td.db.Prepare("SELECT result FROM results WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var b []byte
	err = stmt.QueryRowContext(ctx, id).Scan(&b)
	if err != nil {
		return nil, err
	}
	var run openapi.TestRunResult

	err = json.Unmarshal(b, &run)
	if err != nil {
		return nil, err
	}
	return &run, nil
}

func (td *TestDB) GetResultsByTraceID(ctx context.Context, testID, traceID string) (openapi.TestRunResult, error) {
	stmt, err := td.db.Prepare("SELECT result FROM results WHERE test_id = $1 AND result ->> 'traceId' = $2")
	if err != nil {
		return openapi.TestRunResult{}, err
	}
	defer stmt.Close()

	var b []byte
	err = stmt.QueryRowContext(ctx, testID, traceID).Scan(&b)
	if err != nil {
		return openapi.TestRunResult{}, err
	}
	var run openapi.TestRunResult

	err = json.Unmarshal(b, &run)
	if err != nil {
		return openapi.TestRunResult{}, err
	}
	return run, nil
}

func (td *TestDB) GetResultsByTestID(ctx context.Context, testID string) ([]openapi.TestRunResult, error) {
	stmt, err := td.db.Prepare("SELECT result FROM results WHERE test_id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, testID) //.Scan(&b)
	if err != nil {
		return nil, err
	}
	var run []openapi.TestRunResult

	for rows.Next() {
		var b []byte
		if err := rows.Scan(&b); err != nil {
			return nil, err
		}
		var result openapi.TestRunResult
		err = json.Unmarshal(b, &result)
		if err != nil {
			return nil, err
		}

		run = append(run, result)
	}

	return run, nil
}

func (td *TestDB) UpdateResult(ctx context.Context, run *openapi.TestRunResult) error {
	stmt, err := td.db.Prepare("UPDATE results SET result = $2 WHERE id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	b, err := json.Marshal(run)
	if err != nil {
		return fmt.Errorf("json Marshal: %w", err)
	}
	_, err = stmt.Exec(run.ResultId, b)
	if err != nil {
		return fmt.Errorf("sql exec: %w", err)
	}

	return nil
}
