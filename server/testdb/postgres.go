package testdb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kubeshop/tracetest/id"
	"github.com/kubeshop/tracetest/model"
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

func Postgres(options ...PostgresOption) (model.Repository, error) {
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

func (td *postgresDB) Drop() error {
	return dropTables(td, "runs", "definitions", "tests", "schema_migrations")
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
