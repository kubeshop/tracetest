/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type EnvironmentValue struct {
	Key string `json:"key,omitempty"`

	Value string `json:"value,omitempty"`
}

// AssertEnvironmentValueRequired checks if the required fields are not zero-ed
func AssertEnvironmentValueRequired(obj EnvironmentValue) error {
	return nil
}

// AssertRecurseEnvironmentValueRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of EnvironmentValue (e.g. [][]EnvironmentValue), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseEnvironmentValueRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aEnvironmentValue, ok := obj.(EnvironmentValue)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertEnvironmentValueRequired(aEnvironmentValue)
	})
}
