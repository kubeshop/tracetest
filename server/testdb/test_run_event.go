package testdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
)

const insertTestRunEventQuery = `
	INSERT INTO test_run_events (
		"test_id",
		"run_id",
		"type",
		"stage",
		"title",
		"description",
		"created_at",
		"data_store_connection",
		"polling",
		"outputs"
	) VALUES (
		$1, -- test_id
		$2, -- run_id
		$3, -- type
		$4, -- stage
		$5, -- title
		$6, -- description
		$7, -- created_at
		$8, -- data_store_connection
		$9, -- polling
		$10  -- outputs
	)
	RETURNING "id"
`

func (td *postgresDB) CreateTestRunEvent(ctx context.Context, event model.TestRunEvent) error {
	dataStoreConnectionJSON, err := json.Marshal(event.DataStoreConnection)
	if err != nil {
		return fmt.Errorf("could not marshal data store connection into JSON: %w", err)
	}

	pollingJSON, err := json.Marshal(event.Polling)
	if err != nil {
		return fmt.Errorf("could not marshal polling into JSON: %w", err)
	}

	outputsJSON, err := json.Marshal(event.Outputs)
	if err != nil {
		return fmt.Errorf("could not marshal outputs into JSON: %w", err)
	}

	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now()
	}

	err = td.db.QueryRowContext(
		ctx,
		insertTestRunEventQuery,
		event.TestID,
		event.RunID,
		event.Type,
		event.Stage,
		event.Title,
		event.Description,
		event.CreatedAt,
		dataStoreConnectionJSON,
		pollingJSON,
		outputsJSON,
	).Scan(&event.ID)

	if err != nil {
		return fmt.Errorf("could not insert event into database: %w", err)
	}

	return nil
}

const getTestRunEventsQuery = `
	SELECT
		"id",
		"test_id",
		"run_id",
		"type",
		"stage",
		"title",
		"description",
		"created_at",
		"data_store_connection",
		"polling",
		"outputs"
	FROM test_run_events WHERE "test_id" = $1 AND "run_id" = $2 ORDER BY "created_at" ASC;
`

func (td *postgresDB) GetTestRunEvents(ctx context.Context, testID id.ID, runID int) ([]model.TestRunEvent, error) {
	rows, err := td.db.QueryContext(ctx, getTestRunEventsQuery, testID, runID)
	if err != nil {
		return []model.TestRunEvent{}, fmt.Errorf("could not query test runs: %w", err)
	}

	events := make([]model.TestRunEvent, 0)

	for rows.Next() {
		event, err := readTestRunEventFromRows(rows)
		if err != nil {
			return []model.TestRunEvent{}, fmt.Errorf("could not parse row: %w", err)
		}

		events = append(events, event)
	}

	return events, nil
}

func readTestRunEventFromRows(rows *sql.Rows) (model.TestRunEvent, error) {
	var dataStoreConnectionBytes, pollingBytes, outputsBytes []byte
	event := model.TestRunEvent{}

	err := rows.Scan(
		&event.ID,
		&event.TestID,
		&event.RunID,
		&event.Type,
		&event.Stage,
		&event.Title,
		&event.Description,
		&event.CreatedAt,
		&dataStoreConnectionBytes,
		&pollingBytes,
		&outputsBytes,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.TestRunEvent{}, ErrNotFound
		}

		return model.TestRunEvent{}, fmt.Errorf("could not scan event: %w", err)
	}

	var dataStoreConnection model.ConnectionResult
	var polling model.PollingInfo
	outputs := make([]model.OutputInfo, 0)

	err = json.Unmarshal(dataStoreConnectionBytes, &dataStoreConnection)
	if err != nil {
		return model.TestRunEvent{}, fmt.Errorf("could not unmarshal data store connection: %w", err)
	}

	err = json.Unmarshal(pollingBytes, &polling)
	if err != nil {
		return model.TestRunEvent{}, fmt.Errorf("could not unmarshal polling information: %w", err)
	}

	err = json.Unmarshal(outputsBytes, &outputs)
	if err != nil {
		return model.TestRunEvent{}, fmt.Errorf("could not unmarshal outputs: %w", err)
	}

	event.DataStoreConnection = dataStoreConnection
	event.Polling = polling
	event.Outputs = outputs

	return event, nil
}
