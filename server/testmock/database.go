package testmock

import (
	"database/sql"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
)

const baseDatabaseName = "tracetest"

var singletonTestDatabaseEnvironment *testDatabaseEnvironment

type testDatabaseEnvironment struct {
	container      *gnomock.Container
	mainConnection *sql.DB

	mutex sync.Mutex
}

func getTestDatabaseEnvironment() *testDatabaseEnvironment {
	if singletonTestDatabaseEnvironment == nil {
		panic(fmt.Errorf("testing database environment not started"))
	}

	return singletonTestDatabaseEnvironment
}

func StartTestEnvironment() {
	if singletonTestDatabaseEnvironment != nil {
		return // Already started
	}

	db := &testDatabaseEnvironment{
		mutex: sync.Mutex{},
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()

	container, err := getPostgresContainer()
	if err != nil {
		panic(err)
	}
	db.container = container

	connection, err := getMainDatabaseConnection(db.container)
	if err != nil {
		panic(err)
	}
	db.mainConnection = connection

	// Starts this singleton only here, to guarantee that we
	// will only initiate this singleton when starting the environment
	singletonTestDatabaseEnvironment = db
}

func StopTestEnvironment() {
	db := getTestDatabaseEnvironment()

	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Close main connection
	if db.mainConnection != nil {
		err := db.mainConnection.Close()
		if err != nil {
			panic(err)
		}

		db.mainConnection = nil
	}

	if db.container != nil {
		err := gnomock.Stop(db.container)
		if err != nil {
			panic(err)
		}
	}
}

func GetTestingDatabase() model.Repository {
	dbConnection := GetRawTestingDatabase()
	return GetTestingDatabaseFromRawDB(dbConnection)
}

func GetTestingDatabaseFromRawDB(db *sql.DB) model.Repository {
	testingDatabase, err := testdb.Postgres(testdb.WithDB(db))
	if err != nil {
		panic(err)
	}

	return testingDatabase
}

func CreateMigratedDatabase() *sql.DB {
	newConn, err := createRandomDatabaseForTest(baseDatabaseName)
	if err != nil {
		panic(err)
	}

	// migrate DB
	_, err = testdb.Postgres(testdb.WithDB(newConn))
	if err != nil {
		panic(err)
	}

	return newConn
}

func GetRawTestingDatabase() *sql.DB {
	newDbConnection, err := createRandomDatabaseForTest(baseDatabaseName)

	if err != nil {
		panic(err)
	}

	return newDbConnection
}

func createRandomDatabaseForTest(baseDatabase string) (*sql.DB, error) {
	db := getTestDatabaseEnvironment()

	newDatabaseName := fmt.Sprintf("%s_%d%d%d", baseDatabase, randomInt(), randomInt(), randomInt())
	_, err := db.mainConnection.Exec(fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s", newDatabaseName, baseDatabase))
	if err != nil {
		return nil, fmt.Errorf("could not create database %s: %w", newDatabaseName, err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s  dbname=%s sslmode=disable",
		db.container.Host, db.container.DefaultPort(), "tracetest", "tracetest", newDatabaseName,
	)

	return sql.Open("postgres", connStr)
}

func getPostgresContainer() (*gnomock.Container, error) {
	preset := postgres.Preset(
		postgres.WithUser("tracetest", "tracetest"),
		postgres.WithDatabase("tracetest"),
	)

	dbContainer, err := gnomock.Start(preset)
	if err != nil {
		return nil, fmt.Errorf("could not start postgres container")
	}

	return dbContainer, nil
}

func getMainDatabaseConnection(container *gnomock.Container) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s  dbname=%s sslmode=disable",
		container.Host, container.DefaultPort(), "tracetest", "tracetest", "postgres",
	)

	return sql.Open("postgres", connStr)
}

func randomInt() int {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 1000000
	return rand.Intn(max-min) + min
}
