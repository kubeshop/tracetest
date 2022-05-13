package testdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/kubeshop/tracetest/openapi"
)

type postgresDB struct {
	db               *sql.DB
	migrationsFolder string
}

func Postgres(options ...PostgresOption) (Repository, error) {
	postgres := &postgresDB{}
	postgres.migrationsFolder = "file://./migrations"
	for _, option := range options {
		err := option(postgres)
		if err != nil {
			return nil, err
		}
	}

	err := postgres.ensureLatestMigration()
	if err != nil {
		return nil, fmt.Errorf("could not execute migrations: %w", err)
	}

	return postgres, nil
}

func (p *postgresDB) ensureLatestMigration() error {
	driver, err := postgres.WithInstance(p.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not get driver from postgres connection: %w", err)
	}
	migrateClient, err := migrate.NewWithDatabaseInstance(p.migrationsFolder, "tracetest", driver)
	if err != nil {
		return fmt.Errorf("could not get migration client: %w", err)
	}

	err = migrateClient.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	return nil
}

func (td *postgresDB) CreateTest(ctx context.Context, test *openapi.Test) (string, error) {
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

func (td *postgresDB) UpdateTest(ctx context.Context, test *openapi.Test) error {
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

func (td *postgresDB) DeleteTest(ctx context.Context, test *openapi.Test) error {
	queries := []string{
		"DELETE FROM tests WHERE id = $1",
		"DELETE FROM assertions WHERE test_id = $1",
		"DELETE FROM results WHERE test_id = $1",
	}

	for _, sql := range queries {
		fmt.Println("aca", sql)
		_, err := td.db.Exec(sql, test.TestId)
		if err != nil {
			return fmt.Errorf("sql prepare: %w", err)
		}
	}

	return nil
}

func (td *postgresDB) GetTest(ctx context.Context, id string) (*openapi.Test, error) {
	stmt, err := td.db.Prepare("SELECT test FROM tests WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("prepare: %w", err)
	}
	defer stmt.Close()

	var b []byte
	err = stmt.QueryRowContext(ctx, id).Scan(&b)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, ErrNotFound
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

func (td *postgresDB) GetTests(ctx context.Context, take, skip int32) ([]openapi.Test, error) {
	stmt, err := td.db.Prepare("SELECT test FROM tests LIMIT $1 OFFSET $2")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, take, skip)
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

func (td *postgresDB) CreateAssertion(ctx context.Context, testid string, assertion *openapi.Assertion) (string, error) {
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

func (td *postgresDB) UpdateAssertion(ctx context.Context, testID, assertionID string, assertion openapi.Assertion) error {
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

func (td *postgresDB) DeleteAssertion(ctx context.Context, testID, assertionID string) error {
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

func (td *postgresDB) GetAssertion(ctx context.Context, id string) (*openapi.Assertion, error) {
	stmt, err := td.db.Prepare("SELECT assertion FROM assertions WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var b []byte
	err = stmt.QueryRowContext(ctx, id).Scan(&b)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, ErrNotFound
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

func (td *postgresDB) GetAssertionsByTestID(ctx context.Context, testID string) ([]openapi.Assertion, error) {
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

func (td *postgresDB) Drop() error {
	return dropTables(td, "results", "assertions", "tests", "schema_migrations")
}

func dropTables(td *postgresDB, tables ...string) error {
	for _, table := range tables {
		_, err := td.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s;", table))
		if err != nil {
			return err
		}
	}

	return nil
}

func (td *postgresDB) CreateResult(ctx context.Context, testid string, run *openapi.TestRun) error {
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

func (td *postgresDB) GetResult(ctx context.Context, id string) (*openapi.TestRun, error) {
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
	var run openapi.TestRun

	err = json.Unmarshal(b, &run)
	if err != nil {
		return nil, err
	}
	return &run, nil
}

func (td *postgresDB) GetResultByTraceID(ctx context.Context, testID, traceID string) (openapi.TestRun, error) {
	stmt, err := td.db.Prepare("SELECT result FROM results WHERE test_id = $1 AND result ->> 'traceId' = $2")
	if err != nil {
		return openapi.TestRun{}, err
	}
	defer stmt.Close()

	var b []byte
	err = stmt.QueryRowContext(ctx, testID, traceID).Scan(&b)
	if err != nil {
		return openapi.TestRun{}, err
	}
	var run openapi.TestRun

	err = json.Unmarshal(b, &run)
	if err != nil {
		return openapi.TestRun{}, err
	}
	return run, nil
}

func (td *postgresDB) GetResultsByTestID(ctx context.Context, testID string, take, skip int32) ([]openapi.TestRun, error) {
	stmt, err := td.db.Prepare("SELECT result FROM results WHERE test_id = $1 ORDER BY result ->> 'createdAt' DESC LIMIT $2 OFFSET $3")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, testID, take, skip) //.Scan(&b)
	if err != nil {
		return nil, err
	}
	var run []openapi.TestRun

	for rows.Next() {
		var b []byte
		if err := rows.Scan(&b); err != nil {
			return nil, err
		}
		var result openapi.TestRun
		err = json.Unmarshal(b, &result)
		if err != nil {
			return nil, err
		}

		run = append(run, result)
	}

	return run, nil
}

func (td *postgresDB) UpdateResult(ctx context.Context, run *openapi.TestRun) error {
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
