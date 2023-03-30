package model_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/stretchr/testify/assert"
)

func TestTestRunEvent_ResourceID(t *testing.T) {
	testID := id.NewRandGenerator().ID()
	runID := 1

	event := model.TestRunEvent{TestID: testID, RunID: runID}

	assert.Equal(t, event.ResourceID(), fmt.Sprintf("test/%s/run/%d/event", testID, runID))
}

func TestNewTestRunEvent_CorrectEvent(t *testing.T) {
	testID := id.NewRandGenerator().ID()
	runID := 1

	event := model.NewTestRunEvent(model.StageTest, "TEST_SPECS_RUN_START", testID, runID)

	assert.NotEqual(t, event.ID, 0)
	assert.LessOrEqual(t, event.CreatedAt, time.Now())

	assert.Equal(t, event.TestID, testID)
	assert.Equal(t, event.RunID, runID)

	assert.NotEmpty(t, event.Title)
	assert.NotEmpty(t, event.Description)
}

func TestNewTestRunEventWithArgs_CorrectEvent(t *testing.T) {
	testID := id.NewRandGenerator().ID()
	runID := 1

	event := model.NewTestRunEventWithArgs(model.StageTrace, "POLLING_ITERATION_INFO", testID, runID, map[string]string{
		"NUMBER_OF_SPANS":  "3",
		"ITERATION_NUMBER": "2",
		"ITERATION_REASON": "More spans found",
	})

	assert.NotEqual(t, event.ID, 0)
	assert.LessOrEqual(t, event.CreatedAt, time.Now())

	assert.Equal(t, event.TestID, testID)
	assert.Equal(t, event.RunID, runID)

	assert.Equal(t, event.Description, "A polling iteration has been executed, 3 spans - iteration 2 - reason of next iteration: More spans found")
}

func TestNewTestRunEvent_BaseEventImmutable(t *testing.T) {
	testID := id.NewRandGenerator().ID()
	runID := 1

	event := model.NewTestRunEvent(model.StageTest, "TEST_SPECS_RUN_START", testID, runID)
	event.Title = "some title"
	event.Description = "some description"

	anotherEvent := model.NewTestRunEvent(model.StageTest, "TEST_SPECS_RUN_START", testID, runID)

	assert.NotEqual(t, event.Title, anotherEvent.Title)
	assert.NotEqual(t, event.Description, anotherEvent.Description)
}

func TestNewTestRunEvent_EventNotFound(t *testing.T) {
	testID := id.NewRandGenerator().ID()
	runID := 1

	event := model.NewTestRunEvent(model.StageTest, "MY_UNKNOWN_TYPE", testID, runID)

	assert.Empty(t, event.ID)
	assert.Empty(t, event.TestID.String())
	assert.Empty(t, event.RunID)
	assert.Empty(t, event.Title)
	assert.Empty(t, event.Description)
}
