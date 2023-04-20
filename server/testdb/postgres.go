package testdb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/kubeshop/tracetest/server/migrations"
	"github.com/kubeshop/tracetest/server/pkg/id"
)

type postgresDB struct {
	db *sql.DB
}

var (
	IDGen       = id.NewRandGenerator()
	ErrNotFound = errors.New("record not found")
)

type scanner interface {
	Scan(dest ...interface{}) error
}

func Postgres(options ...PostgresOption) (*postgresDB, error) {
	ps := &postgresDB{}
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
	sourceDriver, err := iofs.New(migrations.Migrations, ".")
	if err != nil {
		log.Fatal(err)
	}

	migrateClient, err := migrate.NewWithInstance("iofs", sourceDriver, "tracetest", driver)
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
