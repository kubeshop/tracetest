package migrations_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMigrations(t *testing.T) {
	db, err := testmock.GetRawTestingDatabase()
	require.NoError(t, err)

	_, err = testdb.Postgres(testdb.WithMigrations("file://../migrations"), testdb.WithDB(db))
	require.NoError(t, err, "postgres migrations up should not fail")

	err = rollback(db)
	assert.NoError(t, err, "rollback should not fail")
}

func rollback(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not get driver from postgres connection: %w", err)
	}

	migrateClient, err := migrate.NewWithDatabaseInstance("file://../migrations/", "tracetest", driver)
	if err != nil {
		return fmt.Errorf("could not get migration client: %w", err)
	}

	err = migrateClient.Down()

	return err
}
