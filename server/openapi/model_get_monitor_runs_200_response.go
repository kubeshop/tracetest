/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type GetMonitorRuns200Response struct {
	Items []MonitorRun `json:"items,omitempty"`

	Count int32 `json:"count,omitempty"`
}

// AssertGetMonitorRuns200ResponseRequired checks if the required fields are not zero-ed
func AssertGetMonitorRuns200ResponseRequired(obj GetMonitorRuns200Response) error {
	for _, el := range obj.Items {
		if err := AssertMonitorRunRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertRecurseGetMonitorRuns200ResponseRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of GetMonitorRuns200Response (e.g. [][]GetMonitorRuns200Response), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseGetMonitorRuns200ResponseRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aGetMonitorRuns200Response, ok := obj.(GetMonitorRuns200Response)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertGetMonitorRuns200ResponseRequired(aGetMonitorRuns200Response)
	})
}
