/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type TestSpecsSpecs struct {
	Name *string `json:"name,omitempty"`

	Selector Selector `json:"selector,omitempty"`

	Assertions []Assertion `json:"assertions,omitempty"`
}

// AssertTestSpecsSpecsRequired checks if the required fields are not zero-ed
func AssertTestSpecsSpecsRequired(obj TestSpecsSpecs) error {
	if err := AssertSelectorRequired(obj.Selector); err != nil {
		return err
	}
	for _, el := range obj.Assertions {
		if err := AssertAssertionRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertRecurseTestSpecsSpecsRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of TestSpecsSpecs (e.g. [][]TestSpecsSpecs), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseTestSpecsSpecsRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aTestSpecsSpecs, ok := obj.(TestSpecsSpecs)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertTestSpecsSpecsRequired(aTestSpecsSpecs)
	})
}
