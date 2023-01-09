package testdb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kubeshop/tracetest/server/id"
)

type postgresDB struct {
	db               *sql.DB
	migrationsFolder string
}

var (
	IDGen       = id.NewRandGenerator()
	ErrNotFound = errors.New("record not found")
)

type scanner interface {
	Scan(dest ...interface{}) error
}

func Postgres(options ...PostgresOption) (*postgresDB, error) {
	ps := &postgresDB{
		migrationsFolder: "file://./migrations",
	}
	for _, option := range options {
		err := option(ps)
		if err != nil {
			return nil, err
		}
	}

	err := ps.ensureLatestMigration()
	if err != nil {
		return nil, fmt.Errorf("could not execute migrations: %w", err)
	}

	return ps, nil
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

func (td *postgresDB) Drop() error {
	return dropTables(
		td,
		"transaction_run_steps",
		"transaction_runs",
		"transaction_steps",
		"transactions",
		"test_runs",
		"tests",
		"environments",
		"data_stores",
		"server",
		"schema_migrations",
	)
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

func (td *postgresDB) Close() error {
	return td.db.Close()
}
