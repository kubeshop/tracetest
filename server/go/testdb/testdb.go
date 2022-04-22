package testdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/j2gg0s/otsql"
	"github.com/j2gg0s/otsql/hook/trace"
	openapi "github.com/kubeshop/tracetest/server/go"
	pq "github.com/lib/pq"
)

type TestDB struct {
	db *sql.DB
}

func New(dsn string) (*TestDB, error) {
	connector, err := pq.NewConnector(dsn)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}
	db := sql.OpenDB(
		otsql.WrapConnector(connector,
			otsql.WithHooks(trace.New(
				trace.WithQuery(true),
				trace.WithQueryParams(true),
				trace.WithRowsAffected(true),
			))))

	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS tests  (
	id UUID NOT NULL PRIMARY KEY,
	test json NOT NULL
);
`)
	if err != nil {
		return nil, fmt.Errorf("create table tests: %w", err)
	}
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS results  (
	id UUID NOT NULL PRIMARY KEY,
	test_id UUID NOT NULL,
	result json NOT NULL
);
`)
	if err != nil {
		return nil, fmt.Errorf("create table results: %w", err)
	}
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS assertions  (
	id UUID NOT NULL PRIMARY KEY,
	test_id UUID NOT NULL,
	assertion json NOT NULL
);
`)
	if err != nil {
		return nil, fmt.Errorf("create table assertions: %w", err)
	}

	return &TestDB{
		db: db,
	}, nil
}

func (td *TestDB) CreateTest(ctx context.Context, test *openapi.Test) (string, error) {
	stmt, err := td.db.Prepare("INSERT INTO tests(id, test) VALUES( $1, $2 )")
	if err != nil {
		return "", fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()
	id := uuid.New().String()
	test.TestId = id
	b, err := json.Marshal(test)
	if err != nil {
		return "", fmt.Errorf("json Marshal: %w", err)
	}
	_, err = stmt.ExecContext(ctx, id, b)
	if err != nil {
		return "", fmt.Errorf("sql exec: %w", err)
	}

	return id, nil
}

func (td *TestDB) UpdateTest(ctx context.Context, test *openapi.Test) error {
	stmt, err := td.db.Prepare("UPDATE tests SET test = $2 WHERE id = $1")
	if err != nil {
		return fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()
	b, err := json.Marshal(test)
	if err != nil {
		return fmt.Errorf("json Marshal: %w", err)
	}
	_, err = stmt.Exec(test.TestId, b)
	if err != nil {
		return fmt.Errorf("sql exec: %w", err)
	}

	return nil
}

func (td *TestDB) GetTest(ctx context.Context, id string) (*openapi.Test, error) {
	stmt, err := td.db.Prepare("SELECT test FROM tests WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	var b []byte
	err = stmt.QueryRowContext(ctx, id).Scan(&b)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, openapi.ErrNotFound
		}
		return nil, fmt.Errorf("query: %w", err)
	}
	var test openapi.Test

	err = json.Unmarshal(b, &test)
	if err != nil {
		return nil, err
	}
	as, err := td.GetAssertionsByTestID(ctx, id)
	if err != nil {
		return nil, err
	}

	test.Assertions = as

	return &test, nil
}

func (td *TestDB) GetTests(ctx context.Context) ([]openapi.Test, error) {
	stmt, err := td.db.Prepare("SELECT test FROM tests")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	var tests []openapi.Test
	for rows.Next() {
		var b []byte
		if err := rows.Scan(&b); err != nil {
			return nil, err
		}
		var test openapi.Test
		err = json.Unmarshal(b, &test)
		if err != nil {
			return nil, err
		}

		as, err := td.GetAssertionsByTestID(ctx, test.TestId)
		if err != nil {
			return nil, err
		}

		test.Assertions = as

		tests = append(tests, test)
	}
	return tests, nil
}

func (td *TestDB) CreateAssertion(ctx context.Context, testid string, assertion *openapi.Assertion) (string, error) {
	for i := range assertion.SpanAssertions {
		assertion.SpanAssertions[i].SpanAssertionId = ensureUUID(assertion.SpanAssertions[i].SpanAssertionId)
	}

	stmt, err := td.db.Prepare("INSERT INTO assertions(id, test_id, assertion) VALUES( $1, $2, $3 )")
	if err != nil {
		return "", fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()
	id := uuid.New().String()
	assertion.AssertionId = id
	b, err := json.Marshal(assertion)
	if err != nil {
		return "", fmt.Errorf("json Marshal: %w", err)
	}
	_, err = stmt.ExecContext(ctx, id, testid, b)
	if err != nil {
		return "", fmt.Errorf("sql exec: %w", err)
	}

	return id, nil
}

func ensureUUID(in string) string {
	if in == "" {
		return uuid.New().String()
	}

	return in
}

func (td *TestDB) UpdateAssertion(ctx context.Context, testID, assertionID string, assertion openapi.Assertion) error {
	for i := range assertion.SpanAssertions {
		assertion.SpanAssertions[i].SpanAssertionId = ensureUUID(assertion.SpanAssertions[i].SpanAssertionId)
	}

	stmt, err := td.db.Prepare("UPDATE assertions SET assertion = $3 WHERE id = $1 AND test_id = $2")
	if err != nil {
		return fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()
	b, err := json.Marshal(assertion)
	if err != nil {
		return fmt.Errorf("json Marshal: %w", err)
	}
	_, err = stmt.ExecContext(ctx, assertionID, testID, b)
	if err != nil {
		return fmt.Errorf("sql exec: %w", err)
	}

	return nil
}

func (td *TestDB) DeleteAssertion(ctx context.Context, testID, assertionID string) error {
	stmt, err := td.db.Prepare("DELETE FROM assertions WHERE id = $1 AND test_id = $2")
	if err != nil {
		return fmt.Errorf("sql prepare: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, assertionID, testID)
	if err != nil {
		return fmt.Errorf("sql exec: %w", err)
	}

	return nil
}

func (td *TestDB) GetAssertion(ctx context.Context, id string) (*openapi.Assertion, error) {
	stmt, err := td.db.Prepare("SELECT assertion FROM assertions WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var b []byte
	err = stmt.QueryRowContext(ctx, id).Scan(&b)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, openapi.ErrNotFound
		}
		return nil, fmt.Errorf("query: %w", err)
	}
	var assertion openapi.Assertion

	err = json.Unmarshal(b, &assertion)
	if err != nil {
		return nil, err
	}
	return &assertion, nil
}

func (td *TestDB) GetAssertionsByTestID(ctx context.Context, testID string) ([]openapi.Assertion, error) {
	stmt, err := td.db.Prepare("SELECT assertion FROM assertions WHERE test_id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, testID)
	if err != nil {
		return nil, err
	}
	var assertions []openapi.Assertion
	for rows.Next() {
		var b []byte
		if err := rows.Scan(&b); err != nil {
			return nil, err
		}
		var assertion openapi.Assertion
		err = json.Unmarshal(b, &assertion)
		if err != nil {
			return nil, err
		}
		assertions = append(assertions, assertion)
	}
	return assertions, nil
}

func (td *TestDB) Drop() error {
	_, err := td.db.Exec(`
DROP TABLE IF EXISTS tests;
`)
	if err != nil {
		return err
	}
	_, err = td.db.Exec(`
DROP TABLE IF EXISTS results;
`)
	if err != nil {
		return err
	}

	return nil
}
