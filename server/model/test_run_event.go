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
	TestRunEventType  string
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
	TriggerEventType_CreatedInfo                    TestRunEventType = "CREATED_INFO"
	TriggerEventType_ResolveError                   TestRunEventType = "RESOLVE_ERROR"
	TriggerEventType_ResolveSuccess                 TestRunEventType = "RESOLVE_SUCCESS"
	TriggerEventType_ResolveStart                   TestRunEventType = "RESOLVE_START"
	TriggerEventType_ExecutionStart                 TestRunEventType = "EXECUTION_START"
	TriggerEventType_ExecutionSuccess               TestRunEventType = "EXECUTION_SUCCESS"
	TriggerEventType_HTTPUnreachableHostError       TestRunEventType = "HTTP_UNREACHABLE_HOST_ERROR"
	TriggerEventType_DockerComposeHostMismatchError TestRunEventType = "DOCKER_COMPOSE_HOST_MISMATCH_ERROR"
	TriggerEventType_gRPCUnreachableHostError       TestRunEventType = "GRPC_UNREACHABLE_HOST_ERROR"

	TraceEventType_FetchingStart           TestRunEventType = "FETCHING_START"
	TraceEventType_QueuedInfo              TestRunEventType = "QUEUED_INFO"
	TraceEventType_DataStoreConnectionInfo TestRunEventType = "DATA_STORE_CONNECTION_INFO"
	TraceEventType_PollingStart            TestRunEventType = "POLLING_START"
	TraceEventType_PollingIterationInfo    TestRunEventType = "POLLING_ITERATION_INFO"
	TraceEventType_PollingSuccess          TestRunEventType = "POLLING_SUCCESS"
	TraceEventType_PollingError            TestRunEventType = "POLLING_ERROR"
	TraceEventType_FetchingSuccess         TestRunEventType = "FETCHING_SUCCESS"
	TraceEventType_FetchingError           TestRunEventType = "FETCHING_ERROR"
	TraceEventType_StoppedInfo             TestRunEventType = "STOPPED_INFO"

	TestEventType_OutputGenerationWarning TestRunEventType = "OUTPUT_GENERATION_WARNING"
	TestEventType_ResolveStart            TestRunEventType = "RESOLVE_START"
	TestEventType_ResolveSuccess          TestRunEventType = "RESOLVE_SUCCESS"
	TestEventType_ResolveError            TestRunEventType = "RESOLVE_ERROR"
	TestEventType_TestSpecsRunSuccess     TestRunEventType = "TEST_SPECS_RUN_SUCCESS"
	TestEventType_TestSpecsRunError       TestRunEventType = "TEST_SPECS_RUN_ERROR"
	TestEventType_TestSpecsRunStart       TestRunEventType = "TEST_SPECS_RUN_START"
)

var (
	eventTable map[TestRunEventStage]map[TestRunEventType]TestRunEvent = map[TestRunEventStage]map[TestRunEventType]TestRunEvent{
		StageTrigger: map[TestRunEventType]TestRunEvent{
			TriggerEventType_CreatedInfo:                    TestRunEvent{Stage: StageTrigger, Type: TriggerEventType_CreatedInfo, Title: "Trigger Run has been created", Description: "Trigger Run has been created"},
			TriggerEventType_ResolveError:                   TestRunEvent{Stage: StageTrigger, Type: TriggerEventType_ResolveError, Title: "Resolving trigger details failed", Description: "Resolving trigger details failed"},
			TriggerEventType_ResolveSuccess:                 TestRunEvent{Stage: StageTrigger, Type: TriggerEventType_ResolveSuccess, Title: "Successful resolving of trigger details", Description: "Successful resolving of trigger details"},
			TriggerEventType_ResolveStart:                   TestRunEvent{Stage: StageTrigger, Type: TriggerEventType_ResolveStart, Title: "Resolving trigger details based on environment variables", Description: "Resolving trigger details based on environment variables"},
			TriggerEventType_ExecutionStart:                 TestRunEvent{Stage: StageTrigger, Type: TriggerEventType_ExecutionStart, Title: "Initial trigger execution", Description: "Initial trigger execution"},
			TriggerEventType_ExecutionSuccess:               TestRunEvent{Stage: StageTrigger, Type: TriggerEventType_ExecutionSuccess, Title: "Successful trigger execution", Description: "Successful trigger execution"},
			TriggerEventType_HTTPUnreachableHostError:       TestRunEvent{Stage: StageTrigger, Type: TriggerEventType_HTTPUnreachableHostError, Title: "Tracetest could not reach the defined host in the trigger", Description: "Tracetest could not reach the defined host in the trigger"},
			TriggerEventType_DockerComposeHostMismatchError: TestRunEvent{Stage: StageTrigger, Type: TriggerEventType_DockerComposeHostMismatchError, Title: "Tracetest is running inside a Docker container", Description: "We identified Tracetest is running inside a docker compose container, so if you are trying to access your local host machine please use the host.docker.internal hostname. For more information, see https://docs.docker.com/docker-for-mac/networking/#use-cases-and-workarounds"},
			TriggerEventType_gRPCUnreachableHostError:       TestRunEvent{Stage: StageTrigger, Type: TriggerEventType_gRPCUnreachableHostError, Title: "Tracetest could not reach the defined host in the trigger", Description: "Tracetest could not reach the defined host in the trigger"},
		},
		StageTrace: map[TestRunEventType]TestRunEvent{
			TraceEventType_FetchingStart:           TestRunEvent{Stage: StageTrace, Type: TraceEventType_FetchingStart, Title: "Starting the trace fetching process", Description: "Starting the trace fetching process"},
			TraceEventType_QueuedInfo:              TestRunEvent{Stage: StageTrace, Type: TraceEventType_QueuedInfo, Title: "Trace Run has been queued to start the fetching process", Description: "Trace Run has been queued to start the fetching process"},
			TraceEventType_DataStoreConnectionInfo: TestRunEvent{Stage: StageTrace, Type: TraceEventType_DataStoreConnectionInfo, Title: "A Data store connection request has been executed,test connection result information", Description: "A Data store connection request has been executed,test connection result information"},
			TraceEventType_PollingStart:            TestRunEvent{Stage: StageTrace, Type: TraceEventType_PollingStart, Title: "Starting the trace polling process", Description: "Starting the trace polling process"},
			TraceEventType_PollingIterationInfo:    TestRunEvent{Stage: StageTrace, Type: TraceEventType_PollingIterationInfo, Title: "A polling iteration has been executed", Description: "A polling iteration has been executed, {{NUMBER_OF_SPANS}} spans - iteration {{ITERATION_NUMBER}} - reason of next iteration: {{ITERATION_REASON}}"},
			TraceEventType_PollingSuccess:          TestRunEvent{Stage: StageTrace, Type: TraceEventType_PollingSuccess, Title: "The polling strategy has succeeded in fetching the trace from the Data Store", Description: "The polling strategy has succeeded in fetching the trace from the Data Store"},
			TraceEventType_PollingError:            TestRunEvent{Stage: StageTrace, Type: TraceEventType_PollingError, Title: "The polling strategy has failed to fetch the trace", Description: "The polling strategy has failed to fetch the trace"},
			TraceEventType_FetchingSuccess:         TestRunEvent{Stage: StageTrace, Type: TraceEventType_FetchingSuccess, Title: "The trace was successfully processed by the backend", Description: "The trace was successfully processed by the backend"},
			TraceEventType_FetchingError:           TestRunEvent{Stage: StageTrace, Type: TraceEventType_FetchingError, Title: "The trace was not able to be fetched", Description: "The trace was not able to be fetched"},
			TraceEventType_StoppedInfo:             TestRunEvent{Stage: StageTrace, Type: TraceEventType_StoppedInfo, Title: "The test run was stopped during its execution", Description: "The test run was stopped during its execution"},
		},
		StageTest: map[TestRunEventType]TestRunEvent{
			TestEventType_OutputGenerationWarning: TestRunEvent{Stage: StageTest, Type: TestEventType_OutputGenerationWarning, Title: "Output {{OUTPUT_NAME}} not be generated", Description: "The value for output {{OUTPUT_NAME}} could not be generated"},
			TestEventType_ResolveStart:            TestRunEvent{Stage: StageTest, Type: TestEventType_ResolveStart, Title: "Resolving test specs details start", Description: "Resolving test specs details start"},
			TestEventType_ResolveSuccess:          TestRunEvent{Stage: StageTest, Type: TestEventType_ResolveSuccess, Title: "Resolving test specs details success", Description: "Resolving test specs details success"},
			TestEventType_ResolveError:            TestRunEvent{Stage: StageTest, Type: TestEventType_ResolveError, Title: "An error ocurred while parsing the test specs", Description: "An error ocurred while parsing the test specs"},
			TestEventType_TestSpecsRunSuccess:     TestRunEvent{Stage: StageTest, Type: TestEventType_TestSpecsRunSuccess, Title: "Test Specs were successfully executed", Description: "Test Specs were successfully executed"},
			TestEventType_TestSpecsRunError:       TestRunEvent{Stage: StageTest, Type: TestEventType_TestSpecsRunError, Title: "Test specs execution error", Description: "Test specs execution error"},
			TestEventType_TestSpecsRunStart:       TestRunEvent{Stage: StageTest, Type: TestEventType_TestSpecsRunStart, Title: "Test specs execution start", Description: "Test specs execution start"},
		},
	}
)

type TestRunEvent struct {
	ID                  int64
	Type                TestRunEventType
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

func NewTestRunEvent(stage TestRunEventStage, eventType TestRunEventType, testID id.ID, runID int) TestRunEvent {
	return NewTestRunEventWithArgs(stage, eventType, testID, runID, map[string]string{})
}

func NewTestRunEventWithArgs(stage TestRunEventStage, eventType TestRunEventType, testID id.ID, runID int, titleAndDescriptionArgs map[string]string) TestRunEvent {
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
