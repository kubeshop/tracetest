package events

import (
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/id"
	"github.com/kubeshop/tracetest/server/model"
)

func TriggerCreatedInfo(testID id.ID, runID int) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrigger,
		Type:                "CREATED_INFO",
		Title:               "Trigger Run has been created",
		Description:         "Trigger Run has been created",
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TriggerResolveError(testID id.ID, runID int, err error) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrigger,
		Type:                "RESOLVE_ERROR",
		Title:               "Resolving trigger details failed",
		Description:         fmt.Sprintf("Resolving trigger details failed: %s", err.Error()),
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TriggerResolveSuccess(testID id.ID, runID int) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrigger,
		Type:                "RESOLVE_SUCCESS",
		Title:               "Successful resolving of trigger details",
		Description:         "Successful resolving of trigger details",
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TriggerResolveStart(testID id.ID, runID int) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrigger,
		Type:                "RESOLVE_START",
		Title:               "Resolving trigger details based on environment variables",
		Description:         "Resolving trigger details based on environment variables",
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TriggerExecutionStart(testID id.ID, runID int) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrigger,
		Type:                "EXECUTION_START",
		Title:               "Initial trigger execution",
		Description:         "Initial trigger execution",
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TriggerExecutionSuccess(testID id.ID, runID int) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrigger,
		Type:                "EXECUTION_SUCCESS",
		Title:               "Successful trigger execution",
		Description:         "Successful trigger execution",
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TriggerExecutionError(testID id.ID, runID int, err error) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrigger,
		Type:                "EXECUTION_ERROR",
		Title:               "Failed to trigger execution",
		Description:         fmt.Sprintf("Failed to trigger execution: %s", err.Error()),
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TriggerHTTPUnreachableHostError(testID id.ID, runID int, err error) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrigger,
		Type:                "HTTP_UNREACHABLE_HOST_ERROR",
		Title:               "Tracetest could not reach the defined host in the trigger",
		Description:         fmt.Sprintf("Tracetest could not reach the defined host in the trigger: %s", err.Error()),
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TriggerDockerComposeHostMismatchError(testID id.ID, runID int) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrigger,
		Type:                "DOCKER_COMPOSE_HOST_MISMATCH_ERROR",
		Title:               "Tracetest is running inside a Docker container",
		Description:         "We identified Tracetest is running inside a docker compose container, so if you are trying to access your local host machine please use the host.docker.internal hostname. For more information, see https://docs.docker.com/docker-for-mac/networking/#use-cases-and-workarounds",
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TriggergRPCUnreachableHostError(testID id.ID, runID int, err error) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrigger,
		Type:                "GRPC_UNREACHABLE_HOST_ERROR",
		Title:               "Tracetest could not reach the defined host in the trigger",
		Description:         fmt.Sprintf("Tracetest could not reach the defined host in the trigger: %s", err.Error()),
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TraceDataStoreConnectionInfo(testID id.ID, runID int, connectionResult model.ConnectionResult) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrace,
		Type:                "DATA_STORE_CONNECTION_INFO",
		Title:               "A Data store connection request has been executed",
		Description:         "A Data store connection request has been executed,test connection result information",
		CreatedAt:           time.Now(),
		DataStoreConnection: connectionResult,
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TracePollingStart(testID id.ID, runID int) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrace,
		Type:                "POLLING_START",
		Title:               "Starting the trace polling process",
		Description:         "Starting the trace polling process",
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling: model.PollingInfo{
			Type:       model.PollingTypePeriodic,
			IsComplete: false,
			Periodic: &model.PeriodicPollingConfig{
				NumberSpans:      0,
				NumberIterations: 0,
			},
		},
		Outputs: []model.OutputInfo{},
	}
}

func TracePollingIterationInfo(testID id.ID, runID, numberOfSpans, iteration int, isComplete bool) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrace,
		Type:                "POLLING_ITERATION_INFO",
		Title:               "A polling iteration has been executed",
		Description:         fmt.Sprintf("A polling iteration has been executed, %d spans collected - iteration %d", numberOfSpans, iteration),
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling: model.PollingInfo{
			Type:       model.PollingTypePeriodic,
			IsComplete: isComplete,
			Periodic: &model.PeriodicPollingConfig{
				NumberSpans:      numberOfSpans,
				NumberIterations: iteration,
			},
		},
		Outputs: []model.OutputInfo{},
	}
}

func TracePollingSuccess(testID id.ID, runID int) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrace,
		Type:                "POLLING_SUCCESS",
		Title:               "The polling strategy has succeeded in fetching the trace from the Data Store",
		Description:         "The polling strategy has succeeded in fetching the trace from the Data Store",
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TracePollingError(testID id.ID, runID int, reason string, err error) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrace,
		Type:                "POLLING_ERROR",
		Title:               "The polling strategy has failed to fetch the trace",
		Description:         fmt.Sprintf("The polling strategy has failed to fetch the trace. Reason: %s Error: %s", reason, err.Error()),
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TraceFetchingStart(testID id.ID, runID int) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrace,
		Type:                "FETCHING_START",
		Title:               "Starting the trace fetching process",
		Description:         "Starting the trace fetching process",
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TraceFetchingSuccess(testID id.ID, runID int) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrace,
		Type:                "FETCHING_SUCCESS",
		Title:               "The trace was successfully processed by the backend",
		Description:         "The trace was successfully processed by the backend",
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TraceFetchingError(testID id.ID, runID int, connectionResult model.ConnectionResult, err error) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrace,
		Type:                "FETCHING_ERROR",
		Title:               "The trace was not able to be fetched",
		Description:         fmt.Sprintf("The trace was not able to be fetched from the DataStore. Error: %s", err),
		CreatedAt:           time.Now(),
		DataStoreConnection: connectionResult,
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TraceStoppedInfo(testID id.ID, runID int) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrace,
		Type:                "STOPPED_INFO",
		Title:               "The test run was stopped during its execution",
		Description:         "The test run was stopped during its execution",
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TestOutputGenerationWarning(testID id.ID, runID int, output string) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTest,
		Type:                "OUTPUT_GENERATION_WARNING",
		Title:               fmt.Sprintf("Output %s not be generated", output),
		Description:         fmt.Sprintf("The value for output %s could not be generated", output),
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TestSpecsRunSuccess(testID id.ID, runID int) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTest,
		Type:                "TEST_SPECS_RUN_SUCCESS",
		Title:               "Test Specs were successfully executed",
		Description:         "Test Specs were successfully executed",
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TestSpecsRunError(testID id.ID, runID int, err error) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTest,
		Type:                "TEST_SPECS_RUN_ERROR",
		Title:               "Test specs execution error",
		Description:         fmt.Sprintf("An error happened when trying to run test specs. Error: %s", err.Error()),
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TestSpecsRunPersistenceError(testID id.ID, runID int, err error) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTest,
		Type:                "TEST_SPECS_RUN_PERSISTENCE_ERROR",
		Title:               "Test Specs persistence error",
		Description:         fmt.Sprintf("Test specs were succesfully executed, however an error happened when trying to persist them. Error: %s", err.Error()),
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TestSpecsRunStart(testID id.ID, runID int) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTest,
		Type:                "TEST_SPECS_RUN_START",
		Title:               "Test specs execution start",
		Description:         "Test specs execution start",
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TestSpecsAssertionError(testID id.ID, runID int, err error, spanID string, assertion string) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTest,
		Type:                "TEST_SPECS_ASSERTION_ERROR",
		Title:               "Assertion execution failed",
		Description:         fmt.Sprintf(`Assertion '%s' returned an error on span %s: %s`, assertion, spanID, err.Error()),
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}
