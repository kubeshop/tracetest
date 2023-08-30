package testdb_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/testmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunEvents(t *testing.T) {
	rawDB := testmock.GetRawTestingDatabase()
	db := testmock.GetTestingDatabaseFromRawDB(rawDB)
	defer rawDB.Close()

	testRepo := test.NewRepository(rawDB)
	testRunRepo := test.NewRunRepository(rawDB, test.NewCache("test"))

	test1 := createTestWithName(t, testRepo, "test 1")

	run1 := createRun(t, testRunRepo, test1)
	run2 := createRun(t, testRunRepo, test1)

	events := []model.TestRunEvent{
		{TestID: test1.ID, RunID: run1.ID, Type: "EVENT_1", Stage: model.StageTrigger, Title: "OP 1", Description: "This happened"},
		{TestID: test1.ID, RunID: run1.ID, Type: "EVENT_2", Stage: model.StageTrigger, Title: "OP 2", Description: "That happened now"},

		{TestID: test1.ID, RunID: run2.ID, Type: "EVENT_1", Stage: model.StageTrigger, Title: "OP 1", Description: "That happened", DataStoreConnection: model.ConnectionResult{
			PortCheck: model.ConnectionTestStep{
				Passed:  true,
				Status:  model.StatusPassed,
				Message: "Should pass",
				Error:   nil,
			},
		}},
		{TestID: test1.ID, RunID: run2.ID, Type: "EVENT_2_FAILED", Stage: model.StageTrigger, Title: "OP 2", Description: "That happened, but failed", Polling: model.PollingInfo{
			Type: model.PollingTypePeriodic,
			Periodic: &model.PeriodicPollingConfig{
				NumberSpans:      3,
				NumberIterations: 1,
			},
		}},
		{TestID: test1.ID, RunID: run2.ID, Type: "ANOTHER_EVENT", Stage: model.StageTrigger, Title: "OP 3", Description: "Clean up after error", Outputs: []model.OutputInfo{
			{LogLevel: model.LogLevelWarn, Message: "INVALID SYNTAX", OutputName: "my_output"},
		}},
	}

	for _, event := range events {
		err := db.CreateTestRunEvent(context.Background(), event)
		require.NoError(t, err)
	}

	events, err := db.GetTestRunEvents(context.Background(), test1.ID, run1.ID)
	require.NoError(t, err)

	assert.Len(t, events, 2)
	assert.LessOrEqual(t, events[0].CreatedAt, events[1].CreatedAt)
	assert.Equal(t, "OP 1", events[0].Title)
	assert.Equal(t, "OP 2", events[1].Title)

	eventsFromRun2, err := db.GetTestRunEvents(context.Background(), test1.ID, run2.ID)
	require.NoError(t, err)

	assert.Len(t, eventsFromRun2, 3)
	assert.LessOrEqual(t, eventsFromRun2[0].CreatedAt, eventsFromRun2[1].CreatedAt)
	assert.LessOrEqual(t, eventsFromRun2[1].CreatedAt, eventsFromRun2[2].CreatedAt)

	// assert events from run 2 have fields that were stored as JSON
	// data store connection
	assert.Equal(t, true, eventsFromRun2[0].DataStoreConnection.PortCheck.Passed)
	assert.Equal(t, model.StatusPassed, eventsFromRun2[0].DataStoreConnection.PortCheck.Status)
	assert.Equal(t, "Should pass", eventsFromRun2[0].DataStoreConnection.PortCheck.Message)
	assert.Nil(t, eventsFromRun2[0].DataStoreConnection.PortCheck.Error)

	// polling
	assert.Equal(t, model.PollingTypePeriodic, eventsFromRun2[1].Polling.Type)
	assert.Equal(t, 3, eventsFromRun2[1].Polling.Periodic.NumberSpans)
	assert.Equal(t, 1, eventsFromRun2[1].Polling.Periodic.NumberIterations)

	// outputs
	assert.Len(t, eventsFromRun2[2].Outputs, 1)
	assert.Equal(t, model.LogLevelWarn, eventsFromRun2[2].Outputs[0].LogLevel)
	assert.Equal(t, "INVALID SYNTAX", eventsFromRun2[2].Outputs[0].Message)
	assert.Equal(t, "my_output", eventsFromRun2[2].Outputs[0].OutputName)
}
