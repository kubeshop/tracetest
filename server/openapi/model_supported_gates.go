/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type SupportedGates string

// List of SupportedGates
const (
	ANALYZER_SCORE SupportedGates = "analyzer-score"
	ANALYZER_RULES SupportedGates = "analyzer-rules"
	TEST_SPECS     SupportedGates = "test-specs"
)

// AssertSupportedGatesRequired checks if the required fields are not zero-ed
func AssertSupportedGatesRequired(obj SupportedGates) error {
	return nil
}

// AssertRecurseSupportedGatesRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of SupportedGates (e.g. [][]SupportedGates), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseSupportedGatesRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aSupportedGates, ok := obj.(SupportedGates)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertSupportedGatesRequired(aSupportedGates)
	})
}
