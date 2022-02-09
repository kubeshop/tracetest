package testdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	openapi "github.com/GIT_USER_ID/GIT_REPO_ID/go"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type TestDB struct {
	db *sql.DB
}

func New(connStr string) (*TestDB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS tests  (
	id UUID NOT NULL PRIMARY KEY,
	test json NOT NULL
);
`)
	if err != nil {
		return nil, err
	}
	return &TestDB{
		db: db,
	}, nil
}

func (td *TestDB) CreateTest(ctx context.Context, test *openapi.Test) (string, error) {
	stmt, err := td.db.Prepare("INSERT INTO tests(id, test) VALUES( $1, $2 )")
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	id := uuid.New().String()
	test.Id = id
	b, err := json.Marshal(test)
	if err != nil {
		return "", err
	}
	_, err = stmt.ExecContext(ctx, id, b)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (td *TestDB) GetTest(ctx context.Context, id string) (*openapi.Test, error) {
	stmt, err := td.db.Prepare("SELECT test FROM tests WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var b []byte
	err = stmt.QueryRowContext(ctx, id).Scan(&b)
	if err != nil {
		return nil, err
	}
	var test openapi.Test

	err = json.Unmarshal(b, &test)
	if err != nil {
		return nil, err
	}
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
		tests = append(tests, test)
	}
	return tests, nil
}

func (td *TestDB) Drop() error {
	_, err := td.db.Exec(`
DROP TABLE IF EXISTS tests;
`)
	if err != nil {
		return err
	}
	return nil
}
