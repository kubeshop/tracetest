package testdb

import (
	"database/sql"
	"fmt"

	"github.com/j2gg0s/otsql"
	"github.com/j2gg0s/otsql/hook/trace"
	pq "github.com/lib/pq"
)

type PostgresOption func(*postgresDB) error

func WithDSN(dsn string) PostgresOption {
	return func(pd *postgresDB) error {
		connector, err := pq.NewConnector(dsn)
		if err != nil {
			return fmt.Errorf("sql open: %w", err)
		}
		db := sql.OpenDB(
			otsql.WrapConnector(connector,
				otsql.WithHooks(
					trace.New(
						trace.WithQuery(true),
						trace.WithQueryParams(true),
						trace.WithRowsAffected(true),
					),
				),
			),
		)
		pd.db = db
		return nil
	}
}

func WithDB(db *sql.DB) PostgresOption {
	return func(pd *postgresDB) error {
		pd.db = db
		return nil
	}
}

func WithMigrations(migrationFolder string) PostgresOption {
	return func(pd *postgresDB) error {
		pd.migrationsFolder = migrationFolder

		return nil
	}
}
