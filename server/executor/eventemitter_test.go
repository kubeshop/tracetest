package executor_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/kubeshop/tracetest/server/executor"
	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
	"github.com/kubeshop/tracetest/server/subscription"
	"github.com/kubeshop/tracetest/server/test"
	"github.com/kubeshop/tracetest/server/test/trigger"
	"github.com/kubeshop/tracetest/server/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventEmitter_SuccessfulScenario(t *testing.T) {
	// Given I have a test run event

	run := test.NewRun()

	testObj := test.Test{
		ID: id.ID("some-test"),
		Trigger: trigger.Trigger{
			Type: trigger.TriggerTypeHTTP,
		},
	}

	testRunEvent := model.TestRunEvent{
		TestID:      testObj.ID,
		RunID:       run.ID,
		Type:        "EVENT_1",
		Stage:       model.StageTrigger,
		Title:       "OP 1",
		Description: "This happened",
	}

	// When I emit this event successfully
	repository := getTestRunEventRepositoryMock(t, false)
	subscriptionManager, subscriber := getSubscriptionManagerMock(t, testRunEvent)

	eventEmitter := executor.NewEventEmitter(repository, subscriptionManager)

	err := eventEmitter.Emit(context.Background(), testRunEvent)
	require.NoError(t, err)

	// Then I expect that it was persisted
	assert.Len(t, repository.events, 1)
	assert.Equal(t, testRunEvent.Title, repository.events[0].Title)
	assert.Equal(t, testRunEvent.Stage, repository.events[0].Stage)
	assert.Equal(t, testRunEvent.Description, repository.events[0].Description)

	// And that it was sent to subscribers
	assert.Len(t, subscriber.events, 1)
	assert.Equal(t, testRunEvent.Title, subscriber.events[0].Title)
	assert.Equal(t, testRunEvent.Stage, subscriber.events[0].Stage)
	assert.Equal(t, testRunEvent.Description, subscriber.events[0].Description)
}

func TestEventEmitter_FailedScenario(t *testing.T) {
	// Given I have a test run event

	run := test.NewRun()

	testObj := test.Test{
		ID: id.ID("some-test"),
		Trigger: trigger.Trigger{
			Type: trigger.TriggerTypeHTTP,
		},
	}

	testRunEvent := model.TestRunEvent{
		TestID:      testObj.ID,
		RunID:       run.ID,
		Type:        "EVENT_1",
		Stage:       model.StageTrigger,
		Title:       "OP 1",
		Description: "This happened",
	}

	// When I emit this event and it fails
	repository := getTestRunEventRepositoryMock(t, true)
	subscriptionManager, subscriber := getSubscriptionManagerMock(t, testRunEvent)

	eventEmitter := executor.NewEventEmitter(repository, subscriptionManager)

	err := eventEmitter.Emit(context.Background(), testRunEvent)
	require.Error(t, err)

	// Then I expect that it was not persisted
	assert.Len(t, repository.events, 0)

	// And that it was not sent to subscribers
	assert.Len(t, subscriber.events, 0)
}

// TestRunEventRepository
type testRunEventRepositoryMock struct {
	testdb.MockRepository
	events      []model.TestRunEvent
	returnError bool
	// ...
}

func (m *testRunEventRepositoryMock) CreateTestRunEvent(ctx context.Context, event model.TestRunEvent) error {
	if m.returnError {
		return errors.New("error on persistence")
	}

	m.events = append(m.events, event)
	return nil
}

func getTestRunEventRepositoryMock(t *testing.T, returnError bool) *testRunEventRepositoryMock {
	t.Helper()

	mock := new(testRunEventRepositoryMock)
	mock.T = t
	mock.Test(t)

	mock.events = []model.TestRunEvent{}
	mock.returnError = returnError

	return mock
}

// TestRunEventSubscriber
type testRunEventSubscriber struct {
	events []model.TestRunEvent
}

func (s *testRunEventSubscriber) ID() string {
	return "some-id"
}

func (s *testRunEventSubscriber) Notify(m subscription.Message) error {
	event := model.TestRunEvent{}
	err := m.DecodeContent(&event)
	if err != nil {
		panic(fmt.Errorf("cannot read testRunEvent: %w", err))
	}
	s.events = append(s.events, event)
	return nil
}

func getSubscriptionManagerMock(t *testing.T, event model.TestRunEvent) (subscription.Manager, *testRunEventSubscriber) {
	t.Helper()

	subscriptionManager := subscription.NewManager()
	subscriber := &testRunEventSubscriber{
		events: []model.TestRunEvent{},
	}

	subscriptionManager.Subscribe(event.ResourceID(), subscriber)

	return subscriptionManager, subscriber
}
