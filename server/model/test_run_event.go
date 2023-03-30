package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/kubeshop/tracetest/server/id"
)

type (
	Protocol          string
	Status            string
	TestRunEventStage string
	PollingType       string
	LogLevel          string
)

var (
	ProtocolHTTP Protocol = "http"
	ProtocolGRPC Protocol = "grpc"
)

var (
	StatusPassed  Status = "passed"
	StatusWarning Status = "warning"
	StatusFailed  Status = "failed"
)

var (
	LogLevelWarn  LogLevel = "warning"
	LogLevelError LogLevel = "error"
)

var (
	StageTrigger TestRunEventStage = "trigger"
	StageTrace   TestRunEventStage = "trace"
	StageTest    TestRunEventStage = "test"
)

var (
	PollingTypePeriodic PollingType = "periodic"
)

var (
	eventTable map[TestRunEventStage]map[string]TestRunEvent = map[TestRunEventStage]map[string]TestRunEvent{
		StageTrigger: map[string]TestRunEvent{
			"CREATED_INFO":                       TestRunEvent{Stage: StageTrigger, Type: "CREATED_INFO", Title: "Trigger Run has been created", Description: "Trigger Run has been created"},
			"RESOLVE_ERROR":                      TestRunEvent{Stage: StageTrigger, Type: "RESOLVE_ERROR", Title: "Resolving trigger details failed", Description: "Resolving trigger details failed"},
			"RESOLVE_SUCCESS":                    TestRunEvent{Stage: StageTrigger, Type: "RESOLVE_SUCCESS", Title: "Successful resolving of trigger details", Description: "Successful resolving of trigger details"},
			"RESOLVE_START":                      TestRunEvent{Stage: StageTrigger, Type: "RESOLVE_START", Title: "Resolving trigger details based on environment variables", Description: "Resolving trigger details based on environment variables"},
			"EXECUTION_START":                    TestRunEvent{Stage: StageTrigger, Type: "EXECUTION_START", Title: "Initial trigger execution", Description: "Initial trigger execution"},
			"EXECUTION_SUCCESS":                  TestRunEvent{Stage: StageTrigger, Type: "EXECUTION_SUCCESS", Title: "Successful trigger execution", Description: "Successful trigger execution"},
			"HTTP_UNREACHABLE_HOST_ERROR":        TestRunEvent{Stage: StageTrigger, Type: "HTTP_UNREACHABLE_HOST_ERROR", Title: "Tracetest could not reach the defined host in the trigger", Description: "Tracetest could not reach the defined host in the trigger"},
			"DOCKER_COMPOSE_HOST_MISMATCH_ERROR": TestRunEvent{Stage: StageTrigger, Type: "DOCKER_COMPOSE_HOST_MISMATCH_ERROR", Title: "Tracetest is running inside a Docker container", Description: "We identified Tracetest is running inside a docker compose container, so if you are trying to access your local host machine please use the host.docker.internal hostname. For more information, see https://docs.docker.com/docker-for-mac/networking/#use-cases-and-workarounds"},
			"GRPC_UNREACHABLE_HOST_ERROR":        TestRunEvent{Stage: StageTrigger, Type: "GRPC_UNREACHABLE_HOST_ERROR", Title: "Tracetest could not reach the defined host in the trigger", Description: "Tracetest could not reach the defined host in the trigger"},
		},
		StageTrace: map[string]TestRunEvent{
			"FETCHING_START":             TestRunEvent{Stage: StageTrace, Type: "FETCHING_START", Title: "Starting the trace fetching process", Description: "Starting the trace fetching process"},
			"QUEUED_INFO":                TestRunEvent{Stage: StageTrace, Type: "QUEUED_INFO", Title: "Trace Run has been queued to start the fetching process", Description: "Trace Run has been queued to start the fetching process"},
			"DATA_STORE_CONNECTION_INFO": TestRunEvent{Stage: StageTrace, Type: "DATA_STORE_CONNECTION_INFO", Title: "A Data store connection request has been executed,test connection result information", Description: "A Data store connection request has been executed,test connection result information"},
			"POLLING_START":              TestRunEvent{Stage: StageTrace, Type: "POLLING_START", Title: "Starting the trace polling process", Description: "Starting the trace polling process"},
			"POLLING_ITERATION_INFO":     TestRunEvent{Stage: StageTrace, Type: "POLLING_ITERATION_INFO", Title: "A polling iteration has been executed", Description: "A polling iteration has been executed, {{NUMBER_OF_SPANS}} spans - iteration {{ITERATION_NUMBER}} - reason of next iteration: {{ITERATION_REASON}}"},
			"POLLING_SUCCESS":            TestRunEvent{Stage: StageTrace, Type: "POLLING_SUCCESS", Title: "The polling strategy has succeeded in fetching the trace from the Data Store", Description: "The polling strategy has succeeded in fetching the trace from the Data Store"},
			"POLLING_ERROR":              TestRunEvent{Stage: StageTrace, Type: "POLLING_ERROR", Title: "The polling strategy has failed to fetch the trace", Description: "The polling strategy has failed to fetch the trace"},
			"FETCHING_SUCCESS":           TestRunEvent{Stage: StageTrace, Type: "FETCHING_SUCCESS", Title: "The trace was successfully processed by the backend", Description: "The trace was successfully processed by the backend"},
			"FETCHING_ERROR":             TestRunEvent{Stage: StageTrace, Type: "FETCHING_ERROR", Title: "The trace was not able to be fetched", Description: "The trace was not able to be fetched"},
			"STOPPED_INFO":               TestRunEvent{Stage: StageTrace, Type: "STOPPED_INFO", Title: "The test run was stopped during its execution", Description: "The test run was stopped during its execution"},
		},
		StageTest: map[string]TestRunEvent{
			"OUTPUT_GENERATION_WARNING": TestRunEvent{Stage: StageTest, Type: "OUTPUT_GENERATION_WARNING", Title: "Output {{OUTPUT_NAME}} not be generated", Description: "The value for output {{OUTPUT_NAME}} could not be generated"},
			"RESOLVE_START":             TestRunEvent{Stage: StageTest, Type: "RESOLVE_START", Title: "Resolving test specs details start", Description: "Resolving test specs details start"},
			"RESOLVE_SUCCESS":           TestRunEvent{Stage: StageTest, Type: "RESOLVE_SUCCESS", Title: "Resolving test specs details success", Description: "Resolving test specs details success"},
			"RESOLVE_ERROR":             TestRunEvent{Stage: StageTest, Type: "RESOLVE_ERROR", Title: "An error ocurred while parsing the test specs", Description: "An error ocurred while parsing the test specs"},
			"TEST_SPECS_RUN_SUCCESS":    TestRunEvent{Stage: StageTest, Type: "TEST_SPECS_RUN_SUCCESS", Title: "Test Specs were successfully executed", Description: "Test Specs were successfully executed"},
			"TEST_SPECS_RUN_ERROR":      TestRunEvent{Stage: StageTest, Type: "TEST_SPECS_RUN_ERROR", Title: "Test specs execution error", Description: "Test specs execution error"},
			"TEST_SPECS_RUN_START":      TestRunEvent{Stage: StageTest, Type: "TEST_SPECS_RUN_START", Title: "Test specs execution start", Description: "Test specs execution start"},
		},
	}
)

type TestRunEvent struct {
	ID                  int64
	Type                string
	Stage               TestRunEventStage
	Title               string
	Description         string
	CreatedAt           time.Time
	TestID              id.ID
	RunID               int
	DataStoreConnection ConnectionResult
	Polling             PollingInfo
	Outputs             []OutputInfo
}

func (e TestRunEvent) ResourceID() string {
	return fmt.Sprintf("test/%s/run/%d/event", e.TestID, e.RunID)
}

func (e TestRunEvent) WithDataStoreConnection(dataStoreConnection ConnectionResult) TestRunEvent {
	e.DataStoreConnection = dataStoreConnection
	return e
}

func (e TestRunEvent) WithPolling(pollingInfo PollingInfo) TestRunEvent {
	e.Polling = pollingInfo
	return e
}

func (e TestRunEvent) WithOutputs(outputs ...OutputInfo) TestRunEvent {
	e.Outputs = outputs
	return e
}

func NewTestRunEvent(stage TestRunEventStage, eventType string, testID id.ID, runID int) TestRunEvent {
	return NewTestRunEventWithArgs(stage, eventType, testID, runID, map[string]string{})
}

func NewTestRunEventWithArgs(stage TestRunEventStage, eventType string, testID id.ID, runID int, titleAndDescriptionArgs map[string]string) TestRunEvent {
	eventTypeTable, found := eventTable[stage]
	if !found {
		return TestRunEvent{}
	}

	event, found := eventTypeTable[eventType]
	if !found {
		return TestRunEvent{}
	}

	event.ID = time.Now().Unix() // TODO: think how we can generate these IDs
	event.CreatedAt = time.Now()
	event.TestID = testID
	event.RunID = runID

	if len(titleAndDescriptionArgs) > 0 {
		event.Title = replaceArgsInString(event.Title, titleAndDescriptionArgs)
		event.Description = replaceArgsInString(event.Description, titleAndDescriptionArgs)
	}

	return event
}

func replaceArgsInString(baseString string, args map[string]string) string {
	finalString := baseString

	for key, value := range args {
		keyAsPlaceholder := fmt.Sprintf("{{%s}}", key)
		finalString = strings.ReplaceAll(finalString, keyAsPlaceholder, value)
	}

	return finalString
}

type PollingInfo struct {
	Type                PollingType
	ReasonNextIteration string
	IsComplete          bool
	Periodic            *PeriodicPollingConfig
}

type PeriodicPollingConfig struct {
	NumberSpans      int
	NumberIterations int
}

type OutputInfo struct {
	LogLevel   LogLevel
	Message    string
	OutputName string
}
