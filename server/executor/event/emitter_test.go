package event_test

import (
	"context"
	"testing"

	"github.com/kubeshop/tracetest/server/executor/event"
	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/stretchr/testify/require"
)

func TestEventEmitter(t *testing.T) {
	repository := getTestRunEventRepositoryMock(t)
	subscriptionManager := subscription.NewManager()

	eventEmitter := event.NewEmitter(repository, subscriptionManager)

	run := model.NewRun()

	test := model.Test{
		ID: id.ID("some-test"),
		ServiceUnderTest: model.Trigger{
			Type: model.TriggerTypeHTTP,
		},
	}

	event := model.TestRunEvent{
		TestID:      test.ID,
		RunID:       run.ID,
		Type:        "EVENT_1",
		Stage:       model.StageTrigger,
		Title:       "OP 1",
		Description: "This happened",
	}

	err := eventEmitter.Emit(context.TODO(), event)
	require.NoError(t, err)

	//TODO: check if event was persisted
	//TODO: check if event was sent to subscribers
}

// TestRunEventRepository
type testRunEventRepositoryMock struct {
	testdb.MockRepository
	events []model.TestRunEvent
	// ...
}

func (m testRunEventRepositoryMock) CreateTestRunEvent(ctx context.Context, event model.TestRunEvent) error {
	m.events = append(m.events, event)
	return nil
}

func getTestRunEventRepositoryMock(t *testing.T) model.Repository {
	t.Helper()

	mock := new(testRunEventRepositoryMock)
	mock.T = t
	mock.Test(t)

	mock.events = []model.TestRunEvent{}

	return mock
}
