/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type AssertionResults struct {
	AllPassed bool `json:"allPassed,omitempty"`

	Results []AssertionResultsResults `json:"results,omitempty"`
}

// AssertAssertionResultsRequired checks if the required fields are not zero-ed
func AssertAssertionResultsRequired(obj AssertionResults) error {
	for _, el := range obj.Results {
		if err := AssertAssertionResultsResultsRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertRecurseAssertionResultsRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of AssertionResults (e.g. [][]AssertionResults), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseAssertionResultsRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aAssertionResults, ok := obj.(AssertionResults)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertAssertionResultsRequired(aAssertionResults)
	})
}
