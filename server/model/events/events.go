package events

import (
	"fmt"
	"time"

	"github.com/kubeshop/tracetest/server/model"
	"github.com/kubeshop/tracetest/server/pkg/id"
)

func TriggerCreatedInfo(testID id.ID, runID int) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrigger,
		Type:                "CREATED_INFO",
		Title:               "Trigger created",
		Description:         "The trigger run has been created",
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
		Description:         fmt.Sprintf("The resolution of trigger details has failed. Error: %s", err.Error()),
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
		Title:               "Resolving trigger details succeeded",
		Description:         "The resolution of trigger details was executed successfully",
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
		Title:               "Resolving trigger details started",
		Description:         "The resolution of trigger details based on variables has started",
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
		Title:               "Trigger execution started",
		Description:         "The execution of the trigger has started",
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
		Title:               "Trigger execution succeeded",
		Description:         "The execution of the trigger was performed successfully",
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
		Title:               "Trigger execution failed",
		Description:         fmt.Sprintf("The execution of the trigger has failed. Error: %s", err.Error()),
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
		Title:               "Unreachable host in the trigger",
		Description:         fmt.Sprintf("Tracetest could not reach the defined host in the trigger. Error: %s", err.Error()),
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
		Title:               "Docker compose mismatch error",
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
		Title:               "Unreachable host in the trigger",
		Description:         fmt.Sprintf("Tracetest could not reach the defined host in the trigger. Error: %s", err.Error()),
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
		Title:               "Data store test connection executed",
		Description:         "A data store test connection has been executed with the following results",
		CreatedAt:           time.Now(),
		DataStoreConnection: connectionResult,
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TracePollingStart(testID id.ID, runID int, dsType, endpoints string) model.TestRunEvent {
	endpointsDescription := ""
	if endpoints != "" {
		endpointsDescription = fmt.Sprintf(" with the following endpoints: %s", endpoints)
	}
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrace,
		Type:                "POLLING_START",
		Title:               "Trace polling started",
		Description:         fmt.Sprintf("The trace polling process has started using %s %s", dsType, endpointsDescription),
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

func TracePollingIterationInfo(testID id.ID, runID, numberOfSpans, iteration int, isComplete bool, reason string) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrace,
		Type:                "POLLING_ITERATION_INFO",
		Title:               "Trace polling iteration executed",
		Description:         fmt.Sprintf("A trace polling iteration has been executed. Reason: %s", reason),
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

func TracePollingSuccess(testID id.ID, runID int, reason string) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrace,
		Type:                "POLLING_SUCCESS",
		Title:               "Trace polling strategy succeeded",
		Description:         fmt.Sprintf("The polling strategy has succeeded in fetching the trace from the data store. Reason: %s", reason),
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
		Title:               "Trace polling strategy failed",
		Description:         fmt.Sprintf("The polling strategy has failed to fetch the trace. Reason: %s - Error: %s", reason, err.Error()),
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
		Title:               "Trace fetching started",
		Description:         "The trace fetching process has started",
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
		Title:               "Trace fetching succeeded",
		Description:         "The trace fetching process was performed successfully",
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TraceFetchingError(testID id.ID, runID int, err error) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrace,
		Type:                "FETCHING_ERROR",
		Title:               "Trace fetching failed",
		Description:         fmt.Sprintf("The trace was not able to be fetched from the data store. Error: %s", err),
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
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
		Title:               "Test run stopped",
		Description:         "The test run was stopped during its execution",
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TestOutputGenerationWarning(testID id.ID, runID int, err error, output string) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTest,
		Type:                "OUTPUT_GENERATION_WARNING",
		Title:               fmt.Sprintf(`Output '%s' not generated`, output),
		Description:         fmt.Sprintf(`The output '%s' returned an error. Error: %s`, output, err.Error()),
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
		Title:               "Test specs succeeded",
		Description:         "The execution of the test specs were performed successfully",
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
		Title:               "Test specs failed",
		Description:         fmt.Sprintf("The execution of the test specs has failed. Error: %s", err.Error()),
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
		Title:               "Test specs persistence error",
		Description:         fmt.Sprintf("The execution of the test specs were performed successfully, however an error happened when trying to persist them. Error: %s", err.Error()),
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
		Title:               "Test specs started",
		Description:         "The execution of the test specs has started",
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TestSpecsAssertionWarning(testID id.ID, runID int, err error, spanID string, assertion string) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTest,
		Type:                "TEST_SPECS_ASSERTION_WARNING",
		Title:               fmt.Sprintf(`Assertion '%s' failed`, assertion),
		Description:         fmt.Sprintf(`The assertion '%s' returned an error on span %s. Error: %s`, assertion, spanID, err.Error()),
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TraceOtlpServerReceivedSpans(testID id.ID, runID, spanCount int, requestType string) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrace,
		Type:                "OTLP_SERVER_RECEIVED_SPANS",
		Title:               fmt.Sprintf("%s OTLP server endpoint received spans", requestType),
		Description:         fmt.Sprintf("The Tracetest %s OTLP endpoint server received %d spans", requestType, spanCount),
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TraceLinterStart(testID id.ID, runID int) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrace,
		Type:                "TRACE_LINTER_START",
		Title:               "Trace linter started",
		Description:         "The trace linter process has started",
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TraceLinterSkip(testID id.ID, runID int) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrace,
		Type:                "TRACE_LINTER_SKIPPED",
		Title:               "Trace linter skipped",
		Description:         "The trace linter process has been skipped",
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TraceLinterSuccess(testID id.ID, runID int) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrace,
		Type:                "TRACE_LINTER_SUCCESS",
		Title:               "Trace linter succeeded",
		Description:         "The trace linter process was performed successfully",
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}

func TraceLinterError(testID id.ID, runID int, err error) model.TestRunEvent {
	return model.TestRunEvent{
		TestID:              testID,
		RunID:               runID,
		Stage:               model.StageTrace,
		Type:                "TRACE_LINTER_ERROR",
		Title:               "Trace linter error",
		Description:         fmt.Sprintf("The trace linter encountered fatal errors. Error: %s", err),
		CreatedAt:           time.Now(),
		DataStoreConnection: model.ConnectionResult{},
		Polling:             model.PollingInfo{},
		Outputs:             []model.OutputInfo{},
	}
}
