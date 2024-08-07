/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"time"
)

type MonitorRun struct {
	Id int32 `json:"id,omitempty"`

	MonitorId string `json:"monitorId,omitempty"`

	MonitorVersion int32 `json:"monitorVersion,omitempty"`

	RunGroupId string `json:"runGroupId,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`

	CompletedAt time.Time `json:"completedAt,omitempty"`

	ExecutionType string `json:"executionType,omitempty"`

	LastError string `json:"lastError,omitempty"`

	State string `json:"state,omitempty"`

	VariableSet VariableSet `json:"variableSet,omitempty"`

	Metadata map[string]string `json:"metadata,omitempty"`

	TestRunsCount int32 `json:"testRunsCount,omitempty"`

	TestSuiteRunsCount int32 `json:"testSuiteRunsCount,omitempty"`

	// list of test runs of the Monitor Run
	TestRuns []TestRun `json:"testRuns,omitempty"`

	// list of test suite runs of the Monitor Run
	TestSuiteRuns []TestSuiteRun `json:"testSuiteRuns,omitempty"`

	Alerts []AlertResult `json:"alerts,omitempty"`
}

// AssertMonitorRunRequired checks if the required fields are not zero-ed
func AssertMonitorRunRequired(obj MonitorRun) error {
	if err := AssertVariableSetRequired(obj.VariableSet); err != nil {
		return err
	}
	for _, el := range obj.TestRuns {
		if err := AssertTestRunRequired(el); err != nil {
			return err
		}
	}
	for _, el := range obj.TestSuiteRuns {
		if err := AssertTestSuiteRunRequired(el); err != nil {
			return err
		}
	}
	for _, el := range obj.Alerts {
		if err := AssertAlertResultRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertRecurseMonitorRunRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of MonitorRun (e.g. [][]MonitorRun), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseMonitorRunRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aMonitorRun, ok := obj.(MonitorRun)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertMonitorRunRequired(aMonitorRun)
	})
}
