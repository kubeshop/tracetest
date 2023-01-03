/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type Variables struct {
	Environment []EnvironmentValue `json:"environment,omitempty"`

	Variables []string `json:"variables,omitempty"`

	Missing []string `json:"missing,omitempty"`
}

// AssertVariablesRequired checks if the required fields are not zero-ed
func AssertVariablesRequired(obj Variables) error {
	for _, el := range obj.Environment {
		if err := AssertEnvironmentValueRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertRecurseVariablesRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of Variables (e.g. [][]Variables), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseVariablesRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aVariables, ok := obj.(Variables)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertVariablesRequired(aVariables)
	})
}
