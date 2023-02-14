package testmock

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/config"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
)

var pgContainer *gnomock.Container

func GetTestingDatabase() (model.Repository, error) {
	db, err := GetRawTestingDatabase()
	if err != nil {
		return nil, err
	}

	return testdb.Postgres(testdb.WithDB(db))
}

func ConfigureDB(cfg *config.Config) error {
	pgContainer, err := getPostgresContainer()
	if err != nil {
		return err
	}

	cfg.Set("postgres.host", pgContainer.Host)
	cfg.Set("postgres.user", "tracetest")
	cfg.Set("postgres.password", "tracetest")
	cfg.Set("postgres.dbname", "postgres")
	cfg.Set("postgres.port", pgContainer.DefaultPort())

	return nil

}

func GetRawTestingDatabase() (*sql.DB, error) {
	pgContainer, err := getPostgresContainer()
	if err != nil {
		return nil, err
	}
	db, err := getMainDatabaseConnection(pgContainer)
	if err != nil {
		return nil, err
	}
	newDbConnection, err := createRandomDatabaseForTest(db, "tracetest")

	if err != nil {
		return nil, err
	}

	return newDbConnection, nil
}

func getMainDatabaseConnection(container *gnomock.Container) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s  dbname=%s sslmode=disable",
		container.Host, container.DefaultPort(), "tracetest", "tracetest", "postgres",
	)

	return sql.Open("postgres", connStr)
}

func createRandomDatabaseForTest(db *sql.DB, baseDatabase string) (*sql.DB, error) {
	epoch := time.Now().UnixNano()
	newDatabaseName := fmt.Sprintf("%s_%d", baseDatabase, epoch)
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s", newDatabaseName, baseDatabase))
	if err != nil {
		return nil, fmt.Errorf("could not create database %s: %w", newDatabaseName, err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s  dbname=%s sslmode=disable",
		pgContainer.Host, pgContainer.DefaultPort(), "tracetest", "tracetest", newDatabaseName,
	)

	return sql.Open("postgres", connStr)
}

func getPostgresContainer() (*gnomock.Container, error) {
	if pgContainer != nil {
		return pgContainer, nil
	}

	preset := postgres.Preset(
		postgres.WithUser("tracetest", "tracetest"),
		postgres.WithDatabase("tracetest"),
	)

	dbContainer, err := gnomock.Start(preset)
	if err != nil {
		return nil, fmt.Errorf("could not start postgres container")
	}

	pgContainer = dbContainer

	return dbContainer, nil
}
