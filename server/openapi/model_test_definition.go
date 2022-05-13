/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

// TestDefinition - Map using selector query as key, and an array of assertions as value
type TestDefinition struct {
	Definitions map[string][]Assertion `json:"definitions,omitempty"`
}

// AssertTestDefinitionRequired checks if the required fields are not zero-ed
func AssertTestDefinitionRequired(obj TestDefinition) error {
	return nil
}

// AssertRecurseTestDefinitionRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of TestDefinition (e.g. [][]TestDefinition), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseTestDefinitionRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aTestDefinition, ok := obj.(TestDefinition)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertTestDefinitionRequired(aTestDefinition)
	})
}
