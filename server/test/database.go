package test

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
)

var container *gnomock.Container

func GetTestingDatabase() (*sql.DB, error) {
	container, err := getPostgresContainer()
	db, err := getMainDatabaseConnection(container)
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
		container.Host, container.DefaultPort(), "tracetest", "tracetest", newDatabaseName,
	)

	return sql.Open("postgres", connStr)
}

func getPostgresContainer() (*gnomock.Container, error) {
	if container != nil {
		return container, nil
	}

	preset := postgres.Preset(
		postgres.WithUser("tracetest", "tracetest"),
		postgres.WithDatabase("tracetest"),
	)

	dbContainer, err := gnomock.Start(preset)
	if err != nil {
		return nil, fmt.Errorf("could not start postgres container")
	}

	container = dbContainer

	return dbContainer, nil
}
