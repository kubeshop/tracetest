package migrations_test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/kubeshop/tracetest/server/migrations"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrations(t *testing.T) {
	db := testmock.GetRawTestingDatabase()

	t.Run("applying migrations", func(t *testing.T) {
		_, err := testdb.Postgres(testdb.WithDB(db))
		require.NoError(t, err, "postgres migrations up should not fail")
	})

	t.Run("rolling back migrations", func(t *testing.T) {
		err := rollback(db)
		assert.NoError(t, err, "rollback should not fail")
	})
}

func rollback(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
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

	err = migrateClient.Down()

	return err
}
