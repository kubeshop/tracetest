/*
 * TraceTest
 *
 * OpenAPI definition for TraceTest endpoint and resources
 *
 * API version: 0.2.1
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type EnvironmentResourceList struct {
	Count int32 `json:"count,omitempty"`

	Items []EnvironmentResource `json:"items,omitempty"`
}

// AssertEnvironmentResourceListRequired checks if the required fields are not zero-ed
func AssertEnvironmentResourceListRequired(obj EnvironmentResourceList) error {
	for _, el := range obj.Items {
		if err := AssertEnvironmentResourceRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertRecurseEnvironmentResourceListRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of EnvironmentResourceList (e.g. [][]EnvironmentResourceList), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseEnvironmentResourceListRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aEnvironmentResourceList, ok := obj.(EnvironmentResourceList)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertEnvironmentResourceListRequired(aEnvironmentResourceList)
	})
}
