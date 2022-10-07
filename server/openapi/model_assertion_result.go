/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type AssertionResult struct {
	Assertion string `json:"assertion,omitempty"`

	AllPassed bool `json:"allPassed,omitempty"`

	SpanResults []AssertionSpanResult `json:"spanResults,omitempty"`
}

// AssertAssertionResultRequired checks if the required fields are not zero-ed
func AssertAssertionResultRequired(obj AssertionResult) error {
	for _, el := range obj.SpanResults {
		if err := AssertAssertionSpanResultRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertRecurseAssertionResultRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of AssertionResult (e.g. [][]AssertionResult), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseAssertionResultRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aAssertionResult, ok := obj.(AssertionResult)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertAssertionResultRequired(aAssertionResult)
	})
}
