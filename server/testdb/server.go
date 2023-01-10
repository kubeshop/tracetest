package testdb

import (
	"database/sql"
	"fmt"

	"github.com/kubeshop/tracetest/server/config"
)

func (td *postgresDB) ServerID() (id string, isNew bool, err error) {
	isNew = false
	id = ""

	err = td.db.
		QueryRow(`SELECT id FROM "server" LIMIT 1`).
		Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		err = fmt.Errorf("could not get serverID from DB: %w", err)
		return
	}

	if id != "" {
		return
	}

	// no id, let's creat it
	isNew = true
	id = config.GetMachineID()

	stmt, err := td.db.Prepare(`INSERT INTO "server" (id) VALUES ($1)`)
	if err != nil {
		err = fmt.Errorf("could not prepare stmt: %w", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		err = fmt.Errorf("could not save serverID into DB: %w", err)
		return
	}

	return

}
