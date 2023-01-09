package testdb

import (
	"database/sql"
	"fmt"

	"github.com/denisbrodbeck/machineid"
)

func (td *postgresDB) ServerID() (string, bool, error) {
	isNew := false
	id := ""

	err := td.db.
		QueryRow(`SELECT id FROM "server" LIMIT 1`).
		Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return id, isNew, fmt.Errorf("could not get serverID from DB: %w", err)
	}

	if id != "" {
		return id, isNew, err
	}

	// there is no ID, let's create a new one and register on database
	isNew = true

	id, err = generateUniqueServerID()
	if err != nil {
		return id, isNew, fmt.Errorf("could not generate serverID: %w", err)
	}

	stmt, err := td.db.Prepare(`INSERT INTO "server" (id) VALUES ($1)`)
	if err != nil {
		return id, isNew, fmt.Errorf("could not prepare stmt: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return id, isNew, fmt.Errorf("could not save serverID into DB: %w", err)
	}

	return id, isNew, nil

}

func generateUniqueServerID() (string, error) {
	id, err := machineid.ProtectedID("tracetest")

	if err != nil {
		return "", fmt.Errorf("could not get machineID: %w", err)
	}

	id = id[:10] // limit lenght to avoid issues with GA

	return id, nil
}
