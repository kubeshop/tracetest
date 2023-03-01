/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type ConnectionTestStep struct {
	Passed bool `json:"passed,omitempty"`

	Status string `json:"status,omitempty"`

	Message string `json:"message,omitempty"`

	Error string `json:"error,omitempty"`
}

// AssertConnectionTestStepRequired checks if the required fields are not zero-ed
func AssertConnectionTestStepRequired(obj ConnectionTestStep) error {
	return nil
}

// AssertRecurseConnectionTestStepRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of ConnectionTestStep (e.g. [][]ConnectionTestStep), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseConnectionTestStepRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aConnectionTestStep, ok := obj.(ConnectionTestStep)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertConnectionTestStepRequired(aConnectionTestStep)
	})
}
